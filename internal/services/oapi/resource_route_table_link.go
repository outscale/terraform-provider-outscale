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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
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
	Client *osc.Client
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
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
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
	data.RouteTableId = to.String(routeTableId)
	data.LinkRouteTableId = to.String(linkRouteTableId)
	data.Id = to.String(routeTableId + "_" + linkRouteTableId)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.LinkRouteTableRequest{
		RouteTableId: data.RouteTableId.ValueString(),
		SubnetId:     data.SubnetId.ValueString(),
	}

	createResp, err := r.Client.LinkRouteTable(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Route Table Link.",
			err.Error(),
		)
		return
	}
	linkRouteTableId := ptr.From(createResp.LinkRouteTableId)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.LinkRouteTableId = to.String(linkRouteTableId)
	data.Id = to.String(linkRouteTableId)

	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Route Table Link state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRouteTableLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteTableLinkModel

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
			"Unable to set Route Table Link API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.UnlinkRouteTableRequest{
		LinkRouteTableId: data.LinkRouteTableId.ValueString(),
	}

	_, err := r.Client.UnlinkRouteTable(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Route Table Link.",
			err.Error(),
		)
	}
}

func (r *resourceRouteTableLink) read(ctx context.Context, data RouteTableLinkModel) (RouteTableLinkModel, error) {
	routeTableLinkFilters := osc.FiltersRouteTable{
		RouteTableIds: &[]string{data.RouteTableId.ValueString()},
	}
	readReq := osc.ReadRouteTablesRequest{
		Filters: &routeTableLinkFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'route table link' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadRouteTables(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	if readResp.RouteTables == nil || len(*readResp.RouteTables) == 0 {
		return data, ErrResourceEmpty
	}

	routeTableLink, _ := lo.Find((*readResp.RouteTables)[0].LinkRouteTables, func(elem osc.LinkRouteTable) bool {
		return elem.LinkRouteTableId == data.LinkRouteTableId.ValueString()
	})

	data.LinkRouteTableId = to.String(routeTableLink.LinkRouteTableId)
	data.Main = to.Bool(routeTableLink.Main)
	data.NetId = to.String(routeTableLink.NetId)
	data.RouteTableId = to.String(routeTableLink.RouteTableId)
	data.SubnetId = to.String(routeTableLink.SubnetId)
	data.Id = to.String(routeTableLink.LinkRouteTableId)

	return data, nil
}
