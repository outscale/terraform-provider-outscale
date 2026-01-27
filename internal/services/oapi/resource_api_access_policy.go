package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                   = &resourceApiAccessPolicy{}
	_ resource.ResourceWithConfigure      = &resourceApiAccessPolicy{}
	_ resource.ResourceWithValidateConfig = &resourceApiAccessPolicy{}
)

type apiAccessPolicyModel struct {
	MaxAccessKeyExpirationSeconds types.Int64    `tfsdk:"max_access_key_expiration_seconds"`
	RequireTrustedEnv             types.Bool     `tfsdk:"require_trusted_env"`
	Id                            types.String   `tfsdk:"id"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
	RequestId                     types.String   `tfsdk:"request_id"`
}

type resourceApiAccessPolicy struct {
	Client *osc.APIClient
}

func NewResourceApiAccessPolicy() resource.Resource {
	return &resourceApiAccessPolicy{}
}

func (r *resourceApiAccessPolicy) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
}

func (r *resourceApiAccessPolicy) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_access_policy"
}

func (r *resourceApiAccessPolicy) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data apiAccessPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if fwhelpers.IsSet(data.RequireTrustedEnv) && data.RequireTrustedEnv.ValueBool() {
		if fwhelpers.IsSet(data.MaxAccessKeyExpirationSeconds) && data.MaxAccessKeyExpirationSeconds.ValueInt64() == 0 {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"When require_trusted_env is true, max_access_key_expiration_seconds must be greater than 0",
			)
		}
	}
}

func (r *resourceApiAccessPolicy) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"max_access_key_expiration_seconds": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 3153600000),
				},
			},
			"require_trusted_env": schema.BoolAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *resourceApiAccessPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apiAccessPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.UpdateApiAccessPolicyRequest{
		MaxAccessKeyExpirationSeconds: data.MaxAccessKeyExpirationSeconds.ValueInt64(),
		RequireTrustedEnv:             data.RequireTrustedEnv.ValueBool(),
	}

	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.ApiAccessPolicyApi.UpdateApiAccessPolicy(ctx).UpdateApiAccessPolicyRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to modify Api Access Policy",
			err.Error(),
		)
		return
	}
	data.Id = to.String(id.UniqueId())

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Api Access Policy response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceApiAccessPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiAccessPolicyModel
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
			"Unable to get Api Access Policy response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceApiAccessPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData apiAccessPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	updateReq := osc.UpdateApiAccessPolicyRequest{}

	if fwhelpers.HasChange(planData.MaxAccessKeyExpirationSeconds, stateData.MaxAccessKeyExpirationSeconds) {
		updateReq.SetMaxAccessKeyExpirationSeconds(planData.MaxAccessKeyExpirationSeconds.ValueInt64())
	}
	if fwhelpers.HasChange(planData.RequireTrustedEnv, stateData.RequireTrustedEnv) {
		updateReq.SetRequireTrustedEnv(planData.RequireTrustedEnv.ValueBool())
	}

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.ApiAccessPolicyApi.UpdateApiAccessPolicy(ctx).UpdateApiAccessPolicyRequest(updateReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Api Access Policy",
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
			"Unable to get Api Access Policy response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *resourceApiAccessPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *resourceApiAccessPolicy) read(ctx context.Context, timeout time.Duration, data apiAccessPolicyModel) (apiAccessPolicyModel, error) {
	req := osc.ReadApiAccessPolicyRequest{}

	var resp osc.ReadApiAccessPolicyResponse
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.ApiAccessPolicyApi.ReadApiAccessPolicy(ctx).ReadApiAccessPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return data, err
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)
	policy := resp.GetApiAccessPolicy()

	data.MaxAccessKeyExpirationSeconds = to.Int64(policy.MaxAccessKeyExpirationSeconds)
	data.RequireTrustedEnv = to.Bool(policy.RequireTrustedEnv)

	return data, nil
}
