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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
)

var (
	_ resource.Resource                = &resourceInternetServiceLink{}
	_ resource.ResourceWithConfigure   = &resourceInternetServiceLink{}
	_ resource.ResourceWithImportState = &resourceInternetServiceLink{}
	_ resource.ResourceWithModifyPlan  = &resourceInternetServiceLink{}
)

const (
	internetServiceLinkErrCreate    = "Unable to link Internet Service"
	internetServiceLinkErrDelete    = "Unable to unlink Internet Service"
	internetServiceLinkErrReadNics  = "Unable to read NICs linked to Internet Service net"
	internetServiceLinkErrUnmapIPs  = "Unable to unlink Public IPs from Internet Service net NICs"
	internetServiceLinkErrRetryLink = "Unable to unlink Internet Service after unmapping Public IPs"
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

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.LinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	createResp, err := r.Client.LinkInternetService(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(internetServiceLinkErrCreate, err.Error())
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)

	data.InternetServiceId = to.String(data.InternetServiceId.ValueString())
	data.Id = to.String(data.InternetServiceId.ValueString())
	// The API response does not contain enough information to set the state directly, which would cause an error.
	// The next read will fill the state

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceInternetServiceLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceInternetServiceLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceInternetServiceLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InternetServiceLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	unlinkReq := osc.UnlinkInternetServiceRequest{
		InternetServiceId: data.InternetServiceId.ValueString(),
		NetId:             data.NetId.ValueString(),
	}

	var hasIPs, hasLBU bool
	_, err := r.Client.UnlinkInternetService(ctx, unlinkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		oscErr := oapihelpers.GetError(err)
		switch oscErr.Code {
		case "1004":
			// 409 with code 1004 is returned when the Net of the Internet Service has mapped Public IPs
			// In this case, we unlink them before unlinking the Internet Service
			hasIPs = true
		case "1005":
			// 409 with code 1005 is returned when the Net of the Internet Service has a Load Balancer
			hasLBU = true
		default:
			resp.Diagnostics.AddError(internetServiceLinkErrDelete, err.Error())
			return
		}
	}

	if hasIPs {
		var publicIps []osc.LinkPublicIp
		respNics, err := r.Client.ReadNics(ctx, osc.ReadNicsRequest{
			Filters: &osc.FiltersNic{
				NetIds: &[]string{data.NetId.ValueString()},
			},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(internetServiceLinkErrReadNics, err.Error())
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
			}, options.WithRetryTimeout(timeout))
			if err != nil {
				switch {
				case osc.HasErrorCode(err, []string{"5026"}):
					// 400 with code 5026 is returned when Public IP is already destroyed
					// In this case, we can skip it
				default:
					resp.Diagnostics.AddError(internetServiceLinkErrUnmapIPs, err.Error())
				}
			}
		}

		_, err = r.Client.UnlinkInternetService(ctx, unlinkReq, options.WithRetryTimeout(timeout))
		if err != nil {
			oscErr := oapihelpers.GetError(err)
			// 424 with code 1007 is returned when Internet Service is already unlinked
			// In this case, the link resource can be destroyed
			if oscErr.Code != "1007" {
				resp.Diagnostics.AddError(internetServiceLinkErrRetryLink, err.Error())
				return
			}
		}
	}

	if hasLBU {
		// We retry on the 1005 error rather than reading the LBU state to decide.
		// During a terraform destroy, resources are deleted in parallel, so the LBU on this Net
		// can be in any transient state (reloading, reconfiguring, deleting, etc.).
		// Checking only for specific states would miss cases like a concurrent
		// backend vms unlink putting the LBU in "reloading" state (like in TF-2 integration test)
		_, err := oapihelpers.RetryOnCodes(ctx, []string{"1005"}, func() (resp any, err error) {
			return r.Client.UnlinkInternetService(ctx, unlinkReq, options.WithRetryTimeout(timeout))
		}, timeout)
		if err != nil {
			resp.Diagnostics.AddError(internetServiceLinkErrDelete, err.Error())
		}
	}
}

func (r *resourceInternetServiceLink) read(ctx context.Context, timeout time.Duration, data InternetServiceLinkModel) (InternetServiceLinkModel, error) {
	internetServiceFilters := osc.FiltersInternetService{
		InternetServiceIds: &[]string{data.InternetServiceId.ValueString()},
	}
	readReq := osc.ReadInternetServicesRequest{
		Filters: &internetServiceFilters,
	}

	readResp, err := r.Client.ReadInternetServices(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if readResp.InternetServices == nil || len(*readResp.InternetServices) == 0 {
		return data, ErrResourceEmpty
	}

	internetService := (*readResp.InternetServices)[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, internetService.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
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
