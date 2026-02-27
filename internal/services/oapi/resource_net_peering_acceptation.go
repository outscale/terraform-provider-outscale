package oapi

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
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithConfigure   = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithImportState = &resourceNetPeeringAcceptation{}
	_ resource.ResourceWithModifyPlan  = &resourceNetPeeringAcceptation{}
)

type NetPeeringAcceptationModel struct {
	AccepterNet        types.List     `tfsdk:"accepter_net"`
	ExpirationDate     types.String   `tfsdk:"expiration_date"`
	NetPeeringId       types.String   `tfsdk:"net_peering_id"`
	SourceNet          types.List     `tfsdk:"source_net"`
	State              types.List     `tfsdk:"state"`
	AccepterOwnerId    types.String   `tfsdk:"accepter_owner_id"`
	SourceNetAccountId types.String   `tfsdk:"source_net_account_id"`
	AccepterNetId      types.String   `tfsdk:"accepter_net_id"`
	SourceNetId        types.String   `tfsdk:"source_net_id"`
	RequestId          types.String   `tfsdk:"request_id"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	TagsComputedModel
}

type resourceNetPeeringAcceptation struct {
	Client *osc.Client
}

func NewResourceNetPeeringAcceptation() resource.Resource {
	return &resourceNetPeeringAcceptation{}
}

func (r *resourceNetPeeringAcceptation) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
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
	data.NetPeeringId = to.String(netPeeringId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts

	data.State = types.ListNull(types.ObjectType{AttrTypes: stateAttrTypes})
	data.SourceNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})
	data.AccepterNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})
	data.Tags = ComputedTagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
			"accepter_net": PeeringconnectionOptionsSchema(),
			"source_net":   PeeringconnectionOptionsSchema(),
			"state": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":    types.StringType,
						"message": types.StringType,
					},
				},
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (r *resourceNetPeeringAcceptation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetPeeringAcceptationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.AcceptNetPeeringRequest{
		NetPeeringId: data.NetPeeringId.ValueString(),
	}

	createResp, err := r.Client.AcceptNetPeering(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Net Peering accepter.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	netPeering := ptr.From(createResp.NetPeering)
	data.NetPeeringId = to.String(netPeering.NetPeeringId)
	data.Id = to.String(netPeering.NetPeeringId)

	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering Acceptation state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetPeeringAcceptation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetPeeringAcceptationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.read(ctx, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
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
}

func (r *resourceNetPeeringAcceptation) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *resourceNetPeeringAcceptation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *resourceNetPeeringAcceptation) read(ctx context.Context, data NetPeeringAcceptationModel) (NetPeeringAcceptationModel, error) {
	netPeeringFilters := osc.FiltersNetPeering{
		NetPeeringIds: &[]string{data.NetPeeringId.ValueString()},
	}
	readReq := osc.ReadNetPeeringsRequest{
		Filters: &netPeeringFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net peering' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadNetPeerings(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil || readResp.NetPeerings == nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	if len(*readResp.NetPeerings) == 0 {
		return data, ErrResourceEmpty
	}

	netPeering := (*readResp.NetPeerings)[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, netPeering.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert tags to the schema list: %v", diags.Errors())
	}
	sourceNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, SourceNetToList(netPeering.SourceNet))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert source net to the schema list: %v", diags.Errors())
	}
	accepterNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, AccepterNetToList(netPeering.AccepterNet))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert accepter net to the schema list: %v", diags.Errors())
	}
	state, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: stateAttrTypes}, NetPeerStateToList(netPeering.State))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert state to the schema list: %v", diags.Errors())
	}

	data.ExpirationDate = to.String(from.ISO8601(netPeering.ExpirationDate))
	data.Tags = tags
	data.AccepterNet = accepterNet
	data.NetPeeringId = to.String(netPeering.NetPeeringId)
	data.Id = to.String(netPeering.NetPeeringId)
	data.SourceNet = sourceNet
	data.State = state
	data.AccepterOwnerId = to.String(ptr.From(netPeering.AccepterNet.AccountId))
	data.SourceNetAccountId = to.String(ptr.From(netPeering.SourceNet.AccountId))
	data.AccepterNetId = to.String(ptr.From(netPeering.AccepterNet.NetId))
	data.SourceNetId = to.String(ptr.From(netPeering.SourceNet.NetId))

	return data, nil
}
