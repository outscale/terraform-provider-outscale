package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                     = &snapshotResource{}
	_ resource.ResourceWithConfigure        = &snapshotResource{}
	_ resource.ResourceWithImportState      = &snapshotResource{}
	_ resource.ResourceWithConfigValidators = &snapshotResource{}
)

const (
	snapshotErrCreate = "Unable to create Snapshot"
	snapshotErrRead   = "Unable to read Snapshot"
	snapshotErrDelete = "Unable to delete Snapshot"
	snapshotErrState  = "Unable to set Snapshot state"
	snapshotErrWait   = "Unable to wait for Snapshot to be available"

	snapshotCreateTimeout = 40 * time.Minute
)

type snapshotModel struct {
	Description               types.String   `tfsdk:"description"`
	SnapshotSize              types.Int64    `tfsdk:"snapshot_size"`
	FileLocation              types.String   `tfsdk:"file_location"`
	SourceRegionName          types.String   `tfsdk:"source_region_name"`
	SourceSnapshotId          types.String   `tfsdk:"source_snapshot_id"`
	VolumeId                  types.String   `tfsdk:"volume_id"`
	AccountAlias              types.String   `tfsdk:"account_alias"`
	AccountId                 types.String   `tfsdk:"account_id"`
	CreationDate              types.String   `tfsdk:"creation_date"`
	PermissionsToCreateVolume types.List     `tfsdk:"permissions_to_create_volume"`
	Progress                  types.Int64    `tfsdk:"progress"`
	SnapshotId                types.String   `tfsdk:"snapshot_id"`
	State                     types.String   `tfsdk:"state"`
	VolumeSize                types.Int64    `tfsdk:"volume_size"`
	Id                        types.String   `tfsdk:"id"`
	RequestId                 types.String   `tfsdk:"request_id"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type snapshotPermissionsModel struct {
	GlobalPermission types.Bool   `tfsdk:"global_permission"`
	AccountId        types.String `tfsdk:"account_id"`
}

var snapshotPermissionsAttrTypes = fwhelpers.GetAttrTypes(snapshotPermissionsModel{})

type snapshotResource struct {
	Client *osc.Client
}

func NewResourceSnapshot() resource.Resource {
	return &snapshotResource{}
}

func (r *snapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *snapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Snapshot identifier. Got: %v", req.ID),
		)
		return
	}

	var data snapshotModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.SnapshotId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal
	data.Tags = TagsNull()
	data.PermissionsToCreateVolume = types.ListNull(types.ObjectType{AttrTypes: snapshotPermissionsAttrTypes})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *snapshotResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("volume_id"),
			path.MatchRoot("snapshot_size"),
			path.MatchRoot("source_snapshot_id"),
		),
	}
}

func (r *snapshotResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_size": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"file_location": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_region_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_snapshot_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_alias": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"permissions_to_create_volume": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: snapshotPermissionsAttrTypes,
				},
			},
			"progress": schema.Int64Attribute{
				Computed: true,
			},
			"snapshot_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"volume_size": schema.Int64Attribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data snapshotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := data.Timeouts.Create(ctx, snapshotCreateTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreateSnapshotRequest{}

	if fwhelpers.IsSet(data.Description) {
		createReq.Description = data.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.FileLocation) {
		createReq.FileLocation = data.FileLocation.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.VolumeId) {
		createReq.VolumeId = data.VolumeId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.SnapshotSize) {
		createReq.SnapshotSize = data.SnapshotSize.ValueInt64Pointer()
	}
	if fwhelpers.IsSet(data.SourceSnapshotId) {
		createReq.SourceSnapshotId = data.SourceSnapshotId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.SourceRegionName) {
		createReq.SourceRegionName = data.SourceRegionName.ValueStringPointer()
	}

	createResp, err := r.Client.CreateSnapshot(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(snapshotErrCreate, err.Error())
		return
	}

	snapshotId := createResp.Snapshot.SnapshotId
	data.Id = to.String(snapshotId)
	data.SnapshotId = to.String(snapshotId)

	stateConf := &stateconf.StateChangeConf[osc.SnapshotState]{
		Pending: stateconf.States(
			osc.SnapshotStatePending,
			osc.SnapshotStateInQueue,
		),
		Target:  stateconf.States(osc.SnapshotStateCompleted),
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(snapshotId),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(snapshotErrWait, err.Error())
		return
	}

	diag = createOAPITagsFW(ctx, r.Client, timeout, data.Tags, snapshotId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(snapshotErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data snapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
		resp.Diagnostics.AddError(snapshotErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData snapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	diag = updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.Id.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(snapshotErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data snapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	_, err := r.Client.DeleteSnapshot(ctx, osc.DeleteSnapshotRequest{
		SnapshotId: data.Id.ValueString(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(snapshotErrDelete, err.Error())
	}
}

func (r *snapshotResource) read(ctx context.Context, timeout time.Duration, data snapshotModel) (snapshotModel, error) {
	readReq := osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{data.Id.ValueString()}},
	}

	resp, err := r.Client.ReadSnapshots(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Snapshots == nil || len(*resp.Snapshots) == 0 {
		return data, ErrResourceEmpty
	}

	snapshot := (*resp.Snapshots)[0]

	tags, diag := flattenOAPITagsFW(ctx, ptr.From(snapshot.Tags))
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	permissions, diag := to.ListObject(ctx, r.flattenPermissions(snapshot.PermissionsToCreateVolume), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Description = to.String(ptr.From(snapshot.Description))
	data.SnapshotId = to.String(snapshot.SnapshotId)
	data.VolumeId = to.String(snapshot.VolumeId)
	data.AccountAlias = to.String(ptr.From(snapshot.AccountAlias))
	data.AccountId = to.String(snapshot.AccountId)
	data.CreationDate = to.String(from.ISO8601(snapshot.CreationDate))
	data.PermissionsToCreateVolume = permissions
	data.Progress = to.Int64(ptr.From(snapshot.Progress))
	data.State = to.String(snapshot.State)
	data.VolumeSize = to.Int64(snapshot.VolumeSize)
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Id = to.String(snapshot.SnapshotId)
	data.Tags = tags

	return data, nil
}

func (r *snapshotResource) flattenPermissions(p *osc.PermissionsOnResource) []snapshotPermissionsModel {
	if p == nil || p.AccountIds == nil {
		return nil
	}

	return lo.Map(*p.AccountIds, func(id string, _ int) snapshotPermissionsModel {
		return snapshotPermissionsModel{
			AccountId:        to.String(id),
			GlobalPermission: to.Bool(ptr.From(p.GlobalPermission)),
		}
	})
}

func (r *snapshotResource) stateRefreshFunc(id string) func(context.Context) (any, osc.SnapshotState, error) {
	return func(ctx context.Context) (any, osc.SnapshotState, error) {
		req := osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{id}},
		}
		resp, err := r.Client.ReadSnapshots(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp.Snapshots == nil || len(*resp.Snapshots) == 0 {
			return nil, "", ErrResourceEmpty
		}

		snapshot := (*resp.Snapshots)[0]
		return resp, snapshot.State, nil
	}
}
