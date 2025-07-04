package outscale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var (
	_ resource.Resource               = &resourceLoadBalancer{}
	_ resource.ResourceWithConfigure  = &resourceLoadBalancer{}
	_ resource.ResourceWithModifyPlan = &resourceLoadBalancer{}
)

type LoadBalancerModel struct {
	LoadBalancerStickyCookiePolicies types.Set      `tfsdk:"load_balancer_sticky_cookie_policies"`
	ApplicationStickyCookiePolicies  types.Set      `tfsdk:"application_sticky_cookie_policies"`
	SourceSecurityGroup              types.Set      `tfsdk:"source_security_group"`
	LoadBalancerName                 types.String   `tfsdk:"load_balancer_name"`
	LoadBalancerType                 types.String   `tfsdk:"load_balancer_type"`
	SecuredCookies                   types.Bool     `tfsdk:"secured_cookies"`
	SecurityGroups                   types.Set      `tfsdk:"security_groups"`
	SubregionNames                   types.Set      `tfsdk:"subregion_names"`
	BackendVmIds                     types.Set      `tfsdk:"backend_vm_ids"`
	HealthCheck                      types.Set      `tfsdk:"health_check"`
	BackendIps                       types.Set      `tfsdk:"backend_ips"`
	Listeners                        []Listeners    `tfsdk:"listeners"`
	AccessLog                        types.Set      `tfsdk:"access_log"`
	RequestId                        types.String   `tfsdk:"request_id"`
	PublicIp                         types.String   `tfsdk:"public_ip"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
	DnsName                          types.String   `tfsdk:"dns_name"`
	Subnets                          types.Set      `tfsdk:"subnets"`
	NetId                            types.String   `tfsdk:"net_id"`
	Tags                             []ResourceTag  `tfsdk:"tags"`
	Id                               types.String   `tfsdk:"id"`
}

type resourceLoadBalancer struct {
	Client *oscgo.APIClient
}

type Listeners struct {
	BackendPort          types.Int32  `tfsdk:"backend_port"`
	BackendProtocol      types.String `tfsdk:"backend_protocol"`
	LoadBalancerPort     types.Int32  `tfsdk:"load_balancer_port"`
	LoadBalancerProtocol types.String `tfsdk:"load_balancer_protocol"`
	ServerCertificateId  types.String `tfsdk:"server_certificate_id"`
	PolicyNames          types.Set    `tfsdk:"policy_names"`
}
type ComparableListener struct {
	BackendPort          types.Int32  `tfsdk:"backend_port"`
	BackendProtocol      types.String `tfsdk:"backend_protocol"`
	LoadBalancerPort     types.Int32  `tfsdk:"load_balancer_port"`
	LoadBalancerProtocol types.String `tfsdk:"load_balancer_protocol"`
	ServerCertificateId  types.String `tfsdk:"server_certificate_id"`
	PolicyNames          types.String `tfsdk:"policy_names"`
}
type BlockHealthCheck struct {
	CheckInterval      types.Int32  `tfsdk:"check_interval"`
	HealthyThreshold   types.Int32  `tfsdk:"healthy_threshold"`
	Path               types.String `tfsdk:"path"`
	Port               types.Int32  `tfsdk:"port"`
	Protocol           types.String `tfsdk:"protocol"`
	Timeout            types.Int32  `tfsdk:"timeout"`
	UnhealthyThreshold types.Int32  `tfsdk:"unhealthy_threshold"`
}

type BlockAccessLog struct {
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	OsuBucketName       types.String `tfsdk:"osu_bucket_name"`
	OsuBucketPrefix     types.String `tfsdk:"osu_bucket_prefix"`
	PublicationInterval types.Int32  `tfsdk:"publication_interval"`
}
type BlockSourceSG struct {
	SecurityGroupName      types.String `tfsdk:"security_group_name"`
	SecurityGroupAccountId types.String `tfsdk:"security_group_account_id"`
}
type BlockAppStickyCookiePolicy struct {
	CookieName types.String `tfsdk:"cookie_name"`
	PolicyName types.String `tfsdk:"policy_name"`
}
type BlockLBUStickyCookiePolicy struct {
	CookieExpirationPeriod types.Int32  `tfsdk:"cookie_expiration_period"`
	PolicyName             types.String `tfsdk:"policy_name"`
}

func NewResourceLoadBalancer() resource.Resource {
	return &resourceLoadBalancer{}
}

func (r *resourceLoadBalancer) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClient_fw)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
}

func (r *resourceLoadBalancer) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data LoadBalancerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, listener := range data.Listeners {
		if !listener.ServerCertificateId.IsNull() {
			protocolNeedCerticate := []string{"https", "ssl"}
			if !slices.Contains(protocolNeedCerticate, strings.ToLower(listener.BackendProtocol.ValueString())) &&
				!slices.Contains(protocolNeedCerticate, strings.ToLower(listener.LoadBalancerProtocol.ValueString())) {
				resp.Diagnostics.AddError(
					"Invalide Listener Attributes Configuration",
					"LBU Listener: server_certificate_id may be set only when protocol is 'HTTPS' or 'SSL'",
				)
				return
			}
		}
	}
}
func (r *resourceLoadBalancer) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will only unlink backend vms from load_balancer.",
		)
		return
	}

	if req.Plan.Raw.IsKnown() && req.State.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource 'outscale_load_balancer_vms' Considerations",
			"You have to apply twice or run 'terraform refesh' after the fist apply to get"+
				" the 'backend_ips' or 'backend_vm_ids' block values in load_balancer resource state",
		)
		return
	}
}

func (r *resourceLoadBalancer) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer"
}

func ListenersSchema() schema.SetNestedBlock {
	return schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"backend_port": schema.Int32Attribute{
					Required: true,
					Validators: []validator.Int32{
						int32validator.Between(1, 65535),
					},
				},
				"backend_protocol": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"HTTP", "HTTPS", "TCP", "SSL"}...),
					},
				},
				"load_balancer_port": schema.Int32Attribute{
					Required: true,
					Validators: []validator.Int32{
						int32validator.Between(1, 65535),
					},
				},
				"load_balancer_protocol": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"HTTP", "HTTPS", "TCP", "SSL"}...),
					},
				},
				"server_certificate_id": schema.StringAttribute{
					Optional: true,
				},
				"policy_names": schema.SetAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
	}
}

func SourceSGSchema() schema.SetAttribute {
	return schema.SetAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"security_group_name":       types.StringType,
				"security_group_account_id": types.StringType,
			},
		},
	}
}

func HealthCheckSchema() schema.SetAttribute {
	return schema.SetAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"check_interval":      types.Int32Type,
				"healthy_threshold":   types.Int32Type,
				"port":                types.Int32Type,
				"protocol":            types.StringType,
				"path":                types.StringType,
				"timeout":             types.Int32Type,
				"unhealthy_threshold": types.Int32Type,
			},
		},
	}
}

func AccessLogSchema() schema.SetAttribute {
	return schema.SetAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"is_enabled":           types.BoolType,
				"osu_bucket_name":      types.StringType,
				"osu_bucket_prefix":    types.StringType,
				"publication_interval": types.Int32Type,
			},
		},
	}
}

func (r *resourceLoadBalancer) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags":      TagsSchema(),
			"listeners": ListenersSchema(),
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"load_balancer_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"internet-facing", "internal"}...),
				},
			},
			"subnets": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"subregion_names": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("subnets"),
					}...),
				},
			},
			"public_ip": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"security_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"secured_cookies": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"access_log":   AccessLogSchema(),
			"health_check": HealthCheckSchema(),
			"source_security_group": schema.SetAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"security_group_name":       types.StringType,
						"security_group_account_id": types.StringType,
					},
				},
			},
			"application_sticky_cookie_policies": schema.SetAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"cookie_name": types.StringType,
						"policy_name": types.StringType,
					},
				},
			},
			"load_balancer_sticky_cookie_policies": schema.SetAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"cookie_expiration_period": types.Int32Type,
						"policy_name":              types.StringType,
					},
				},
			},
			"backend_vm_ids": schema.SetAttribute{
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
				Computed:    true,
			},
			"backend_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_name": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *resourceLoadBalancer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lbuName := data.LoadBalancerName.ValueString()
	createReq := oscgo.NewCreateLoadBalancerRequest(getListenersConfig(data.Listeners), data.LoadBalancerName.ValueString())
	if !data.LoadBalancerType.IsUnknown() {
		createReq.SetLoadBalancerType(data.LoadBalancerType.ValueString())
	}
	if !data.PublicIp.IsUnknown() {
		createReq.SetPublicIp(data.PublicIp.ValueString())
	}
	if !data.Subnets.IsUnknown() {
		subnets, diags := utils.GetSliceFromFwtypeSet(ctx, data.Subnets)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.SetSubnets(subnets)
	}
	if !data.SecurityGroups.IsUnknown() {
		securityGroups, diags := utils.GetSliceFromFwtypeSet(ctx, data.SecurityGroups)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.SetSecurityGroups(securityGroups)
	}
	if !data.SubregionNames.IsUnknown() {
		subregionNames, diags := utils.GetSliceFromFwtypeSet(ctx, data.SubregionNames)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.SetSubregionNames(subregionNames)
	}
	if tags := tagsToOSCResourceTag(data.Tags); len(tags) != 0 {
		createReq.SetTags(tags)
	}

	var createResp oscgo.CreateLoadBalancerResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.LoadBalancerApi.CreateLoadBalancer(ctx).CreateLoadBalancerRequest(*createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create load_balancer resource",
			err.Error(),
		)
		return
	}

	if data.SecuredCookies.ValueBool() {
		req := oscgo.UpdateLoadBalancerRequest{
			LoadBalancerName: lbuName,
		}
		req.SetSecuredCookies(data.SecuredCookies.ValueBool())
		err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.UpdateLoadBalancer(ctx).UpdateLoadBalancerRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to set load_balancer secured_cookies paremeter",
				err.Error(),
			)
			return
		}
	}
	data.RequestId = types.StringValue(*createResp.ResponseContext.RequestId)
	data.Id = types.StringValue(lbuName)
	data.LoadBalancerName = types.StringValue(lbuName)
	err = setLoadBalancerState(ctx, r, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set load_balancer state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLoadBalancer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data LoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := setLoadBalancerState(ctx, r, &data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set load_balancer API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLoadBalancer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var dataPlan, dataState LoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateTimeout, diags := dataPlan.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if dataPlan.SecurityGroups.IsUnknown() {
		sgReq := oscgo.UpdateLoadBalancerRequest{
			LoadBalancerName: dataPlan.LoadBalancerName.ValueString(),
		}
		sgToAdd, _, diags := utils.GetSlicesFromTypesSetForUpdating(ctx, dataState.SecurityGroups, dataPlan.SecurityGroups)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		sgReq.SetSecurityGroups(sgToAdd)
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.UpdateLoadBalancer(ctx).UpdateLoadBalancerRequest(sgReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Load Balancer security_groups",
				err.Error(),
			)
			return
		}
	}
	if !reflect.DeepEqual(dataPlan.Listeners, dataState.Listeners) {
		listenersToAdd, listenersToRemove := getListenersForUpdate(&dataPlan, &dataState)
		//if len(listenersToAdd) != 0 || len(listenersToRemove) != 0 {
		if len(listenersToRemove) > 0 {
			lbuPorts := []int32{}
			for _, listener := range listenersToRemove {
				lbuPorts = append(lbuPorts, listener.GetLoadBalancerPort())
			}

			removeReq := oscgo.DeleteLoadBalancerListenersRequest{
				LoadBalancerName:  dataPlan.LoadBalancerName.ValueString(),
				LoadBalancerPorts: lbuPorts,
			}
			err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
				_, httpResp, err := r.Client.ListenerApi.DeleteLoadBalancerListeners(
					ctx).DeleteLoadBalancerListenersRequest(removeReq).Execute()

				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})

			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to remove Load Balancer listeners",
					err.Error(),
				)
				return
			}
		}

		if len(listenersToAdd) > 0 {
			AddReq := oscgo.CreateLoadBalancerListenersRequest{
				LoadBalancerName: dataPlan.LoadBalancerName.ValueString(),
				Listeners:        listenersToAdd,
			}

			err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
				_, httpResp, err := r.Client.ListenerApi.CreateLoadBalancerListeners(ctx).CreateLoadBalancerListenersRequest(AddReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to add Load Balancer listeners",
					err.Error(),
				)
				return
			}
		}
	}
	if dataPlan.SecuredCookies.ValueBool() != dataState.SecuredCookies.ValueBool() {
		req := oscgo.UpdateLoadBalancerRequest{
			LoadBalancerName: dataPlan.LoadBalancerName.ValueString(),
		}
		req.SetSecuredCookies(dataPlan.SecuredCookies.ValueBool())
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.UpdateLoadBalancer(ctx).UpdateLoadBalancerRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update load_balancer secured_cookies paremeter",
				err.Error(),
			)
			return
		}
	}

	err := setLoadBalancerState(ctx, r, &dataPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set LBU backend vms state after updating.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &dataPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLoadBalancer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteReq := oscgo.NewDeleteLoadBalancerRequest(data.LoadBalancerName.ValueString())
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.LoadBalancerApi.DeleteLoadBalancer(ctx).DeleteLoadBalancerRequest(*deleteReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete load_balancer",
			err.Error(),
		)
		return
	}
}

func setLoadBalancerState(ctx context.Context, r *resourceLoadBalancer, data *LoadBalancerModel) error {
	lbuFilters := oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{data.LoadBalancerName.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse read timeout value. Error: %v: ", diags.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	readReq := oscgo.ReadLoadBalancersRequest{
		Filters: &lbuFilters,
	}
	var readResp oscgo.ReadLoadBalancersResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.LoadBalancerApi.ReadLoadBalancers(context.Background()).ReadLoadBalancersRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if len(readResp.GetLoadBalancers()) == 0 {
		return errors.New("Empty")
	}

	data.RequestId = types.StringValue(*readResp.ResponseContext.RequestId)
	lbu := readResp.GetLoadBalancers()[0]
	data.NetId = types.StringValue(lbu.GetNetId())
	data.DnsName = types.StringValue(lbu.GetDnsName())
	data.PublicIp = types.StringValue(lbu.GetPublicIp())
	data.SecuredCookies = types.BoolValue(lbu.GetSecuredCookies())
	data.LoadBalancerName = types.StringValue(lbu.GetLoadBalancerName())
	data.LoadBalancerType = types.StringValue(lbu.GetLoadBalancerType())
	data.Id = types.StringValue(lbu.GetLoadBalancerName())
	data.Tags = getTagsFromApiResponse(lbu.GetTags())

	if lbu.HasBackendVmIds() {
		data.BackendVmIds, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetBackendVmIds())
		if diags.HasError() {
			return fmt.Errorf("unable to set LBU backend_vm_ids: %v", diags.Errors())
		}
	}
	if lbu.HasBackendIps() {
		data.BackendIps, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetBackendIps())
		if diags.HasError() {
			return fmt.Errorf("unable to set LBU backend_ips: %v", diags.Errors())
		}
	}
	data.SubregionNames, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetSubregionNames())
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU SubregionNames: %v", diags.Errors())
	}
	data.SecurityGroups, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetSecurityGroups())
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU SecurityGroups: %v", diags.Errors())
	}
	data.Subnets, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetSubnets())
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU Subnets: %v", diags.Errors())
	}
	data.SourceSecurityGroup, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockSourceSG{})}, getSourceSGFromApiResponse(lbu.GetSourceSecurityGroup()))
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU SourceSecurityGroup: %v", diags.Errors())
	}
	data.HealthCheck, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockHealthCheck{})}, getHealthCheckFromApiResponse(lbu.GetHealthCheck()))
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU HealthCheck: %v", diags.Errors())
	}

	if lbu.HasAccessLog() {
		data.AccessLog, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockAccessLog{})}, getAccessLogFromApiResponse(lbu.GetAccessLog()))
		if diags.HasError() {
			return fmt.Errorf("unable to set LBU AccessLog: %v", diags.Errors())
		}
	}
	data.Listeners, diags = getListenersFromApiResponse(ctx, lbu.GetListeners())
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU Listeners: %v", diags.Errors())
	}
	data.ApplicationStickyCookiePolicies, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockAppStickyCookiePolicy{})}, getAppStickyCookiePoliciesFromApiResponse(lbu.GetApplicationStickyCookiePolicies()))
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU ApplicationStickyCookiePolicies: %v", diags.Errors())
	}
	data.LoadBalancerStickyCookiePolicies, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockLBUStickyCookiePolicy{})}, getLBUStickyCookiePoliciesFromApiResponse(lbu.GetLoadBalancerStickyCookiePolicies()))
	if diags.HasError() {
		return fmt.Errorf("unable to set LBU LoadBalancerStickyCookiePolicies: %v", diags.Errors())
	}
	return nil
}

func getListenersForUpdate(dataPlan, dataState *LoadBalancerModel) ([]oscgo.ListenerForCreation, []oscgo.ListenerForCreation) {
	var (
		listenersToAdd, listenersToRemove []oscgo.ListenerForCreation
	)

	stateListeners := set.NewSet[ComparableListener]()
	planListeners := set.NewSet[ComparableListener]()
	for _, sListeners := range dataState.Listeners {
		stateListeners.Add(listenersToComparable(sListeners))
	}
	for _, pListeners := range dataPlan.Listeners {
		planListeners.Add(listenersToComparable(pListeners))
	}

	ltnToAdd := planListeners.Difference(stateListeners)
	ltnToRemove := stateListeners.Difference(planListeners)
	if len(ltnToAdd.ToSlice()) != 0 {
		listenersToAdd = getComparableListenersForUpdate(ltnToAdd.ToSlice())
	}
	if len(ltnToRemove.ToSlice()) != 0 {
		listenersToRemove = getComparableListenersForUpdate(ltnToRemove.ToSlice())
	}

	return listenersToAdd, listenersToRemove
}

func getListenersConfig(confListeners []Listeners) []oscgo.ListenerForCreation {
	result := make([]oscgo.ListenerForCreation, 0, len(confListeners))
	for _, listener := range confListeners {
		oscListener := oscgo.NewListenerForCreation(listener.BackendPort.ValueInt32(), listener.LoadBalancerPort.ValueInt32(), listener.LoadBalancerProtocol.ValueString())

		if !listener.BackendProtocol.IsUnknown() && !listener.BackendProtocol.IsNull() {
			oscListener.SetBackendProtocol(listener.BackendProtocol.ValueString())
		}
		if !listener.ServerCertificateId.IsUnknown() && !listener.ServerCertificateId.IsNull() {
			oscListener.SetServerCertificateId(listener.ServerCertificateId.ValueString())
		}
		result = append(result, *oscListener)
	}
	return result
}

func getComparableListenersForUpdate(confListeners []ComparableListener) []oscgo.ListenerForCreation {
	result := make([]oscgo.ListenerForCreation, 0, len(confListeners))
	for _, listener := range confListeners {
		oscListener := oscgo.NewListenerForCreation(listener.BackendPort.ValueInt32(), listener.LoadBalancerPort.ValueInt32(), listener.LoadBalancerProtocol.ValueString())

		if !listener.BackendProtocol.IsUnknown() && !listener.BackendProtocol.IsNull() {
			oscListener.SetBackendProtocol(listener.BackendProtocol.ValueString())
		}
		if !listener.ServerCertificateId.IsUnknown() && !listener.ServerCertificateId.IsNull() {
			oscListener.SetServerCertificateId(listener.ServerCertificateId.ValueString())
		}
		result = append(result, *oscListener)
	}
	return result
}
func getListenersFromApiResponse(ctx context.Context, listenerResp []oscgo.Listener) (listeners []Listeners, diag diag.Diagnostics) {

	pNames := types.SetNull(types.StringType)
	for _, listener := range listenerResp {
		rlistener := Listeners{
			BackendPort:          types.Int32Value(listener.GetBackendPort()),
			BackendProtocol:      types.StringValue(listener.GetBackendProtocol()),
			LoadBalancerPort:     types.Int32Value(listener.GetLoadBalancerPort()),
			LoadBalancerProtocol: types.StringValue(listener.GetLoadBalancerProtocol()),
			ServerCertificateId:  types.StringValue(listener.GetServerCertificateId()),
		}
		if rlistener.ServerCertificateId.ValueString() == "" {
			rlistener.ServerCertificateId = types.StringNull()
		}
		if listener.HasPolicyNames() {
			pNames, diag = types.SetValueFrom(ctx, types.StringType, listener.GetPolicyNames())
			if diag.HasError() {
				return listeners, diag
			}
		}
		rlistener.PolicyNames = pNames
		listeners = append(listeners, rlistener)
	}
	return
}

func getAccessLogFromApiResponse(accLog oscgo.AccessLog) []BlockAccessLog {
	accessLog := make([]BlockAccessLog, 0, 1)
	accesslog := BlockAccessLog{
		IsEnabled:           types.BoolValue(accLog.GetIsEnabled()),
		OsuBucketName:       types.StringValue(accLog.GetOsuBucketName()),
		OsuBucketPrefix:     types.StringValue(accLog.GetOsuBucketPrefix()),
		PublicationInterval: types.Int32Value(accLog.GetPublicationInterval()),
	}
	accessLog = append(accessLog, accesslog)
	return accessLog
}

func getHealthCheckFromApiResponse(hCheck oscgo.HealthCheck) []BlockHealthCheck {

	healthCheck := make([]BlockHealthCheck, 0, 1)
	path := types.StringNull()
	if hCheck.HasPath() {
		path = types.StringValue(hCheck.GetPath())
	}

	check := BlockHealthCheck{
		CheckInterval:      types.Int32Value(hCheck.GetCheckInterval()),
		HealthyThreshold:   types.Int32Value(hCheck.GetHealthyThreshold()),
		UnhealthyThreshold: types.Int32Value(hCheck.GetUnhealthyThreshold()),
		Port:               types.Int32Value(hCheck.GetPort()),
		Protocol:           types.StringValue(hCheck.GetProtocol()),
		Path:               path,
		Timeout:            types.Int32Value(hCheck.GetTimeout()),
	}

	healthCheck = append(healthCheck, check)
	return healthCheck
}

func getSourceSGFromApiResponse(sourceSG oscgo.SourceSecurityGroup) []BlockSourceSG {
	sourceSecurityGroup := make([]BlockSourceSG, 0, 1)
	sSgroup := BlockSourceSG{
		SecurityGroupName:      types.StringValue(sourceSG.GetSecurityGroupName()),
		SecurityGroupAccountId: types.StringValue(sourceSG.GetSecurityGroupAccountId()),
	}
	sourceSecurityGroup = append(sourceSecurityGroup, sSgroup)
	return sourceSecurityGroup
}
func getLBUStickyCookiePoliciesFromApiResponse(lbuSCPolicies []oscgo.LoadBalancerStickyCookiePolicy) (lbuStickyCookiePolicy []BlockLBUStickyCookiePolicy) {
	for _, lbuSCPolicy := range lbuSCPolicies {
		lbuPolicy := BlockLBUStickyCookiePolicy{
			CookieExpirationPeriod: types.Int32Value(lbuSCPolicy.GetCookieExpirationPeriod()),
			PolicyName:             types.StringValue(lbuSCPolicy.GetPolicyName()),
		}
		lbuStickyCookiePolicy = append(lbuStickyCookiePolicy, lbuPolicy)
	}
	return
}

func getAppStickyCookiePoliciesFromApiResponse(appSCPolicies []oscgo.ApplicationStickyCookiePolicy) (appStickyCookiePolicy []BlockAppStickyCookiePolicy) {
	for _, appSCPolicy := range appSCPolicies {
		lbuPolicy := BlockAppStickyCookiePolicy{
			CookieName: types.StringValue(appSCPolicy.GetCookieName()),
			PolicyName: types.StringValue(appSCPolicy.GetPolicyName()),
		}
		appStickyCookiePolicy = append(appStickyCookiePolicy, lbuPolicy)
	}
	return
}

func listenersToComparable(listener Listeners) ComparableListener {
	var policyNamesStr string
	if !listener.PolicyNames.IsNull() && !listener.PolicyNames.IsUnknown() {
		elements := listener.PolicyNames.Elements()
		policies := make([]string, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				policies = append(policies, strVal.ValueString())
			}
		}
		// Trier pour assurer la comparabilité
		sort.Strings(policies)
		policyBytes, _ := json.Marshal(policies)
		policyNamesStr = string(policyBytes)
	}
	return ComparableListener{
		BackendPort:          listener.BackendPort,
		BackendProtocol:      listener.BackendProtocol,
		LoadBalancerPort:     listener.LoadBalancerPort,
		LoadBalancerProtocol: listener.LoadBalancerProtocol,
		ServerCertificateId:  listener.ServerCertificateId,
		PolicyNames:          types.StringValue(policyNamesStr),
	}
}
