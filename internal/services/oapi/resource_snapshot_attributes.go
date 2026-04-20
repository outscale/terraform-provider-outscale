package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                     = &snapshotAttributesResource{}
	_ resource.ResourceWithConfigure        = &snapshotAttributesResource{}
	_ resource.ResourceWithConfigValidators = &snapshotAttributesResource{}
	_ resource.ResourceWithModifyPlan       = &snapshotAttributesResource{}
)

const (
	snapAttrErrCreate = "Unable to create Snapshot Attributes"
	snapAttrErrRead   = "Unable to read Snapshot Attributes"
	snapAttrErrState  = "Unable to set Snapshot Attributes state"
)

type snapshotAttributesModel struct {
	Additions  types.List     `tfsdk:"permissions_to_create_volume_additions"`
	Removals   types.List     `tfsdk:"permissions_to_create_volume_removals"`
	SnapshotId types.String   `tfsdk:"snapshot_id"`
	AccountId  types.String   `tfsdk:"account_id"`
	Id         types.String   `tfsdk:"id"`
	RequestId  types.String   `tfsdk:"request_id"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

type snapshotAttributesPermModel struct {
	AccountIds       types.List `tfsdk:"account_ids"`
	GlobalPermission types.Bool `tfsdk:"global_permission"`
}

var snapshotAttributesPermAttrTypes = types.ObjectType{AttrTypes: map[string]attr.Type{
	"account_ids":       types.ListType{ElemType: types.StringType},
	"global_permission": types.BoolType,
}}

type snapshotAttributesResource struct {
	Client *osc.Client
}

func NewResourceSnapshotAttributes() resource.Resource {
	return &snapshotAttributesResource{}
}

func (r *snapshotAttributesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *snapshotAttributesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_attributes"
}

func (r *snapshotAttributesResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("permissions_to_create_volume_additions"),
			path.MatchRoot("permissions_to_create_volume_removals"),
		),
	}
}

func (r *snapshotAttributesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	permBlock := schema.NestedBlockObject{
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"account_ids": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"global_permission": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}

	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
			"permissions_to_create_volume_additions": schema.ListNestedBlock{
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: permBlock,
			},
			"permissions_to_create_volume_removals": schema.ListNestedBlock{
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: permBlock,
			},
		},
		Attributes: map[string]schema.Attribute{
			"snapshot_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

func (r *snapshotAttributesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data snapshotAttributesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	snapshotId := data.SnapshotId.ValueString()
	updateReq := osc.UpdateSnapshotRequest{
		SnapshotId: snapshotId,
	}

	perms := osc.PermissionsOnResourceCreation{}
	buildPermissions := func(list types.List) (*osc.PermissionsOnResource, diag.Diagnostics) {
		if !fwhelpers.IsSet(list) {
			return nil, nil
		}

		models, diag := to.Slice[snapshotAttributesPermModel](ctx, list)
		if diag.HasError() {
			return nil, diag
		}
		if len(models) == 0 {
			return nil, nil
		}

		perm := osc.PermissionsOnResource{}
		if fwhelpers.IsSet(models[0].AccountIds) {
			accountIds, diag := to.Slice[string](ctx, models[0].AccountIds)
			if diag.HasError() {
				return nil, diag
			}
			perm.AccountIds = &accountIds
		}
		if fwhelpers.IsSet(models[0].GlobalPermission) {
			gp := models[0].GlobalPermission.ValueBool()
			perm.GlobalPermission = &gp
		}

		return &perm, diag
	}

	addition, diag := buildPermissions(data.Additions)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	perms.Additions = addition

	removal, diag := buildPermissions(data.Removals)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	perms.Removals = removal

	updateReq.PermissionsToCreateVolume = perms

	_, err := r.Client.UpdateSnapshot(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(snapAttrErrCreate, err.Error())
		return
	}

	data.Id = to.String(snapshotId)

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(snapAttrErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *snapshotAttributesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data snapshotAttributesModel
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
		resp.Diagnostics.AddError(snapAttrErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *snapshotAttributesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	if req.State.Raw.IsNull() {
		return
	}

	var stateData snapshotAttributesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planData snapshotAttributesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if fwhelpers.IsSet(stateData.Additions) && fwhelpers.IsSet(stateData.Removals) {
		if planData.Additions.IsNull() || planData.Additions.IsUnknown() {
			planData.Additions = stateData.Additions
		}
		if planData.Removals.IsNull() || planData.Removals.IsUnknown() {
			planData.Removals = stateData.Removals
		}
	}

	additions, diags := normalizeSnapshotAttributesPermissionsList(ctx, planData.Additions)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	planData.Additions = additions

	removals, diags := normalizeSnapshotAttributesPermissionsList(ctx, planData.Removals)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	planData.Removals = removals

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
}

func (r *snapshotAttributesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Every attribute is RequiresReplace, Update is never called
}

func (r *snapshotAttributesResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func normalizeSnapshotAttributesPermissionsList(ctx context.Context, list types.List) (types.List, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return list, nil
	}

	perms, diags := to.Slice[snapshotAttributesPermModel](ctx, list)
	if diags.HasError() {
		return list, diags
	}

	for i := range perms {
		if perms[i].AccountIds.IsNull() || perms[i].AccountIds.IsUnknown() {
			perms[i].AccountIds = types.ListValueMust(types.StringType, []attr.Value{})
		}
		if perms[i].GlobalPermission.IsNull() || perms[i].GlobalPermission.IsUnknown() {
			perms[i].GlobalPermission = types.BoolValue(false)
		}
	}

	return types.ListValueFrom(ctx, snapshotAttributesPermAttrTypes, perms)
}

func (r *snapshotAttributesResource) read(ctx context.Context, timeout time.Duration, data snapshotAttributesModel) (snapshotAttributesModel, error) {
	resp, err := r.Client.ReadSnapshots(ctx, osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{data.Id.ValueString()}},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Snapshots == nil || len(*resp.Snapshots) == 0 {
		return data, ErrResourceEmpty
	}

	snapshot := (*resp.Snapshots)[0]

	data.Id = to.String(snapshot.SnapshotId)
	data.AccountId = to.String(snapshot.AccountId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return data, nil
}
