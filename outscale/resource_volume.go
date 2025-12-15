package outscale

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
	"github.com/outscale/osc-sdk-go/v2"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/framework/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceVolume{}
	_ resource.ResourceWithConfigure   = &resourceVolume{}
	_ resource.ResourceWithImportState = &resourceVolume{}
	_ resource.ResourceWithModifyPlan  = &resourceVolume{}
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

type resourceVolume struct {
	Client *osc.APIClient
}

func NewResourceVolume() resource.Resource {
	return &resourceVolume{}
}

func (r *resourceVolume) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
			utils.VolumeIOPSError,
		)
		return
	}
	if data.VolumeType.ValueString() == "io1" {
		if data.Iops.IsUnknown() || data.Iops.IsNull() || data.Iops.ValueInt32() < utils.MinIops || data.Iops.ValueInt32() > utils.MaxIops {
			resp.Diagnostics.AddError(
				"Setting IOPS Considerations",
				fmt.Sprintf("iops parameter is required for 'io1' volume and must be between %d and %d inclusive, got: %d", utils.MinIops, utils.MaxIops, data.Iops.ValueInt32()),
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
	data.VolumeId = types.StringValue(volumeId)
	data.Id = types.StringValue(volumeId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.LinkedVolumes = types.SetNull(types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockLinkedVolumes{})})
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
					fwmodifyplan.ForceNewFramework(),
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

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := osc.NewCreateVolumeRequest(data.SubregionName.ValueString())
	createReq.SetVolumeType(data.VolumeType.ValueString())
	if !data.Size.IsUnknown() && !data.Size.IsNull() {
		createReq.SetSize(data.Size.ValueInt32())
	}
	if !data.Iops.IsUnknown() && !data.Iops.IsNull() {
		createReq.SetIops(data.Iops.ValueInt32())
	}
	if !data.SnapshotId.IsUnknown() && !data.SnapshotId.IsNull() {
		createReq.SetSnapshotId(data.SnapshotId.ValueString())
	}

	var createResp osc.CreateVolumeResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.VolumeApi.CreateVolume(ctx).CreateVolumeRequest(*createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create volume resource",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(*createResp.ResponseContext.RequestId)
	volumeId := createResp.Volume.GetVolumeId()
	data.VolumeId = types.StringValue(volumeId)
	data.Id = types.StringValue(volumeId)

	diag := createOAPITagsFW(ctx, r.Client, data.Tags, volumeId)
	if utils.CheckDiags(resp, diag) {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    getVolumeStateRefreshFunc(ctx, r.Client, createTimeout, volumeId),
		Timeout:    createTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Waiting for volume state to become available",
			fmt.Sprintf("Unexpected volume (%s) state: '%s' ", volumeId, err.Error()),
		)
		return
	}

	err = setVolumeState(ctx, r, &data)
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

func (r *resourceVolume) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setVolumeState(ctx, r, &data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set volume API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVolume) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData VolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := planData.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeId := stateData.VolumeId.ValueString()
	stateData.TerminationSnapshotName = planData.TerminationSnapshotName

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, volumeId)
	if utils.CheckDiags(resp, diag) {
		return
	}

	updateReq := osc.NewUpdateVolumeRequest(volumeId)
	shouldUpdate := false

	if utils.IsSet(planData.Size) && !planData.Size.Equal(stateData.Size) {
		updateReq.Size = planData.Size.ValueInt32Pointer()
		shouldUpdate = true
	}

	// When updating a gp2 volume to io1 while keeping the iops value at 100 (default gp2 iops value),
	// Terraform will not detect the change between config and state
	if utils.IsSet(planData.Iops) && (!planData.Iops.Equal(stateData.Iops) || (planData.VolumeType.ValueString() == "io1" && stateData.VolumeType.ValueString() != "io1")) {
		updateReq.Iops = planData.Iops.ValueInt32Pointer()
		shouldUpdate = true
	}

	if utils.IsSet(planData.VolumeType) && !planData.VolumeType.Equal(stateData.VolumeType) {
		updateReq.VolumeType = planData.VolumeType.ValueStringPointer()
		shouldUpdate = true
	}

	if shouldUpdate {
		var volume osc.UpdateVolumeResponse
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.VolumeApi.UpdateVolume(ctx).UpdateVolumeRequest(*updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			volume = rp
			stateData.RequestId = types.StringValue(*rp.ResponseContext.RequestId)
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update volume resource",
				err.Error(),
			)
			return
		}

		if volume.GetVolume().TaskId.IsSet() {
			err := WaitForVolumeTasks(ctx, updateTimeout, []string{*volume.GetVolume().TaskId.Get()}, r.Client)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error waiting for volume update task to complete",
					err.Error(),
				)
				return
			}
		}
	}

	err := setVolumeState(ctx, r, &stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set volume state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func WaitForVolumeTasks(ctx context.Context, timeout time.Duration, tasksIds []string, client *osc.APIClient) error {
	failed := []string{"failed", "canceled"}
	running := []string{"pending", "active"}
	stateConf := &retry.StateChangeConf{
		Pending: []string{"running"},
		Target:  []string{"finished"},
		Refresh: func() (any, string, error) {
			var resp osc.ReadVolumeUpdateTasksResponse

			err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				rp, httpResp, err := client.VolumeApi.ReadVolumeUpdateTasks(ctx).ReadVolumeUpdateTasksRequest(osc.ReadVolumeUpdateTasksRequest{
					Filters: &osc.FiltersUpdateVolumeTask{
						TaskIds: &tasksIds,
					},
				}).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				resp = rp
				return nil
			})
			if err != nil {
				return nil, "", err
			}
			if len(resp.GetVolumeUpdateTasks()) == 0 {
				return nil, "", fmt.Errorf("tasks %v not found", tasksIds)
			}

			if lo.ContainsBy(resp.GetVolumeUpdateTasks(), func(task osc.VolumeUpdateTask) bool {
				return lo.Contains(running, task.GetState())
			}) {
				return resp, "running", nil
			}

			failedTasks := lo.FilterMap(resp.GetVolumeUpdateTasks(), func(task osc.VolumeUpdateTask, _ int) (error, bool) {
				return fmt.Errorf("task (%s) did not complete and ended with state: %s. Comment: %s", task.GetTaskId(), task.GetState(), task.GetComment()), lo.Contains(failed, task.GetState())
			})
			if len(failedTasks) > 0 {
				return resp, "failed", fmt.Errorf("volume update tasks failed: %w", errors.Join(failedTasks...))
			}
			return resp, "finished", nil
		},
		Timeout: timeout,
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

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
		err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
			var err error
			r, httpResp, err := r.Client.SnapshotApi.CreateSnapshot(ctx).CreateSnapshotRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			snapshotId = *r.GetSnapshot().SnapshotId
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to create snapshot from volume '%s' ", volumeId),
				err.Error(),
			)
			return
		}

		tags := oscgo.ResourceTag{
			Key:   "Name",
			Value: data.TerminationSnapshotName.String(),
		}
		err = createOAPITags(ctx, r.Client, []oscgo.ResourceTag{tags}, snapshotId)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to create tags for snapshot '%s' ", snapshotId),
				err.Error(),
			)
			return
		}
	}

	delReq := osc.DeleteVolumeRequest{
		VolumeId: data.VolumeId.ValueString(),
	}
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.VolumeApi.DeleteVolume(ctx).DeleteVolumeRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete volume",
			err.Error(),
		)
		return
	}
}

