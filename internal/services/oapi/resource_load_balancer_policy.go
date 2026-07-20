package oapi

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
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
	_ resource.Resource              = &loadBalancerPolicyResource{}
	_ resource.ResourceWithConfigure = &loadBalancerPolicyResource{}
)

const (
	loadBalancerPolicyErrCreate = "Unable to create Load Balancer Policy"
	loadBalancerPolicyErrDelete = "Unable to delete Load Balancer Policy"
)

type loadBalancerPolicyModel struct {
	PolicyName                       types.String   `tfsdk:"policy_name"`
	AccessLog                        types.List     `tfsdk:"access_log"`
	HealthCheck                      types.List     `tfsdk:"health_check"`
	ApplicationStickyCookiePolicies  types.List     `tfsdk:"application_sticky_cookie_policies"`
	LoadBalancerStickyCookiePolicies types.List     `tfsdk:"load_balancer_sticky_cookie_policies"`
	Listeners                        types.List     `tfsdk:"listeners"`
	SourceSecurityGroup              types.List     `tfsdk:"source_security_group"`
	PublicIp                         types.String   `tfsdk:"public_ip"`
	SecuredCookies                   types.Bool     `tfsdk:"secured_cookies"`
	NetId                            types.String   `tfsdk:"net_id"`
	BackendVmIds                     types.List     `tfsdk:"backend_vm_ids"`
	SubregionNames                   types.List     `tfsdk:"subregion_names"`
	LoadBalancerType                 types.String   `tfsdk:"load_balancer_type"`
	SecurityGroups                   types.Set      `tfsdk:"security_groups"`
	Subnets                          types.List     `tfsdk:"subnets"`
	DnsName                          types.String   `tfsdk:"dns_name"`
	PolicyType                       types.String   `tfsdk:"policy_type"`
	LoadBalancerName                 types.String   `tfsdk:"load_balancer_name"`
	CookieExpirationPeriod           types.Int64    `tfsdk:"cookie_expiration_period"`
	CookieName                       types.String   `tfsdk:"cookie_name"`
	RequestId                        types.String   `tfsdk:"request_id"`
	Id                               types.String   `tfsdk:"id"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
	TagsComputedModel
}

type loadBalancerPolicyResource struct {
	loadBalancerCommon
}

func NewResourceLoadBalancerPolicy() resource.Resource {
	return &loadBalancerPolicyResource{}
}

func (r *loadBalancerPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *loadBalancerPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_policy"
}

func ifApplicationPolicy(_ context.Context, req validator.StringRequest) bool {
	return fwhelpers.IsSet(req.ConfigValue) && req.ConfigValue.ValueString() == "app"
}

func (r *loadBalancerPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"policy_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9A-Za-z-]+$`), "only alphanumeric characters and hyphens allowed"),
				},
			},
			"access_log": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.accessLogComputedAttributes(),
			},
			"health_check": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.healthCheckComputedAttributes(),
			},
			"application_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.appStickyPolicyComputedAttributes(),
			},
			"load_balancer_sticky_cookie_policies": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.stickyPolicyComputedAttributes(),
			},
			"listeners": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.listenersComputedAttributes(),
			},
			"source_security_group": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: r.sourceSecurityGroupAttributes(),
			},
			"public_ip": schema.StringAttribute{
				Computed: true,
			},
			"secured_cookies": schema.BoolAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
			},
			"backend_vm_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"subregion_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
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
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("app", "load_balancer"),
					validatorstring.AlsoRequiresIf(
						path.MatchRoot("cookie_name"),
						ifApplicationPolicy,
						"Missing cookie name",
						"'cookie_name' is required when 'policy_type' is 'app'.",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cookie_expiration_period": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cookie_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (c *loadBalancerCommon) listenersComputedAttributes() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"backend_port": schema.Int64Attribute{
				Computed: true,
			},
			"backend_protocol": schema.StringAttribute{
				Computed: true,
			},
			"load_balancer_port": schema.Int64Attribute{
				Computed: true,
			},
			"load_balancer_protocol": schema.StringAttribute{
				Computed: true,
			},
			"server_certificate_id": schema.StringAttribute{
				Computed: true,
			},
			"policy_names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *loadBalancerPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data loadBalancerPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateLoadBalancerPolicyRequest{
		LoadBalancerName: data.LoadBalancerName.ValueString(),
		PolicyName:       data.PolicyName.ValueString(),
		PolicyType:       data.PolicyType.ValueString(),
	}
	if fwhelpers.IsSet(data.CookieName) {
		createReq.CookieName = data.CookieName.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.CookieExpirationPeriod) {
		createReq.CookieExpirationPeriod = new(int(data.CookieExpirationPeriod.ValueInt64()))
	}

	createResp, err := r.Client.CreateLoadBalancerPolicy(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerPolicyErrCreate, err.Error())
		return
	}
	data.Id = to.String(id.UniqueId())
	data.RequestId = to.String(createResp.ResponseContext.RequestId)

	stateData, err := r.flatten(ctx, data, *createResp.LoadBalancer)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data loadBalancerPolicyModel
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

