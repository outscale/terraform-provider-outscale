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
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
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
	Client *osc.Client
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
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
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
	data.InternetServiceId = to.String(internetServiceId)
	data.Id = to.String(internetServiceId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = ComputedTagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.LinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	createResp, err := r.Client.LinkInternetService(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to link Internet Service resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)

	data.InternetServiceId = to.String(data.InternetServiceId.ValueString())
	data.Id = to.String(data.InternetServiceId.ValueString())

	data, err = setInternetServiceLinkState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Internet Service Link state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceInternetServiceLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setInternetServiceLinkState(ctx, r, data)
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
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	unlinkReq := osc.UnlinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	_, err := r.Client.UnlinkInternetService(ctx, unlinkReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		oscErr := oapihelpers.GetError(err)
		// 409 with code 1004 is returned when the Net of the Internet Service has mapped Public IPs
		// In this case, we unlink them before unlinking the Internet Service
		if oscErr.Code != "1004" {
			resp.Diagnostics.AddError(
				"Unable to unlink Internet Service resource.",
				err.Error(),
			)
			return
		}
	}

	var publicIps []osc.LinkPublicIp
	respNics, err := r.Client.ReadNics(ctx, osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{
			NetIds: &[]string{data.NetId.ValueString()},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Nics linked to the net of the Internet Service.",
			err.Error(),
		)
	}
	if len(ptr.From(respNics.Nics)) > 0 {
		for _, nic := range ptr.From(respNics.Nics) {
			if nic.LinkPublicIp == nil {
				continue
			}
			publicIps = append(publicIps, *nic.LinkPublicIp)
		}
	}
	for _, ip := range publicIps {
		_, err := r.Client.UnlinkPublicIp(ctx, osc.UnlinkPublicIpRequest{
			LinkPublicIpId: &ip.LinkPublicIpId,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to unlink Public IP from the NICs linked to the Net of the Internet Service.",
				err.Error(),
			)
		}
	}

	_, err = r.Client.UnlinkInternetService(ctx, unlinkReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		oscErr := oapihelpers.GetError(err)
		// 424 with code 1007 is returned when Internet Service is already unlinked
		// In this case, the link resource can be destroyed
		if oscErr.Code != "1007" {
			resp.Diagnostics.AddError(
				"Unable to unlink Internet Service resource after unmapping Public IPs.",
				err.Error(),
			)
			return
		}
	}
}

func setInternetServiceLinkState(ctx context.Context, r *resourceInternetServiceLink, data InternetServiceLinkModel) (InternetServiceLinkModel, error) {
	internetServiceFilters := osc.FiltersInternetService{
		InternetServiceIds: &[]string{data.InternetServiceId.ValueString()},
	}
	readReq := osc.ReadInternetServicesRequest{
		Filters: &internetServiceFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'internet service' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadInternetServices(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	if readResp.InternetServices == nil || len(*readResp.InternetServices) == 0 {
		return data, ErrResourceEmpty
	}

	internetService := (*readResp.InternetServices)[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, internetService.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	data.NetId = to.String(internetService.NetId)
	data.State = to.String(internetService.State)
	data.Id = to.String(internetService.InternetServiceId)
	data.InternetServiceId = to.String(internetService.InternetServiceId)

	return data, nil
}

func ResourceInternetServiceLinkStateRefreshFunc(ctx context.Context, r *resourceInternetServiceLink, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		readReq := osc.ReadInternetServicesRequest{Filters: &osc.FiltersInternetService{InternetServiceIds: &[]string{id}}}
		resp, err := r.Client.ReadInternetServices(ctx, readReq, options.WithRetryTimeout(ReadDefaultTimeout))
		if err != nil {
			return resp, "error", err
		}
		internetService := (*resp.InternetServices)[0]

		return resp, internetService.State, nil
	}
}
