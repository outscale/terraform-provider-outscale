package outscale

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var (
	_ resource.Resource                = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithConfigure   = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithImportState = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithModifyPlan  = &resourceNetPeeringAcceptation{}
)

type NetPeeringAcceptationModel struct {
	AccepterNet        types.List   `tfsdk:"accepter_net"`
	ExpirationDate     types.String `tfsdk:"expiration_date"`
	NetPeeringId       types.String `tfsdk:"net_peering_id"`
	SourceNet          types.List   `tfsdk:"source_net"`
	State              types.List   `tfsdk:"state"`
	AccepterOwnerId    types.String `tfsdk:"accepter_owner_id"`
	SourceNetAccountId types.String `tfsdk:"source_net_account_id"`
	AccepterNetId      types.String `tfsdk:"accepter_net_id"`
	SourceNetId        types.String `tfsdk:"source_net_id"`
	Tags               types.List   `tfsdk:"tags"`

	RequestId types.String   `tfsdk:"request_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
	Id        types.String   `tfsdk:"id"`
}

type resourceNetPeeringAcceptation struct {
	Client *oscgo.APIClient
}

func NewResourceNetPeeringAcceptation() resource.Resource {
	return &resourceNetPeeringAcceptation{}
}

func (r *resourceNetPeeringAcceptation) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
}

func (r *resourceNetPeeringAcceptation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	netPeeringId := req.ID

	if netPeeringId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Net Peering identifier, got: %v", req.ID),
		)
		return
	}

	var data NetPeeringAcceptationModel
	var timeouts timeouts.Value
	data.NetPeeringId = types.StringValue(netPeeringId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts

	data.Tags = types.ListNull(types.ObjectType{AttrTypes: tagAttrTypes})
	data.State = types.ListNull(types.ObjectType{AttrTypes: stateAttrTypes})
	data.SourceNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})
	data.AccepterNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeeringAcceptation) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net_peering_acceptation"
}

func (r *resourceNetPeeringAcceptation) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will remove resource from terraform state but not fully delete. It will be gone when deleting the Net Peering.",
		)
	}
}

func (r *resourceNetPeeringAcceptation) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"tags": TagsSchemaComputedAttr(),
			"net_peering_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_net_account_id": schema.StringAttribute{
				Computed: true,
			},
			"accepter_owner_id": schema.StringAttribute{
				Computed: true,
			},
			"accepter_net_id": schema.StringAttribute{
				Computed: true,
			},
			"source_net_id": schema.StringAttribute{
				Computed: true,
			},
			"expiration_date": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"accepter_net": PeeringConnectionOptionsSchema(),
			"source_net":   PeeringConnectionOptionsSchema(),
			"state": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":    types.StringType,
						"message": types.StringType,
					},
				},
			},
		},
	}
}

func TagsSchemaComputedAttr() *schema.ListAttribute {
	return &schema.ListAttribute{
		Computed: true,
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		},
	}
}

var tagAttrTypes = utils.GetAttrTypes(ResourceTag{})

func (r *resourceNetPeeringAcceptation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetPeeringAcceptationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.AcceptNetPeeringRequest{
		NetPeeringId: data.NetPeeringId.ValueString(),
	}

	var createResp oscgo.AcceptNetPeeringResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetPeeringApi.AcceptNetPeering(ctx).AcceptNetPeeringRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Net Peering accepter.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	netPeering := createResp.GetNetPeering()
	data.NetPeeringId = types.StringValue(netPeering.GetNetPeeringId())
	data.Id = types.StringValue(netPeering.GetNetPeeringId())

	data, err = r.setNetPeeringAcceptationState(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering Acceptation state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeeringAcceptation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetPeeringAcceptationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.setNetPeeringAcceptationState(ctx, data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set subnet API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeeringAcceptation) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *resourceNetPeeringAcceptation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *resourceNetPeeringAcceptation) setNetPeeringAcceptationState(ctx context.Context, data NetPeeringAcceptationModel) (NetPeeringAcceptationModel, error) {
	netPeeringFilters := oscgo.FiltersNetPeering{
		NetPeeringIds: &[]string{data.NetPeeringId.ValueString()},
	}
	readReq := oscgo.ReadNetPeeringsRequest{
		Filters: &netPeeringFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net peering' read timeout value. Error: %v: ", diags.Errors())
	}

	var readResp oscgo.ReadNetPeeringsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetPeeringApi.ReadNetPeerings(ctx).ReadNetPeeringsRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}
	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	if len(readResp.GetNetPeerings()) == 0 {
		return data, errors.New("Empty")
	}

	netPeering := readResp.GetNetPeerings()[0]

	tags, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: tagAttrTypes}, getTagsFromApiResponse(netPeering.GetTags()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Tags to the schema List. Error: %v: ", diags.Errors())
	}
	sourceNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, SourceNetToList(netPeering.GetSourceNet()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Source Net to the schema List. Error: %v: ", diags.Errors())
	}
	accepterNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, AccepterNetToList(netPeering.GetAccepterNet()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Accepter Net to the schema List. Error: %v: ", diags.Errors())
	}
	state, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: stateAttrTypes}, NetPeerStateToList(netPeering.GetState()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert State to the schema List. Error: %v: ", diags.Errors())
	}

	data.ExpirationDate = types.StringValue(netPeering.GetExpirationDate())
	data.Tags = tags
	data.AccepterNet = accepterNet
	data.NetPeeringId = types.StringValue(netPeering.GetNetPeeringId())
	data.Id = types.StringValue(netPeering.GetNetPeeringId())
	data.SourceNet = sourceNet
	data.State = state
	data.AccepterOwnerId = types.StringValue(netPeering.AccepterNet.GetAccountId())
	data.SourceNetAccountId = types.StringValue(netPeering.SourceNet.GetAccountId())
	data.AccepterNetId = types.StringValue(netPeering.AccepterNet.GetNetId())
	data.SourceNetId = types.StringValue(netPeering.SourceNet.GetNetId())

	return data, nil
}
