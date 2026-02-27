package oapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                   = &resourceRoute{}
	_ resource.ResourceWithConfigure      = &resourceRoute{}
	_ resource.ResourceWithImportState    = &resourceRoute{}
	_ resource.ResourceWithModifyPlan     = &resourceRoute{}
	_ resource.ResourceWithValidateConfig = &resourceRoute{}
)

type RouteCoreModel struct {
	CreationMethod       types.String `tfsdk:"creation_method"`
	DestinationIpRange   types.String `tfsdk:"destination_ip_range"`
	DestinationServiceId types.String `tfsdk:"destination_service_id"`
	GatewayId            types.String `tfsdk:"gateway_id"`
	NatServiceId         types.String `tfsdk:"nat_service_id"`
	NetAccessPointId     types.String `tfsdk:"net_access_point_id"`
	NetPeeringId         types.String `tfsdk:"net_peering_id"`
	NicId                types.String `tfsdk:"nic_id"`
	State                types.String `tfsdk:"state"`
	VmAccountId          types.String `tfsdk:"vm_account_id"`
	VmId                 types.String `tfsdk:"vm_id"`
}

type RouteModel struct {
	RouteCoreModel

	RouteTableId     types.String `tfsdk:"route_table_id"`
	AwaitActiveState types.Bool   `tfsdk:"await_active_state"`

	Timeouts  timeouts.Value `tfsdk:"timeouts"`
	RequestId types.String   `tfsdk:"request_id"`
	Id        types.String   `tfsdk:"id"`
}

type resourceRoute struct {
	Client *osc.Client
}

func NewResourceRoute() resource.Resource {
	return &resourceRoute{}
}

func (r *resourceRoute) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceRoute) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	parts := strings.SplitN(req.ID, "_", 2)
	if len(parts) != 2 || req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Route identifier in the format {route_table_id}_{destination_ip_range}, got: %v", req.ID),
		)
		return
	}
	routeTableId := parts[0]
	destinationIpRange := parts[1]

	var data RouteModel
	var timeouts timeouts.Value
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.RouteTableId = to.String(routeTableId)
	data.DestinationIpRange = to.String(destinationIpRange)
	data.Id = to.String(routeTableId + "_" + destinationIpRange)
	data.AwaitActiveState = to.Bool(AwaitActiveStateDefaultValue)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceRoute) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *resourceRoute) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceRoute) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	configAttrCount := 0
	exclusiveTargets := []string{
		"gateway_id",
		"nat_service_id",
		"net_peering_id",
		"nic_id",
		"vm_id",
	}
	for _, target := range exclusiveTargets {
		var valConfig types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(target), &valConfig)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !valConfig.IsNull() {
			configAttrCount++
		}
	}
	if configAttrCount != 1 {
		resp.Diagnostics.AddError(
			"Attribute Configuration",
			fmt.Sprintf("Exactly one of %v should be set.", exclusiveTargets),
		)
	}
}

func (r *resourceRoute) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"creation_method": schema.StringAttribute{
				Computed: true,
			},
			"destination_ip_range": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_service_id": schema.StringAttribute{
				Computed: true,
			},
			"gateway_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"nat_service_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"net_access_point_id": schema.StringAttribute{
				Computed: true,
			},
			"net_peering_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"nic_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"vm_account_id": schema.StringAttribute{
				Computed: true,
			},
			"vm_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"route_table_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"await_active_state": schema.BoolAttribute{
				Computed:           true,
				Optional:           true,
				Default:            booldefault.StaticBool(AwaitActiveStateDefaultValue),
				DeprecationMessage: "Route's state is always active. The attribute will be removed in the next major version.",
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

func (r *resourceRoute) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouteModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateRouteRequest{
		DestinationIpRange: data.DestinationIpRange.ValueString(),
		RouteTableId:       data.RouteTableId.ValueString(),
	}

	if !data.GatewayId.IsUnknown() && !data.GatewayId.IsNull() {
		createReq.GatewayId = data.GatewayId.ValueStringPointer()
	}
	if !data.NatServiceId.IsUnknown() && !data.NatServiceId.IsNull() {
		createReq.NatServiceId = data.NatServiceId.ValueStringPointer()
	}
	if !data.VmId.IsUnknown() && !data.VmId.IsNull() {
		createReq.VmId = data.VmId.ValueStringPointer()
	}
	if !data.NicId.IsUnknown() && !data.NicId.IsNull() {
		createReq.NicId = data.NicId.ValueStringPointer()
	}
	if !data.NetPeeringId.IsUnknown() && !data.NetPeeringId.IsNull() {
		createReq.NetPeeringId = data.NetPeeringId.ValueStringPointer()
	}

	createResp, err := r.Client.CreateRoute(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route.",
			err.Error(),
		)
		return
	}
	rt := ptr.From(createResp.RouteTable)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(rt.RouteTableId + "_" + data.DestinationIpRange.ValueString())

	if data.AwaitActiveState.ValueBool() {
		stateConf := &retry.StateChangeConf{
			Target:  []string{"active"},
			Refresh: ResourceRouteStateRefreshFunc(ctx, r, "blackhole", rt.RouteTableId, data.DestinationIpRange.ValueString()),
			Timeout: createTimeout,
		}

		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error waiting for Route (%s) to become active.",
					data.Id.ValueString(),
				),
				err.Error(),
			)
			return
		}
	}

	stateData, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func getRouteFromRouteTable(routeTable osc.RouteTable, destinationIpRange string) (osc.Route, error) {
	for _, route := range routeTable.Routes {
		if route.DestinationIpRange == destinationIpRange {
			return route, nil
		}
	}
	return osc.Route{}, fmt.Errorf("unable to find matching route for route table (%s) "+
		"and destination CIDR block (%s)", routeTable.RouteTableId, destinationIpRange)
}

