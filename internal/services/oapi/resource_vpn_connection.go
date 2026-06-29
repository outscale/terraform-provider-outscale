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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &vpnConnectionResource{}
	_ resource.ResourceWithConfigure   = &vpnConnectionResource{}
	_ resource.ResourceWithImportState = &vpnConnectionResource{}
)

const (
	vpnConnectionErrCreate = "Unable to create VPN Connection"
	vpnConnectionErrDelete = "Unable to delete VPN Connection"
	vpnConnectionErrWait   = "Unable to wait for VPN Connection state"

	vpnConnectionCreateTimeout = 15 * time.Minute
)

type vpnConnectionModel struct {
	ClientGatewayId            types.String   `tfsdk:"client_gateway_id"`
	VirtualGatewayId           types.String   `tfsdk:"virtual_gateway_id"`
	ConnectionType             types.String   `tfsdk:"connection_type"`
	StaticRoutesOnly           types.Bool     `tfsdk:"static_routes_only"`
	ClientGatewayConfiguration types.String   `tfsdk:"client_gateway_configuration"`
	VpnConnectionId            types.String   `tfsdk:"vpn_connection_id"`
	State                      types.String   `tfsdk:"state"`
	Routes                     types.List     `tfsdk:"routes"`
	VgwTelemetries             types.List     `tfsdk:"vgw_telemetries"`
	RequestId                  types.String   `tfsdk:"request_id"`
	Id                         types.String   `tfsdk:"id"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type vpnConnectionRouteLightModel struct {
	DestinationIpRange types.String `tfsdk:"destination_ip_range"`
	RouteType          types.String `tfsdk:"route_type"`
	State              types.String `tfsdk:"state"`
}

type vpnConnectionVgwTelemetryModel struct {
	AcceptedRouteCount  types.Int64  `tfsdk:"accepted_route_count"`
	LastStateChangeDate types.String `tfsdk:"last_state_change_date"`
	OutsideIpAddress    types.String `tfsdk:"outside_ip_address"`
	State               types.String `tfsdk:"state"`
	StateDescription    types.String `tfsdk:"state_description"`
}

var (
	vpnConnectionRouteLightAttrTypes   = fwhelpers.GetAttrTypes(vpnConnectionRouteLightModel{})
	vpnConnectionVgwTelemetryAttrTypes = fwhelpers.GetAttrTypes(vpnConnectionVgwTelemetryModel{})
)

type vpnConnectionResource struct {
	Client *osc.Client
}

func NewResourceVPNConnection() resource.Resource {
	return &vpnConnectionResource{}
}

func (r *vpnConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpnConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_connection"
}

func (r *vpnConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import VPN connection identifier. Got: %v", req.ID),
		)
		return
	}

	var data vpnConnectionModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.VpnConnectionId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal
	data.Tags = TagsNull()
	data.Routes = types.ListNull(to.ObjType(vpnConnectionRouteLightAttrTypes))
	data.VgwTelemetries = types.ListNull(to.ObjType(vpnConnectionVgwTelemetryAttrTypes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpnConnectionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"client_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtual_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"static_routes_only": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"client_gateway_configuration": schema.StringAttribute{
				Computed: true,
			},
			"vpn_connection_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"routes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination_ip_range": schema.StringAttribute{
							Computed: true,
						},
						"route_type": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"vgw_telemetries": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accepted_route_count": schema.Int64Attribute{
							Computed: true,
						},
						"last_state_change_date": schema.StringAttribute{
							Computed: true,
						},
						"outside_ip_address": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"state_description": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *vpnConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data vpnConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, vpnConnectionCreateTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateVpnConnectionRequest{
		ClientGatewayId:  data.ClientGatewayId.ValueString(),
		VirtualGatewayId: data.VirtualGatewayId.ValueString(),
		ConnectionType:   data.ConnectionType.ValueString(),
	}

	if fwhelpers.IsSet(data.StaticRoutesOnly) {
		createReq.StaticRoutesOnly = data.StaticRoutesOnly.ValueBoolPointer()
	}

	createResp, err := r.Client.CreateVpnConnection(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		oscErr := oapihelpers.GetError(err)
		if oscErr.Code == "6008" {
			resp.Diagnostics.AddError(vpnConnectionErrCreate, err.Error())
			return
		}
		resp.Diagnostics.AddError(vpnConnectionErrCreate, err.Error())
		return
	}
	vpn := createResp.VpnConnection
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(vpn.VpnConnectionId)
	data.VpnConnectionId = to.String(vpn.VpnConnectionId)

	stateData, err := r.flatten(ctx, data, vpn)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, vpn.VpnConnectionId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err = r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *vpnConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vpnConnectionModel
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
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *vpnConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData vpnConnectionModel
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
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *vpnConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vpnConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	id := data.Id.ValueString()

	_, err := r.Client.DeleteVpnConnection(ctx, osc.DeleteVpnConnectionRequest{
		VpnConnectionId: id,
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionErrDelete, err.Error())
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.VpnConnectionState]{
		Pending: stateconf.States(osc.VpnConnectionStateDeleting),
		Target:  stateconf.States(osc.VpnConnectionStateDeleted),
		Timeout: timeout,
		Refresh: r.refreshFunc(id),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(vpnConnectionErrWait, err.Error())
	}
}

func (r *vpnConnectionResource) read(ctx context.Context, timeout time.Duration, data vpnConnectionModel) (vpnConnectionModel, error) {
	conf := &stateconf.StateChangeConf[osc.VpnConnectionState]{
		Pending: stateconf.States(osc.VpnConnectionStatePending),
		Target:  stateconf.States(osc.VpnConnectionStateAvailable),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.Id.ValueString()),
	}
	respAny, err := conf.WaitForStateContext(ctx)
	if err != nil {
		return data, err
	}

	resp := respAny.(*osc.ReadVpnConnectionsResponse)
	vpn := (*resp.VpnConnections)[0]

	data.RequestId = to.String(resp.ResponseContext.RequestId)

	if vpn.State == osc.VpnConnectionStateDeleted {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return r.flatten(ctx, data, &vpn)
}

func (r *vpnConnectionResource) flatten(ctx context.Context, data vpnConnectionModel, vpn *osc.VpnConnection) (vpnConnectionModel, error) {
	tags, diag := flattenOAPITagsFW(ctx, vpn.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	routes, diag := to.ListFromAttrType(ctx, r.flattenRoutes(vpn.Routes), to.ObjType(vpnConnectionRouteLightAttrTypes), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	vgwTelemetries, diag := to.ListFromAttrType(ctx, r.flattenVgwTelemetries(vpn.VgwTelemetries), to.ObjType(vpnConnectionVgwTelemetryAttrTypes), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Routes = routes
	data.Tags = tags
	data.VgwTelemetries = vgwTelemetries
	data.ClientGatewayConfiguration = to.String(ptr.From(vpn.ClientGatewayConfiguration))
	data.VpnConnectionId = to.String(vpn.VpnConnectionId)
	data.State = to.String(vpn.State)
	data.StaticRoutesOnly = to.Bool(vpn.StaticRoutesOnly)
	data.VirtualGatewayId = to.String(vpn.VirtualGatewayId)
	data.ConnectionType = to.String(vpn.ConnectionType)
	data.ClientGatewayId = to.String(vpn.ClientGatewayId)
	data.Id = to.String(vpn.VpnConnectionId)

	return data, nil
}

func (r *vpnConnectionResource) flattenRoutes(routes []osc.RouteLight) []vpnConnectionRouteLightModel {
	return lo.Map(routes, func(route osc.RouteLight, _ int) vpnConnectionRouteLightModel {
		return vpnConnectionRouteLightModel{
			DestinationIpRange: to.String(route.DestinationIpRange),
			RouteType:          to.String(route.RouteType),
			State:              to.String(route.State),
		}
	})
}

func (r *vpnConnectionResource) flattenVgwTelemetries(telemetries []osc.VgwTelemetry) []vpnConnectionVgwTelemetryModel {
	return lo.Map(telemetries, func(t osc.VgwTelemetry, _ int) vpnConnectionVgwTelemetryModel {
		return vpnConnectionVgwTelemetryModel{
			AcceptedRouteCount:  to.Int64(ptr.From(t.AcceptedRouteCount)),
			LastStateChangeDate: to.String(from.ISO8601(t.LastStateChangeDate)),
			OutsideIpAddress:    to.String(ptr.From(t.OutsideIpAddress)),
			State:               to.String(ptr.From(t.State)),
			StateDescription:    to.String(ptr.From(t.StateDescription)),
		}
	})
}

func (r *vpnConnectionResource) refreshFunc(id string) stateconf.StateRefreshFunc[osc.VpnConnectionState] {
	return func(ctx context.Context) (any, osc.VpnConnectionState, error) {
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
		vpnConnection := (*resp.VpnConnections)[0]

		return resp, vpnConnection.State, nil
	}
}
