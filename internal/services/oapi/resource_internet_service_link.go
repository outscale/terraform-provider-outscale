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
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                = &resourceInternetServiceLink{}
	_ resource.ResourceWithConfigure   = &resourceInternetServiceLink{}
	_ resource.ResourceWithImportState = &resourceInternetServiceLink{}
	_ resource.ResourceWithModifyPlan  = &resourceInternetServiceLink{}
)

type InternetServiceLinkModel struct {
	InternetServiceId types.String   `tfsdk:"internet_service_id"`
	NetId             types.String   `tfsdk:"net_id"`
	State             types.String   `tfsdk:"state"`
	RequestId         types.String   `tfsdk:"request_id"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	TagsComputedModel
}

type resourceInternetServiceLink struct {
	Client *oscgo.APIClient
}

func NewResourceInternetServiceLink() resource.Resource {
	return &resourceInternetServiceLink{}
}

func (r *resourceInternetServiceLink) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceInternetServiceLink) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	internetServiceId := req.ID

	if internetServiceId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import internet service identifier. Got: %v", req.ID),
		)
		return
	}

	var data InternetServiceLinkModel
	var timeouts timeouts.Value
	data.InternetServiceId = types.StringValue(internetServiceId)
	data.Id = types.StringValue(internetServiceId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = ComputedTagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetServiceLink) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_internet_service_link"
}

func (r *resourceInternetServiceLink) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will destroy the state and unlink the Internet Service but not fully delete it. It will be gone when deleting the Internet Service.",
		)
	}
}

func (r *resourceInternetServiceLink) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"state": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"internet_service_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (r *resourceInternetServiceLink) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.LinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	var createResp oscgo.LinkInternetServiceResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.InternetServiceApi.LinkInternetService(ctx).LinkInternetServiceRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to link Internet Service resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())

	data.InternetServiceId = types.StringValue(data.InternetServiceId.ValueString())
	data.Id = types.StringValue(data.InternetServiceId.ValueString())

	data, err = setInternetServiceLinkState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Internet Service Link state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceInternetServiceLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setInternetServiceLinkState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
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

func (r *resourceInternetServiceLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceInternetServiceLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	unlinkReq := oscgo.UnlinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.InternetServiceApi.UnlinkInternetService(ctx).UnlinkInternetServiceRequest(unlinkReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink Internet Service resource.",
			err.Error(),
		)
		return
	}
}

func setInternetServiceLinkState(ctx context.Context, r *resourceInternetServiceLink, data InternetServiceLinkModel) (InternetServiceLinkModel, error) {
	internetServiceFilters := oscgo.FiltersInternetService{
		InternetServiceIds: &[]string{data.InternetServiceId.ValueString()},
	}
	readReq := oscgo.ReadInternetServicesRequest{
		Filters: &internetServiceFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'internet service' read timeout value. Error: %v: ", diags.Errors())
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
		return data, errors.New("Empty")
	}

	internetService := readResp.GetInternetServices()[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, internetService.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	data.NetId = types.StringValue(internetService.GetNetId())
	data.State = types.StringValue(internetService.GetState())
	data.Id = types.StringValue(internetService.GetInternetServiceId())
	data.InternetServiceId = types.StringValue(internetService.GetInternetServiceId())
	return data, nil
}

func ResourceInternetServiceLinkStateRefreshFunc(ctx context.Context, r *resourceInternetServiceLink, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := oscgo.ReadInternetServicesRequest{Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{id}}}
		var resp oscgo.ReadInternetServicesResponse

		err := retry.RetryContext(ctx, ReadDefaultTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.InternetServiceApi.ReadInternetServices(ctx).ReadInternetServicesRequest(readReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return resp, "error", err
		}
		internetService := resp.GetInternetServices()[0]

		return resp, internetService.GetState(), nil
	}
}
