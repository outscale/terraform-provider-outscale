package oapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
)

var (
	_ resource.Resource                = &nicLinkResource{}
	_ resource.ResourceWithConfigure   = &nicLinkResource{}
	_ resource.ResourceWithImportState = &nicLinkResource{}
)

const (
	nicLinkErrCreate = "Unable to create NIC Link"
	nicLinkErrRead   = "Unable to read NIC Link"
	nicLinkErrDelete = "Unable to delete NIC Link"
)

type nicLinkModel struct {
	DeviceNumber       types.Int64    `tfsdk:"device_number"`
	VmId               types.String   `tfsdk:"vm_id"`
	NicId              types.String   `tfsdk:"nic_id"`
	State              types.String   `tfsdk:"state"`
	DeleteOnVmDeletion types.Bool     `tfsdk:"delete_on_vm_deletion"`
	VmAccountId        types.String   `tfsdk:"vm_account_id"`
	LinkNicId          types.String   `tfsdk:"link_nic_id"`
	Id                 types.String   `tfsdk:"id"`
	RequestId          types.String   `tfsdk:"request_id"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

type nicLinkResource struct {
	Client *osc.Client
}

func NewResourceNicLink() resource.Resource {
	return &nicLinkResource{}
}

func (r *nicLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *nicLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nic_link"
}

func (r *nicLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import NIC Link identifier. Got: %v", req.ID),
		)
		return
	}

	var data nicLinkModel
	var timeoutsVal timeouts.Value
	data.NicId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *nicLinkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"device_number": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nic_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"delete_on_vm_deletion": schema.BoolAttribute{
				Computed: true,
			},
			"vm_account_id": schema.StringAttribute{
				Computed: true,
			},
			"link_nic_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *nicLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data nicLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.LinkNicRequest{
		DeviceNumber: int(data.DeviceNumber.ValueInt64()),
		VmId:         data.VmId.ValueString(),
		NicId:        data.NicId.ValueString(),
	}

	createResp, err := r.Client.LinkNic(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(nicLinkErrCreate, err.Error())
		return
	}

	data.Id = to.String(createResp.LinkNicId)
	data.LinkNicId = to.String(createResp.LinkNicId)

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(nicLinkErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data nicLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	switch {
	case errors.Is(err, ErrResourceEmpty):
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(nicLinkErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *nicLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data nicLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	_, err := r.Client.UnlinkNic(ctx, osc.UnlinkNicRequest{
		LinkNicId: data.Id.ValueString(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(nicLinkErrDelete, err.Error())
		return
	}

	nicId := data.NicId.ValueString()
	stateConf := &stateconf.StateChangeConf[osc.LinkNicState]{
		Pending: stateconf.States(osc.LinkNicStateDetaching),
		Target:  stateconf.States(osc.LinkNicStateDetached),
		Timeout: timeout,
		Refresh: r.nicLinkRefreshFunc(nicId),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	switch {
	case errors.Is(err, ErrResourceEmpty):
	case err != nil:
		resp.Diagnostics.AddError(nicLinkErrDelete, err.Error())
	}
}

func (r *nicLinkResource) read(ctx context.Context, timeout time.Duration, data nicLinkModel) (nicLinkModel, error) {
	nicId := data.NicId.ValueString()

	stateConf := &stateconf.StateChangeConf[osc.LinkNicState]{
		Pending: stateconf.States(osc.LinkNicStateAttaching, osc.LinkNicStateDetaching),
		Target:  stateconf.States(osc.LinkNicStateAttached, osc.LinkNicStateDetached),
		Timeout: timeout,
		Refresh: r.nicLinkRefreshFunc(nicId),
	}

	respAny, err := stateConf.WaitForStateContext(ctx)
	switch {
	case errors.Is(err, ErrResourceEmpty):
		return data, err
	case err != nil:
		return data, fmt.Errorf("error waiting for NIC (%s): %w", nicId, err)
	}

	resp := respAny.(*osc.ReadNicsResponse)
	linkNic := (*resp.Nics)[0].LinkNic

	data.DeviceNumber = to.Int64(linkNic.DeviceNumber)
	data.VmId = to.String(linkNic.VmId)
	data.State = to.String(linkNic.State)
	data.DeleteOnVmDeletion = to.Bool(linkNic.DeleteOnVmDeletion)
	data.VmAccountId = to.String(linkNic.VmAccountId)
	data.LinkNicId = to.String(linkNic.LinkNicId)
	data.Id = to.String(linkNic.LinkNicId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return data, nil
}

func (r *nicLinkResource) nicLinkRefreshFunc(nicId string) func(context.Context) (any, osc.LinkNicState, error) {
	return func(ctx context.Context) (any, osc.LinkNicState, error) {
		resp, err := r.Client.ReadNics(ctx, osc.ReadNicsRequest{
			Filters: &osc.FiltersNic{NicIds: &[]string{nicId}},
		})
		if err != nil {
			return nil, "", err
		}
		if resp.Nics == nil || len(*resp.Nics) == 0 || (*resp.Nics)[0].LinkNic == nil || (*resp.Nics)[0].LinkNic.LinkNicId == "" {
			return nil, "", ErrResourceEmpty
		}

		linkNic := ptr.From((*resp.Nics)[0].LinkNic)
		return resp, linkNic.State, nil
	}
}

// nicLinkRefreshFunc is a standalone SDKv2-compatible refresh function used by
// resource_outscale_nic.go which is still on SDKv2.
func nicLinkRefreshFunc(ctx context.Context, client *osc.Client, nicID string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (any, string, error) {
		req := osc.ReadNicsRequest{
			Filters: &osc.FiltersNic{
				NicIds: &[]string{nicID},
			},
		}

		resp, err := client.ReadNics(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "failed", err
		}
		if resp.Nics == nil || len(*resp.Nics) < 1 {
			return nil, "failed", fmt.Errorf("error to find the nic(%s): %#v", nicID, resp.Nics)
		}

		linkNic := ptr.From((*resp.Nics)[0].LinkNic)
		if reflect.DeepEqual(linkNic, osc.LinkNic{}) {
			return resp, "detached", nil
		}

		return resp, string(linkNic.State), nil
	}
}
