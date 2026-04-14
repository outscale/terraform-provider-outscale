package oapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
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
	volumeCommon
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
	r.Batcher = client.VolumeBatcher
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

	volStateConf := &stateconf.StateChangeConf[osc.VolumeState]{
		Pending: stateconf.States(osc.VolumeStateCreating),
		Target:  stateconf.States(osc.VolumeStateAvailable),
		Timeout: createTimeout,
		Refresh: r.stateRefreshFunc(data.VolumeId.ValueString()),
	}
	if _, err := volStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid volume state to be linked",
			fmt.Sprintf("Expected state 'available', got: '%s' ", err.Error()),
		)
		return
	}

	vmStateConf := &stateconf.StateChangeConf[osc.VmState]{
		Pending: stateconf.States(osc.VmStatePending),
		Target:  stateconf.States(osc.VmStateRunning, osc.VmStateStopped),
		Timeout: createTimeout,
		Refresh: r.vmStateRefreshFunc(data.VmId.ValueString()),
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

	linkedStateConf := &stateconf.StateChangeConf[osc.LinkedVolumeState]{
		Pending: stateconf.States(osc.LinkedVolumeStateAttaching),
		Target:  stateconf.States(osc.LinkedVolumeStateAttached),
		Timeout: createTimeout,
		Refresh: getvolumeAttachmentStateRefreshFunc(r.Client, data.VmId.ValueString(), data.VolumeId.ValueString()),
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
	unLinkedStateConf := &stateconf.StateChangeConf[osc.LinkedVolumeState]{
		Pending: stateconf.States(osc.LinkedVolumeStateDetaching),
		Target:  stateconf.States(osc.LinkedVolumeStateDetached),
		Timeout: deleteTimeout,
		Refresh: getvolumeAttachmentStateRefreshFunc(r.Client, data.VmId.ValueString(), data.VolumeId.ValueString()),
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
	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'volume_link' read timeout value: %v", diags.Errors())
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	vol, err := r.Batcher.Read(ctxWithTimeout, data.VolumeId.ValueString())
	if err != nil {
		if errors.Is(err, batch.ErrNotFound) {
			return ErrResourceEmpty
		}
		return err
	}

	isForceUnlink := false
	linkedVol := vol.LinkedVolumes[0]
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

func (r *resourceVolumeLink) vmStateRefreshFunc(vmId string) stateconf.StateRefreshFunc[osc.VmState] {
	return func(ctx context.Context) (any, osc.VmState, error) {
		req := osc.ReadVmsStateRequest{
			AllVms: new(true),
			Filters: &osc.FiltersVmsState{
				VmIds: &[]string{vmId},
			},
		}
		resp, err := r.Client.ReadVmsState(ctx, req)
		if err != nil {
			return nil, "", err
		}
		vmStatus := (*resp.VmStates)[0]
		return vmStatus, vmStatus.VmState, nil
	}
}

func getvolumeAttachmentStateRefreshFunc(client *osc.Client, vmId string, volumeId string) stateconf.StateRefreshFunc[osc.LinkedVolumeState] {
	return func(ctx context.Context) (any, osc.LinkedVolumeState, error) {
		request := osc.ReadVolumesRequest{
			Filters: &osc.FiltersVolume{
				VolumeIds:       &[]string{volumeId},
				LinkVolumeVmIds: &[]string{vmId},
			},
		}

		resp, err := client.ReadVolumes(ctx, request)
		if err != nil {
			return nil, "", err
		}

		if len(ptr.From(resp.Volumes)) > 0 && len((*resp.Volumes)[0].LinkedVolumes) > 0 {
			linkedVolume := (*resp.Volumes)[0].LinkedVolumes[0]
			return linkedVolume, linkedVolume.State, nil
		}

		return resp.Volumes, osc.LinkedVolumeStateDetached, nil
	}
}