func (r *loadBalancerPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *loadBalancerPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data loadBalancerPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.DeleteLoadBalancerPolicyRequest{
		LoadBalancerName: data.LoadBalancerName.ValueString(),
		PolicyName:       data.PolicyName.ValueString(),
	}

	_, err := r.Client.DeleteLoadBalancerPolicy(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerPolicyErrDelete, err.Error())
	}
}

func (r *loadBalancerPolicyResource) read(ctx context.Context, timeout time.Duration, data loadBalancerPolicyModel) (loadBalancerPolicyModel, error) {
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
	readResp := respAny.(*osc.ReadLoadBalancersResponse)
	lb := (*readResp.LoadBalancers)[0]

	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	return r.flatten(ctx, data, lb)
}

func (r *loadBalancerPolicyResource) flatten(ctx context.Context, data loadBalancerPolicyModel, lb osc.LoadBalancer) (loadBalancerPolicyModel, error) {
	appPolicy, foundAppPolicy := lo.Find(lb.ApplicationStickyCookiePolicies, func(v osc.ApplicationStickyCookiePolicy) bool {
		return ptr.From(v.PolicyName) == data.PolicyName.ValueString()
	})
	lbuPolicy, foundLbuPolicy := lo.Find(lb.LoadBalancerStickyCookiePolicies, func(v osc.LoadBalancerStickyCookiePolicy) bool {
		return ptr.From(v.PolicyName) == data.PolicyName.ValueString()
	})
	if !foundAppPolicy && !foundLbuPolicy {
		return data, ErrResourceEmpty
	}

	tags, diag := flattenOAPIComputedTagsFW(ctx, lb.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	accessLogModel := r.flattenAccessLog(lb.AccessLog)
	accessLog, diag := to.ListObject(ctx, accessLogModel, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	healthCheckModel := r.flattenHealthCheck(lb.HealthCheck)
	healthCheck, diag := to.ListObject(ctx, healthCheckModel, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
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

	subregionNames, diag := to.List(ctx, lb.SubregionNames, to.ZeroValueAsEmpty)
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

	data.Tags = tags
	data.AccessLog = accessLog
	data.HealthCheck = healthCheck
	data.Listeners = listeners
	data.SourceSecurityGroup = sourceSecurityGroup
	data.ApplicationStickyCookiePolicies = appPolicies
	data.LoadBalancerStickyCookiePolicies = stickyPolicies
	data.BackendVmIds = backendVmIds
	data.SubregionNames = subregionNames
	data.Subnets = subnets
	data.SecurityGroups = securityGroups
	data.LoadBalancerName = to.String(lb.LoadBalancerName)
	data.LoadBalancerType = to.String(lb.LoadBalancerType)
	data.DnsName = to.String(lb.DnsName)
	data.PublicIp = to.String(ptr.From(lb.PublicIp))
	data.SecuredCookies = to.Bool(lb.SecuredCookies)
	data.NetId = to.String(ptr.From(lb.NetId))
	data.CookieName = to.String(appPolicy.CookieName)
	data.CookieExpirationPeriod = to.Int64(lbuPolicy.CookieExpirationPeriod)

	return data, nil
}
