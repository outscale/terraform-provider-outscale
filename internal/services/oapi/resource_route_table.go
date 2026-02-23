package oapi

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
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
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

var (
	linkRouteTableAttrTypes                 = fwhelpers.GetAttrTypes(RouteTableLinkCoreModel{})
	routePropagatingVirtualGatewayAttrTypes = fwhelpers.GetAttrTypes(RoutePropagatingVirtualGatewayModel{})
	routeAttrTypes                          = fwhelpers.GetAttrTypes(RouteCoreModel{})
)

func RoutesToModel(routes []osc.Route) []RouteCoreModel {
	routeModels := []RouteCoreModel{}

	for _, r := range routes {
		route := RouteCoreModel{
			CreationMethod:       to.String(r.CreationMethod),
			DestinationIpRange:   to.String(r.DestinationIpRange),
			DestinationServiceId: to.String(r.DestinationServiceId),
			GatewayId:            to.String(r.GatewayId),
			NatServiceId:         to.String(r.NatServiceId),
			NetAccessPointId:     to.String(r.NetAccessPointId),
			NetPeeringId:         to.String(r.NetPeeringId),
			NicId:                to.String(r.NicId),
			State:                to.String(r.State),
			VmAccountId:          to.String(r.VmAccountId),
			VmId:                 to.String(r.VmId),
		}
		routeModels = append(routeModels, route)
	}
	return routeModels
}

func LinkRouteTablesToModel(linkRouteTables []osc.LinkRouteTable) []RouteTableLinkCoreModel {
	linkRouteTableModels := []RouteTableLinkCoreModel{}

	for _, lrt := range linkRouteTables {
		link := RouteTableLinkCoreModel{
			LinkRouteTableId: to.String(lrt.LinkRouteTableId),
			Main:             to.Bool(lrt.Main),
			NetId:            to.String(lrt.NetId),
			RouteTableId:     to.String(lrt.RouteTableId),
			SubnetId:         to.String(lrt.SubnetId),
		}
		linkRouteTableModels = append(linkRouteTableModels, link)
	}
	return linkRouteTableModels
}

func RoutePropagatingVirtualGatewaysToModel(routePropagatingVirtualGateways []osc.RoutePropagatingVirtualGateway) []RoutePropagatingVirtualGatewayModel {
	virtualGatewaysModels := []RoutePropagatingVirtualGatewayModel{}

	for _, vgw := range routePropagatingVirtualGateways {
		virtualGateway := RoutePropagatingVirtualGatewayModel{
			VirtualGatewayId: to.String(vgw.VirtualGatewayId),
		}
		virtualGatewaysModels = append(virtualGatewaysModels, virtualGateway)
	}
	return virtualGatewaysModels
}

type resourceRouteTable struct {
	Client *osc.Client
}

func NewResourceRouteTable() resource.Resource {
	return &resourceRouteTable{}
}

func (r *resourceRouteTable) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	data.RouteTableId = to.String(RouteTableId)
	data.Id = to.String(RouteTableId)
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

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateRouteTableRequest{
		NetId: data.NetId.ValueString(),
	}

	createResp, err := r.Client.CreateRouteTable(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route Table resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	routeTable := ptr.From(createResp.RouteTable)

	diag := createOAPITagsFW(ctx, r.Client, createTimeout, data.Tags, routeTable.RouteTableId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data.RouteTableId = to.String(routeTable.RouteTableId)
	data.Id = to.String(routeTable.RouteTableId)
	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRouteTable) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteTableModel

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
			"Unable to set Route Table API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRouteTable) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.RouteTableId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := r.read(ctx, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRouteTable) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteTableModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteRouteTableRequest{
		RouteTableId: data.RouteTableId.ValueString(),
	}

	_, err := r.Client.DeleteRouteTable(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Route Table.",
			err.Error(),
		)
	}
}

func (r *resourceRouteTable) read(ctx context.Context, data RouteTableModel) (RouteTableModel, error) {
	routeTableFilters := osc.FiltersRouteTable{
		RouteTableIds: &[]string{data.RouteTableId.ValueString()},
	}
	readReq := osc.ReadRouteTablesRequest{
		Filters: &routeTableFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'route table' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadRouteTables(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	if readResp.RouteTables == nil || len(*readResp.RouteTables) == 0 {
		return data, ErrResourceEmpty
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	routeTable := (*readResp.RouteTables)[0]
	tags, diag := flattenOAPITagsFW(ctx, routeTable.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	linkRouteTables, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: linkRouteTableAttrTypes}, LinkRouteTablesToModel(routeTable.LinkRouteTables))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert link route tables to the schema set: %v", diags.Errors())
	}
	routePropagatingVirtualGateways, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: routePropagatingVirtualGatewayAttrTypes}, RoutePropagatingVirtualGatewaysToModel(routeTable.RoutePropagatingVirtualGateways))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert route propagating virtual gateways to the schema set: %v", diags.Errors())
	}
	routes, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: routeAttrTypes}, RoutesToModel(routeTable.Routes))
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert routes to the schema set: %v", diags.Errors())
	}

	data.LinkRouteTables = linkRouteTables
	data.NetId = to.String(routeTable.NetId)
	data.RoutePropagatingVirtualGateways = routePropagatingVirtualGateways
	data.RouteTableId = to.String(routeTable.RouteTableId)
	data.Routes = routes
	data.Id = to.String(routeTable.RouteTableId)

	return data, nil
}
