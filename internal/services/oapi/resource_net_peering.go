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
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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

func AccepterNetToList(n oscgo.AccepterNet) []NetPeerModel {
	return []NetPeerModel{
		{
			NetId:     types.StringValue(n.GetNetId()),
			IpRange:   types.StringValue(n.GetIpRange()),
			AccountId: types.StringValue(n.GetAccountId()),
		},
	}
}

func SourceNetToList(n oscgo.SourceNet) []NetPeerModel {
	return []NetPeerModel{
		{
			NetId:     types.StringValue(n.GetNetId()),
			IpRange:   types.StringValue(n.GetIpRange()),
			AccountId: types.StringValue(n.GetAccountId()),
		},
	}
}

func NetPeerStateToList(s oscgo.NetPeeringState) []NetPeeringState {
	return []NetPeeringState{
		{
			Message: types.StringValue(s.GetMessage()),
			Name:    types.StringValue(s.GetName()),
		},
	}
}

type resourceNetPeering struct {
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *oscgo.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.NetPeeringId = types.StringValue(netPeeringId)
	data.Id = types.StringValue(netPeeringId)
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
	if resp.Diagnostics.HasError() {
		return
	}
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
		return
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
					fwmodifyplan.ForceNewFramework(),
				},
			},
			"accepter_owner_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					fwmodifyplan.ForceNewFramework(),
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
			"accepter_net": PeeringConnectionOptionsSchema(),
			"source_net":   PeeringConnectionOptionsSchema(),
			"state": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: stateAttrTypes,
				},
			},
		},
	}
}

func PeeringConnectionOptionsSchema() *schema.ListAttribute {
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateNetPeeringRequest{
		AccepterNetId: data.AccepterNetId.ValueString(),
		SourceNetId:   data.SourceNetId.ValueString(),
	}

	if !data.AccepterOwnerId.IsUnknown() && !data.AccepterOwnerId.IsNull() {
		createReq.SetAccepterOwnerId(data.AccepterOwnerId.ValueString())
	} else {
		ids := []string{createReq.AccepterNetId, createReq.SourceNetId}
		filters := oscgo.FiltersNet{
			NetIds: &ids,
		}
		req := oscgo.ReadNetsRequest{
			Filters: &filters,
		}

		var readResp oscgo.ReadNetsResponse
		err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.NetApi.ReadNets(ctx).ReadNetsRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			readResp = rp
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read associated nets resource.",
				err.Error(),
			)
			return
		}
		if len(readResp.GetNets()) != 2 {
			resp.Diagnostics.AddError(
				"Your 'accept_net_id' and 'source_net_id' are on different accounts, so the 'accept_owner_id' parameter is mandatory.",
				"Accepter owner id is mandatory.",
			)
			return
		}
	}

	var createResp oscgo.CreateNetPeeringResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetPeeringApi.CreateNetPeering(ctx).CreateNetPeeringRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create net peering resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	netPeering := createResp.GetNetPeering()

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, netPeering.GetNetPeeringId())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"initiating-request", "provisioning", "pending"},
		Target:     []string{"pending-acceptance", "active"},
		Refresh:    ResourceNetPeeringConnectionStateRefreshFunc(ctx, createTimeout, r, netPeering.GetNetPeeringId()),
		Timeout:    createTimeout,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error waiting for Net Peering (%s) to become available.",
				netPeering.GetNetPeeringId()),
			err.Error(),
		)
		return
	}

	data.NetPeeringId = types.StringValue(netPeering.GetNetPeeringId())
	data.Id = types.StringValue(netPeering.GetNetPeeringId())
	data, err = setNetPeeringState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeering) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetPeeringModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setNetPeeringState(ctx, r, data)
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeering) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData NetPeeringModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.NetPeeringId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := setNetPeeringState(ctx, r, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Peering state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetPeering) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetPeeringModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteNetPeeringRequest{
		NetPeeringId: data.NetPeeringId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.NetPeeringApi.DeleteNetPeering(ctx).DeleteNetPeeringRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Net Peering.",
			err.Error(),
		)
		return
	}
}

func setNetPeeringState(ctx context.Context, r *resourceNetPeering, data NetPeeringModel) (NetPeeringModel, error) {
	netPeeringFilters := oscgo.FiltersNetPeering{
		NetPeeringIds: &[]string{data.NetPeeringId.ValueString()},
	}
	readReq := oscgo.ReadNetPeeringsRequest{
		Filters: &netPeeringFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net peering' read timeout value: %v", diags.Errors())
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
		return data, ErrResourceEmpty
	}

	netPeering := readResp.GetNetPeerings()[0]
	tags, diag := flattenOAPITagsFW(ctx, netPeering.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	sourceNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, SourceNetToList(netPeering.GetSourceNet()))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert source net to the schema list: %v", diags.Errors())
	}
	accepterNet, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: netAttrTypes}, AccepterNetToList(netPeering.GetAccepterNet()))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert accepter net to the schema list: %v", diags.Errors())
	}
	state, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: stateAttrTypes}, NetPeerStateToList(netPeering.GetState()))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert state to the schema list: %v", diags.Errors())
	}

	data.ExpirationDate = types.StringValue(netPeering.GetExpirationDate())
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

func ResourceNetPeeringConnectionStateRefreshFunc(ctx context.Context, to time.Duration, r *resourceNetPeering, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := oscgo.ReadNetPeeringsRequest{Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{id}}}
		var resp oscgo.ReadNetPeeringsResponse

		err := retry.RetryContext(ctx, to, func() *retry.RetryError {
			rp, httpResp, err := r.Client.NetPeeringApi.ReadNetPeerings(ctx).ReadNetPeeringsRequest(readReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
				return nil, "", nil
			}
			return resp, "error", err
		}
		netPeering := resp.GetNetPeerings()[0]

		if netPeering.State.GetName() == "failed" {
			return nil, "failed", errors.New(*netPeering.GetState().Message)
		}

		return resp, *netPeering.GetState().Name, nil
	}
}
