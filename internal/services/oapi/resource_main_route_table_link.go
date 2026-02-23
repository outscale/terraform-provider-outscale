package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	_ resource.Resource               = &resourceMainRouteTableLink{}
	_ resource.ResourceWithConfigure  = &resourceMainRouteTableLink{}
	_ resource.ResourceWithModifyPlan = &resourceMainRouteTableLink{}
)

type MainRouteTableLinkModel struct {
	RouteTableLinkCoreModel
	DefaultRouteTableId types.String `tfsdk:"default_route_table_id"`

	RequestId types.String   `tfsdk:"request_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
	Id        types.String   `tfsdk:"id"`
}

type resourceMainRouteTableLink struct {
	Client *osc.Client
}

func NewResourceMainRouteTableLink() resource.Resource {
	return &resourceMainRouteTableLink{}
}

func (r *resourceMainRouteTableLink) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceMainRouteTableLink) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_main_route_table_link"
}

func (r *resourceMainRouteTableLink) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Destroying a `main_route_table_link` resets the original route table as the main for the Net. Ensure the additional route table remains intact for this operation to succeed (see: https://registry.terraform.io/providers/outscale/outscale/latest/docs/resources/main_route_table_link).",
		)
	}
}

func (r *resourceMainRouteTableLink) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"route_table_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Computed: true,
			},
			"default_route_table_id": schema.StringAttribute{
				Computed: true,
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

func (r *resourceMainRouteTableLink) GetAssociatedRouteTable(ctx context.Context, to time.Duration, data MainRouteTableLinkModel) (*osc.ReadRouteTablesResponse, error) {
	routeTableFilters := osc.FiltersRouteTable{
		NetIds:             &[]string{data.NetId.ValueString()},
		LinkRouteTableMain: new(true),
	}
	readReq := osc.ReadRouteTablesRequest{
		Filters: &routeTableFilters,
	}

	readResp, err := r.Client.ReadRouteTables(ctx, readReq, options.WithRetryTimeout(to))
	if err != nil {
		return readResp, err
	}

	return readResp, nil
}

func (r *resourceMainRouteTableLink) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MainRouteTableLinkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	routeTableResp, err := r.GetAssociatedRouteTable(ctx, createTimeout, data)
	if err != nil {
		return
	}
	routeTable := ptr.From(routeTableResp.RouteTables)[0]
	oldLinkRouteTableId := routeTable.LinkRouteTables[0].LinkRouteTableId
	defaultRouteTableId := routeTable.LinkRouteTables[0].RouteTableId

	createReq := osc.UpdateRouteTableLinkRequest{
		RouteTableId:     data.RouteTableId.ValueString(),
		LinkRouteTableId: oldLinkRouteTableId,
	}

	createResp, err := r.Client.UpdateRouteTableLink(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set the Main Route Table.",
			err.Error(),
		)
		return
	}
	data.DefaultRouteTableId = to.String(defaultRouteTableId)

	linkRouteTableId := ptr.From(createResp.LinkRouteTableId)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.LinkRouteTableId = to.String(linkRouteTableId)
	data.Id = to.String(linkRouteTableId)

	data, err = r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Main Route Table Link state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceMainRouteTableLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MainRouteTableLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	data, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Main Route Table Link API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceMainRouteTableLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceMainRouteTableLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MainRouteTableLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.UpdateRouteTableLinkRequest{
		LinkRouteTableId: data.LinkRouteTableId.ValueString(),
		RouteTableId:     data.DefaultRouteTableId.ValueString(),
	}

	_, err := r.Client.UpdateRouteTableLink(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink the Main Route Table.",
			err.Error(),
		)
	}
}

func (r *resourceMainRouteTableLink) read(ctx context.Context, timeout time.Duration, data MainRouteTableLinkModel) (MainRouteTableLinkModel, error) {
	routeTableResp, err := r.GetAssociatedRouteTable(ctx, timeout, data)
	if err != nil {
		return data, err
	}
	if routeTableResp.RouteTables == nil || len(*routeTableResp.RouteTables) == 0 {
		return data, ErrResourceEmpty
	}
	routeTable := (*routeTableResp.RouteTables)[0]

	data.RequestId = to.String(routeTableResp.ResponseContext.RequestId)

	var mainRouteTableLink osc.LinkRouteTable
	for _, elem := range routeTable.LinkRouteTables {
		if elem.LinkRouteTableId == data.LinkRouteTableId.ValueString() {
			mainRouteTableLink = elem
		}
	}
	data.LinkRouteTableId = to.String(mainRouteTableLink.LinkRouteTableId)
	data.Main = to.Bool(mainRouteTableLink.Main)
	data.NetId = to.String(mainRouteTableLink.NetId)
	data.RouteTableId = to.String(mainRouteTableLink.RouteTableId)
	data.SubnetId = to.String(mainRouteTableLink.SubnetId)
	data.Id = to.String(mainRouteTableLink.LinkRouteTableId)

	return data, nil
}
