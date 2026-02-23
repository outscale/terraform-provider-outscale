package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                = &resoureCa{}
	_ resource.ResourceWithConfigure   = &resoureCa{}
	_ resource.ResourceWithImportState = &resoureCa{}
)

type caModel struct {
	CaPem         types.String   `tfsdk:"ca_pem"`
	CaFingerprint types.String   `tfsdk:"ca_fingerprint"`
	CaId          types.String   `tfsdk:"ca_id"`
	Description   types.String   `tfsdk:"description"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	RequestId     types.String   `tfsdk:"request_id"`
}

type resoureCa struct {
	Client *osc.Client
}

func NewResourceCa() resource.Resource {
	return &resoureCa{}
}

func (r *resoureCa) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resoureCa) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resoureCa) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ca"
}

func (r *resoureCa) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"ca_pem": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ca_fingerprint": schema.StringAttribute{
				Computed: true,
			},
			"ca_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
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

func (r *resoureCa) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data caModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreateCaRequest{
		CaPem: data.CaPem.ValueString(),
	}

	if fwhelpers.IsSet(data.Description) {
		createReq.Description = data.Description.ValueStringPointer()
	}

	createResp, err := r.Client.CreateCa(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Ca",
			err.Error(),
		)
		return
	}
	data.Id = to.String(createResp.Ca.CaId)

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Ca state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resoureCa) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data caModel
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
			"Unable to set Ca API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resoureCa) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData caModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	updateReq := osc.UpdateCaRequest{
		CaId: stateData.Id.ValueString(),
	}

	if fwhelpers.HasChange(planData.Description, stateData.Description) {
		updateReq.Description = planData.Description.ValueStringPointer()
	}

	_, err := r.Client.UpdateCa(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Ca",
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
			"Unable to set Ca API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *resoureCa) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data caModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	deleteReq := osc.DeleteCaRequest{
		CaId: data.Id.ValueString(),
	}

	_, err := r.Client.DeleteCa(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Ca",
			err.Error(),
		)
	}
}

func (r *resoureCa) read(ctx context.Context, timeout time.Duration, data caModel) (caModel, error) {
	req := osc.ReadCasRequest{
		Filters: &osc.FiltersCa{
			CaIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadCas(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Cas == nil || len(*resp.Cas) == 0 {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)
	ca := (*resp.Cas)[0]

	data.Id = to.String(ca.CaId)
	data.CaFingerprint = to.String(ptr.From(ca.CaFingerprint))
	data.CaId = to.String(ca.CaId)
	data.Description = to.String(ptr.From(ca.Description))

	return data, nil
}
