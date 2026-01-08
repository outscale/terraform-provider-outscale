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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                = &resourceInternetService{}
	_ resource.ResourceWithConfigure   = &resourceInternetService{}
	_ resource.ResourceWithImportState = &resourceInternetService{}
	_ resource.ResourceWithModifyPlan  = &resourceInternetService{}
)

type InternetServiceModel struct {
	InternetServiceId types.String   `tfsdk:"internet_service_id"`
	NetId             types.String   `tfsdk:"net_id"`
	State             types.String   `tfsdk:"state"`
	RequestId         types.String   `tfsdk:"request_id"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	TagsModel
}

type resourceInternetService struct {
	Client *oscgo.APIClient
}

func NewResourceInternetService() resource.Resource {
	return &resourceInternetService{}
}

func (r *resourceInternetService) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceInternetService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	internetServiceId := req.ID

	if internetServiceId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import internet service identifier. Got: %v", req.ID),
		)
		return
	}

	var data InternetServiceModel
	var timeouts timeouts.Value
	data.InternetServiceId = types.StringValue(internetServiceId)
	data.Id = types.StringValue(internetServiceId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_internet_service"
}

func (r *resourceInternetService) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceInternetService) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"state": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
			},
			"internet_service_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

func (r *resourceInternetService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InternetServiceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateInternetServiceRequest{}

	var createResp oscgo.CreateInternetServiceResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.InternetServiceApi.CreateInternetService(ctx).CreateInternetServiceRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Internet Service resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	internetService := createResp.GetInternetService()

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, internetService.GetInternetServiceId())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data.InternetServiceId = types.StringValue(internetService.GetInternetServiceId())
	data.Id = types.StringValue(internetService.GetInternetServiceId())

	data, err = setInternetServiceState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Internet Service state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InternetServiceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setInternetServiceState(ctx, r, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Internet Service API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData InternetServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.InternetServiceId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := setInternetServiceState(ctx, r, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Internet Service state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InternetServiceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.InternetServiceApi.DeleteInternetService(ctx).DeleteInternetServiceRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Internet Service.",
			err.Error(),
		)
		return
	}
}

func setInternetServiceState(ctx context.Context, r *resourceInternetService, data InternetServiceModel) (InternetServiceModel, error) {
	internetServiceFilters := oscgo.FiltersInternetService{
		InternetServiceIds: &[]string{data.InternetServiceId.ValueString()},
	}
	readReq := oscgo.ReadInternetServicesRequest{
		Filters: &internetServiceFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'internet service' read timeout value: %v", diags.Errors())
	}

	var readResp oscgo.ReadInternetServicesResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.InternetServiceApi.ReadInternetServices(ctx).ReadInternetServicesRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}
	if len(readResp.GetInternetServices()) == 0 {
		return data, ErrResourceEmpty
	}

	internetService := readResp.GetInternetServices()[0]

	tags, diag := flattenOAPITagsFW(ctx, internetService.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags
	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	data.NetId = types.StringValue(internetService.GetNetId())
	data.State = types.StringValue(internetService.GetState())
	data.Id = types.StringValue(internetService.GetInternetServiceId())
	data.InternetServiceId = types.StringValue(internetService.GetInternetServiceId())
	return data, nil
}
