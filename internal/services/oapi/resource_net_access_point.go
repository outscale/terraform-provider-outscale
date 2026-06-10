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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceNetAccessPoint{}
	_ resource.ResourceWithConfigure   = &resourceNetAccessPoint{}
	_ resource.ResourceWithImportState = &resourceNetAccessPoint{}
	_ resource.ResourceWithModifyPlan  = &resourceNetAccessPoint{}
)

const (
	netAccessPointErrCreate = "Unable to create Net Access Point"
	netAccessPointErrUpdate = "Unable to update Net Access Point"
	netAccessPointErrDelete = "Unable to delete Net Access Point"
	netAccessPointErrWait   = "Unable to wait for Net Access Point state"
)

type NetAccessPointModel struct {
	NetAccessPointId types.String   `tfsdk:"net_access_point_id"`
	NetId            types.String   `tfsdk:"net_id"`
	RouteTableIds    types.Set      `tfsdk:"route_table_ids"`
	ServiceName      types.String   `tfsdk:"service_name"`
	State            types.String   `tfsdk:"state"`
	RequestId        types.String   `tfsdk:"request_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id               types.String   `tfsdk:"id"`
	TagsModel
}

type resourceNetAccessPoint struct {
	Client *osc.Client
}

func NewResourceNetAccessPoint() resource.Resource {
	return &resourceNetAccessPoint{}
}

func (r *resourceNetAccessPoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
}

func (r *resourceNetAccessPoint) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	netAccessPointId := req.ID

	if netAccessPointId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Net Access Point identifier, got: %v", req.ID),
		)
		return
	}

	var data NetAccessPointModel
	var timeouts timeouts.Value
	data.NetAccessPointId = to.String(netAccessPointId)
	data.Id = to.String(netAccessPointId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.RouteTableIds = types.SetNull(types.StringType)
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceNetAccessPoint) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net_access_point"
}

func (r *resourceNetAccessPoint) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceNetAccessPoint) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"net_access_point_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"route_table_ids": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *resourceNetAccessPoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetAccessPointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateNetAccessPointRequest{
		NetId:       data.NetId.ValueString(),
		ServiceName: data.ServiceName.ValueString(),
	}

	if !data.RouteTableIds.IsUnknown() && !data.RouteTableIds.IsNull() {
		var rtIds []string
		diags = data.RouteTableIds.ElementsAs(ctx, &rtIds, false)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		createReq.RouteTableIds = &rtIds
	}

	createResp, err := r.Client.CreateNetAccessPoint(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(netAccessPointErrCreate, err.Error())
		return
	}

	netAccessPoint := *createResp.NetAccessPoint
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.NetAccessPointId = to.String(netAccessPoint.NetAccessPointId)
	data.Id = to.String(netAccessPoint.NetAccessPointId)

	stateData, err := r.flatten(ctx, data, netAccessPoint)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, netAccessPoint.NetAccessPointId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.NetAccessPointState]{
		Pending: stateconf.States(osc.NetAccessPointStatePending),
		Target:  stateconf.States(osc.NetAccessPointStateAvailable),
		Timeout: timeout,
		Refresh: r.refreshFunc(netAccessPoint.NetAccessPointId),
	}

	napAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(netAccessPointErrWait, err.Error())
		return
	}

	// We set the last read response to the state
	stateData, err = r.flatten(ctx, data, napAny.(osc.NetAccessPoint))
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceNetAccessPoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetAccessPointModel

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

func (r *resourceNetAccessPoint) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData NetAccessPointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diags = updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.NetAccessPointId.ValueString())
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if !planData.RouteTableIds.IsUnknown() && !planData.RouteTableIds.IsNull() {
		planSlice, diags := to.Slice[string](ctx, planData.RouteTableIds)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}
		stateSlice, diags := to.Slice[string](ctx, stateData.RouteTableIds)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}
		addIds, removeIds := lo.Difference(planSlice, stateSlice)

		updateReq := osc.UpdateNetAccessPointRequest{
			NetAccessPointId: stateData.NetAccessPointId.ValueString(),
		}
		if len(addIds) > 0 {
			updateReq.AddRouteTableIds = &addIds
		}
		if len(removeIds) > 0 {
			updateReq.RemoveRouteTableIds = &removeIds
		}
		_, err := r.Client.UpdateNetAccessPoint(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(netAccessPointErrUpdate, err.Error())
			return
		}
	}

	stateData.Timeouts = planData.Timeouts
	data, err := r.read(ctx, timeout, stateData)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *resourceNetAccessPoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetAccessPointModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteNetAccessPointRequest{
		NetAccessPointId: data.NetAccessPointId.ValueString(),
	}

	_, err := r.Client.DeleteNetAccessPoint(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(netAccessPointErrDelete, err.Error())
	}
}

func (r *resourceNetAccessPoint) read(ctx context.Context, timeout time.Duration, data NetAccessPointModel) (NetAccessPointModel, error) {
	readReq := osc.ReadNetAccessPointsRequest{
		Filters: &osc.FiltersNetAccessPoint{
			NetAccessPointIds: &[]string{data.NetAccessPointId.ValueString()},
		},
	}

	readResp, err := r.Client.ReadNetAccessPoints(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	if readResp.NetAccessPoints == nil || len(*readResp.NetAccessPoints) == 0 {
		return data, ErrResourceEmpty
	}
	netAccessPoint := (*readResp.NetAccessPoints)[0]

	return r.flatten(ctx, data, netAccessPoint)
}

func (r *resourceNetAccessPoint) flatten(ctx context.Context, data NetAccessPointModel, netAccessPoint osc.NetAccessPoint) (NetAccessPointModel, error) {
	routeTablesIds, diags := types.SetValueFrom(ctx, types.StringType, netAccessPoint.RouteTableIds)
	if diags.HasError() {
		return data, from.Diag(diags)
	}
	tags, diag := flattenOAPITagsFW(ctx, netAccessPoint.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.Tags = tags

	data.RouteTableIds = routeTablesIds
	data.NetId = to.String(netAccessPoint.NetId)
	data.NetAccessPointId = to.String(netAccessPoint.NetAccessPointId)
	data.Id = to.String(netAccessPoint.NetAccessPointId)
	data.ServiceName = to.String(netAccessPoint.ServiceName)
	data.State = to.String(netAccessPoint.State)

	return data, nil
}

func (r *resourceNetAccessPoint) refreshFunc(id string) func(ctx context.Context) (any, osc.NetAccessPointState, error) {
	return func(ctx context.Context) (any, osc.NetAccessPointState, error) {
		readReq := osc.ReadNetAccessPointsRequest{Filters: &osc.FiltersNetAccessPoint{NetAccessPointIds: &[]string{id}}}

		resp, err := r.Client.ReadNetAccessPoints(ctx, readReq)
		if err != nil {
			return resp, "", err
		}
		if resp.NetAccessPoints == nil || len(*resp.NetAccessPoints) == 0 {
			return resp, "", ErrResourceEmpty
		}

		nap := (*resp.NetAccessPoints)[0]
		return nap, nap.State, nil
	}
}
