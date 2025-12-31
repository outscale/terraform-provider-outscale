package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	set "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                = &resourceNetAccessPoint{}
	_ resource.ResourceWithConfigure   = &resourceNetAccessPoint{}
	_ resource.ResourceWithImportState = &resourceNetAccessPoint{}
	_ resource.ResourceWithModifyPlan  = &resourceNetAccessPoint{}
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
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.NetAccessPointId = types.StringValue(netAccessPointId)
	data.Id = types.StringValue(netAccessPointId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.RouteTableIds = types.SetNull(types.StringType)
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateNetAccessPointRequest{
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

	var createResp oscgo.CreateNetAccessPointResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetAccessPointApi.CreateNetAccessPoint(ctx).CreateNetAccessPointRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Net Access Point resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	netAccessPoint := createResp.GetNetAccessPoint()

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, netAccessPoint.GetNetAccessPointId())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ResourceNetAccessPointStateRefreshFunc(ctx, createTimeout, r, netAccessPoint.GetNetAccessPointId()),
		Timeout:    createTimeout,
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error waiting for Net Access Point (%s) to become available.",
				netAccessPoint.GetNetAccessPointId()),
			err.Error(),
		)
		return
	}

	data.NetAccessPointId = types.StringValue(netAccessPoint.GetNetAccessPointId())
	data.Id = types.StringValue(netAccessPoint.GetNetAccessPointId())
	data, err = setNetAccessPointState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Access Point state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAccessPoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetAccessPointModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setNetAccessPointState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Net Access Point API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.NetAccessPointId.ValueString())
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if !planData.RouteTableIds.IsUnknown() && !planData.RouteTableIds.IsNull() {
		extractSet := func(ctx context.Context, ids types.Set) (set.Set[string], diag.Diagnostics) {
			rtIds, diag := to.Slice[string](ctx, ids)
			if diag.HasError() {
				return nil, diag
			}
			setIds := set.NewSet[string]()
			setIds.Append(rtIds...)
			return setIds, nil
		}

		planSet, diags := extractSet(ctx, planData.RouteTableIds)
		if diags.HasError() {
			return
		}
		stateSet, diags := extractSet(ctx, stateData.RouteTableIds)
		if diags.HasError() {
			return
		}

		addIds := planSet.Difference(stateSet).ToSlice()
		removeIds := stateSet.Difference(planSet).ToSlice()

		updateReq := oscgo.UpdateNetAccessPointRequest{
			NetAccessPointId: stateData.NetAccessPointId.ValueString(),
		}
		if len(addIds) > 0 {
			updateReq.AddRouteTableIds = &addIds
		}
		if len(removeIds) > 0 {
			updateReq.RemoveRouteTableIds = &removeIds
		}
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.NetAccessPointApi.UpdateNetAccessPoint(ctx).UpdateNetAccessPointRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Net Access Point resource.",
				err.Error(),
			)
			return
		}
	}

	data, err := setNetAccessPointState(ctx, r, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Access Point state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAccessPoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetAccessPointModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteNetAccessPointRequest{
		NetAccessPointId: data.NetAccessPointId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.NetAccessPointApi.DeleteNetAccessPoint(ctx).DeleteNetAccessPointRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Net Access Point.",
			err.Error(),
		)
		return
	}
}

func setNetAccessPointState(ctx context.Context, r *resourceNetAccessPoint, data NetAccessPointModel) (NetAccessPointModel, error) {
	readReq := oscgo.ReadNetAccessPointsRequest{
		Filters: &oscgo.FiltersNetAccessPoint{
			NetAccessPointIds: &[]string{data.NetAccessPointId.ValueString()},
		},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net access point' read timeout value. Error: %v: ", diags.Errors())
	}

	var readResp oscgo.ReadNetAccessPointsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetAccessPointApi.ReadNetAccessPoints(ctx).ReadNetAccessPointsRequest(readReq).Execute()
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

	if len(readResp.GetNetAccessPoints()) == 0 {
		return data, errors.New("Empty")
	}
	netAccessPoint := readResp.GetNetAccessPoints()[0]

	routeTablesIds, diags := types.SetValueFrom(ctx, types.StringType, netAccessPoint.GetRouteTableIds())
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert Route Tables Ids into a Set. Error: %v: ", diags.Errors())
	}
	tags, diag := flattenOAPITagsFW(ctx, netAccessPoint.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	data.RouteTableIds = routeTablesIds
	data.NetId = types.StringValue(netAccessPoint.GetNetId())
	data.NetAccessPointId = types.StringValue(netAccessPoint.GetNetAccessPointId())
	data.Id = types.StringValue(netAccessPoint.GetNetAccessPointId())
	data.ServiceName = types.StringValue(netAccessPoint.GetServiceName())
	data.State = types.StringValue(netAccessPoint.GetState())

	return data, nil
}

func ResourceNetAccessPointStateRefreshFunc(ctx context.Context, to time.Duration, r *resourceNetAccessPoint, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := oscgo.ReadNetAccessPointsRequest{Filters: &oscgo.FiltersNetAccessPoint{NetAccessPointIds: &[]string{id}}}
		var resp oscgo.ReadNetAccessPointsResponse

		err := retry.RetryContext(ctx, to, func() *retry.RetryError {
			rp, httpResp, err := r.Client.NetAccessPointApi.ReadNetAccessPoints(ctx).ReadNetAccessPointsRequest(readReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return resp, "error", err
		}
		nap := resp.GetNetAccessPoints()[0]

		return resp, nap.GetState(), nil
	}
}
