package oapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                = &resourceRouteTableLink{}
	_ resource.ResourceWithConfigure   = &resourceRouteTableLink{}
	_ resource.ResourceWithImportState = &resourceRouteTableLink{}
	_ resource.ResourceWithModifyPlan  = &resourceRouteTableLink{}
)

type RouteTableLinkCoreModel struct {
	LinkRouteTableId types.String `tfsdk:"link_route_table_id"`
	Main             types.Bool   `tfsdk:"main"`
	NetId            types.String `tfsdk:"net_id"`
	RouteTableId     types.String `tfsdk:"route_table_id"`
	SubnetId         types.String `tfsdk:"subnet_id"`
}

type RouteTableLinkModel struct {
	RouteTableLinkCoreModel

	RequestId types.String   `tfsdk:"request_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
	Id        types.String   `tfsdk:"id"`
}

type resourceRouteTableLink struct {
	Client *oscgo.APIClient
}

func NewResourceRouteTableLink() resource.Resource {
	return &resourceRouteTableLink{}
}

func (r *resourceRouteTableLink) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceRouteTableLink) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	parts := strings.SplitN(req.ID, "_", 2)
	if len(parts) != 2 || req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Route Table Link identifier in the format {route_table_id}_{link_route_table_id}, got: %v", req.ID),
		)
		return
	}
	routeTableId := parts[0]
	linkRouteTableId := parts[1]

	var data RouteTableLinkModel
	var timeouts timeouts.Value
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.RouteTableId = types.StringValue(routeTableId)
	data.LinkRouteTableId = types.StringValue(linkRouteTableId)
	data.Id = types.StringValue(routeTableId + "_" + linkRouteTableId)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTableLink) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_table_link"
}

func (r *resourceRouteTableLink) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceRouteTableLink) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"link_route_table_id": schema.StringAttribute{
				Computed: true,
			},
			"main": schema.BoolAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
			},
			"route_table_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (r *resourceRouteTableLink) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouteTableLinkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.LinkRouteTableRequest{
		RouteTableId: data.RouteTableId.ValueString(),
		SubnetId:     data.SubnetId.ValueString(),
	}

	var createResp oscgo.LinkRouteTableResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteTableApi.LinkRouteTable(ctx).LinkRouteTableRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route Table Link.",
			err.Error(),
		)
		return
	}
	linkRouteTableId := createResp.GetLinkRouteTableId()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.LinkRouteTableId = types.StringValue(linkRouteTableId)
	data.Id = types.StringValue(linkRouteTableId)

	data, err = setRouteTableLinkState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table Link state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTableLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteTableLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setRouteTableLinkState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Route Table Link API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceRouteTableLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceRouteTableLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteTableLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.UnlinkRouteTableRequest{
		LinkRouteTableId: data.LinkRouteTableId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.RouteTableApi.UnlinkRouteTable(ctx).UnlinkRouteTableRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Route Table Link.",
			err.Error(),
		)
		return
	}
}

func setRouteTableLinkState(ctx context.Context, r *resourceRouteTableLink, data RouteTableLinkModel) (RouteTableLinkModel, error) {
	routeTableLinkFilters := oscgo.FiltersRouteTable{
		RouteTableIds: &[]string{data.RouteTableId.ValueString()},
	}
	readReq := oscgo.ReadRouteTablesRequest{
		Filters: &routeTableLinkFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'Route Table Link' read timeout value. Error: %v: ", diags.Errors())
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

	var routeTableLink oscgo.LinkRouteTable
	for _, elem := range readResp.GetRouteTables()[0].GetLinkRouteTables() {
		if elem.GetLinkRouteTableId() == data.LinkRouteTableId.ValueString() {
			routeTableLink = elem
		}
	}
	data.LinkRouteTableId = types.StringValue(routeTableLink.GetLinkRouteTableId())
	data.Main = types.BoolValue(routeTableLink.GetMain())
	data.NetId = types.StringValue(routeTableLink.GetNetId())
	data.RouteTableId = types.StringValue(routeTableLink.GetRouteTableId())
	data.SubnetId = types.StringValue(routeTableLink.GetSubnetId())
	data.Id = types.StringValue(routeTableLink.GetLinkRouteTableId())

	return data, nil
}
