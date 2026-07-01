package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource              = &virtualGatewayRoutePropagationResource{}
	_ resource.ResourceWithConfigure = &virtualGatewayRoutePropagationResource{}
)

const (
	virtualGatewayRoutePropagationErrModify = "Unable to modify Virtual Gateway Route Propagation"
	virtualGatewayRoutePropagationErrDelete = "Unable to disable Virtual Gateway Route Propagation"
)

type virtualGatewayRoutePropagationModel struct {
	VirtualGatewayId types.String   `tfsdk:"virtual_gateway_id"`
	RouteTableId     types.String   `tfsdk:"route_table_id"`
	Enable           types.Bool     `tfsdk:"enable"`
	RequestId        types.String   `tfsdk:"request_id"`
	Id               types.String   `tfsdk:"id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type virtualGatewayRoutePropagationResource struct {
	Client *osc.Client
}

func NewResourceVirtualGatewayRoutePropagation() resource.Resource {
	return &virtualGatewayRoutePropagationResource{}
}

func (r *virtualGatewayRoutePropagationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *virtualGatewayRoutePropagationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_gateway_route_propagation"
}

func (r *virtualGatewayRoutePropagationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"virtual_gateway_id": schema.StringAttribute{
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
			"enable": schema.BoolAttribute{
				Required: true,
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

func (r *virtualGatewayRoutePropagationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data virtualGatewayRoutePropagationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, diags := r.modifyRoutePropagation(ctx, data, timeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *virtualGatewayRoutePropagationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data virtualGatewayRoutePropagationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	data, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *virtualGatewayRoutePropagationResource) modifyRoutePropagation(ctx context.Context, data virtualGatewayRoutePropagationModel, timeout time.Duration) (virtualGatewayRoutePropagationModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	updateReq := osc.UpdateRoutePropagationRequest{
		VirtualGatewayId: data.VirtualGatewayId.ValueString(),
		RouteTableId:     data.RouteTableId.ValueString(),
		Enable:           data.Enable.ValueBool(),
	}
	updateResp, err := r.Client.UpdateRoutePropagation(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		diags.AddError(virtualGatewayRoutePropagationErrModify, err.Error())
		return data, diags
	}

	data.RequestId = to.String(updateResp.ResponseContext.RequestId)
	data.Id = to.String(data.VirtualGatewayId.ValueString() + "_" + data.RouteTableId.ValueString())

	return data, diags
}

func (r *virtualGatewayRoutePropagationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data virtualGatewayRoutePropagationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, diags := r.modifyRoutePropagation(ctx, data, timeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *virtualGatewayRoutePropagationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data virtualGatewayRoutePropagationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateRoutePropagationRequest{
		VirtualGatewayId: data.VirtualGatewayId.ValueString(),
		RouteTableId:     data.RouteTableId.ValueString(),
		Enable:           false,
	}
	_, err := r.Client.UpdateRoutePropagation(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(virtualGatewayRoutePropagationErrDelete, err.Error())
	}
}

func (r *virtualGatewayRoutePropagationResource) read(ctx context.Context, timeout time.Duration, data virtualGatewayRoutePropagationModel) (virtualGatewayRoutePropagationModel, error) {
	readResp, err := r.Client.ReadRouteTables(ctx, osc.ReadRouteTablesRequest{
		Filters: &osc.FiltersRouteTable{
			RouteTableIds: &[]string{data.RouteTableId.ValueString()},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if len(ptr.From(readResp.RouteTables)) == 0 {
		return data, ErrResourceEmpty
	}

	data.Id = to.String(data.VirtualGatewayId.ValueString() + "_" + data.RouteTableId.ValueString())

	return data, nil
}
