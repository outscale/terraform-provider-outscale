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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                = &resourceVolumeLink{}
	_ resource.ResourceWithConfigure   = &resourceVolumeLink{}
	_ resource.ResourceWithModifyPlan  = &resourceVolumeLink{}
	_ resource.ResourceWithImportState = &resourceVolumeLink{}
)

type VolumeLinkModel struct {
	DeleteOnVmDeletion types.Bool     `tfsdk:"delete_on_vm_deletion"`
	ForceUnlink        types.Bool     `tfsdk:"force_unlink"`
	DeviceName         types.String   `tfsdk:"device_name"`
	RequestId          types.String   `tfsdk:"request_id"`
	VolumeId           types.String   `tfsdk:"volume_id"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	State              types.String   `tfsdk:"state"`
	VmId               types.String   `tfsdk:"vm_id"`
	Id                 types.String   `tfsdk:"id"`
}
type resourceVolumeLink struct {
	Client *osc.Client
}

func NewResourceVolumeLink() resource.Resource {
	return &resourceVolumeLink{}
}

func (r *resourceVolumeLink) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceVolumeLink) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsKnown() && req.State.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource 'outscale_volume_link' Considerations",
			"You have to apply twice or run 'terraform refesh' after the fist apply to get"+
				" the 'linked_volumes' block values in volume resource state",
		)
	}
}

func (r *resourceVolumeLink) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	volumeId := req.ID
	if volumeId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import keypair identifier Got: %v", req.ID),
		)
		return
	}

	var data VolumeLinkModel
	var timeouts timeouts.Value
	data.VolumeId = to.String(volumeId)
	data.Id = to.String(volumeId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceVolumeLink) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_link"
}

func (r *resourceVolumeLink) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"volume_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force_unlink": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"delete_on_vm_deletion": schema.BoolAttribute{
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

func (r *resourceVolumeLink) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VolumeLinkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	volStateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "updating"},
		Target:  []string{"available"},
		Timeout: createTimeout,
		Refresh: getVolumeStateRefreshFunc(ctx, r.Client, createTimeout, data.VolumeId.ValueString()),
	}
	if _, err := volStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid volume state to be linked",
			fmt.Sprintf("Expected state 'available', got: '%s' ", err.Error()),
		)
		return
	}

	vmStateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"running", "stopped"},
		Timeout: createTimeout,
		Refresh: getVmStateRefreshFunc(ctx, r.Client, createTimeout, data.VmId.ValueString()),
	}
	if _, err := vmStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid vm state to be linked",
			fmt.Sprintf("Expected state 'running', got: '%s' ", err.Error()),
		)
		return
	}

	createReq := osc.LinkVolumeRequest{
		DeviceName: data.DeviceName.ValueString(),
		VmId:       data.VmId.ValueString(),
		VolumeId:   data.VolumeId.ValueString(),
	}
	linkResp, err := r.Client.LinkVolume(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Link volume resource",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(*linkResp.ResponseContext.RequestId)

	linkedStateConf := &retry.StateChangeConf{
		Pending: []string{"attaching"},
		Target:  []string{"attached"},
		Timeout: createTimeout,
		Refresh: getvolumeAttachmentStateRefreshFunc(ctx, r.Client, createTimeout, data.VmId.ValueString(), data.VolumeId.ValueString()),
	}
	if _, err = linkedStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid linked volume state",
			fmt.Sprintf("Expected state 'attached', got: '%s' ", err.Error()),
		)
		return
	}

	err = setLinkedVolumeState(ctx, r, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set volume state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceVolumeLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setLinkedVolumeState(ctx, r, &data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set volume API response values.",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceVolumeLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataPlan, dataState VolumeLinkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dataState.ForceUnlink = dataPlan.ForceUnlink

	err := setLinkedVolumeState(ctx, r, &dataState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set volume state after link volume updating.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataState)...)
}

func (r *resourceVolumeLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VolumeLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.UnlinkVolumeRequest{
		VolumeId:    data.VolumeId.ValueString(),
		ForceUnlink: data.ForceUnlink.ValueBoolPointer(),
	}

	_, err := r.Client.UnlinkVolume(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink volume",
			err.Error(),
		)
		return
	}
	unLinkedStateConf := &retry.StateChangeConf{
		Pending: []string{"detaching"},
		Target:  []string{"detached"},
		Timeout: deleteTimeout,
		Refresh: getvolumeAttachmentStateRefreshFunc(ctx, r.Client, deleteTimeout, data.VmId.ValueString(), data.VolumeId.ValueString()),
	}
	if _, err = unLinkedStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid unlinked volume state",
			fmt.Sprintf("Expected state 'detached', got: '%s' ", err.Error()),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func setLinkedVolumeState(ctx context.Context, r *resourceVolumeLink, data *VolumeLinkModel) error {
	volumeFilters := osc.FiltersVolume{
		VolumeIds: &[]string{data.VolumeId.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'volume_link' read timeout value: %v", diags.Errors())
	}

	readReq := osc.ReadVolumesRequest{
		Filters: &volumeFilters,
	}
	readResp, err := r.Client.ReadVolumes(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return err
	}
	if readResp.Volumes == nil || len(*readResp.Volumes) == 0 || len((*readResp.Volumes)[0].LinkedVolumes) == 0 {
		return ErrResourceEmpty
	}

	isForceUnlink := false
	linkedVol := (*readResp.Volumes)[0].LinkedVolumes[0]
	if data.ForceUnlink.ValueBool() {
		isForceUnlink = data.ForceUnlink.ValueBool()
	}
	data.ForceUnlink = to.Bool(isForceUnlink)
	data.DeleteOnVmDeletion = to.Bool(linkedVol.DeleteOnVmDeletion)
	data.VolumeId = to.String(linkedVol.VolumeId)
	data.DeviceName = to.String(linkedVol.DeviceName)
	data.State = to.String(linkedVol.State)
	data.VmId = to.String(linkedVol.VmId)
	data.Id = to.String(linkedVol.VolumeId)

	return nil
}

func getVmStateRefreshFunc(ctx context.Context, client *osc.Client, timeout time.Duration, vmId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := osc.ReadVmsStateRequest{
			AllVms: new(true),
			Filters: &osc.FiltersVmsState{
				VmIds: &[]string{vmId},
			},
		}
		resp, err := client.ReadVmsState(ctx, req)
		if err != nil {
			return nil, "", err
		}
		vmStatus := (*resp.VmStates)[0]
		return vmStatus, string(vmStatus.VmState), nil
	}
}

func getvolumeAttachmentStateRefreshFunc(ctx context.Context, client *osc.Client, timeOut time.Duration, vmId string, volumeId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := osc.ReadVolumesRequest{
			Filters: &osc.FiltersVolume{
				VolumeIds:       &[]string{volumeId},
				LinkVolumeVmIds: &[]string{vmId},
			},
		}

		resp, err := client.ReadVolumes(ctx, request)
		if err != nil {
			return nil, "failed", err
		}

		if len(*resp.Volumes) > 0 && len((*resp.Volumes)[0].LinkedVolumes) > 0 {
			linkedVolume := (*resp.Volumes)[0].LinkedVolumes[0]
			return linkedVolume, string(linkedVolume.State), nil
		}

		return resp.Volumes, "detached", nil
	}
}
