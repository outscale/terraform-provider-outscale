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
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &vpnConnectionRouteResource{}
	_ resource.ResourceWithConfigure   = &vpnConnectionRouteResource{}
	_ resource.ResourceWithImportState = &vpnConnectionRouteResource{}
)

const (
	vpnConnectionRouteErrCreate = "Unable to create VPN Connection Route"
	vpnConnectionRouteErrRead   = "Unable to read VPN Connection Route"
	vpnConnectionRouteErrDelete = "Unable to delete VPN Connection Route"
	vpnConnectionRouteErrWait   = "Unable to wait for VPN Connection Route state"
)

type vpnConnectionRouteModel struct {
	DestinationIpRange types.String   `tfsdk:"destination_ip_range"`
	VpnConnectionId    types.String   `tfsdk:"vpn_connection_id"`
	RequestId          types.String   `tfsdk:"request_id"`
	Id                 types.String   `tfsdk:"id"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

type vpnConnectionRouteResource struct {
	Client *osc.Client
}

func NewResourceVPNConnectionRoute() resource.Resource {
	return &vpnConnectionRouteResource{}
}

func (r *vpnConnectionRouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpnConnectionRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_connection_route"
}

func (r *vpnConnectionRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "_", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"To import a VPN Connection Route, use the format {vpn_connection_id}_{destination_ip_range}. Got: "+req.ID,
		)
		return
	}

	vpnConnectionID := parts[0]
	destinationIPRange := parts[1]

	var data vpnConnectionRouteModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(fmt.Sprintf("%s:%s", destinationIPRange, vpnConnectionID))
	data.VpnConnectionId = to.String(vpnConnectionID)
	data.DestinationIpRange = to.String(destinationIPRange)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpnConnectionRouteResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"destination_ip_range": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpn_connection_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *vpnConnectionRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data vpnConnectionRouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	ipRange := data.DestinationIpRange.ValueString()
	vpnId := data.VpnConnectionId.ValueString()

	createReq := osc.CreateVpnConnectionRouteRequest{
		DestinationIpRange: data.DestinationIpRange.ValueString(),
		VpnConnectionId:    data.VpnConnectionId.ValueString(),
	}
	_, err := r.Client.CreateVpnConnectionRoute(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionRouteErrCreate, err.Error())
		return
	}

	data.Id = to.String(fmt.Sprintf("%s:%s", ipRange, vpnId))

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionRouteErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *vpnConnectionRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vpnConnectionRouteModel
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
		resp.Diagnostics.AddError(vpnConnectionRouteErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *vpnConnectionRouteResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *vpnConnectionRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vpnConnectionRouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	ipRange, connectionID := r.parseConnection(data.Id.ValueString())

	deleteReq := osc.DeleteVpnConnectionRouteRequest{
		DestinationIpRange: ipRange,
		VpnConnectionId:    connectionID,
	}
	_, err := r.Client.DeleteVpnConnectionRoute(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionRouteErrDelete, err.Error())
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.RouteLightState]{
		Pending: stateconf.States(osc.RouteLightStateDeleting),
		Target:  stateconf.States(osc.RouteLightStateDeleted),
		Timeout: timeout,
		Refresh: r.refreshFunc(ipRange, connectionID),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionRouteErrWait, err.Error())
	}
}

func (r *vpnConnectionRouteResource) read(ctx context.Context, timeout time.Duration, data vpnConnectionRouteModel) (vpnConnectionRouteModel, error) {
	ipRange, connectionID := r.parseConnection(data.Id.ValueString())

	stateConf := &stateconf.StateChangeConf[osc.RouteLightState]{
		Pending: stateconf.States(osc.RouteLightStatePending),
		Target:  stateconf.States(osc.RouteLightStateAvailable),
		Timeout: timeout,
		Refresh: r.refreshFunc(ipRange, connectionID),
	}

	respAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return data, err
	}
	resp := respAny.(*osc.ReadVpnConnectionsResponse)

	data.VpnConnectionId = to.String((*resp.VpnConnections)[0].VpnConnectionId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return data, nil
}

func (r *vpnConnectionRouteResource) refreshFunc(destinationIPRange, id string) stateconf.StateRefreshFunc[osc.RouteLightState] {
	return func(ctx context.Context) (any, osc.RouteLightState, error) {
		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				VpnConnectionIds: &[]string{id},
			},
		}

		resp, err := r.Client.ReadVpnConnections(ctx, filter)
		if err != nil {
			return nil, "", err
		}
		if resp.VpnConnections == nil || len(*resp.VpnConnections) == 0 {
			return nil, "", ErrResourceEmpty
		}
		vpn := (*resp.VpnConnections)[0]

		route, found := lo.Find(vpn.Routes, func(route osc.RouteLight) bool {
			return route.DestinationIpRange == destinationIPRange
		})
		if !found {
			return nil, "", ErrResourceEmpty
		}

		return resp, route.State, nil
	}
}

func (r *vpnConnectionRouteResource) parseConnection(id string) (ipRange string, connectionID string) {
	parts := strings.SplitN(id, ":", 2)
	return parts[0], parts[1]
}
