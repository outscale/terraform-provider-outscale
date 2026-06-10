package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var _ ephemeral.EphemeralResource = &resourceEphemeralKeypair{}

const (
	ephemeralKeypairErrCreate     = "Unable to create Ephemeral Keypair"
	ephemeralKeypairErrExistCheck = "Unable to verify if the Keypair already exists"
)

type resourceEphemeralKeypair struct {
	keypairCommon
}

type EphemeralKeypairModel struct {
	KeypairModel
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

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	resp.RenewAt = time.Now().Add(timeout)

	createReq := osc.CreateKeypairRequest{}
	createReq.KeypairName = data.KeypairName.ValueString()

	if !data.PublicKey.IsUnknown() && !data.PublicKey.IsNull() {
		createReq.PublicKey = data.PublicKey.ValueStringPointer()
	}

	isExit, err := e.isResourceExist(ctx, timeout, &data)
	if err != nil {
		resp.Diagnostics.AddError(ephemeralKeypairErrExistCheck, err.Error())
		return
	}
	if !isExit {
		createResp, err := e.Client.CreateKeypair(ctx, createReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(ephemeralKeypairErrCreate, err.Error())
			return
		}

		data.RequestId = to.String(*createResp.ResponseContext.RequestId)
		keypair := ptr.From(createResp.Keypair)
		data.Id = to.String(keypair.KeypairId)
		data.KeypairName = to.String(keypair.KeypairName)
		if createReq.PublicKey != nil {
			data.PublicKey = to.String(createReq.PublicKey)
		}

		stateData := e.flattenCreate(data.KeypairModel, keypair)
		diag := resp.Result.Set(ctx, &stateData)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		diag = createOAPITagsFW(ctx, e.Client, timeout, data.Tags, *keypair.KeypairId)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	stateData, err := e.read(ctx, timeout, data.KeypairModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set ephemeral keypair state",
			err.Error(),
		)
		return
	}
	// Private state is not needed here: ephemeral resources are not persisted to state or plan files,
	// and private state only flows in-memory between Open → Renew → Close lifecycle methods.
	// Since neither Renew nor Close are implemented, private state would be immediately discarded.
	// See: https://developer.hashicorp.com/terraform/plugin/framework/ephemeral-resources/renew
	//
	// privateData, err := json.Marshal(data)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to convert EphemeralKeypairModel to private state",
	// 		err.Error(),
	// 	)
	// 	return
	// }
	// resp.Private.SetKey(ctx, "ephemKeypairData", privateData)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	resp.Diagnostics.Append(resp.Result.Set(ctx, &stateData)...)
}

func (r *resourceEphemeralKeypair) isResourceExist(ctx context.Context, timeout time.Duration, data *EphemeralKeypairModel) (bool, error) {
	keypairFilters := osc.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}

	readReq := osc.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	rp, err := r.Client.ReadKeypairs(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return false, err
	}

	return len(ptr.From(rp.Keypairs)) > 0, err
}
