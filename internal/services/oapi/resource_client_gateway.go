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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                = &clientGatewayResource{}
	_ resource.ResourceWithConfigure   = &clientGatewayResource{}
	_ resource.ResourceWithImportState = &clientGatewayResource{}
)

const (
	clientGwErrCreate = "Unable to create Client Gateway"
	clientGwErrRead   = "Unable to read Client Gateway"
	clientGwErrDelete = "Unable to delete Client Gateway"
	clientGwErrState  = "Unable to set Client Gateway state"
)

type clientGatewayModel struct {
	BgpAsn          types.Int64    `tfsdk:"bgp_asn"`
	ConnectionType  types.String   `tfsdk:"connection_type"`
	PublicIp        types.String   `tfsdk:"public_ip"`
	ClientGatewayId types.String   `tfsdk:"client_gateway_id"`
	State           types.String   `tfsdk:"state"`
	Id              types.String   `tfsdk:"id"`
	RequestId       types.String   `tfsdk:"request_id"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type clientGatewayResource struct {
	Client *osc.Client
}

func NewResourceClientGateway() resource.Resource {
	return &clientGatewayResource{}
}

func (r *clientGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *clientGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_gateway"
}

func (r *clientGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import client gateway identifier. Got: %v", req.ID),
		)
		return
	}

	var data clientGatewayModel
	var timeouts timeouts.Value
	data.Id = to.String(id)
	data.ClientGatewayId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clientGatewayResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bgp_asn": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"connection_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_ip": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_gateway_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *clientGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data clientGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateClientGatewayRequest{
		BgpAsn:         int(data.BgpAsn.ValueInt64()),
		ConnectionType: data.ConnectionType.ValueString(),
		PublicIp:       data.PublicIp.ValueString(),
	}

	createResp, err := r.Client.CreateClientGateway(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(clientGwErrCreate, err.Error())
		return
	}

	data.Id = to.String(createResp.ClientGateway.ClientGatewayId)
	data.ClientGatewayId = to.String(createResp.ClientGateway.ClientGatewayId)

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, createResp.ClientGateway.ClientGatewayId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(ctx, timeout, createResp.ClientGateway.ClientGatewayId),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(clientGwErrState, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(clientGwErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *clientGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data clientGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
		resp.Diagnostics.AddError(clientGwErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *clientGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData clientGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.Id.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(clientGwErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *clientGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data clientGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.DeleteClientGatewayRequest{
		ClientGatewayId: data.Id.ValueString(),
	}

	_, err := r.Client.DeleteClientGateway(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(clientGwErrDelete, err.Error())
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "failed"},
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(ctx, timeout, data.Id.ValueString()),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		switch {
		case errors.Is(err, ErrResourceEmpty):
		default:
			resp.Diagnostics.AddError(clientGwErrDelete, err.Error())
		}
	}
}

func (r *clientGatewayResource) read(ctx context.Context, timeout time.Duration, data clientGatewayModel) (clientGatewayModel, error) {
	req := osc.ReadClientGatewaysRequest{
		Filters: &osc.FiltersClientGateway{
			ClientGatewayIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadClientGateways(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	if resp.ClientGateways == nil || len(*resp.ClientGateways) == 0 || (*resp.ClientGateways)[0].State == "deleted" {
		return data, ErrResourceEmpty
	}

	gw := (*resp.ClientGateways)[0]

	tags, diag := flattenOAPITagsFW(ctx, gw.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("%v", diag.Errors())
	}

	data.Tags = tags
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Id = to.String(gw.ClientGatewayId)
	data.ClientGatewayId = to.String(gw.ClientGatewayId)
	data.BgpAsn = to.Int64(gw.BgpAsn)
	data.ConnectionType = to.String(gw.ConnectionType)
	data.PublicIp = to.String(gw.PublicIp)
	data.State = to.String(gw.State)

	return data, nil
}

func (r *clientGatewayResource) stateRefreshFunc(ctx context.Context, timeout time.Duration, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		req := osc.ReadClientGatewaysRequest{
			Filters: &osc.FiltersClientGateway{
				ClientGatewayIds: &[]string{id},
			},
		}
		resp, err := r.Client.ReadClientGateways(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "", err
		}

		if len(ptr.From(resp.ClientGateways)) == 0 {
			return nil, "", ErrResourceEmpty
		}

		gw := (*resp.ClientGateways)[0]

		return resp, gw.State, nil
	}
}
