package outscale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var _ ephemeral.EphemeralResource = &resourceEphemeralKeypair{}

type resourceEphemeralKeypair struct {
	Client *oscgo.APIClient
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
	Tags               []ResourceTag  `tfsdk:"tags"`
	Id                 types.String   `tfsdk:"id"`
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
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
			"tags": TagsSchema(),
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
	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.RenewAt = time.Now().Add(createTimeout)

	createReq := oscgo.CreateKeypairRequest{}
	createReq.SetKeypairName(data.KeypairName.ValueString())

	if !data.PublicKey.IsUnknown() && !data.PublicKey.IsNull() {
		createReq.SetPublicKey(data.PublicKey.ValueString())
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
		var createResp oscgo.CreateKeypairResponse
		err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			rp, httpResp, err := e.Client.KeypairApi.CreateKeypair(ctx).CreateKeypairRequest(createReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			createResp = rp
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create ephemeral keypair resource",
				err.Error(),
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
		if len(data.Tags) > 0 {
			err = createFrameworkTags(ctx, e.Client, tagsToOSCResourceTag(data.Tags), keypair.GetKeypairId())
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to add Tags on ephemeral keypair resource",
					err.Error(),
				)
				return
			}
		}
		data.PrivateKey = types.StringValue(keypair.GetPrivateKey())
		if createReq.HasPublicKey() {
			data.PublicKey = types.StringValue(createReq.GetPublicKey())
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func setEphKeypairState(ctx context.Context, r *resourceEphemeralKeypair, data *EphemeralKeypairModel) error {
	keypairFilters := oscgo.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'keypair' read timeout value. Error: %v: ", diags.Errors())
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
		return utils.GetErrorResponse(err)
	}
	if len(readResp.GetKeypairs()) == 0 {
		return errors.New("Empty")
	}

	keypair := readResp.GetKeypairs()[0]
	data.Tags = getTagsFromApiResponse(keypair.GetTags())
	data.KeypairFingerprint = types.StringValue(keypair.GetKeypairFingerprint())
	data.KeypairName = types.StringValue(keypair.GetKeypairName())
	data.KeypairType = types.StringValue(keypair.GetKeypairType())
	data.KeypairId = types.StringValue(keypair.GetKeypairId())
	return nil
}

func isResourceExist(ctx context.Context, data *EphemeralKeypairModel, e *resourceEphemeralKeypair) (bool, error) {
	keypairFilters := oscgo.FiltersKeypair{
		KeypairNames: &[]string{data.KeypairName.ValueString()},
	}
	isKeypairExist := false
	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return isKeypairExist, fmt.Errorf("unable to parse 'ephemeral keypair' read timeout value. Error: %v: ", diags.Errors())
	}

	readReq := oscgo.ReadKeypairsRequest{
		Filters: &keypairFilters,
	}
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := e.Client.KeypairApi.ReadKeypairs(ctx).ReadKeypairsRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		if len(rp.GetKeypairs()) > 0 {
			isKeypairExist = true
		}
		return nil
	})
	return isKeypairExist, err
}
