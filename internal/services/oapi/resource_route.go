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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.RouteTableId = types.StringValue(routeTableId)
	data.DestinationIpRange = types.StringValue(destinationIpRange)
	data.Id = types.StringValue(routeTableId + "_" + destinationIpRange)
	data.AwaitActiveState = types.BoolValue(AwaitActiveStateDefaultValue)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateRouteRequest{
		DestinationIpRange: data.DestinationIpRange.ValueString(),
		RouteTableId:       data.RouteTableId.ValueString(),
	}

	if !data.GatewayId.IsUnknown() && !data.GatewayId.IsNull() {
		createReq.SetGatewayId(data.GatewayId.ValueString())
	}
	if !data.NatServiceId.IsUnknown() && !data.NatServiceId.IsNull() {
		createReq.SetNatServiceId(data.NatServiceId.ValueString())
	}
	if !data.VmId.IsUnknown() && !data.VmId.IsNull() {
		createReq.SetVmId(data.VmId.ValueString())
	}
	if !data.NicId.IsUnknown() && !data.NicId.IsNull() {
		createReq.SetNicId(data.NicId.ValueString())
	}
	if !data.NetPeeringId.IsUnknown() && !data.NetPeeringId.IsNull() {
		createReq.SetNetPeeringId(data.NetPeeringId.ValueString())
	}

	var createResp oscgo.CreateRouteResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteApi.CreateRoute(ctx).CreateRouteRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route.",
			err.Error(),
		)
		return
	}
	rt := createResp.GetRouteTable()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.Id = types.StringValue(rt.GetRouteTableId() + "_" + data.DestinationIpRange.ValueString())

	if data.AwaitActiveState.ValueBool() {
		stateConf := &retry.StateChangeConf{
			Target:     []string{"active"},
			Refresh:    ResourceRouteStateRefreshFunc(ctx, r, "blackhole", rt.GetRouteTableId(), data.DestinationIpRange.ValueString()),
			Timeout:    ReadDefaultTimeout,
			MinTimeout: 3 * time.Second,
			Delay:      2 * time.Second,
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

	stateData, err := r.setRouteState(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetRouteFromRouteTable(routeTable oscgo.RouteTable, destinationIpRange string) (oscgo.Route, error) {
	for _, route := range routeTable.GetRoutes() {
		if route.GetDestinationIpRange() == destinationIpRange {
			return route, nil
		}
	}
	return oscgo.Route{}, fmt.Errorf("unable to find matching route for Route Table (%s) "+
		"and destination CIDR block (%s)", routeTable.GetRouteTableId(), destinationIpRange)
}

func ResourceRouteStateRefreshFunc(ctx context.Context, r *resourceRoute, failState string, routeTableId string, destinationIpRange string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := oscgo.ReadRouteTablesRequest{Filters: &oscgo.FiltersRouteTable{
			RouteDestinationIpRanges: &[]string{destinationIpRange},
			RouteTableIds:            &[]string{routeTableId},
		}}
		var resp oscgo.ReadRouteTablesResponse

		err := retry.RetryContext(ctx, ReadDefaultTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.RouteTableApi.ReadRouteTables(ctx).ReadRouteTablesRequest(readReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return resp, "error", err
		}
		route, err := GetRouteFromRouteTable(resp.GetRouteTables()[0], destinationIpRange)
		if err != nil {
			return resp, "error", err
		}
		if route.GetState() == failState {
			return resp, failState, fmt.Errorf("failed to reach target state. Route is in '%v' failing state", failState)
		}

		return resp, route.GetState(), nil
	}
}

func (r *resourceRoute) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.setRouteState(ctx, data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	updateReq := oscgo.UpdateRouteRequest{
		RouteTableId:       stateData.RouteTableId.ValueString(),
		DestinationIpRange: stateData.DestinationIpRange.ValueString(),
	}
	if !planData.GatewayId.IsUnknown() && !planData.GatewayId.IsNull() && !planData.GatewayId.Equal(stateData.GatewayId) {
		updateReq.SetGatewayId(planData.GatewayId.ValueString())
		updateCall = true
	}
	if !planData.NatServiceId.IsUnknown() && !planData.NatServiceId.IsNull() && !planData.NatServiceId.Equal(stateData.NatServiceId) {
		updateReq.SetNatServiceId(planData.NatServiceId.ValueString())
		updateCall = true
	}
	if !planData.VmId.IsUnknown() && !planData.VmId.IsNull() && !planData.VmId.Equal(stateData.VmId) {
		updateReq.SetVmId(planData.VmId.ValueString())
		updateCall = true
	}
	if !planData.NicId.IsUnknown() && !planData.NicId.IsNull() && !planData.NicId.Equal(stateData.NicId) {
		updateReq.SetNicId(planData.NicId.ValueString())
		updateCall = true
	}
	if !planData.NetPeeringId.IsUnknown() && !planData.NetPeeringId.IsNull() && !planData.NetPeeringId.Equal(stateData.NetPeeringId) {
		updateReq.SetNetPeeringId(planData.NetPeeringId.ValueString())
		updateCall = true
	}

	if updateCall {
		updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var createResp oscgo.UpdateRouteResponse
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.RouteApi.UpdateRoute(ctx).UpdateRouteRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			createResp = rp
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Route resource.",
				err.Error(),
			)
			return
		}
		stateData.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())

		if stateData.AwaitActiveState.ValueBool() {
			stateConf := &retry.StateChangeConf{
				Target:     []string{"active"},
				Refresh:    ResourceRouteStateRefreshFunc(ctx, r, "blackhole", stateData.RouteTableId.ValueString(), stateData.DestinationIpRange.ValueString()),
				Timeout:    ReadDefaultTimeout,
				MinTimeout: 3 * time.Second,
				Delay:      2 * time.Second,
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

	data, err := r.setRouteState(ctx, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRoute) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteRouteRequest{
		DestinationIpRange: data.DestinationIpRange.ValueString(),
		RouteTableId:       data.RouteTableId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.RouteApi.DeleteRoute(ctx).DeleteRouteRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Route.",
			err.Error(),
		)
		return
	}
}

func (r *resourceRoute) setRouteState(ctx context.Context, data RouteModel) (RouteModel, error) {
	readReq := oscgo.ReadRouteTablesRequest{Filters: &oscgo.FiltersRouteTable{
		RouteDestinationIpRanges: &[]string{data.DestinationIpRange.ValueString()},
		RouteTableIds:            &[]string{data.RouteTableId.ValueString()},
	}}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'route' read timeout value. Error: %v: ", diags.Errors())
	}

	var readResp oscgo.ReadRouteTablesResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteTableApi.ReadRouteTables(ctx).ReadRouteTablesRequest(readReq).Execute()
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

	routeTable := readResp.GetRouteTables()[0]
	route, err := GetRouteFromRouteTable(routeTable, data.DestinationIpRange.ValueString())
	if err != nil {
		return data, errors.New("Empty")
	}

	data.GatewayId = types.StringValue(route.GetGatewayId())
	data.NatServiceId = types.StringValue(route.GetNatServiceId())
	data.NetPeeringId = types.StringValue(route.GetNetPeeringId())
	data.VmId = types.StringValue(route.GetVmId())
	data.NicId = types.StringValue(route.GetNicId())
	data.CreationMethod = types.StringValue(route.GetCreationMethod())
	data.DestinationIpRange = types.StringValue(route.GetDestinationIpRange())
	data.DestinationServiceId = types.StringValue(route.GetDestinationServiceId())
	data.NetAccessPointId = types.StringValue(route.GetNetAccessPointId())
	data.State = types.StringValue(route.GetState())
	data.VmAccountId = types.StringValue(route.GetVmAccountId())
	data.RouteTableId = types.StringValue(routeTable.GetRouteTableId())

	return data, nil
}
