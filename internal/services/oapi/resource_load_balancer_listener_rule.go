package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
)

var (
	_ resource.Resource                     = &loadBalancerListenerRuleResource{}
	_ resource.ResourceWithConfigure        = &loadBalancerListenerRuleResource{}
	_ resource.ResourceWithConfigValidators = &loadBalancerListenerRuleResource{}
)

const (
	loadBalancerListenerRuleErrCreate = "Unable to create Load Balancer Listener Rule"
	loadBalancerListenerRuleErrUpdate = "Unable to update Load Balancer Listener Rule"
	loadBalancerListenerRuleErrDelete = "Unable to delete Load Balancer Listener Rule"
)

type loadBalancerListenerRuleModel struct {
	VmIds        types.Set      `tfsdk:"vm_ids"`
	Listener     types.List     `tfsdk:"listener"`
	RequestId    types.String   `tfsdk:"request_id"`
	ListenerRule types.List     `tfsdk:"listener_rule"`
	Id           types.String   `tfsdk:"id"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

type loadBalancerListenerRuleListenerModel struct {
	LoadBalancerName types.String `tfsdk:"load_balancer_name"`
	LoadBalancerPort types.Int64  `tfsdk:"load_balancer_port"`
}

type loadBalancerListenerRuleDetailModel struct {
	Action           types.String `tfsdk:"action"`
	HostNamePattern  types.String `tfsdk:"host_name_pattern"`
	ListenerRuleName types.String `tfsdk:"listener_rule_name"`
	ListenerRuleId   types.Int64  `tfsdk:"listener_rule_id"`
	ListenerId       types.Int64  `tfsdk:"listener_id"`
	PathPattern      types.String `tfsdk:"path_pattern"`
	Priority         types.Int64  `tfsdk:"priority"`
}

type loadBalancerListenerRuleResource struct {
	loadBalancerCommon
}

func NewResourceLoadBalancerListenerRule() resource.Resource {
	return &loadBalancerListenerRuleResource{}
}

func (r *loadBalancerListenerRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *loadBalancerListenerRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_listener_rule"
}

func (r *loadBalancerListenerRuleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("listener_rule").AtAnyListIndex().AtName("host_name_pattern"),
			path.MatchRoot("listener_rule").AtAnyListIndex().AtName("path_pattern"),
		),
	}
}

func (r *loadBalancerListenerRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"listener": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeBetween(1, 1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"load_balancer_name": schema.StringAttribute{
							Required: true,
						},
						"load_balancer_port": schema.Int64Attribute{
							Required: true,
						},
					},
				},
			},
			"listener_rule": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeBetween(1, 1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
						},
						"host_name_pattern": schema.StringAttribute{
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								fwtypes.EmptyStringAsNull(),
							},
							Default: stringdefault.StaticString(""),
						},
						"listener_rule_name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"listener_rule_id": schema.Int64Attribute{
							Computed: true,
						},
						"listener_id": schema.Int64Attribute{
							Computed: true,
						},
						"path_pattern": schema.StringAttribute{
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								fwtypes.EmptyStringAsNull(),
							},
							Default: stringdefault.StaticString(""),
						},
						"priority": schema.Int64Attribute{
							Required: true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"vm_ids": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
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

func (r *loadBalancerListenerRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data loadBalancerListenerRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	vmIds, diag := to.Slice[string](ctx, data.VmIds)
	diags.Append(diag...)

	listeners, diag := to.Slice[loadBalancerListenerRuleListenerModel](ctx, data.Listener)
	diags.Append(diag...)

	rules, diag := to.Slice[loadBalancerListenerRuleDetailModel](ctx, data.ListenerRule)
	diags.Append(diag...)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateListenerRuleRequest{
		VmIds: vmIds,
		Listener: osc.LoadBalancerLight{
			LoadBalancerName: listeners[0].LoadBalancerName.ValueString(),
			LoadBalancerPort: int(listeners[0].LoadBalancerPort.ValueInt64()),
		},
		ListenerRule: osc.ListenerRuleForCreation{
			ListenerRuleName: rules[0].ListenerRuleName.ValueString(),
			Priority:         int(rules[0].Priority.ValueInt64()),
		},
	}

	if fwhelpers.IsSet(rules[0].Action) {
		createReq.ListenerRule.Action = rules[0].Action.ValueStringPointer()
	}
	if fwhelpers.IsSet(rules[0].HostNamePattern) {
		createReq.ListenerRule.HostNamePattern = rules[0].HostNamePattern.ValueStringPointer()
	}
	if fwhelpers.IsSet(rules[0].PathPattern) {
		createReq.ListenerRule.PathPattern = rules[0].PathPattern.ValueStringPointer()
	}

	createResp, err := r.Client.CreateListenerRule(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerListenerRuleErrCreate, err.Error())
		return
	}
	data.Id = to.String(createResp.ListenerRule.ListenerRuleName)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)

	stateData, err := r.flatten(ctx, data, *createResp.ListenerRule)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *loadBalancerListenerRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data loadBalancerListenerRuleModel
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

func (r *loadBalancerListenerRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state loadBalancerListenerRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	planRules, diags := to.Slice[loadBalancerListenerRuleDetailModel](ctx, plan.ListenerRule)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	stateRules, diags := to.Slice[loadBalancerListenerRuleDetailModel](ctx, state.ListenerRule)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateListenerRuleRequest{
		ListenerRuleName: state.Id.ValueString(),
	}
	if fwhelpers.HasChange(planRules[0].HostNamePattern, stateRules[0].HostNamePattern) {
		updateReq.HostPattern = planRules[0].HostNamePattern.ValueStringPointer()
	}
	if fwhelpers.HasChange(planRules[0].PathPattern, stateRules[0].PathPattern) {
		updateReq.PathPattern = planRules[0].PathPattern.ValueStringPointer()
	}

	updateResp, err := r.Client.UpdateListenerRule(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerListenerRuleErrUpdate, err.Error())
		return
	}

	plan.RequestId = to.String(updateResp.ResponseContext.RequestId)

	newData, err := r.flatten(ctx, plan, *updateResp.ListenerRule)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *loadBalancerListenerRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data loadBalancerListenerRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.DeleteListenerRuleRequest{
		ListenerRuleName: data.Id.ValueString(),
	}
	_, err := r.Client.DeleteListenerRule(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(loadBalancerListenerRuleErrDelete, err.Error())
		return
	}
}

func (r *loadBalancerListenerRuleResource) read(ctx context.Context, timeout time.Duration, data loadBalancerListenerRuleModel) (loadBalancerListenerRuleModel, error) {
	req := osc.ReadListenerRulesRequest{
		Filters: &osc.FiltersListenerRule{
			ListenerRuleNames: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadListenerRules(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	rules := ptr.From(resp.ListenerRules)
	if len(rules) == 0 {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return r.flatten(ctx, data, rules[0])
}

func (r *loadBalancerListenerRuleResource) flatten(ctx context.Context, data loadBalancerListenerRuleModel, rule osc.ListenerRule) (loadBalancerListenerRuleModel, error) {
	oldRules, diags := to.Slice[loadBalancerListenerRuleDetailModel](ctx, data.ListenerRule)
	if diags.HasError() {
		return data, from.Diag(diags)
	}

	// The API returns nil after a pattern is cleared (setting "" in the configuration). Preserve an explicit "" from
	// the plan or state so Terraform remains consistent; non-nil API values still win
	hostNamePattern := to.String(rule.HostNamePattern)
	pathPattern := to.String(rule.PathPattern)
	if len(oldRules) == 1 {
		oldHostNamePattern := oldRules[0].HostNamePattern
		if hostNamePattern.IsNull() && fwhelpers.IsSet(oldHostNamePattern) && oldHostNamePattern.ValueString() == "" {
			hostNamePattern = oldHostNamePattern
		}

		oldPathPattern := oldRules[0].PathPattern
		if pathPattern.IsNull() && fwhelpers.IsSet(oldPathPattern) && oldPathPattern.ValueString() == "" {
			pathPattern = oldPathPattern
		}
	}

	model := []loadBalancerListenerRuleDetailModel{{
		Action:           to.String(rule.Action),
		HostNamePattern:  hostNamePattern,
		ListenerRuleName: to.String(rule.ListenerRuleName),
		ListenerRuleId:   to.Int64(rule.ListenerRuleId),
		ListenerId:       to.Int64(rule.ListenerId),
		PathPattern:      pathPattern,
		Priority:         to.Int64(rule.Priority),
	}}

	listenerRule, diags := to.ListObject(ctx, model, to.ZeroValueAsEmpty)
	if diags.HasError() {
		return data, from.Diag(diags)
	}

	vmIds, diags := to.Set(ctx, ptr.From(rule.VmIds), to.ZeroValueAsEmpty)
	if diags.HasError() {
		return data, from.Diag(diags)
	}

	data.ListenerRule = listenerRule
	data.VmIds = vmIds
	data.Id = to.String(rule.ListenerRuleName)

	return data, nil
}