func ResourceRouteStateRefreshFunc(ctx context.Context, r *resourceRoute, failState string, routeTableId string, destinationIpRange string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := osc.ReadRouteTablesRequest{Filters: &osc.FiltersRouteTable{
			RouteDestinationIpRanges: &[]string{destinationIpRange},
			RouteTableIds:            &[]string{routeTableId},
		}}

		resp, err := r.Client.ReadRouteTables(ctx, readReq, options.WithRetryTimeout(ReadDefaultTimeout))
		if err != nil {
			return resp, "error", err
		}
		route, err := getRouteFromRouteTable((*resp.RouteTables)[0], destinationIpRange)
		if err != nil {
			return resp, "error", err
		}
		if route.State == failState {
			return resp, failState, fmt.Errorf("failed to reach target state: route is in '%v' failing state", failState)
		}

		return resp, route.State, nil
	}
}

func (r *resourceRoute) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteModel

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
			"Unable to set Route API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRoute) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData RouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.AwaitActiveState.IsUnknown() && !planData.AwaitActiveState.IsNull() {
		stateData.AwaitActiveState = planData.AwaitActiveState
	}

	updateCall := false
	updateReq := osc.UpdateRouteRequest{
		RouteTableId:       stateData.RouteTableId.ValueString(),
		DestinationIpRange: stateData.DestinationIpRange.ValueString(),
	}
	if !planData.GatewayId.IsUnknown() && !planData.GatewayId.IsNull() && !planData.GatewayId.Equal(stateData.GatewayId) {
		updateReq.GatewayId = planData.GatewayId.ValueStringPointer()
		updateCall = true
	}
	if !planData.NatServiceId.IsUnknown() && !planData.NatServiceId.IsNull() && !planData.NatServiceId.Equal(stateData.NatServiceId) {
		updateReq.NatServiceId = planData.NatServiceId.ValueStringPointer()
		updateCall = true
	}
	if !planData.VmId.IsUnknown() && !planData.VmId.IsNull() && !planData.VmId.Equal(stateData.VmId) {
		updateReq.VmId = planData.VmId.ValueStringPointer()
		updateCall = true
	}
	if !planData.NicId.IsUnknown() && !planData.NicId.IsNull() && !planData.NicId.Equal(stateData.NicId) {
		updateReq.NicId = planData.NicId.ValueStringPointer()
		updateCall = true
	}
	if !planData.NetPeeringId.IsUnknown() && !planData.NetPeeringId.IsNull() && !planData.NetPeeringId.Equal(stateData.NetPeeringId) {
		updateReq.NetPeeringId = planData.NetPeeringId.ValueStringPointer()
		updateCall = true
	}

	if updateCall {
		updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		createResp, err := r.Client.UpdateRoute(ctx, updateReq, options.WithRetryTimeout(updateTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Route resource.",
				err.Error(),
			)
			return
		}
		stateData.RequestId = to.String(createResp.ResponseContext.RequestId)

		if stateData.AwaitActiveState.ValueBool() {
			stateConf := &retry.StateChangeConf{
				Target:  []string{"active"},
				Timeout: updateTimeout,
				Refresh: ResourceRouteStateRefreshFunc(ctx, r, "blackhole", stateData.RouteTableId.ValueString(), stateData.DestinationIpRange.ValueString()),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error waiting for Route (%s) to become active.",
						stateData.Id.ValueString(),
					),
					err.Error(),
				)
				return
			}
		}
	}

	data, err := r.read(ctx, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRoute) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteRouteRequest{
		DestinationIpRange: data.DestinationIpRange.ValueString(),
		RouteTableId:       data.RouteTableId.ValueString(),
	}

	_, err := r.Client.DeleteRoute(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Route.",
			err.Error(),
		)
	}
}

func (r *resourceRoute) read(ctx context.Context, data RouteModel) (RouteModel, error) {
	readReq := osc.ReadRouteTablesRequest{Filters: &osc.FiltersRouteTable{
		RouteDestinationIpRanges: &[]string{data.DestinationIpRange.ValueString()},
		RouteTableIds:            &[]string{data.RouteTableId.ValueString()},
	}}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'route' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadRouteTables(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	routeTable := (ptr.From(readResp.RouteTables))[0]
	route, err := getRouteFromRouteTable(routeTable, data.DestinationIpRange.ValueString())
	if err != nil {
		return data, ErrResourceEmpty
	}

	data.GatewayId = to.String(ptr.From(route.GatewayId))
	data.NatServiceId = to.String(ptr.From(route.NatServiceId))
	data.NetPeeringId = to.String(ptr.From(route.NetPeeringId))
	data.VmId = to.String(ptr.From(route.VmId))
	data.NicId = to.String(ptr.From(route.NicId))
	data.CreationMethod = to.String(route.CreationMethod)
	data.DestinationIpRange = to.String(route.DestinationIpRange)
	data.DestinationServiceId = to.String(ptr.From(route.DestinationServiceId))
	data.NetAccessPointId = to.String(ptr.From(route.NetAccessPointId))
	data.State = to.String(route.State)
	data.VmAccountId = to.String(ptr.From(route.VmAccountId))
	data.RouteTableId = to.String(routeTable.RouteTableId)

	return data, nil
}
