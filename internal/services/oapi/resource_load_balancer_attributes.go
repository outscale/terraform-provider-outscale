package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
)

var (
	_ resource.Resource                = &loadBalancerAttributesResource{}
	_ resource.ResourceWithConfigure   = &loadBalancerAttributesResource{}
	_ resource.ResourceWithImportState = &loadBalancerAttributesResource{}
)

const (
	loadBalancerAttributesErrCreate = "Unable to create Load Balancer Attributes"
	loadBalancerAttributesErrUpdate = "Unable to update Load Balancer Attributes"
	loadBalancerAttributesErrDelete = "Unable to delete Load Balancer Attributes"
	loadBalancerAttributesErrWait   = "Unable to wait for Load Balancer state"
)

type loadBalancerAttributesModel struct {
	AccessLog   types.Object `tfsdk:"access_log"`
	HealthCheck types.Object `tfsdk:"health_check"`
	// AccessLog                        types.List     `tfsdk:"access_log"`
	// HealthCheck                      types.List     `tfsdk:"health_check"`
	Listeners                        types.List     `tfsdk:"listeners"`
	Subnets                          types.List     `tfsdk:"subnets"`
	SubregionNames                   types.List     `tfsdk:"subregion_names"`
	LoadBalancerPort                 types.Int64    `tfsdk:"load_balancer_port"`
	DnsName                          types.String   `tfsdk:"dns_name"`
	SecurityGroups                   types.List     `tfsdk:"security_groups"`
	ServerCertificateId              types.String   `tfsdk:"server_certificate_id"`
	SourceSecurityGroup              types.List     `tfsdk:"source_security_group"`
	BackendVmIds                     types.List     `tfsdk:"backend_vm_ids"`
	ApplicationStickyCookiePolicies  types.List     `tfsdk:"application_sticky_cookie_policies"`
	LoadBalancerStickyCookiePolicies types.List     `tfsdk:"load_balancer_sticky_cookie_policies"`
	LoadBalancerName                 types.String   `tfsdk:"load_balancer_name"`
	LoadBalancerType                 types.String   `tfsdk:"load_balancer_type"`
	PolicyNames                      types.List     `tfsdk:"policy_names"`
	RequestId                        types.String   `tfsdk:"request_id"`
	Id                               types.String   `tfsdk:"id"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
	TagsComputedModel
}

type loadBalancerAttributesResource struct {
	loadBalancerCommon
}

func NewResourceLoadBalancerAttributes() resource.Resource {
	return &loadBalancerAttributesResource{}
}

func (r *loadBalancerAttributesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *loadBalancerAttributesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_attributes"
}

func (r *loadBalancerAttributesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import load balancer attributes identifier. Got: %v", req.ID),
		)
		return
	}

	var data loadBalancerAttributesModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(req.ID)
	data.LoadBalancerName = to.String(req.ID)
	data.AccessLog = types.ObjectNull(loadBalancerAccessLogAttrTypes)
	data.HealthCheck = types.ObjectNull(loadBalancerHealthCheckAttrTypes)
	data.Listeners = types.ListNull(to.ObjType(loadBalancerListenerAttrTypes))
	data.Subnets = types.ListNull(types.StringType)
	data.SubregionNames = types.ListNull(types.StringType)
	data.Tags = ComputedTagsNull()
	data.SecurityGroups = types.ListNull(types.StringType)
	data.SourceSecurityGroup = types.ListNull(to.ObjType(loadBalancerSourceSecurityGroupAttrs))
	data.BackendVmIds = types.ListNull(types.StringType)
	data.ApplicationStickyCookiePolicies = types.ListNull(to.ObjType(loadBalancerAppStickyPolicyAttrTypes))
	data.LoadBalancerStickyCookiePolicies = types.ListNull(to.ObjType(loadBalancerStickyPolicyAttrTypes))
	data.PolicyNames = types.ListNull(types.StringType)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *loadBalancerAttributesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			// Default LBU value causes issue with blocks, they cannot receive computed values when not configured
			// "access_log": schema.ListNestedBlock{
			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 	},
			// 	NestedObject: schema.NestedBlockObject{
			// 		PlanModifiers: []planmodifier.Object{
			// 			objectplanmodifier.RequiresReplaceIfConfigured(),
			// 		},
			// 		Attributes: map[string]schema.Attribute{
			// 			"is_enabled": schema.BoolAttribute{
			// 				Optional: true,
			// 			},
			// 			"osu_bucket_name": schema.StringAttribute{
			// 				Optional: true,
			// 			},
			// 			"osu_bucket_prefix": schema.StringAttribute{
			// 				Optional: true,
			// 			},
			// 			"publication_interval": schema.Int64Attribute{
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },
			// "health_check": schema.ListNestedBlock{
			// 	Validators: []validator.List{
			// 		listvalidator.SizeAtMost(1),
			// 	},
			// 	NestedObject: schema.NestedBlockObject{
			// 		PlanModifiers: []planmodifier.Object{
			// 			objectplanmodifier.RequiresReplaceIfConfigured(),
			// 		},
			// 		Attributes: map[string]schema.Attribute{
			// 			"healthy_threshold": schema.Int64Attribute{
			// 				Computed: true,
			// 				Optional: true,
			// 				Validators: []validator.Int64{
			// 					int64validator.Between(2, 10),
			// 				},
			// 			},
			// 			"unhealthy_threshold": schema.Int64Attribute{
			// 				Computed: true,
			// 				Optional: true,
			// 				Validators: []validator.Int64{
			// 					int64validator.Between(2, 10),
			// 				},
			// 			},
			// 			"path": schema.StringAttribute{
			// 				Optional: true,
			// 			},
			// 			"port": schema.Int64Attribute{
			// 				Required: true,
			// 				Validators: []validator.Int64{
			// 					int64validator.Between(1, 65535),
			// 				},
			// 			},
			// 			"protocol": schema.StringAttribute{
			// 				Required: true,
			// 			},
			// 			"check_interval": schema.Int64Attribute{
			// 				Computed: true,
			// 				Optional: true,
			// 				Validators: []validator.Int64{
			// 					int64validator.Between(5, 600),
			// 				},
			// 			},
			// 			"timeout": schema.Int64Attribute{
			// 				Computed: true,
			// 				Optional: true,
			// 				Validators: []validator.Int64{
			// 					int64validator.Between(2, 60),
			// 				},
			// 			},
			// 		},
			// 	},
			// },
		},
		Attributes: map[string]schema.Attribute{
			"access_log": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIfConfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"is_enabled": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"osu_bucket_name": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"osu_bucket_prefix": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"publication_interval": schema.Int64Attribute{
						Computed: true,
						Optional: true,
					},
				},
			},

			"health_check": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIfConfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"healthy_threshold": schema.Int64Attribute{
						Computed: true,
						Optional: true,
						Validators: []validator.Int64{
							int64validator.Between(2, 10),
						},
					},
					"unhealthy_threshold": schema.Int64Attribute{
						Computed: true,
						Optional: true,
						Validators: []validator.Int64{
							int64validator.Between(2, 10),
						},
					},
					"path": schema.StringAttribute{
						Optional: true,
					},
					"port": schema.Int64Attribute{
						Required: true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"protocol": schema.StringAttribute{
						Required: true,
					},
					"check_interval": schema.Int64Attribute{
						Computed: true,
						Optional: true,
						Validators: []validator.Int64{
							int64validator.Between(5, 600),
						},
					},
					"timeout": schema.Int64Attribute{
						Computed: true,
						Optional: true,
						Validators: []validator.Int64{
							int64validator.Between(2, 60),
						},
					},
				},
			},
			"listeners": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.listenersComputedAttributes(),
			},
			"subnets": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"subregion_names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"load_balancer_port": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"dns_name": schema.StringAttribute{
				Computed: true,
			},
			"security_groups": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"server_certificate_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("load_balancer_port")),
				},
			},
			"source_security_group": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.sourceSecurityGroupAttributes(),
			},
			"backend_vm_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"application_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.appStickyPolicyComputedAttributes(),
			},
			"load_balancer_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.stickyPolicyComputedAttributes(),
			},
			"load_balancer_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_type": schema.StringAttribute{
				Computed: true,
			},
			"policy_names": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (r *loadBalancerAttributesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data loadBalancerAttributesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	name := data.LoadBalancerName.ValueString()
	updateReq := osc.UpdateLoadBalancerRequest{
		LoadBalancerName: name,
	}

	if fwhelpers.IsSet(data.LoadBalancerPort) {
		updateReq.LoadBalancerPort = new(int(data.LoadBalancerPort.ValueInt64()))
	}
	if fwhelpers.IsSet(data.ServerCertificateId) {
		updateReq.ServerCertificateId = data.ServerCertificateId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.PolicyNames) {
		policyNames, diag := to.Slice[string](ctx, data.PolicyNames)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		updateReq.PolicyNames = &policyNames
	}
	if fwhelpers.IsSet(data.AccessLog) {
		accessLogs, diag := to.Model[loadBalancerAccessLogModel](ctx, data.AccessLog)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		updateReq.AccessLog = r.expandAccessLog(accessLogs)
	}
	if fwhelpers.IsSet(data.HealthCheck) {
		healthChecks, diag := to.Model[loadBalancerHealthCheckModel](ctx, data.HealthCheck)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		updateReq.HealthCheck = r.expandHealthCheck(healthChecks)
	}

	updateReqs := []osc.UpdateLoadBalancerRequest{updateReq}
	// API returns an error if both AccessLog and HealthCheck are set, so we need to split the update into two requests
	if fwhelpers.IsSet(data.AccessLog) && fwhelpers.IsSet(data.HealthCheck) {
		updateReqs[0].HealthCheck = nil

		updateReqs = append(updateReqs, osc.UpdateLoadBalancerRequest{
			LoadBalancerName: name,
			HealthCheck:      updateReq.HealthCheck,
		})
	}

	var waitResp *osc.ReadLoadBalancersResponse
	for _, updateReq := range updateReqs {
		_, err := r.Client.UpdateLoadBalancer(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(loadBalancerAttributesErrCreate, err.Error())
			return
		}

		stateConf := &stateconf.StateChangeConf[osc.LoadBalancerState]{
			Pending: stateconf.States(osc.LoadBalancerStateStarting, osc.LoadBalancerStateProvisioning, osc.LoadBalancerStateReloading, osc.LoadBalancerStateReconfiguring),
			Target:  stateconf.States(osc.LoadBalancerStateActive),
			Timeout: timeout,
			Refresh: r.refreshFunc(name),
		}
		waitRespAny, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError(loadBalancerAttributesErrWait, err.Error())
			return
		}
		waitResp = waitRespAny.(*osc.ReadLoadBalancersResponse)
	}

	data.Id = to.String(name)
	data.LoadBalancerName = to.String(name)
	data.RequestId = to.String(waitResp.ResponseContext.RequestId)

	stateData, err := r.flatten(ctx, data, (*waitResp.LoadBalancers)[0])
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerAttributesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data loadBalancerAttributesModel
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

func (r *loadBalancerAttributesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data loadBalancerAttributesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(loadBalancerAttributesErrUpdate, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerAttributesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data loadBalancerAttributesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The only configurable attribute which can be "deleted" is AccessLog
	if !fwhelpers.IsSet(data.AccessLog) {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.UpdateLoadBalancerRequest{
		LoadBalancerName: data.LoadBalancerName.ValueString(),
		AccessLog: &osc.AccessLog{
			IsEnabled: false,
		},
	}

	// 409 with code 6031 is returned when the LBU is in an invalid state
	_, err := oapihelpers.RetryOnCodes(ctx, []string{"6031"}, func() (any, error) {
		return r.Client.UpdateLoadBalancer(ctx, deleteReq, options.WithRetryTimeout(timeout))
	}, timeout)
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerAttributesErrDelete, err.Error())
	}
}

func (r *loadBalancerAttributesResource) read(ctx context.Context, timeout time.Duration, data loadBalancerAttributesModel) (loadBalancerAttributesModel, error) {
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

func (r *loadBalancerAttributesResource) flatten(ctx context.Context, data loadBalancerAttributesModel, lb osc.LoadBalancer) (loadBalancerAttributesModel, error) {
	tags, diag := flattenOAPIComputedTagsFW(ctx, lb.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.Tags = tags

	accessLogModel := r.flattenAccessLog(lb.AccessLog)
	accessLog, diags := to.Object(ctx, accessLogModel[0])
	if diags.HasError() {
		return data, from.Diag(diags)
	}

	healthCheckModel := r.flattenHealthCheck(lb.HealthCheck)
	healthCheck, diags := to.Object(ctx, healthCheckModel[0])
	if diags.HasError() {
		return data, from.Diag(diags)
	}

	listenersModel, diag := r.flattenListeners(ctx, lb.Listeners)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	listeners, diag := to.ListFromAttrType(ctx, listenersModel, to.ObjType(loadBalancerListenerAttrTypes), to.ZeroValueAsEmpty)
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

	backendVmIds, diag := to.List(ctx, lb.BackendVmIds, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	securityGroups, diag := to.List(ctx, lb.SecurityGroups, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	subnets, diag := to.List(ctx, lb.Subnets, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	subregionNames, diag := to.List(ctx, lb.SubregionNames, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.AccessLog = accessLog
	data.HealthCheck = healthCheck
	data.Listeners = listeners
	data.SourceSecurityGroup = sourceSecurityGroup
	data.ApplicationStickyCookiePolicies = appPolicies
	data.LoadBalancerStickyCookiePolicies = stickyPolicies
	data.Tags = tags
	data.BackendVmIds = backendVmIds
	data.Subnets = subnets
	data.SubregionNames = subregionNames
	data.SecurityGroups = securityGroups
	data.DnsName = to.String(lb.DnsName)
	data.LoadBalancerName = to.String(lb.LoadBalancerName)
	data.LoadBalancerType = to.String(lb.LoadBalancerType)
	data.Id = to.String(lb.LoadBalancerName)

	return data, nil
}

func (r *loadBalancerAttributesResource) expandAccessLog(model loadBalancerAccessLogModel) *osc.AccessLog {
	access := &osc.AccessLog{}

	if fwhelpers.IsSet(model.IsEnabled) {
		access.IsEnabled = model.IsEnabled.ValueBool()
	}
	if fwhelpers.IsSet(model.PublicationInterval) {
		access.PublicationInterval = new(int(model.PublicationInterval.ValueInt64()))
	}
	if fwhelpers.IsSet(model.OsuBucketName) && model.OsuBucketName.ValueString() != "" {
		access.OsuBucketName = model.OsuBucketName.ValueStringPointer()
	}
	if fwhelpers.IsSet(model.OsuBucketPrefix) && model.OsuBucketPrefix.ValueString() != "" {
		access.OsuBucketPrefix = model.OsuBucketPrefix.ValueStringPointer()
	}

	return access
}

func (r *loadBalancerAttributesResource) expandHealthCheck(model loadBalancerHealthCheckModel) *osc.HealthCheck {
	healthCheck := &osc.HealthCheck{
		Port:     int(model.Port.ValueInt64()),
		Protocol: model.Protocol.ValueString(),
	}

	if fwhelpers.IsSet(model.HealthyThreshold) {
		healthCheck.HealthyThreshold = int(model.HealthyThreshold.ValueInt64())
	}
	if fwhelpers.IsSet(model.UnhealthyThreshold) {
		healthCheck.UnhealthyThreshold = int(model.UnhealthyThreshold.ValueInt64())
	}
	if fwhelpers.IsSet(model.CheckInterval) {
		healthCheck.CheckInterval = int(model.CheckInterval.ValueInt64())
	}
	if fwhelpers.IsSet(model.Timeout) {
		healthCheck.Timeout = int(model.Timeout.ValueInt64())
	}
	if fwhelpers.IsSet(model.Path) {
		healthCheck.Path = model.Path.ValueStringPointer()
	}

	return healthCheck
}
