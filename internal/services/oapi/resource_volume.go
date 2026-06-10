package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceVolume{}
	_ resource.ResourceWithConfigure   = &resourceVolume{}
	_ resource.ResourceWithImportState = &resourceVolume{}
	_ resource.ResourceWithModifyPlan  = &resourceVolume{}
)

const (
	volumeErrCreate   = "Unable to create Volume"
	volumeErrUpdate   = "Unable to update Volume"
	volumeErrDelete   = "Unable to delete Volume"
	volumeErrWait     = "Unable to wait for Volume state"
	volumeErrTask     = "Unable to wait for Volume update task"
	volumeErrSnapshot = "Unable to create snapshot during Volume deletion"
	volumeErrTags     = "Unable to create snapshot tags during Volume deletion"
)

type VolumeModel struct {
	TerminationSnapshotName types.String   `tfsdk:"termination_snapshot_name"`
	LinkedVolumes           types.Set      `tfsdk:"linked_volumes"`
	SubregionName           types.String   `tfsdk:"subregion_name"`
	CreationDate            types.String   `tfsdk:"creation_date"`
	SnapshotId              types.String   `tfsdk:"snapshot_id"`
	VolumeType              types.String   `tfsdk:"volume_type"`
	RequestId               types.String   `tfsdk:"request_id"`
	VolumeId                types.String   `tfsdk:"volume_id"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
	State                   types.String   `tfsdk:"state"`
	Iops                    types.Int32    `tfsdk:"iops"`
	Size                    types.Int32    `tfsdk:"size"`
	Id                      types.String   `tfsdk:"id"`
	TagsModel
}

type BlockLinkedVolumes struct {
	DeleteOnVmDeletion types.Bool   `tfsdk:"delete_on_vm_deletion"`
	VolumeId           types.String `tfsdk:"volume_id"`
	DeviceName         types.String `tfsdk:"device_name"`
	State              types.String `tfsdk:"state"`
	VmId               types.String `tfsdk:"vm_id"`
}

type volumeCommon struct {
	Client  *osc.Client
	Batcher *batch.BatcherByID[osc.Volume]
}

type resourceVolume struct {
	volumeCommon
}

func NewResourceVolume() resource.Resource {
	return &resourceVolume{}
}

func (r *resourceVolume) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceVolume) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
		return
	}

	var data VolumeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Size.IsNull() && data.SnapshotId.IsNull() {
		resp.Diagnostics.AddError(
			"Setting 'size' and 'snapshot_id' Considerations",
			"Volume 'size' parameter is required if the volume is not created from a snapshot (snapshot_id unspecified)",
		)
		return
	}

	if data.VolumeType.ValueString() != "io1" && !data.Iops.IsNull() && !data.Iops.IsUnknown() {
		resp.Diagnostics.AddError(
			"Setting IOPS Considerations",
			ErrResourceInvalidIOPS.Error(),
		)
		return
	}
	if data.VolumeType.ValueString() == "io1" {
		if data.Iops.IsUnknown() || data.Iops.IsNull() || data.Iops.ValueInt32() < MinIops || data.Iops.ValueInt32() > MaxIops {
			resp.Diagnostics.AddError(
				"Setting IOPS Considerations",
				fmt.Sprintf("iops parameter is required for 'io1' volume and must be between %d and %d inclusive, got: %d", MinIops, MaxIops, data.Iops.ValueInt32()),
			)
			return
		}
	}
}

func (r *resourceVolume) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	volumeId := req.ID
	if volumeId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import volume identifier Got: %v", req.ID),
		)
		return
	}

	var data VolumeModel
	var timeouts timeouts.Value
	data.VolumeId = to.String(volumeId)
	data.Id = to.String(volumeId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.LinkedVolumes = types.SetNull(types.ObjectType{AttrTypes: fwhelpers.GetAttrTypes(BlockLinkedVolumes{})})
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceVolume) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func LinkedVolumesSchema() schema.SetAttribute {
	return schema.SetAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"delete_on_vm_deletion": types.BoolType,
				"volume_id":             types.StringType,
				"device_name":           types.StringType,
				"state":                 types.StringType,
				"vm_id":                 types.StringType,
			},
		},
	}
}

func (r *resourceVolume) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"linked_volumes": LinkedVolumesSchema(),
			"termination_snapshot_name": schema.StringAttribute{
				Optional: true,
			},
			"subregion_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("standard"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"gp2", "io1", "standard"}...),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"iops": schema.Int32Attribute{
				Optional: true,
				Computed: true,
			},
			"size": schema.Int32Attribute{
				Optional: true,
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"volume_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceVolume) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateVolumeRequest{
		SubregionName: data.SubregionName.ValueString(),
	}
	createReq.VolumeType = new(osc.VolumeType(data.VolumeType.ValueString()))
	if !data.Size.IsUnknown() && !data.Size.IsNull() {
		createReq.Size = new(int(data.Size.ValueInt32()))
	}
	if !data.Iops.IsUnknown() && !data.Iops.IsNull() {
		createReq.Iops = new(int(data.Iops.ValueInt32()))
	}
	if !data.SnapshotId.IsUnknown() && !data.SnapshotId.IsNull() {
		createReq.SnapshotId = data.SnapshotId.ValueStringPointer()
	}

	createResp, err := r.Client.CreateVolume(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(volumeErrCreate, err.Error())
		return
	}
	volumeId := createResp.Volume.VolumeId
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.VolumeId = to.String(volumeId)
	data.Id = to.String(volumeId)

	stateData, err := r.flatten(ctx, data, *createResp.Volume)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, volumeId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.VolumeState]{
		Pending: stateconf.States(osc.VolumeStateCreating),
		Target:  stateconf.States(osc.VolumeStateAvailable),
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(volumeId),
	}
	volumeAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(volumeErrWait,
			fmt.Sprintf("Unexpected volume (%s) state: '%s' ", volumeId, err.Error()),
		)
		return
	}

	// We set the last read response to the state
	volume := volumeAny.(osc.Volume)
	stateData, err = r.flatten(ctx, data, volume)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceVolume) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeModel

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

func (r *resourceVolume) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData VolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	volumeId := stateData.VolumeId.ValueString()
	stateData.TerminationSnapshotName = planData.TerminationSnapshotName

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, volumeId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	updateReq := osc.UpdateVolumeRequest{
		VolumeId: volumeId,
	}
	shouldUpdate := false

	if fwhelpers.IsSet(planData.Size) && !planData.Size.Equal(stateData.Size) {
		updateReq.Size = new(int(planData.Size.ValueInt32()))
		shouldUpdate = true
	}

	// When updating a gp2 volume to io1 while keeping the iops value at 100 (default gp2 iops value),
	// Terraform will not detect the change between config and state
	if fwhelpers.IsSet(planData.Iops) && (!planData.Iops.Equal(stateData.Iops) || (planData.VolumeType.ValueString() == "io1" && stateData.VolumeType.ValueString() != "io1")) {
		updateReq.Iops = new(int(planData.Iops.ValueInt32()))
		shouldUpdate = true
	}

	if fwhelpers.IsSet(planData.VolumeType) && !planData.VolumeType.Equal(stateData.VolumeType) {
		updateReq.VolumeType = new(osc.VolumeType(planData.VolumeType.ValueString()))
		shouldUpdate = true
	}

	if shouldUpdate {
		volume, err := r.Client.UpdateVolume(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(volumeErrUpdate, err.Error())
			return
		}

		if volume.Volume.TaskId != nil {
			err := WaitForVolumeTasks(ctx, timeout, []string{*volume.Volume.TaskId}, r.Client)
			if err != nil {
				resp.Diagnostics.AddError(volumeErrTask, err.Error())
				return
			}
		}
	}

	stateData.Timeouts = planData.Timeouts
	newData, err := r.read(ctx, timeout, stateData)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func WaitForVolumeTasks(ctx context.Context, timeout time.Duration, tasksIds []string, client *osc.Client) error {
	failed := []string{"failed", "canceled"}
	running := []string{"pending", "active"}
	stateConf := &retry.StateChangeConf{
		Pending: []string{"running"},
		Target:  []string{"finished"},
		Timeout: timeout,
		Refresh: func() (any, string, error) {
			req := osc.ReadVolumeUpdateTasksRequest{
				Filters: &osc.FiltersReadVolumeUpdateTask{
					TaskIds: &tasksIds,
				},
			}
			resp, err := client.ReadVolumeUpdateTasks(ctx, req, options.WithRetryTimeout(timeout))
			if err != nil {
				return nil, "", err
			}
			if resp.VolumeUpdateTasks == nil || len(*resp.VolumeUpdateTasks) == 0 {
				return nil, "", fmt.Errorf("tasks %v not found", tasksIds)
			}

			if lo.ContainsBy(*resp.VolumeUpdateTasks, func(task osc.VolumeUpdateTask) bool {
				return lo.Contains(running, *task.State)
			}) {
				return resp, "running", nil
			}

			failedTasks := lo.FilterMap(*resp.VolumeUpdateTasks, func(task osc.VolumeUpdateTask, _ int) (error, bool) {
				return fmt.Errorf("task (%s) did not complete and ended with state: %s - comment: %s", ptr.From(task.TaskId), ptr.From(task.State), ptr.From(task.Comment)), lo.Contains(failed, *task.State)
			})
			if len(failedTasks) > 0 {
				return resp, "failed", fmt.Errorf("volume update tasks failed: %w", errors.Join(failedTasks...))
			}
			return resp, "finished", nil
		},
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func (r *resourceVolume) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	volumeId := data.VolumeId.ValueString()
	if !data.TerminationSnapshotName.IsNull() {
		description := "created before volume deletion"
		var snapshotId string
		request := osc.CreateSnapshotRequest{
			Description: &description,
			VolumeId:    &volumeId,
		}
		_, err := r.Client.CreateSnapshot(ctx, request, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(volumeErrSnapshot, err.Error())
			return
		}

		tags := osc.ResourceTag{
			Key:   "Name",
			Value: data.TerminationSnapshotName.String(),
		}
		err = createOAPITags(ctx, r.Client, timeout, []osc.ResourceTag{tags}, snapshotId)
		if err != nil {
			resp.Diagnostics.AddError(volumeErrTags, err.Error())
			return
		}
	}

	delReq := osc.DeleteVolumeRequest{
		VolumeId: data.VolumeId.ValueString(),
	}
	_, err := r.Client.DeleteVolume(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			volumeErrDelete,
			err.Error(),
		)
	}
}

func (r *resourceVolume) read(ctx context.Context, timeout time.Duration, data VolumeModel) (VolumeModel, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	volume, err := r.Batcher.Read(ctxWithTimeout, data.VolumeId.ValueString())
	if err != nil {
		if errors.Is(err, batch.ErrNotFound) {
			return data, ErrResourceEmpty
		}
		return data, err
	}

	return r.flatten(ctx, data, *volume)
}

func (r *resourceVolume) flatten(ctx context.Context, data VolumeModel, volume osc.Volume) (VolumeModel, error) {
	tags, diag := flattenOAPITagsFW(ctx, volume.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.Tags = tags

	data.LinkedVolumes, diag = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: fwhelpers.GetAttrTypes(BlockLinkedVolumes{})}, getLinkedVolumesFromApiResponse(volume.LinkedVolumes))
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	data.CreationDate = to.String(from.ISO8601(volume.CreationDate))
	if data.TerminationSnapshotName.IsNull() {
		data.TerminationSnapshotName = types.StringNull()
	}
	data.SubregionName = to.String(volume.SubregionName)
	data.VolumeType = to.String(volume.VolumeType)
	data.VolumeId = to.String(volume.VolumeId)
	data.SnapshotId = to.String(ptr.From(volume.SnapshotId))
	data.State = to.String(volume.State)
	if data.VolumeType.ValueString() != string(osc.VolumeTypeStandard) {
		data.Iops = to.Int32(int32(volume.Iops)) //nolint:gosec
	} else {
		data.Iops = types.Int32Null()
	}
	data.Size = to.Int32(int32(volume.Size)) //nolint:gosec
	data.Id = to.String(volume.VolumeId)

	return data, nil
}

func (r *volumeCommon) stateRefreshFunc(volumeID string) stateconf.StateRefreshFunc[osc.VolumeState] {
	return func(ctx context.Context) (any, osc.VolumeState, error) {
		resp, err := r.Client.ReadVolumes(ctx, osc.ReadVolumesRequest{
			Filters: &osc.FiltersVolume{
				VolumeIds: &[]string{volumeID},
			},
		})
		if err != nil {
			return nil, "", err
		}
		v := ptr.From(resp.Volumes)[0]
		return v, v.State, nil
	}
}

func getLinkedVolumesFromApiResponse(linkedVols []osc.LinkedVolume) []BlockLinkedVolumes {
	linkedVolumes := make([]BlockLinkedVolumes, 0, len(linkedVols))

	for _, linkedVol := range linkedVols {
		linkedVolume := BlockLinkedVolumes{}
		linkedVolume.DeleteOnVmDeletion = to.Bool(linkedVol.DeleteOnVmDeletion)
		linkedVolume.DeviceName = to.String(linkedVol.DeviceName)
		linkedVolume.VmId = to.String(linkedVol.VmId)
		linkedVolume.State = to.String(linkedVol.State)
		linkedVolume.VolumeId = to.String(linkedVol.VolumeId)
		linkedVolumes = append(linkedVolumes, linkedVolume)
	}
	return linkedVolumes
}
