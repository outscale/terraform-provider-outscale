package outscale

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/utils"
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
	Tags                    []ResourceTag  `tfsdk:"tags"`
	Id                      types.String   `tfsdk:"id"`
}

type BlockLinkedVolumes struct {
	DeleteOnVmDeletion types.Bool   `tfsdk:"delete_on_vm_deletion"`
	VolumeId           types.String `tfsdk:"volume_id"`
	DeviceName         types.String `tfsdk:"device_name"`
	State              types.String `tfsdk:"state"`
	VmId               types.String `tfsdk:"vm_id"`
}

type resourceVolume struct {
	Client *oscgo.APIClient
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
			"tags": TagsSchema(),
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
				PlanModifiers: []planmodifier.String{
					fwmodifyplan.ForceNewFramework(),
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

	createReq := oscgo.NewCreateVolumeRequest(data.SubregionName.ValueString())
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

	var createResp oscgo.CreateVolumeResponse
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

	if len(data.Tags) > 0 {
		err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(data.Tags), volumeId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add Tags on volume resource",
				err.Error(),
			)
			return
		}
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
	var dataPlan, dataState VolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateTimeout, diags := dataPlan.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeId := dataState.VolumeId.ValueString()
	dataState.TerminationSnapshotName = dataPlan.TerminationSnapshotName

	if !dataPlan.Size.Equal(dataState.Size) {
		updateReq := oscgo.NewUpdateVolumeRequest(volumeId)
		updateReq.Size = dataPlan.Size.ValueInt32Pointer()
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.VolumeApi.UpdateVolume(ctx).UpdateVolumeRequest(*updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			dataState.RequestId = types.StringValue(*rp.ResponseContext.RequestId)
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update volume size",
				err.Error(),
			)
			return
		}
	}

	if dataPlan.VolumeType.ValueString() == "io1" && !dataPlan.Iops.Equal(dataState.Iops) {
		updateReq := oscgo.NewUpdateVolumeRequest(volumeId)
		updateReq.Iops = dataPlan.Iops.ValueInt32Pointer()
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			rp, httpResp, err := r.Client.VolumeApi.UpdateVolume(ctx).UpdateVolumeRequest(*updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			dataState.RequestId = types.StringValue(*rp.ResponseContext.RequestId)
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update volume resource",
				err.Error(),
			)
			return
		}
	}

	if !reflect.DeepEqual(dataPlan.Tags, dataState.Tags) {
		toRemove, toCreate := diffOSCAPITags(tagsToOSCResourceTag(dataPlan.Tags), tagsToOSCResourceTag(dataState.Tags))
		err := updateFrameworkTags(ctx, r.Client, toCreate, toRemove, volumeId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Tags on volume resource",
				err.Error(),
			)
			return
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"available", "in-use"},
		Refresh:    getVolumeStateRefreshFunc(ctx, r.Client, updateTimeout, volumeId),
		Timeout:    updateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Waiting for volume state to become available",
			fmt.Sprintf("Unexpected volume (%s) state: '%s' ", volumeId, err.Error()),
		)
		return
	}
	err := setVolumeState(ctx, r, &dataState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set volume state after tags updating.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
		request := oscgo.CreateSnapshotRequest{
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

		snapTags := []ResourceTag{
			{
				Key:   types.StringValue("Name"),
				Value: types.StringValue(data.TerminationSnapshotName.String()),
			},
		}
		err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(snapTags), snapshotId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add Tags on snapshot resource",
				err.Error(),
			)
			return
		}
	}

	delReq := oscgo.DeleteVolumeRequest{
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

	volumeFilters := oscgo.FiltersVolume{
		VolumeIds: &[]string{data.VolumeId.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'volume' read timeout value. Error: %v: ", diags.Errors())
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
	if len(readResp.GetVolumes()) == 0 {
		return errors.New("Empty")
	}

	volume := readResp.GetVolumes()[0]
	data.Tags = getTagsFromApiResponse(volume.GetTags())
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
		data.Iops = types.Int32Value(utils.DefaultIops)
	}
	data.Size = types.Int32Value(volume.GetSize())
	data.Id = types.StringValue(volume.GetVolumeId())
	return nil
}

func getVolumeStateRefreshFunc(ctx context.Context, conn *oscgo.APIClient, timeout time.Duration, volumeID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		var resp oscgo.ReadVolumesResponse

		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(ctx).ReadVolumesRequest(oscgo.ReadVolumesRequest{
				Filters: &oscgo.FiltersVolume{
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

func getLinkedVolumesFromApiResponse(linkedVols []oscgo.LinkedVolume) []BlockLinkedVolumes {
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
