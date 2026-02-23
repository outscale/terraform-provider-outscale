package oapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var _ ephemeral.EphemeralResource = &resourceEphemeralKeypair{}

type resourceEphemeralKeypair struct {
	Client *osc.Client
}
type EphemeralKeypairModel struct {
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

func NewKeypairEphemeralResource() ephemeral.EphemeralResource {
	return &resourceEphemeralKeypair{}
}

func (r *resourceEphemeralKeypair) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

func (d *resourceEphemeralKeypair) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypair"
}

func (d resourceEphemeralKeypair) ValidateConfig(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	var data EphemeralKeypairModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If attribute_one is not configured, return without warning.
	if data.KeypairName.IsNull() || data.KeypairName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("keypair_name"),
			"Missing Attribute Configuration",
			"The 'keypair_name' parameter is mandatory.",
		)
	}
}

func (r *resourceEphemeralKeypair) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
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
			},
			"keypair_type": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Optional: true,
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

func (e *resourceEphemeralKeypair) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data EphemeralKeypairModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	resp.RenewAt = time.Now().Add(createTimeout)

	createReq := osc.CreateKeypairRequest{}
	createReq.KeypairName = data.KeypairName.ValueString()

	if !data.PublicKey.IsUnknown() && !data.PublicKey.IsNull() {
		createReq.PublicKey = data.PublicKey.ValueStringPointer()
	}

	isExit, err := isResourceExist(ctx, &data, e)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to check whether the keypair resource already exists",
			err.Error(),
		)
		return
	}
	if !isExit {
		createResp, err := e.Client.CreateKeypair(ctx, createReq, options.WithRetryTimeout(createTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create ephemeral keypair resource",
				err.Error(),
			)
			return
		}

		data.RequestId = to.String(*createResp.ResponseContext.RequestId)
		keypair := ptr.From(createResp.Keypair)
		data.Id = to.String(keypair.KeypairId)
		data.KeypairName = to.String(keypair.KeypairName)
		if createReq.PublicKey != nil {
			data.PublicKey = to.String(createReq.PublicKey)
		}
		data.PrivateKey = to.String(keypair.PrivateKey)

		diag := createOAPITagsFW(ctx, e.Client, createTimeout, data.Tags, *keypair.KeypairId)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	err = setEphKeypairState(ctx, e, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set ephemeral keypair state",
			err.Error(),
		)
		return
	}
	privateData, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to convert EphemeralKeypairModel to private state",
			err.Error(),
		)
		return
	}
	resp.Private.SetKey(ctx, "ephemKeypairData", privateData)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func setEphKeypairState(ctx context.Context, r *resourceEphemeralKeypair, data *EphemeralKeypairModel) error {
	keypairFilters := osc.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'keypair' read timeout value: %v", diags.Errors())
	}

	readReq := osc.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	readResp, err := r.Client.ReadKeypairs(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return err
	}
	if readResp.Keypairs == nil || len(*readResp.Keypairs) == 0 {
		return ErrResourceEmpty
	}

	keypair := (*readResp.Keypairs)[0]

	tags, diag := flattenOAPITagsFW(ctx, ptr.From(keypair.Tags))
	if diag.HasError() {
		return fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags
	data.KeypairFingerprint = to.String(keypair.KeypairFingerprint)
	data.KeypairName = to.String(keypair.KeypairName)
	data.KeypairType = to.String(keypair.KeypairType)
	data.KeypairId = to.String(keypair.KeypairId)

	return nil
}

func isResourceExist(ctx context.Context, data *EphemeralKeypairModel, e *resourceEphemeralKeypair) (bool, error) {
	keypairFilters := osc.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}
	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return false, fmt.Errorf("unable to parse 'ephemeral keypair' read timeout value: %v", diags.Errors())
	}

	readReq := osc.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	rp, err := e.Client.ReadKeypairs(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return false, err
	}
	if rp.Keypairs == nil {
		return false, fmt.Errorf("keypair not found")
	}

	return len(*rp.Keypairs) > 0, err
}
