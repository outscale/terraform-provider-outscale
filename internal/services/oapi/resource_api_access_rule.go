package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
)

var (
	_ resource.Resource                = &resourceApiAccessRule{}
	_ resource.ResourceWithConfigure   = &resourceApiAccessRule{}
	_ resource.ResourceWithImportState = &resourceApiAccessRule{}
)

type apiAccessRuleModel struct {
	ApiAccessRuleId types.String   `tfsdk:"api_access_rule_id"`
	CaIds           types.Set      `tfsdk:"ca_ids"`
	Cns             types.Set      `tfsdk:"cns"`
	Description     types.String   `tfsdk:"description"`
	IpRanges        types.Set      `tfsdk:"ip_ranges"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	Id              types.String   `tfsdk:"id"`
	RequestId       types.String   `tfsdk:"request_id"`
}

type resourceApiAccessRule struct {
	Client *osc.Client
}

func NewResourceApiAccessRule() resource.Resource {
	return &resourceApiAccessRule{}
}

func (r *resourceApiAccessRule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceApiAccessRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceApiAccessRule) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_access_rule"
}

func (r *resourceApiAccessRule) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"api_access_rule_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ca_ids": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"cns": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"ip_ranges": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.Any(
							validatorstring.IsCIDR(),
							validatorstring.IsIP(),
						),
					),
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

func (r *resourceApiAccessRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apiAccessRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreateApiAccessRuleRequest{}

	if fwhelpers.IsSet(data.CaIds) {
		ids, diag := to.Slice[string](ctx, data.CaIds)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.CaIds = &ids
	}
	if fwhelpers.IsSet(data.IpRanges) {
		ips, diag := to.Slice[string](ctx, data.IpRanges)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.IpRanges = &ips
	}
	if fwhelpers.IsSet(data.Cns) {
		cns, diag := to.Slice[string](ctx, data.Cns)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.Cns = &cns
	}
	if fwhelpers.IsSet(data.Description) {
		createReq.Description = data.Description.ValueStringPointer()
	}

	createResp, err := r.Client.CreateApiAccessRule(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Api Access Rule",
			err.Error(),
		)
		return
	}
	data.Id = to.String(createResp.ApiAccessRule.ApiAccessRuleId)

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Api Access Rule state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceApiAccessRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiAccessRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to get Api Access Rule response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceApiAccessRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData apiAccessRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	updateReq := osc.UpdateApiAccessRuleRequest{
		ApiAccessRuleId: stateData.Id.ValueString(),
	}

	if fwhelpers.IsSet(planData.CaIds) {
		caIds, diag := to.Slice[string](ctx, planData.CaIds)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.CaIds = &caIds
	}
	if fwhelpers.IsSet(planData.Cns) {
		cns, diag := to.Slice[string](ctx, planData.Cns)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.Cns = &cns
	}
	if fwhelpers.IsSet(planData.IpRanges) {
		ipRanges, diag := to.Slice[string](ctx, planData.IpRanges)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.IpRanges = &ipRanges
	}
	if fwhelpers.IsSet(planData.Description) {
		updateReq.Description = planData.Description.ValueStringPointer()
	}

	_, err := r.Client.UpdateApiAccessRule(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Api Access Rule",
			err.Error(),
		)
	}

	newStateData, err := r.read(ctx, timeout, stateData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read Api Access Rule response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *resourceApiAccessRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data apiAccessRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	deleteReq := osc.DeleteApiAccessRuleRequest{
		ApiAccessRuleId: data.Id.ValueString(),
	}

	_, err := r.Client.DeleteApiAccessRule(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Api Access Rule",
			err.Error(),
		)
	}
}

func (r *resourceApiAccessRule) read(ctx context.Context, timeout time.Duration, data apiAccessRuleModel) (apiAccessRuleModel, error) {
	req := osc.ReadApiAccessRulesRequest{
		Filters: &osc.FiltersApiAccessRule{
			ApiAccessRuleIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadApiAccessRules(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.ApiAccessRules == nil || len(*resp.ApiAccessRules) == 0 {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)
	rule := (*resp.ApiAccessRules)[0]

	caIds, diag := to.Set(ctx, ptr.From(rule.CaIds))
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert ca_ids into a set: %v", diag.Errors())
	}
	cns, diag := to.Set(ctx, ptr.From(rule.Cns))
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert cns into a set: %v", diag.Errors())
	}
	ipRanges, diag := to.Set(ctx, ptr.From(rule.IpRanges))
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert ip_ranges into a set: %v", diag.Errors())
	}

	data.CaIds = caIds
	data.Cns = cns
	data.IpRanges = ipRanges
	data.Id = to.String(rule.ApiAccessRuleId)
	data.ApiAccessRuleId = to.String(rule.ApiAccessRuleId)
	data.Description = to.String(ptr.From(rule.Description))

	return data, nil
}
