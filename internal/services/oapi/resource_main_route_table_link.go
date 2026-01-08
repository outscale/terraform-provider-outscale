package oapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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

func (r *resourceMainRouteTableLink) GetAssociatedRouteTable(ctx context.Context, data MainRouteTableLinkModel) (oscgo.ReadRouteTablesResponse, error) {
	var readResp oscgo.ReadRouteTablesResponse
	routeTableFilters := oscgo.FiltersRouteTable{
		NetIds:             &[]string{data.NetId.ValueString()},
		LinkRouteTableMain: &[]bool{true}[0],
	}
	readReq := oscgo.ReadRouteTablesRequest{
		Filters: &routeTableFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return readResp, fmt.Errorf("unable to parse 'route table' read timeout value: %v", diags.Errors())
	}

	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteTableApi.ReadRouteTables(ctx).ReadRouteTablesRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	routeTableResp, err := r.GetAssociatedRouteTable(ctx, data)
	if err != nil {
		return
	}
	routeTable := routeTableResp.GetRouteTables()[0]
	oldLinkRouteTableId := routeTable.GetLinkRouteTables()[0].GetLinkRouteTableId()
	defaultRouteTableId := routeTable.GetLinkRouteTables()[0].GetRouteTableId()

	createReq := oscgo.UpdateRouteTableLinkRequest{
		RouteTableId:     data.RouteTableId.ValueString(),
		LinkRouteTableId: oldLinkRouteTableId,
	}

	var createResp oscgo.UpdateRouteTableLinkResponse
	err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.RouteTableApi.UpdateRouteTableLink(ctx).UpdateRouteTableLinkRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set the Main Route Table.",
			err.Error(),
		)
		return
	}
	data.DefaultRouteTableId = types.StringValue(defaultRouteTableId)
	linkRouteTableId := createResp.GetLinkRouteTableId()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.LinkRouteTableId = types.StringValue(linkRouteTableId)
	data.Id = types.StringValue(linkRouteTableId)

	data, err = setMainRouteTableLinkState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Main Route Table Link state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceMainRouteTableLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MainRouteTableLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setMainRouteTableLinkState(ctx, r, data)
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.UpdateRouteTableLinkRequest{
		LinkRouteTableId: data.LinkRouteTableId.ValueString(),
		RouteTableId:     data.DefaultRouteTableId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.RouteTableApi.UpdateRouteTableLink(ctx).UpdateRouteTableLinkRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink the Main Route Table.",
			err.Error(),
		)
		return
	}
}

func setMainRouteTableLinkState(ctx context.Context, r *resourceMainRouteTableLink, data MainRouteTableLinkModel) (MainRouteTableLinkModel, error) {
	routeTableResp, err := r.GetAssociatedRouteTable(ctx, data)
	routeTable := routeTableResp.GetRouteTables()[0]
	if err != nil {
		return data, err
	}
	data.RequestId = types.StringValue(routeTableResp.ResponseContext.GetRequestId())
	if len(routeTableResp.GetRouteTables()) == 0 {
		return data, ErrResourceEmpty
	}

	var mainRouteTableLink oscgo.LinkRouteTable
	for _, elem := range routeTable.GetLinkRouteTables() {
		if elem.GetLinkRouteTableId() == data.LinkRouteTableId.ValueString() {
			mainRouteTableLink = elem
		}
	}
	data.LinkRouteTableId = types.StringValue(mainRouteTableLink.GetLinkRouteTableId())
	data.Main = types.BoolValue(mainRouteTableLink.GetMain())
	data.NetId = types.StringValue(mainRouteTableLink.GetNetId())
	data.RouteTableId = types.StringValue(mainRouteTableLink.GetRouteTableId())
	data.SubnetId = types.StringValue(mainRouteTableLink.GetSubnetId())
	data.Id = types.StringValue(mainRouteTableLink.GetLinkRouteTableId())

	return data, nil
}