func setVolumeState(ctx context.Context, r *resourceVolume, data *VolumeModel) error {

	volumeFilters := osc.FiltersVolume{
		VolumeIds: &[]string{data.VolumeId.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'volume' read timeout value. Error: %v: ", diags.Errors())
	}

	readReq := osc.ReadVolumesRequest{
		Filters: &volumeFilters,
	}
	var readResp osc.ReadVolumesResponse
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
	if len(readResp.GetVolumes()) == 0 {
		return errors.New("Empty")
	}

	volume := readResp.GetVolumes()[0]
	tags, diag := flattenOAPITagsFW(ctx, volume.GetTags())
	if diag.HasError() {
		return fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	data.LinkedVolumes, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: utils.GetAttrTypes(BlockLinkedVolumes{})}, getLinkedVolumesFromApiResponse(volume.GetLinkedVolumes())) //volume.GetLinkedVolumes()) //utils.GetAttrTypes(BlockLinkedVolumes{})}, volume.GetLinkedVolumes())
	if diags.HasError() {
		return fmt.Errorf("unable to set LinkedVolumes block: %v: ", diags.Errors())
	}
	data.CreationDate = types.StringValue(volume.GetCreationDate())
	if data.TerminationSnapshotName.IsNull() {
		data.TerminationSnapshotName = types.StringNull()
	}
	if data.SnapshotId.IsNull() {
		data.SnapshotId = types.StringNull()
	}
	data.SubregionName = types.StringValue(volume.GetSubregionName())
	data.VolumeType = types.StringValue(volume.GetVolumeType())
	data.VolumeId = types.StringValue(volume.GetVolumeId())
	data.SnapshotId = types.StringValue(volume.GetSnapshotId())
	data.State = types.StringValue(volume.GetState())
	if data.VolumeType.ValueString() != "standard" {
		data.Iops = types.Int32Value(volume.GetIops())
	} else {
		data.Iops = types.Int32Null()
	}
	data.Size = types.Int32Value(volume.GetSize())
	data.Id = types.StringValue(volume.GetVolumeId())
	return nil
}

func getVolumeStateRefreshFunc(ctx context.Context, conn *osc.APIClient, timeout time.Duration, volumeID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		var resp osc.ReadVolumesResponse

		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(ctx).ReadVolumesRequest(osc.ReadVolumesRequest{
				Filters: &osc.FiltersVolume{
					VolumeIds: &[]string{volumeID},
				}}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return nil, "", err
		}
		v := resp.GetVolumes()[0]
		return v, v.GetState(), nil
	}
}

func getLinkedVolumesFromApiResponse(linkedVols []osc.LinkedVolume) []BlockLinkedVolumes {
	linkedVolumes := make([]BlockLinkedVolumes, 0, len(linkedVols))

	for _, linkedVol := range linkedVols {
		linkedVolume := BlockLinkedVolumes{}
		linkedVolume.DeleteOnVmDeletion = types.BoolValue(linkedVol.GetDeleteOnVmDeletion())
		linkedVolume.DeviceName = types.StringValue(linkedVol.GetDeviceName())
		linkedVolume.VmId = types.StringValue(linkedVol.GetVmId())
		linkedVolume.State = types.StringValue(linkedVol.GetState())
		linkedVolume.VolumeId = types.StringValue(linkedVol.GetVolumeId())
		linkedVolumes = append(linkedVolumes, linkedVolume)
	}
	return linkedVolumes
}
