package outscale

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	_ resource.Resource                = &resourceRouteTable{}
	_ resource.ResourceWithConfigure   = &resourceRouteTable{}
	_ resource.ResourceWithImportState = &resourceRouteTable{}
	_ resource.ResourceWithModifyPlan  = &resourceRouteTable{}
)

type RouteTableModel struct {
	LinkRouteTables                 types.List     `tfsdk:"link_route_tables"`
	NetId                           types.String   `tfsdk:"net_id"`
	RoutePropagatingVirtualGateways types.List     `tfsdk:"route_propagating_virtual_gateways"`
	RouteTableId                    types.String   `tfsdk:"route_table_id"`
	Routes                          types.List     `tfsdk:"routes"`
	RequestId                       types.String   `tfsdk:"request_id"`
	Timeouts                        timeouts.Value `tfsdk:"timeouts"`
	Id                              types.String   `tfsdk:"id"`
	TagsModel
}

type RoutePropagatingVirtualGatewayModel struct {
	VirtualGatewayId types.String `tfsdk:"virtual_gateway_id"`
}

var linkRouteTableAttrTypes = utils.GetAttrTypes(RouteTableLinkCoreModel{})
var routePropagatingVirtualGatewayAttrTypes = utils.GetAttrTypes(RoutePropagatingVirtualGatewayModel{})
var routeAttrTypes = utils.GetAttrTypes(RouteCoreModel{})

func RoutesToModel(routes []oscgo.Route) []RouteCoreModel {
	routeModels := []RouteCoreModel{}

	for _, r := range routes {
		route := RouteCoreModel{
			CreationMethod:       types.StringValue(r.GetCreationMethod()),
			DestinationIpRange:   types.StringValue(r.GetDestinationIpRange()),
			DestinationServiceId: types.StringValue(r.GetDestinationServiceId()),
			GatewayId:            types.StringValue(r.GetGatewayId()),
			NatServiceId:         types.StringValue(r.GetNatServiceId()),
			NetAccessPointId:     types.StringValue(r.GetNetAccessPointId()),
			NetPeeringId:         types.StringValue(r.GetNetPeeringId()),
			NicId:                types.StringValue(r.GetNicId()),
			State:                types.StringValue(r.GetState()),
			VmAccountId:          types.StringValue(r.GetVmAccountId()),
			VmId:                 types.StringValue(r.GetVmId()),
		}
		routeModels = append(routeModels, route)
	}
	return routeModels
}

func LinkRouteTablesToModel(linkRouteTables []oscgo.LinkRouteTable) []RouteTableLinkCoreModel {
	linkRouteTableModels := []RouteTableLinkCoreModel{}

	for _, lrt := range linkRouteTables {
		link := RouteTableLinkCoreModel{
			LinkRouteTableId: types.StringValue(lrt.GetLinkRouteTableId()),
			Main:             types.BoolValue(lrt.GetMain()),
			NetId:            types.StringValue(lrt.GetNetId()),
			RouteTableId:     types.StringValue(lrt.GetRouteTableId()),
			SubnetId:         types.StringValue(lrt.GetSubnetId()),
		}
		linkRouteTableModels = append(linkRouteTableModels, link)
	}
	return linkRouteTableModels
}

func RoutePropagatingVirtualGatewaysToModel(routePropagatingVirtualGateways []oscgo.RoutePropagatingVirtualGateway) []RoutePropagatingVirtualGatewayModel {
	virtualGatewaysModels := []RoutePropagatingVirtualGatewayModel{}

	for _, vgw := range routePropagatingVirtualGateways {
		virtualGateway := RoutePropagatingVirtualGatewayModel{
			VirtualGatewayId: types.StringValue(vgw.GetVirtualGatewayId()),
		}
		virtualGatewaysModels = append(virtualGatewaysModels, virtualGateway)
	}
	return virtualGatewaysModels
}

type resourceRouteTable struct {
	Client *oscgo.APIClient
}

func NewResourceRouteTable() resource.Resource {
	return &resourceRouteTable{}
}

