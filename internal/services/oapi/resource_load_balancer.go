package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &loadBalancerResource{}
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
)

const (
	loadBalancerErrCreate          = "Unable to create Load Balancer"
	loadBalancerErrUpdate          = "Unable to update Load Balancer"
	loadBalancerErrUpdateListeners = "Unable to update Load Balancer listeners"
	loadBalancerErrDelete          = "Unable to delete Load Balancer"
	loadBalancerErrWait            = "Unable to wait for Load Balancer state"
)

type loadBalancerModel struct {
	SubregionNames                   types.List     `tfsdk:"subregion_names"`
	LoadBalancerName                 types.String   `tfsdk:"load_balancer_name"`
	LoadBalancerType                 types.String   `tfsdk:"load_balancer_type"`
	SecurityGroups                   types.Set      `tfsdk:"security_groups"`
	Subnets                          types.List     `tfsdk:"subnets"`
	DnsName                          types.String   `tfsdk:"dns_name"`
	AccessLog                        types.List     `tfsdk:"access_log"`
	HealthCheck                      types.Set      `tfsdk:"health_check"`
	BackendVmIds                     types.List     `tfsdk:"backend_vm_ids"`
	BackendIps                       types.List     `tfsdk:"backend_ips"`
	Listeners                        types.Set      `tfsdk:"listeners"`
	SourceSecurityGroup              types.List     `tfsdk:"source_security_group"`
	PublicIp                         types.String   `tfsdk:"public_ip"`
	SecuredCookies                   types.Bool     `tfsdk:"secured_cookies"`
	NetId                            types.String   `tfsdk:"net_id"`
	ApplicationStickyCookiePolicies  types.List     `tfsdk:"application_sticky_cookie_policies"`
	LoadBalancerStickyCookiePolicies types.List     `tfsdk:"load_balancer_sticky_cookie_policies"`
	State                            types.String   `tfsdk:"state"`
	RequestId                        types.String   `tfsdk:"request_id"`
	Id                               types.String   `tfsdk:"id"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type loadBalancerAccessLogModel struct {
	IsEnabled           types.Bool   `tfsdk:"is_enabled"`
	OsuBucketName       types.String `tfsdk:"osu_bucket_name"`
	OsuBucketPrefix     types.String `tfsdk:"osu_bucket_prefix"`
	PublicationInterval types.Int64  `tfsdk:"publication_interval"`
}

type loadBalancerHealthCheckModel struct {
	HealthyThreshold   types.Int64  `tfsdk:"healthy_threshold"`
	UnhealthyThreshold types.Int64  `tfsdk:"unhealthy_threshold"`
	Path               types.String `tfsdk:"path"`
	CheckInterval      types.Int64  `tfsdk:"check_interval"`
	Port               types.Int64  `tfsdk:"port"`
	Protocol           types.String `tfsdk:"protocol"`
	Timeout            types.Int64  `tfsdk:"timeout"`
}

type loadBalancerListenerModel struct {
	BackendPort          types.Int64  `tfsdk:"backend_port"`
	BackendProtocol      types.String `tfsdk:"backend_protocol"`
	LoadBalancerPort     types.Int64  `tfsdk:"load_balancer_port"`
	LoadBalancerProtocol types.String `tfsdk:"load_balancer_protocol"`
	ServerCertificateId  types.String `tfsdk:"server_certificate_id"`
	PolicyNames          types.List   `tfsdk:"policy_names"`
}

var loadBalancerListenerAttrTypes = map[string]attr.Type{
	"backend_port":           types.Int64Type,
	"backend_protocol":       types.StringType,
	"load_balancer_port":     types.Int64Type,
	"load_balancer_protocol": types.StringType,
	"server_certificate_id":  types.StringType,
	"policy_names":           types.ListType{ElemType: types.StringType},
}

var (
	loadBalancerAccessLogAttrTypes       = fwhelpers.GetAttrTypes(loadBalancerAccessLogModel{})
	loadBalancerHealthCheckAttrTypes     = fwhelpers.GetAttrTypes(loadBalancerHealthCheckModel{})
	loadBalancerSourceSecurityGroupAttrs = fwhelpers.GetAttrTypes(loadBalancerSourceSecurityGroupModel{})
	loadBalancerAppStickyPolicyAttrTypes = fwhelpers.GetAttrTypes(loadBalancerAppStickyCookiePolicyModel{})
	loadBalancerStickyPolicyAttrTypes    = fwhelpers.GetAttrTypes(loadBalancerStickyCookiePolicyModel{})
)

type loadBalancerSourceSecurityGroupModel struct {
	SecurityGroupName      types.String `tfsdk:"security_group_name"`
	SecurityGroupAccountId types.String `tfsdk:"security_group_account_id"`
}

type loadBalancerAppStickyCookiePolicyModel struct {
	CookieName types.String `tfsdk:"cookie_name"`
	PolicyName types.String `tfsdk:"policy_name"`
}

type loadBalancerStickyCookiePolicyModel struct {
	PolicyName types.String `tfsdk:"policy_name"`
}

type loadBalancerCommon struct {
	Client *osc.Client
}

type loadBalancerResource struct {
	loadBalancerCommon
}

func NewResourceLoadBalancer() resource.Resource {
	return &loadBalancerResource{}
}

func (r *loadBalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Client = client.OSC
}

func (r *loadBalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer"
}

func (r *loadBalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import load balancer identifier. Got: %v", req.ID),
		)
		return
	}

	var data loadBalancerModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.LoadBalancerName = to.String(id)
	data.Tags = TagsNull()
	data.SubregionNames = types.ListNull(types.StringType)
	data.SecurityGroups = types.SetNull(types.StringType)
	data.Subnets = types.ListNull(types.StringType)
	data.AccessLog = types.ListNull(to.ObjType(loadBalancerAccessLogAttrTypes))
	data.HealthCheck = types.SetNull(to.ObjType(loadBalancerHealthCheckAttrTypes))
	data.BackendVmIds = types.ListNull(types.StringType)
	data.BackendIps = types.ListNull(types.StringType)
	data.Listeners = types.SetNull(to.ObjType(loadBalancerListenerAttrTypes))
	data.SourceSecurityGroup = types.ListNull(to.ObjType(loadBalancerSourceSecurityGroupAttrs))
	data.ApplicationStickyCookiePolicies = types.ListNull(to.ObjType(loadBalancerAppStickyPolicyAttrTypes))
	data.LoadBalancerStickyCookiePolicies = types.ListNull(to.ObjType(loadBalancerStickyPolicyAttrTypes))

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *loadBalancerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("subnets"),
			path.MatchRoot("subregion_names"),
		),
	}
}

func ifSecureProtocol(ctx context.Context, req validator.StringRequest) bool {
	var protocol types.String
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("load_balancer_protocol"), &protocol)
	if diags.HasError() {
		return false
	}

	return fwhelpers.IsSet(protocol) && (protocol.ValueString() == "HTTPS" || protocol.ValueString() == "SSL")
}

func (r *loadBalancerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags": TagsSchemaFW(),
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"listeners": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"backend_port": schema.Int64Attribute{
							Required: true,
						},
						"backend_protocol": schema.StringAttribute{
							Required: true,
						},
						"load_balancer_port": schema.Int64Attribute{
							Required: true,
						},
						"load_balancer_protocol": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								validatorstring.AlsoRequiresIf(
									path.MatchRelative().AtParent().AtName("server_certificate_id"),
									ifSecureProtocol,
									"Missing server certificate",
									"'server_certificate_id' is required when 'load_balancer_protocol' is 'HTTPS' or 'SSL'.",
								),
							},
						},
						"server_certificate_id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"policy_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"subregion_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"subnets": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_name": schema.StringAttribute{
				Computed: true,
			},
			"access_log": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.accessLogComputedAttributes(),
			},
			"health_check": schema.SetNestedAttribute{
				Computed:     true,
				NestedObject: r.healthCheckComputedAttributes(),
			},
			"backend_vm_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"backend_ips": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_security_group": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.sourceSecurityGroupAttributes(),
			},
			"public_ip": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"secured_cookies": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
			},
			"application_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.appStickyPolicyComputedAttributes(),
			},
			"load_balancer_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.stickyPolicyComputedAttributes(),
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (c *loadBalancerCommon) accessLogComputedAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"is_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"osu_bucket_name": schema.StringAttribute{
				Computed: true,
			},
			"osu_bucket_prefix": schema.StringAttribute{
				Computed: true,
			},
			"publication_interval": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (c *loadBalancerCommon) healthCheckComputedAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"healthy_threshold": schema.Int64Attribute{
				Computed: true,
			},
			"unhealthy_threshold": schema.Int64Attribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Computed: true,
			},
			"check_interval": schema.Int64Attribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"protocol": schema.StringAttribute{
				Computed: true,
			},
			"timeout": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (c *loadBalancerCommon) sourceSecurityGroupAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"security_group_name": schema.StringAttribute{
				Computed: true,
			},
			"security_group_account_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (c *loadBalancerCommon) appStickyPolicyComputedAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"cookie_name": schema.StringAttribute{
				Computed: true,
			},
			"policy_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (c *loadBalancerCommon) stickyPolicyComputedAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"policy_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *loadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data loadBalancerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	listeners, diags := r.expandListeners(ctx, data.Listeners)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateLoadBalancerRequest{
		LoadBalancerName: data.LoadBalancerName.ValueString(),
		Listeners:        listeners,
	}

	if fwhelpers.IsSet(data.Tags) {
		tagsModel, diags := to.Slice[ResourceTag](ctx, data.Tags)
		if diags.HasError() {
			return
		}
		tags := expandOAPITagsFW(tagsModel)
		createReq.Tags = &tags
	}

	if fwhelpers.IsSet(data.LoadBalancerType) {
		createReq.LoadBalancerType = data.LoadBalancerType.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.PublicIp) {
		createReq.PublicIp = data.PublicIp.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.SecurityGroups) {
		sgs, diags := to.Slice[string](ctx, data.SecurityGroups)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}

		createReq.SecurityGroups = &sgs
	}
	if fwhelpers.IsSet(data.Subnets) {
		subnets, diags := to.Slice[string](ctx, data.Subnets)
		if diags.HasError() {
			return
		}

		createReq.Subnets = &subnets
	}
	if fwhelpers.IsSet(data.SubregionNames) {
		subregions, diags := to.Slice[string](ctx, data.SubregionNames)
		if diags.HasError() {
			return
		}

		createReq.SubregionNames = &subregions
	}

	createResp, err := r.Client.CreateLoadBalancer(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerErrCreate, err.Error())
		return
	}

	data.Id = to.String(createResp.LoadBalancer.LoadBalancerName)
	data.LoadBalancerName = to.String(createResp.LoadBalancer.LoadBalancerName)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)

	stateData, err := r.flatten(ctx, data, *createResp.LoadBalancer)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if fwhelpers.IsSet(data.SecuredCookies) {
		updateReq := osc.UpdateLoadBalancerRequest{
			LoadBalancerName: createReq.LoadBalancerName,
			SecuredCookies:   data.SecuredCookies.ValueBoolPointer(),
		}
		_, err = r.Client.UpdateLoadBalancer(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(loadBalancerErrCreate, err.Error())
			return
		}
	}

	stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
		Pending: stateconf.States(osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring),
		Target:  stateconf.States(osc.LoadBalancerStateActive),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.LoadBalancerName.ValueString()),
	}
	waitRespAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerErrWait, err.Error())
		return
	}
	waitResp := waitRespAny.(*osc.ReadLoadBalancersResponse)
	data.RequestId = to.String(waitResp.ResponseContext.RequestId)

	// We set the last read response to the state
	stateData, err = r.flatten(ctx, data, (*waitResp.LoadBalancers)[0])
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data loadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData loadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	name := stateData.Id.ValueString()
	updateReq := osc.UpdateLoadBalancerRequest{LoadBalancerName: name}
	doUpdate := false

	if fwhelpers.HasChange(planData.SecurityGroups, stateData.SecurityGroups) {
		groups, diag := to.Slice[string](ctx, planData.SecurityGroups)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.SecurityGroups = &groups

		doUpdate = true
	}

	if fwhelpers.HasChange(planData.PublicIp, stateData.PublicIp) && fwhelpers.IsSet(planData.PublicIp) {
		updateReq.PublicIp = planData.PublicIp.ValueStringPointer()

		doUpdate = true
	}

	if fwhelpers.HasChange(planData.SecuredCookies, stateData.SecuredCookies) {
		updateReq.SecuredCookies = planData.SecuredCookies.ValueBoolPointer()

		doUpdate = true
	}

	if doUpdate {
		_, err := r.Client.UpdateLoadBalancer(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(loadBalancerErrUpdate, err.Error())
			return
		}

		stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
			Pending: stateconf.States(osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring),
			Target:  stateconf.States(osc.LoadBalancerStateActive),
			Timeout: timeout,
			Refresh: r.refreshFunc(name),
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError(loadBalancerErrWait, err.Error())
			return
		}
	}

	if !planData.Tags.Equal(stateData.Tags) {
		if err := r.updateTags(ctx, name, stateData.Tags, planData.Tags, timeout); err != nil {
			resp.Diagnostics.AddError(loadBalancerErrUpdate, err.Error())
			return
		}
	}

	if !planData.Listeners.Equal(stateData.Listeners) {
		if err := r.updateListeners(ctx, name, stateData.Listeners, planData.Listeners, timeout); err != nil {
			resp.Diagnostics.AddError(loadBalancerErrUpdateListeners, err.Error())
			return
		}
	}

	stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
		Pending: stateconf.States(osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring),
		Target:  stateconf.States(osc.LoadBalancerStateActive),
		Timeout: timeout,
		Refresh: r.refreshFunc(name),
	}
	waitRespAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerErrWait, err.Error())
		return
	}
	waitResp := waitRespAny.(*osc.ReadLoadBalancersResponse)
	planData.RequestId = to.String(waitResp.ResponseContext.RequestId)

	newData, err := r.flatten(ctx, planData, (*waitResp.LoadBalancers)[0])
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *loadBalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data loadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.DeleteLoadBalancerRequest{
		LoadBalancerName: data.LoadBalancerName.ValueString(),
	}
	_, err := r.Client.DeleteLoadBalancer(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerErrDelete, err.Error())
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
		Pending: stateconf.States(osc.LoadBalancerStateDeleting, osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring),
		Target:  stateconf.States(osc.LoadBalancerStateDeleted),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.LoadBalancerName.ValueString()),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	switch {
	case errors.Is(err, ErrResourceEmpty):
	case err != nil:
		resp.Diagnostics.AddError(loadBalancerErrWait, err.Error())
	}
}

func (r *loadBalancerResource) read(ctx context.Context, timeout time.Duration, data loadBalancerModel) (loadBalancerModel, error) {
	stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
		Pending: stateconf.States(osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring, osc.LoadBalancerStateDeleting),
		Target:  stateconf.States(osc.LoadBalancerStateActive),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.LoadBalancerName.ValueString()),
	}
	respAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return data, err
	}
	resp := respAny.(*osc.ReadLoadBalancersResponse)

	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return r.flatten(ctx, data, (*resp.LoadBalancers)[0])
}

func (r *loadBalancerResource) flatten(ctx context.Context, data loadBalancerModel, lb osc.LoadBalancer) (loadBalancerModel, error) {
	tags, diag := flattenOAPITagsFW(ctx, lb.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.Tags = tags

	accessLogModel := r.flattenAccessLog(lb.AccessLog)
	accessLog, diag := to.ListObject(ctx, accessLogModel, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	healthCheckModel := r.flattenHealthCheck(lb.HealthCheck)
	healthCheck, diag := to.SetObject(ctx, healthCheckModel, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	listenersModel, diag := r.flattenListeners(ctx, lb.Listeners)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	listeners, diag := to.SetFromAttrType(ctx, listenersModel, to.ObjType(loadBalancerListenerAttrTypes), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	sourceSecurityGroup, diag := r.flattenSourceSecurityGroup(ctx, lb.SourceSecurityGroup)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	appPolicies, diag := r.flattenAppStickyPolicies(ctx, lb.ApplicationStickyCookiePolicies)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	stickyPolicies, diag := r.flattenStickyPolicies(ctx, lb.LoadBalancerStickyCookiePolicies)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	subregionNames, diag := to.List(ctx, lb.SubregionNames, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	backendVmIds, diag := to.List(ctx, lb.BackendVmIds, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	backendIps, diag := to.List(ctx, lb.BackendIps, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	securityGroups, diag := to.Set(ctx, lb.SecurityGroups, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	subnets, diag := to.List(ctx, lb.Subnets, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.SubregionNames = subregionNames
	data.LoadBalancerName = to.String(lb.LoadBalancerName)
	data.LoadBalancerType = to.String(lb.LoadBalancerType)
	data.SecurityGroups = securityGroups
	data.Subnets = subnets
	data.Tags = tags
	data.DnsName = to.String(lb.DnsName)
	data.AccessLog = accessLog
	data.HealthCheck = healthCheck
	data.BackendVmIds = backendVmIds
	data.BackendIps = backendIps
	data.Listeners = listeners
	data.SourceSecurityGroup = sourceSecurityGroup
	data.PublicIp = to.String(ptr.From(lb.PublicIp))
	data.SecuredCookies = to.Bool(lb.SecuredCookies)
	data.NetId = to.String(ptr.From(lb.NetId))
	data.ApplicationStickyCookiePolicies = appPolicies
	data.LoadBalancerStickyCookiePolicies = stickyPolicies
	data.State = to.String(lb.State)
	data.Id = to.String(lb.LoadBalancerName)

	return data, nil
}

func (r *loadBalancerCommon) flattenAccessLog(logs osc.AccessLog) []loadBalancerAccessLogModel {
	return []loadBalancerAccessLogModel{
		{
			IsEnabled:           to.Bool(logs.IsEnabled),
			OsuBucketName:       to.String(ptr.From(logs.OsuBucketName)),
			OsuBucketPrefix:     to.String(ptr.From(logs.OsuBucketPrefix)),
			PublicationInterval: to.Int64(ptr.From(logs.PublicationInterval)),
		},
	}
}

func (r *loadBalancerCommon) flattenHealthCheck(healthCheck osc.HealthCheck) []loadBalancerHealthCheckModel {
	return []loadBalancerHealthCheckModel{
		{
			HealthyThreshold:   to.Int64(healthCheck.HealthyThreshold),
			UnhealthyThreshold: to.Int64(healthCheck.UnhealthyThreshold),
			Path:               to.String(ptr.From(healthCheck.Path)),
			CheckInterval:      to.Int64(healthCheck.CheckInterval),
			Port:               to.Int64(healthCheck.Port),
			Protocol:           to.String(healthCheck.Protocol),
			Timeout:            to.Int64(healthCheck.Timeout),
		},
	}
}

func (r *loadBalancerCommon) flattenListeners(ctx context.Context, listeners []osc.Listener) ([]loadBalancerListenerModel, diag.Diagnostics) {
	models := make([]loadBalancerListenerModel, 0, len(listeners))
	for _, listener := range listeners {
		policyNames, diags := to.List(ctx, listener.PolicyNames, to.ZeroValueAsEmpty)
		if diags.HasError() {
			return nil, diags
		}

		models = append(models, loadBalancerListenerModel{
			BackendPort:          to.Int64(listener.BackendPort),
			BackendProtocol:      to.String(listener.BackendProtocol),
			LoadBalancerPort:     to.Int64(listener.LoadBalancerPort),
			LoadBalancerProtocol: to.String(listener.LoadBalancerProtocol),
			ServerCertificateId:  to.String(ptr.From(listener.ServerCertificateId)),
			PolicyNames:          policyNames,
		})
	}

	return models, nil
}

func (r *loadBalancerCommon) flattenSourceSecurityGroup(ctx context.Context, sourceSecurityGroup osc.SourceSecurityGroup) (types.List, diag.Diagnostics) {
	model := []loadBalancerSourceSecurityGroupModel{
		{
			SecurityGroupName:      to.String(ptr.From(sourceSecurityGroup.SecurityGroupName)),
			SecurityGroupAccountId: to.String(ptr.From(sourceSecurityGroup.SecurityGroupAccountId)),
		},
	}

	return to.ListObject(ctx, model, to.ZeroValueAsEmpty)
}

func (r *loadBalancerCommon) flattenAppStickyPolicies(ctx context.Context, policies []osc.ApplicationStickyCookiePolicy) (types.List, diag.Diagnostics) {
	model := lo.Map(policies, func(policy osc.ApplicationStickyCookiePolicy, _ int) loadBalancerAppStickyCookiePolicyModel {
		return loadBalancerAppStickyCookiePolicyModel{
			CookieName: to.String(policy.CookieName),
			PolicyName: to.String(policy.PolicyName),
		}
	})

	return to.ListObject(ctx, model, to.ZeroValueAsEmpty)
}

func (r *loadBalancerCommon) flattenStickyPolicies(ctx context.Context, policies []osc.LoadBalancerStickyCookiePolicy) (types.List, diag.Diagnostics) {
	model := lo.Map(policies, func(policy osc.LoadBalancerStickyCookiePolicy, _ int) loadBalancerStickyCookiePolicyModel {
		return loadBalancerStickyCookiePolicyModel{
			PolicyName: to.String(policy.PolicyName),
		}
	})

	return to.ListObject(ctx, model, to.ZeroValueAsEmpty)
}

func (r *loadBalancerResource) expandListeners(ctx context.Context, listenersSet types.Set) ([]osc.ListenerForCreation, diag.Diagnostics) {
	listeners, diags := to.Slice[loadBalancerListenerModel](ctx, listenersSet)
	if diags.HasError() {
		return nil, diags
	}

	return lo.Map(listeners, func(listener loadBalancerListenerModel, _ int) osc.ListenerForCreation {
		apiListener := osc.ListenerForCreation{
			BackendPort:          int(listener.BackendPort.ValueInt64()),
			BackendProtocol:      listener.BackendProtocol.ValueStringPointer(),
			LoadBalancerPort:     int(listener.LoadBalancerPort.ValueInt64()),
			LoadBalancerProtocol: listener.LoadBalancerProtocol.ValueString(),
		}

		if fwhelpers.IsSet(listener.ServerCertificateId) && listener.ServerCertificateId.ValueString() != "" {
			apiListener.ServerCertificateId = listener.ServerCertificateId.ValueStringPointer()
		}

		return apiListener
	}), nil
}

func (r *loadBalancerResource) updateTags(ctx context.Context, name string, oldSet, newSet types.Set, timeout time.Duration) error {
	toCreate, toRemove, diags := diffOAPITagsFW(ctx, oldSet, newSet)
	if diags.HasError() {
		return from.Diag(diags)
	}

	if len(toRemove) > 0 {
		tags := lo.Map(toRemove, func(tag osc.ResourceTag, _ int) osc.ResourceLoadBalancerTag {
			return osc.ResourceLoadBalancerTag{Key: tag.Key}
		})

		deleteReq := osc.DeleteLoadBalancerTagsRequest{
			LoadBalancerNames: []string{name},
			Tags:              tags,
		}
		_, err := r.Client.DeleteLoadBalancerTags(ctx, deleteReq, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}

	if len(toCreate) > 0 {
		tags := lo.Map(toCreate, func(tag osc.ResourceTag, _ int) osc.ResourceTag {
			return osc.ResourceTag{Key: tag.Key, Value: tag.Value}
		})

		createReq := osc.CreateLoadBalancerTagsRequest{
			LoadBalancerNames: []string{name},
			Tags:              tags,
		}
		_, err := r.Client.CreateLoadBalancerTags(ctx, createReq, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *loadBalancerResource) updateListeners(ctx context.Context, name string, oldSet, newSet types.Set, timeout time.Duration) error {
	toCreateSet, toRemoveSet, diag := fwhelpers.Difference(ctx, oldSet, newSet)
	if diag.HasError() {
		return from.Diag(diag)
	}

	toRemove, diag := r.expandListeners(ctx, toRemoveSet)
	if diag.HasError() {
		return from.Diag(diag)
	}
	toCreate, diag := r.expandListeners(ctx, toCreateSet)
	if diag.HasError() {
		return from.Diag(diag)
	}

	if len(toRemove) > 0 {
		ports := lo.Map(toRemove, func(listener osc.ListenerForCreation, _ int) int {
			return int(listener.LoadBalancerPort)
		})
		deleteReq := osc.DeleteLoadBalancerListenersRequest{
			LoadBalancerName:  name,
			LoadBalancerPorts: ports,
		}

		_, err := r.Client.DeleteLoadBalancerListeners(ctx, deleteReq, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}
	if len(toCreate) > 0 {
		createReq := osc.CreateLoadBalancerListenersRequest{
			LoadBalancerName: name,
			Listeners:        toCreate,
		}

		_, err := r.Client.CreateLoadBalancerListeners(ctx, createReq, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *loadBalancerCommon) refreshFunc(name string) stateconf.StateRefreshFunc[osc.LoadBalancerState] {
	return func(ctx context.Context) (any, osc.LoadBalancerState, error) {
		req := osc.ReadLoadBalancersRequest{
			Filters: &osc.FiltersLoadBalancer{
				LoadBalancerNames: &[]string{name},
			},
		}

		rawResp, err := r.Client.ReadLoadBalancersRaw(ctx, req)
		if err != nil {
			return nil, "", err
		}

		parsed, err := osc.ParseReadLoadBalancersResp(rawResp)
		if err != nil {
			return nil, "", err
		}
		// Retry only non-applicative HTTP 404s
		if parsed.StatusCode() == 404 {
			// We have to return a non-nil state to not return an error and force the retry
			return nil, osc.LoadBalancerStateProvisioning, nil
		}
		resp, err := parsed.Expect()
		if err != nil {
			return nil, "", err
		}

		if len(ptr.From(resp.LoadBalancers)) == 0 {
			return nil, "", ErrResourceEmpty
		}

		lb := (*resp.LoadBalancers)[0]
		if lb.State == osc.LoadBalancerStateDeleted {
			return nil, "", ErrResourceEmpty
		}

		return resp, lb.State, nil
	}
}
