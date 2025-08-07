package outscale

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
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
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
	Client *oscgo.APIClient
}

func NewResourceVolumeLink() resource.Resource {
	return &resourceVolumeLink{}
}

func (r *resourceVolumeLink) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.VolumeId = types.StringValue(volumeId)
	data.Id = types.StringValue(volumeId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volStateConf := &retry.StateChangeConf{
		Pending:      []string{"creating", "updating"},
		Target:       []string{"available"},
		Refresh:      getVolumeStateRefreshFunc(ctx, r.Client, createTimeout, data.VolumeId.ValueString()),
		Timeout:      createTimeout,
		Delay:        2 * time.Second,
		PollInterval: 4 * time.Second,
	}
	if _, err := volStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid volume state to be linked",
			fmt.Sprintf("Expected state 'available', got: '%s' ", err.Error()),
		)
		return
	}

	vmStateConf := &retry.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"running", "stopped"},
		Refresh:      getVmStateRefreshFunc(ctx, r.Client, createTimeout, data.VmId.ValueString()),
		Timeout:      createTimeout,
		Delay:        2 * time.Second,
		PollInterval: 3 * time.Second,
	}
	if _, err := vmStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Invalid vm state to be linked",
			fmt.Sprintf("Expected state 'running', got: '%s' ", err.Error()),
		)
		return
	}

	createReq := oscgo.NewLinkVolumeRequest(data.DeviceName.ValueString(), data.VmId.ValueString(), data.VolumeId.ValueString())
	var linkResp oscgo.LinkVolumeResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.VolumeApi.LinkVolume(ctx).LinkVolumeRequest(*createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		linkResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Link volume resource",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(*linkResp.ResponseContext.RequestId)

	linkedStateConf := &retry.StateChangeConf{
		Pending:    []string{"attaching"},
		Target:     []string{"attached"},
		Refresh:    getvolumeAttachmentStateRefreshFunc(ctx, r.Client, createTimeout, data.VmId.ValueString(), data.VolumeId.ValueString()),
		Timeout:    createTimeout,
		Delay:      3 * time.Second,
		MinTimeout: 4 * time.Second,
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVolumeLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setLinkedVolumeState(ctx, r, &data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVolumeLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VolumeLinkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.UnlinkVolumeRequest{
		VolumeId:    data.VolumeId.ValueString(),
		ForceUnlink: data.ForceUnlink.ValueBoolPointer(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.VolumeApi.UnlinkVolume(ctx).UnlinkVolumeRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink volume",
			err.Error(),
		)
		return
	}
	unLinkedStateConf := &retry.StateChangeConf{
		Pending:    []string{"detaching"},
		Target:     []string{"detached"},
		Refresh:    getvolumeAttachmentStateRefreshFunc(ctx, r.Client, deleteTimeout, data.VmId.ValueString(), data.VolumeId.ValueString()),
		Timeout:    deleteTimeout,
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
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

	volumeFilters := oscgo.FiltersVolume{
		VolumeIds: &[]string{data.VolumeId.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'volume_link' read timeout value. Error: %v: ", diags.Errors())
	}

	readReq := oscgo.ReadVolumesRequest{
		Filters: &volumeFilters,
	}
	var readResp oscgo.ReadVolumesResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.VolumeApi.ReadVolumes(ctx).ReadVolumesRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if len(readResp.GetVolumes()) == 0 || len(readResp.GetVolumes()[0].GetLinkedVolumes()) == 0 {
		return errors.New("Empty")
	}

	isForceUnlink := false
	linkedVol := readResp.GetVolumes()[0].GetLinkedVolumes()[0]
	if data.ForceUnlink.ValueBool() {
		isForceUnlink = data.ForceUnlink.ValueBool()
	}
	data.ForceUnlink = types.BoolValue(isForceUnlink)
	data.DeleteOnVmDeletion = types.BoolValue(linkedVol.GetDeleteOnVmDeletion())
	data.VolumeId = types.StringValue(linkedVol.GetVolumeId())
	data.DeviceName = types.StringValue(linkedVol.GetDeviceName())
	data.State = types.StringValue(linkedVol.GetState())
	data.VmId = types.StringValue(linkedVol.GetVmId())
	data.Id = types.StringValue(linkedVol.GetVolumeId())

	return nil
}

func getVmStateRefreshFunc(ctx context.Context, conn *oscgo.APIClient, timeout time.Duration, vmId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVmsStateResponse
		var err error
		req := oscgo.NewReadVmsStateRequest()
		filters := oscgo.FiltersVmsState{
			VmIds: &[]string{vmId},
		}
		req.SetAllVms(true)
		req.SetFilters(filters)
		err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			rp, httpResp, err := conn.VmApi.ReadVmsState(ctx).ReadVmsStateRequest(*req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return nil, "", err
		}
		vmStatus := resp.GetVmStates()[0]
		return vmStatus, vmStatus.GetVmState(), nil
	}
}
func getvolumeAttachmentStateRefreshFunc(ctx context.Context, conn *oscgo.APIClient, timeOut time.Duration, vmId string, volumeId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds:       &[]string{volumeId},
				LinkVolumeVmIds: &[]string{vmId},
			},
		}

		var resp oscgo.ReadVolumesResponse
		err := retry.RetryContext(ctx, timeOut, func() *retry.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return nil, "failed", err
		}

		if len(resp.GetVolumes()) > 0 && len(resp.GetVolumes()[0].GetLinkedVolumes()) > 0 {
			linkedVolume := resp.GetVolumes()[0].GetLinkedVolumes()[0]
			return linkedVolume, linkedVolume.GetState(), nil
		}

		return resp.GetVolumes(), "detached", nil
	}
}
