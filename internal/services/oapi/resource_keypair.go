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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource               = &resourceKeypair{}
	_ resource.ResourceWithConfigure  = &resourceKeypair{}
	_ resource.ResourceWithModifyPlan = &resourceKeypair{}
)

const (
	keypairErrCreate = "Unable to create Keypair"
	keypairErrDelete = "Unable to delete Keypair"
)

type KeypairModel struct {
	KeypairFingerprint types.String   `tfsdk:"keypair_fingerprint"`
	PrivateKey         types.String   `tfsdk:"private_key"`
	KeypairName        types.String   `tfsdk:"keypair_name"`
	KeypairType        types.String   `tfsdk:"keypair_type"`
	KeypairId          types.String   `tfsdk:"keypair_id"`
	PublicKey          types.String   `tfsdk:"public_key"`
	RequestId          types.String   `tfsdk:"request_id"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	TagsModel
}

type keypairCommon struct {
	Client *osc.Client
}

type resourceKeypair struct {
	keypairCommon
}

func NewResourceKeypair() resource.Resource {
	return &resourceKeypair{}
}

func (r *resourceKeypair) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceKeypair) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource. "+
				"And users will not be able to use this credentials.",
		)
	}
}

func (r *resourceKeypair) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	keypairId := req.ID
	if keypairId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import keypair identifier Got: %v", req.ID),
		)
		return
	}

	var data KeypairModel
	var timeouts timeouts.Value
	data.KeypairId = to.String(keypairId)
	data.Id = to.String(keypairId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceKeypair) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypair"
}

func (r *resourceKeypair) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags": TagsSchemaFW(),
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"keypair_id": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Computed: true,
			},
			"keypair_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keypair_type": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keypair_fingerprint": schema.StringAttribute{
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

func (r *resourceKeypair) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KeypairModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateKeypairRequest{}
	createReq.KeypairName = data.KeypairName.ValueString()

	if !data.PublicKey.IsUnknown() && !data.PublicKey.IsNull() {
		createReq.PublicKey = data.PublicKey.ValueStringPointer()
	}

	createResp, err := r.Client.CreateKeypair(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(keypairErrCreate, err.Error())
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	keypair := ptr.From(createResp.Keypair)

	data.Id = to.String(keypair.KeypairId)
	data.KeypairName = to.String(ptr.From(keypair.KeypairName))
	if createReq.PublicKey != nil {
		data.PublicKey = to.String(createReq.PublicKey)
	}

	stateData := r.flattenCreate(data, keypair)
	diag := resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	diag = createOAPITagsFW(ctx, r.Client, timeout, data.Tags, *keypair.KeypairId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data.PrivateKey = to.String(ptr.From(keypair.PrivateKey))

	stateData, err = r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceKeypair) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeypairModel

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

func (r *resourceKeypair) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData KeypairModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.KeypairId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData.Timeouts = planData.Timeouts
	stateData, err := r.read(ctx, timeout, stateData)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceKeypair) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeypairModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteKeypairRequest{
		KeypairId: data.KeypairId.ValueStringPointer(),
	}

	_, err := r.Client.DeleteKeypair(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(keypairErrDelete, err.Error())
	}
}

func (r *keypairCommon) read(ctx context.Context, timeout time.Duration, data KeypairModel) (KeypairModel, error) {
	keypairFilters := osc.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}

	readReq := osc.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	readResp, err := r.Client.ReadKeypairs(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if readResp.Keypairs == nil || len(*readResp.Keypairs) == 0 {
		return data, ErrResourceEmpty
	}

	keypair := (*readResp.Keypairs)[0]
	tags, diag := flattenOAPITagsFW(ctx, ptr.From(keypair.Tags))
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.Tags = tags
	data.KeypairFingerprint = to.String(ptr.From(keypair.KeypairFingerprint))
	data.KeypairName = to.String(ptr.From(keypair.KeypairName))
	data.KeypairType = to.String(ptr.From(keypair.KeypairType))
	data.KeypairId = to.String(keypair.KeypairId)

	return data, nil
}

func (r *keypairCommon) flattenCreate(data KeypairModel, keypair osc.KeypairCreated) KeypairModel {
	data.KeypairFingerprint = to.String(ptr.From(keypair.KeypairFingerprint))
	data.KeypairId = to.String(keypair.KeypairId)
	data.KeypairName = to.String(ptr.From(keypair.KeypairName))
	data.KeypairType = to.String(ptr.From(keypair.KeypairType))
	data.PrivateKey = to.String(ptr.From(keypair.PrivateKey))

	return data
}
