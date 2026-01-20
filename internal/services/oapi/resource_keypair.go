package oapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource               = &resourceKeypair{}
	_ resource.ResourceWithConfigure  = &resourceKeypair{}
	_ resource.ResourceWithModifyPlan = &resourceKeypair{}
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

type resourceKeypair struct {
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.KeypairId = types.StringValue(keypairId)
	data.Id = types.StringValue(keypairId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateKeypairRequest{}
	createReq.SetKeypairName(data.KeypairName.ValueString())

	if !data.PublicKey.IsUnknown() && !data.PublicKey.IsNull() {
		createReq.SetPublicKey(data.PublicKey.ValueString())
	}

	var createResp oscgo.CreateKeypairResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.KeypairApi.CreateKeypair(ctx).CreateKeypairRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create keypair resource",
			"Error: "+err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(*createResp.ResponseContext.RequestId)
	keypair := createResp.GetKeypair()

	data.Id = types.StringValue(keypair.GetKeypairId())
	data.KeypairName = types.StringValue(keypair.GetKeypairName())
	if createReq.HasPublicKey() {
		data.PublicKey = types.StringValue(createReq.GetPublicKey())
	}

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, keypair.GetKeypairId())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data.PrivateKey = types.StringValue(keypair.GetPrivateKey())
	if createReq.HasPublicKey() {
		data.PublicKey = types.StringValue(createReq.GetPublicKey())
	}
	err = setKeypairState(ctx, r, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set keypair state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceKeypair) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeypairModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setKeypairState(ctx, r, &data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set keypair API response values.",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceKeypair) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData KeypairModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.KeypairId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	err := setKeypairState(ctx, r, &stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set keypair state after tags updating.",
			"Error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceKeypair) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeypairModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteKeypairRequest{
		KeypairId: data.KeypairId.ValueStringPointer(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.KeypairApi.DeleteKeypair(ctx).DeleteKeypairRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete keypair",
			"Error: "+err.Error(),
		)
		return
	}
}

func setKeypairState(ctx context.Context, r *resourceKeypair, data *KeypairModel) error {
	keypairFilters := oscgo.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'keypair' read timeout value: %v", diags.Errors())
	}

	readReq := oscgo.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	var readResp oscgo.ReadKeypairsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.KeypairApi.ReadKeypairs(ctx).ReadKeypairsRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if len(readResp.GetKeypairs()) == 0 {
		return ErrResourceEmpty
	}

	keypair := readResp.GetKeypairs()[0]
	tags, diag := flattenOAPITagsFW(ctx, keypair.GetTags())
	if diag.HasError() {
		return fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags
	data.KeypairFingerprint = types.StringValue(keypair.GetKeypairFingerprint())
	data.KeypairName = types.StringValue(keypair.GetKeypairName())
	data.KeypairType = types.StringValue(keypair.GetKeypairType())
	data.KeypairId = types.StringValue(keypair.GetKeypairId())
	return nil
}