func (r *resourceRouteTable) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceRouteTable) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	RouteTableId := req.ID

	if RouteTableId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Route Table identifier, got: %v", req.ID),
		)
		return
	}

	var data RouteTableModel
	var timeouts timeouts.Value
	data.RouteTableId = types.StringValue(RouteTableId)
	data.Id = types.StringValue(RouteTableId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.LinkRouteTables = types.ListNull(types.ObjectType{AttrTypes: linkRouteTableAttrTypes})
	data.RoutePropagatingVirtualGateways = types.ListNull(types.ObjectType{AttrTypes: routePropagatingVirtualGatewayAttrTypes})
	data.Routes = types.ListNull(types.ObjectType{AttrTypes: routeAttrTypes})
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTable) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_table"
}

func (r *resourceRouteTable) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceRouteTable) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"route_table_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"route_propagating_virtual_gateways": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: routePropagatingVirtualGatewayAttrTypes,
				},
			},
			"link_route_tables": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: linkRouteTableAttrTypes,
				},
			},
			"routes": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: routeAttrTypes,
				},
			},
		},
	}
}

func (r *resourceRouteTable) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouteTableModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateRouteTableRequest{
		NetId: data.NetId.ValueString(),
	}

	var createResp oscgo.CreateRouteTableResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteTableApi.CreateRouteTable(ctx).CreateRouteTableRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route Table resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	routeTable := createResp.GetRouteTable()

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, routeTable.GetRouteTableId())
	if utils.CheckDiags(resp, diag) {
		return
	}

	data.RouteTableId = types.StringValue(routeTable.GetRouteTableId())
	data.Id = types.StringValue(routeTable.GetRouteTableId())
	data, err = setRouteTableState(ctx, r, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTable) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteTableModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setRouteTableState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Route Table API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTable) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.RouteTableId.ValueString())
	if utils.CheckDiags(resp, diag) {
		return
	}

	data, err := setRouteTableState(ctx, r, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTable) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteTableModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteRouteTableRequest{
		RouteTableId: data.RouteTableId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.RouteTableApi.DeleteRouteTable(ctx).DeleteRouteTableRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Route Table.",
			err.Error(),
		)
		return
	}
}

func setRouteTableState(ctx context.Context, r *resourceRouteTable, data RouteTableModel) (RouteTableModel, error) {
	routeTableFilters := oscgo.FiltersRouteTable{
		RouteTableIds: &[]string{data.RouteTableId.ValueString()},
	}
	readReq := oscgo.ReadRouteTablesRequest{
		Filters: &routeTableFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'Route Table' read timeout value. Error: %v: ", diags.Errors())
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
	if len(readResp.GetRouteTables()) == 0 {
		return data, errors.New("Empty")
	}

	routeTable := readResp.GetRouteTables()[0]
	tags, diag := flattenOAPITagsFW(ctx, routeTable.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	linkRouteTables, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: linkRouteTableAttrTypes}, LinkRouteTablesToModel(routeTable.GetLinkRouteTables()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Link Route Tables to the schema Set. Error: %v: ", diags.Errors())
	}
	routePropagatingVirtualGateways, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: routePropagatingVirtualGatewayAttrTypes}, RoutePropagatingVirtualGatewaysToModel(routeTable.GetRoutePropagatingVirtualGateways()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Route Propagating Virtual Gateways to the schema Set. Error: %v: ", diags.Errors())
	}
	routes, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: routeAttrTypes}, RoutesToModel(routeTable.GetRoutes()))
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Routes to the schema Set. Error: %v: ", diags.Errors())
	}

	data.LinkRouteTables = linkRouteTables
	data.NetId = types.StringValue(routeTable.GetNetId())
	data.RoutePropagatingVirtualGateways = routePropagatingVirtualGateways
	data.RouteTableId = types.StringValue(routeTable.GetRouteTableId())
	data.Routes = routes
	data.Id = types.StringValue(routeTable.GetRouteTableId())

	return data, nil
}
