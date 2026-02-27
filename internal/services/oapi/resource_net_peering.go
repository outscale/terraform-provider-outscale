package oapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
)

var (
	_ resource.Resource                = &resourceNetPeering{}
	_ resource.ResourceWithConfigure   = &resourceNetPeering{}
	_ resource.ResourceWithImportState = &resourceNetPeering{}
	_ resource.ResourceWithModifyPlan  = &resourceNetPeering{}
)

type NetPeeringModel struct {
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
	TagsModel
}

type NetPeerModel struct {
	NetId     types.String `tfsdk:"net_id"`
	IpRange   types.String `tfsdk:"ip_range"`
	AccountId types.String `tfsdk:"account_id"`
}

type NetPeeringState struct {
	Message types.String `tfsdk:"message"`
	Name    types.String `tfsdk:"name"`
}

var netAttrTypes = fwhelpers.GetAttrTypes(NetPeerModel{})

var stateAttrTypes = fwhelpers.GetAttrTypes(NetPeeringState{})

func AccepterNetToList(n osc.AccepterNet) []NetPeerModel {
	return []NetPeerModel{
		{
			NetId:     to.String(n.NetId),
			IpRange:   to.String(n.IpRange),
			AccountId: to.String(n.AccountId),
		},
	}
}

func SourceNetToList(n osc.SourceNet) []NetPeerModel {
	return []NetPeerModel{
		{
			NetId:     to.String(n.NetId),
			IpRange:   to.String(n.IpRange),
			AccountId: to.String(n.AccountId),
		},
	}
}

func NetPeerStateToList(s osc.NetPeeringState) []NetPeeringState {
	return []NetPeeringState{
		{
			Message: to.String(s.Message),
			Name:    to.String(s.Name),
		},
	}
}

type resourceNetPeering struct {
	Client *osc.Client
}

func NewResourceNetPeering() resource.Resource {
	return &resourceNetPeering{}
}

func (r *resourceNetPeering) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceNetPeering) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	netPeeringId := req.ID

	if netPeeringId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Net Peering identifier, got: %v", req.ID),
		)
		return
	}

	var data NetPeeringModel
	var timeouts timeouts.Value
	data.NetPeeringId = to.String(netPeeringId)
	data.Id = to.String(netPeeringId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts

	data.State = types.ListNull(types.ObjectType{AttrTypes: stateAttrTypes})
	data.SourceNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})
	data.AccepterNet = types.ListNull(types.ObjectType{AttrTypes: netAttrTypes})
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceNetPeering) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net_peering"
}

func (r *resourceNetPeering) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceNetPeering) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"net_peering_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_net_account_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
			},
			"accepter_owner_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
			},
			"accepter_net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expiration_date": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accepter_net": PeeringconnectionOptionsSchema(),
			"source_net":   PeeringconnectionOptionsSchema(),
			"state": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: stateAttrTypes,
				},
			},
		},
	}
}

func PeeringconnectionOptionsSchema() *schema.ListAttribute {
	return &schema.ListAttribute{
		Computed: true,
		ElementType: types.ObjectType{
			AttrTypes: netAttrTypes,
		},
	}
}

func (r *resourceNetPeering) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetPeeringModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateNetPeeringRequest{
		AccepterNetId: data.AccepterNetId.ValueString(),
		SourceNetId:   data.SourceNetId.ValueString(),
	}

	if !data.AccepterOwnerId.IsUnknown() && !data.AccepterOwnerId.IsNull() {
		createReq.AccepterOwnerId = data.AccepterOwnerId.ValueStringPointer()
	} else {
		ids := []string{createReq.AccepterNetId, createReq.SourceNetId}
		filters := osc.FiltersNet{
			NetIds: &ids,
		}
		req := osc.ReadNetsRequest{
			Filters: &filters,
		}

		readResp, err := r.Client.ReadNets(ctx, req, options.WithRetryTimeout(createTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read associated nets resource.",
				err.Error(),
			)
			return
		}
		if readResp.Nets != nil && len(*readResp.Nets) != 2 {
			resp.Diagnostics.AddError(
				"Your 'accept_net_id' and 'source_net_id' are on different accounts, so the 'accept_owner_id' parameter is mandatory.",
				"Accepter owner id is mandatory.",
			)
			return
		}
	}

	createResp, err := r.Client.CreateNetPeering(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create net peering resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	netPeering := ptr.From(createResp.NetPeering)

	diag := createOAPITagsFW(ctx, r.Client, createTimeout, data.Tags, netPeering.NetPeeringId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{},
		Target:  []string{string(osc.NetPeeringStateNamePendingAcceptance), string(osc.NetPeeringStateNameActive)},
		Timeout: createTimeout,
		Refresh: ResourceNetPeeringconnectionStateRefreshFunc(ctx, createTimeout, r, netPeering.NetPeeringId),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error waiting for Net Peering (%s) to become available.",
				netPeering.NetPeeringId),
			err.Error(),
		)
		return
	}

	data.NetPeeringId = to.String(netPeering.NetPeeringId)
	data.Id = to.String(netPeering.NetPeeringId)
	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetPeering) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetPeeringModel

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

func (r *resourceNetPeering) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData NetPeeringModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.NetPeeringId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := r.read(ctx, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetPeering) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetPeeringModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteNetPeeringRequest{
		NetPeeringId: data.NetPeeringId.ValueString(),
	}

	_, err := r.Client.DeleteNetPeering(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Net Peering.",
			err.Error(),
		)
	}
}

func (r *resourceNetPeering) read(ctx context.Context, data NetPeeringModel) (NetPeeringModel, error) {
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
	if err != nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	if readResp.NetPeerings == nil || len(*readResp.NetPeerings) == 0 {
		return data, ErrResourceEmpty
	}

	netPeering := (*readResp.NetPeerings)[0]
	tags, diag := flattenOAPITagsFW(ctx, netPeering.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

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

func ResourceNetPeeringconnectionStateRefreshFunc(ctx context.Context, to time.Duration, r *resourceNetPeering, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := osc.ReadNetPeeringsRequest{Filters: &osc.FiltersNetPeering{NetPeeringIds: &[]string{id}}}

		resp, err := r.Client.ReadNetPeerings(ctx, readReq, options.WithRetryTimeout(to))
		if err != nil || resp.NetPeerings == nil || len(*resp.NetPeerings) == 0 {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringconnectionID.NotFound") {
				return nil, "", nil
			}
			return resp, "error", err
		}
		netPeering := (*resp.NetPeerings)[0]

		if netPeering.State.Name == "failed" {
			return nil, "failed", errors.New(netPeering.State.Message)
		}

		return resp, string(netPeering.State.Name), nil
	}
}
