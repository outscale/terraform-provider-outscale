package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/fwvalidators"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/to"
)

var (
	_ resource.Resource              = &oksClusterResource{}
	_ resource.ResourceWithConfigure = &oksClusterResource{}
)

type ClusterModel struct {
	AdminLbu              types.Bool     `tfsdk:"admin_lbu"`
	AdminWhitelist        types.Set      `tfsdk:"admin_whitelist"`
	AdmissionFlags        types.Object   `tfsdk:"admission_flags"`
	AutoMaintenances      types.Object   `tfsdk:"auto_maintenances"`
	CidrPods              types.String   `tfsdk:"cidr_pods"`
	CidrService           types.String   `tfsdk:"cidr_service"`
	ClusterDns            types.String   `tfsdk:"cluster_dns"`
	Cni                   types.String   `tfsdk:"cni"`
	ControlPlanes         types.String   `tfsdk:"control_planes"`
	CpMultiAz             types.Bool     `tfsdk:"cp_multi_az"`
	CpSubregions          types.Set      `tfsdk:"cp_subregions"`
	Description           types.String   `tfsdk:"description"`
	DisableApiTermination types.Bool     `tfsdk:"disable_api_termination"`
	ExpectedControlPlanes types.String   `tfsdk:"expected_control_planes"`
	ExpectedVersion       types.String   `tfsdk:"expected_version"`
	Id                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	ProjectId             types.String   `tfsdk:"project_id"`
	Quirks                types.Set      `tfsdk:"quirks"`
	Statuses              types.Object   `tfsdk:"statuses"`
	Version               types.String   `tfsdk:"version"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
	RequestId             types.String   `tfsdk:"request_id"`
	OKSTagsModel
}

type AutoMaintenancesModel struct {
	MinorUpgradeMaintenance *MaintenanceWindowModel `tfsdk:"minor_upgrade_maintenance"`
	PatchUpgradeMaintenance *MaintenanceWindowModel `tfsdk:"patch_upgrade_maintenance"`
}

type MaintenanceWindowModel struct {
	DurationHours types.Int64  `tfsdk:"duration_hours"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	StartHour     types.Int64  `tfsdk:"start_hour"`
	Tz            types.String `tfsdk:"tz"`
	WeekDay       types.String `tfsdk:"week_day"`
}

type AdmissionFlagsModel struct {
	AppliedAdmissionPlugins types.Set `tfsdk:"applied_admission_plugins"`
	DisableAdmissionPlugins types.Set `tfsdk:"disable_admission_plugins"`
	EnableAdmissionPlugins  types.Set `tfsdk:"enable_admission_plugins"`
}

type StatusesModel struct {
	AvailableUpgrade types.String      `tfsdk:"available_upgrade"`
	CreatedAt        timetypes.RFC3339 `tfsdk:"created_at"`
	DeletedAt        timetypes.RFC3339 `tfsdk:"deleted_at"`
	Status           types.String      `tfsdk:"status"`
	UpdatedAt        timetypes.RFC3339 `tfsdk:"updated_at"`
}

type oksClusterResource struct {
	Client *oks.Client
}

func NewResourceCluster() resource.Resource {
	return &oksClusterResource{}
}

func (r *oksClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *oks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Client = client.OKS
}

func (r *oksClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *oksClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_cluster"
}

func (r *oksClusterResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"admin_lbu": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"admin_whitelist": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"admission_flags": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"applied_admission_plugins": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"disable_admission_plugins": schema.SetAttribute{
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
					},
					"enable_admission_plugins": schema.SetAttribute{
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"auto_maintenances": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"minor_upgrade_maintenance": schema.SingleNestedAttribute{
						Computed:   true,
						Optional:   true,
						Attributes: autoMaintenancesSchema,
					},
					"patch_upgrade_maintenance": schema.SingleNestedAttribute{
						Computed:   true,
						Optional:   true,
						Attributes: autoMaintenancesSchema,
					},
				},
			},
			"cidr_pods": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					fwvalidators.IsCIDR(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr_service": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					fwvalidators.IsCIDR(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_dns": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					fwvalidators.IsIP(),
				},
			},
			"cni": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					fwvalidators.IsIP(),
				},
			},
			"control_planes": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"cp_multi_az": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"cp_subregions": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"disable_api_termination": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"expected_control_planes": schema.StringAttribute{
				Computed: true,
			},
			"expected_version": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"quirks": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"statuses": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"available_upgrade": schema.StringAttribute{
						Computed: true,
					},
					"created_at": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"deleted_at": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"status": schema.StringAttribute{
						Computed: true,
					},
					"updated_at": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
				},
			},
			"version": schema.StringAttribute{
				Required: true,
			},
			"tags": OKSTagsSchema(),
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

var autoMaintenancesSchema = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Optional: true,
		Computed: true,
	},
	"duration_hours": schema.Int64Attribute{
		Optional: true,
		Computed: true,
	},
	"start_hour": schema.Int64Attribute{
		Optional: true,
		Computed: true,
	},
	"week_day": schema.StringAttribute{
		Optional: true,
		Computed: true,
	},
	"tz": schema.StringAttribute{
		Optional: true,
		Computed: true,
	},
}

func (r *oksClusterResource) expandOKSAutoMaintenances(data *MaintenanceWindowModel) (auto oks.MaintenanceWindow) {
	if data != nil {
		if utils.IsSet(data.DurationHours) {
			hours := int(data.DurationHours.ValueInt64())
			auto.DurationHours = &hours
		}
		if utils.IsSet(data.Enabled) {
			auto.Enabled = data.Enabled.ValueBoolPointer()
		}
		if utils.IsSet(data.StartHour) {
			hour := int(data.StartHour.ValueInt64())
			auto.StartHour = &hour
		}
		if utils.IsSet(data.Tz) {
			auto.Tz = data.Tz.ValueStringPointer()
		}
		if utils.IsSet(data.WeekDay) {
			auto.WeekDay = (*oks.MaintenanceWindowWeekDay)(data.WeekDay.ValueStringPointer())
		}
	}
	return
}

func (r *oksClusterResource) expandOKSAdmissionFlags(ctx context.Context, data AdmissionFlagsModel) (admissionFlags oks.AdmissionFlagsInput, diags diag.Diagnostics) {
	if utils.IsSet(data.DisableAdmissionPlugins) {
		disablePlugins, diags := to.Slice[string](ctx, data.DisableAdmissionPlugins)
		if diags.HasError() {
			return admissionFlags, diags
		}
		admissionFlags.DisableAdmissionPlugins = &disablePlugins
	}
	if utils.IsSet(data.EnableAdmissionPlugins) {
		enablePlugins, diags := to.Slice[string](ctx, data.EnableAdmissionPlugins)
		if diags.HasError() {
			return admissionFlags, diags
		}
		admissionFlags.EnableAdmissionPlugins = &enablePlugins
	}
	return
}

func (r *oksClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterModel
	diags := req.Plan.Get(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}

	input := oks.ClusterInput{
		Name:      data.Name.ValueString(),
		ProjectId: data.ProjectId.ValueString(),
		Version:   data.Version.ValueString(),
		AutoMaintenances: oks.AutoMaintenances{
			MinorUpgradeMaintenance: oks.MaintenanceWindow{},
			PatchUpgradeMaintenance: oks.MaintenanceWindow{},
		},
	}

	whitelist, diags := to.Slice[string](ctx, data.AdminWhitelist)
	if utils.CheckDiags(resp, diags) {
		return
	}
	input.AdminWhitelist = whitelist

	if utils.IsSet(data.AutoMaintenances) {
		auto, diags := to.Obj[AutoMaintenancesModel](ctx, data.AutoMaintenances)
		if utils.CheckDiags(resp, diags) {
			return
		}

		input.AutoMaintenances.MinorUpgradeMaintenance = r.expandOKSAutoMaintenances(auto.MinorUpgradeMaintenance)
		input.AutoMaintenances.PatchUpgradeMaintenance = r.expandOKSAutoMaintenances(auto.PatchUpgradeMaintenance)
	}

	if utils.IsSet(data.Description) {
		input.Description = data.Description.ValueStringPointer()
	}
	if utils.IsSet(data.CpMultiAz) {
		input.CpMultiAz = data.CpMultiAz.ValueBoolPointer()
	}
	if utils.IsSet(data.CpSubregions) {
		sub, diags := to.Slice[string](ctx, data.CpSubregions)
		if utils.CheckDiags(resp, diags) {
			return
		}
		input.CpSubregions = &sub
	}
	if utils.IsSet(data.AdminLbu) {
		input.AdminLbu = data.AdminLbu.ValueBoolPointer()
	}
	if utils.IsSet(data.AdmissionFlags) {
		admissionModel, diags := to.Obj[AdmissionFlagsModel](ctx, data.AdmissionFlags)
		if utils.CheckDiags(resp, diags) {
			return
		}

		admissionFlags, diags := r.expandOKSAdmissionFlags(ctx, admissionModel)
		if utils.CheckDiags(resp, diags) {
			return
		}
		input.AdmissionFlags = &admissionFlags
	}
	if utils.IsSet(data.CidrPods) {
		input.CidrPods = data.CidrPods.ValueStringPointer()
	}
	if utils.IsSet(data.CidrService) {
		input.CidrService = data.CidrService.ValueStringPointer()
	}
	if utils.IsSet(data.ClusterDns) {
		input.ClusterDns = data.ClusterDns.ValueStringPointer()
	}
	if utils.IsSet(data.ControlPlanes) {
		input.ControlPlanes = data.ControlPlanes.ValueStringPointer()
	}
	if utils.IsSet(data.Quirks) {
		quirks, diags := to.Slice[string](ctx, data.Quirks)
		if utils.CheckDiags(resp, diags) {
			return
		}
		input.Quirks = &quirks
	}
	if utils.IsSet(data.DisableApiTermination) {
		input.DisableApiTermination = data.DisableApiTermination.ValueBoolPointer()
	}

	tags, diags := expandOKSTags(ctx, data.OKSTagsModel)
	if utils.CheckDiags(resp, diags) {
		return
	}
	input.Tags = &tags

	createResp, err := r.Client.CreateCluster(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Cluster",
			"Error: "+err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(createResp.Cluster.Id)

	to, diags := data.Timeouts.Create(ctx, utils.CreateOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	data, err = r.setOKSClusterState(ctx, data, to)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}
}

func (r *oksClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterModel

	diags := req.State.Get(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}

	to, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	data, err := r.setOKSClusterState(ctx, data, to)
	if err != nil {
		if code := oks.StatusCodeHelper(err); code != nil && *code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}
}

func (r *oksClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state ClusterModel
		update      oks.ClusterUpdate
	)
	diags := req.Plan.Get(ctx, &plan)
	if utils.CheckDiags(resp, diags) {
		return
	}
	diags = req.State.Get(ctx, &state)
	if utils.CheckDiags(resp, diags) {
		return
	}

	statuses, diags := to.Obj[StatusesModel](ctx, state.Statuses)
	if utils.CheckDiags(resp, diags) {
		return
	}
	if statuses.Status.ValueString() == "deleting" {
		resp.Diagnostics.AddError(
			"Unable to Update Cluster",
			"Cluster is currently being deleted and cannot be updated. The resource will be removed from the state once deleted.",
		)
		return
	}

	if utils.IsSet(plan.AdminWhitelist) && !plan.AdminWhitelist.Equal(state.AdminWhitelist) {
		admins, diags := to.Slice[string](ctx, plan.AdminWhitelist)
		if utils.CheckDiags(resp, diags) {
			return
		}
		update.AdminWhitelist = &admins
	}
	if utils.IsSet(plan.AdmissionFlags) && !plan.AdmissionFlags.Equal(state.AdmissionFlags) {
		admissionModel, diags := to.Obj[AdmissionFlagsModel](ctx, plan.AdmissionFlags)
		if utils.CheckDiags(resp, diags) {
			return
		}

		admissionFlags, diags := r.expandOKSAdmissionFlags(ctx, admissionModel)
		if utils.CheckDiags(resp, diags) {
			return
		}
		update.AdmissionFlags = &admissionFlags
	}
	if utils.IsSet(plan.AutoMaintenances) && !plan.AutoMaintenances.Equal(state.AutoMaintenances) {
		auto, diags := to.Obj[AutoMaintenancesModel](ctx, plan.AutoMaintenances)
		if utils.CheckDiags(resp, diags) {
			return
		}

		update.AutoMaintenances = &oks.AutoMaintenances{
			MinorUpgradeMaintenance: r.expandOKSAutoMaintenances(auto.MinorUpgradeMaintenance),
			PatchUpgradeMaintenance: r.expandOKSAutoMaintenances(auto.PatchUpgradeMaintenance),
		}
	}
	if utils.IsSet(plan.ControlPlanes) && !plan.ControlPlanes.Equal(state.ControlPlanes) {
		update.ControlPlanes = plan.ControlPlanes.ValueStringPointer()
	}
	if utils.IsSet(plan.Description) && !plan.Description.Equal(state.Description) {
		update.Description = plan.Description.ValueStringPointer()
	}
	if utils.IsSet(plan.DisableApiTermination) && !plan.DisableApiTermination.Equal(state.DisableApiTermination) {
		update.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}
	if utils.IsSet(plan.Quirks) && !plan.Quirks.Equal(state.Quirks) {
		quirks, diags := to.Slice[string](ctx, plan.Quirks)
		if utils.CheckDiags(resp, diags) {
			return
		}
		update.Quirks = &quirks
	}
	if utils.IsSet(plan.Version) && !plan.Version.Equal(state.Version) {
		update.Version = plan.Version.ValueStringPointer()
	}
	tags, diags := cmpOKSTags(ctx, plan.OKSTagsModel, state.OKSTagsModel)
	if utils.CheckDiags(resp, diags) {
		return
	}
	if tags != nil {
		update.Tags = &tags
	}

	updateResp, err := r.Client.UpdateCluster(ctx, state.Id.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Cluster",
			"Error: "+err.Error(),
		)
		return
	}
	state.RequestId = to.String(updateResp.ResponseContext.RequestId)

	to, diags := state.Timeouts.Update(ctx, utils.UpdateOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	data, err := r.setOKSClusterState(ctx, state, to)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}
}

func (r *oksClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterModel

	diags := req.State.Get(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}
	_, err := r.Client.DeleteCluster(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Cluster",
			"Error: "+err.Error(),
		)
		return
	}

	to, diags := data.Timeouts.Update(ctx, utils.DeleteOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	_, err = r.waitForClusterState(ctx, data.Id.ValueString(), []string{"pending", "deleting"}, []string{}, to)
	if err != nil {
		if code := oks.StatusCodeHelper(err); code != nil && *code != 404 {
			resp.Diagnostics.AddError("Unable to wait for Cluster complete deletion.", "Error: "+err.Error())
		}
	}
}

func (r *oksClusterResource) waitForClusterState(ctx context.Context, id string, pending []string, target []string, timeout time.Duration) (*oks.ClusterResponse, error) {
	resp, err := utils.WaitForResource[oks.ClusterResponse](ctx, &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			resp, err := r.Client.GetCluster(ctx, id)
			if err != nil {
				return resp, "", err
			}
			return resp, *resp.Cluster.Statuses.Status, nil
		},
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *oksClusterResource) setOKSAdmissionFlags(ctx context.Context, data *ClusterModel, auto *oks.AdmissionFlags) diag.Diagnostics {
	if auto != nil {
		var model AdmissionFlagsModel
		applied, diags := types.SetValueFrom(ctx, types.StringType, auto.AppliedAdmissionPlugins)
		if diags.HasError() {
			return diags
		}
		disable, diags := types.SetValueFrom(ctx, types.StringType, auto.DisableAdmissionPlugins)
		if diags.HasError() {
			return diags
		}
		enable, diags := types.SetValueFrom(ctx, types.StringType, auto.EnableAdmissionPlugins)
		if diags.HasError() {
			return diags
		}

		model.AppliedAdmissionPlugins = applied
		model.DisableAdmissionPlugins = disable
		model.EnableAdmissionPlugins = enable

		obj, diags := types.ObjectValueFrom(ctx, data.AdmissionFlags.AttributeTypes(ctx), model)
		if diags.HasError() {
			return diags
		}
		data.AdmissionFlags = obj
	}

	return nil
}

func (r *oksClusterResource) setOKSAutoMaintenances(ctx context.Context, data *ClusterModel, auto oks.AutoMaintenances) diag.Diagnostics {
	var model AutoMaintenancesModel
	setAuto := func(window oks.MaintenanceWindow) *MaintenanceWindowModel {
		var model MaintenanceWindowModel
		model.DurationHours = to.Int64(window.DurationHours)
		model.Enabled = to.Bool(window.Enabled)
		model.StartHour = to.Int64(window.StartHour)
		model.Tz = to.String(window.Tz)
		model.WeekDay = to.String((*string)(window.WeekDay))

		return &model
	}
	model.MinorUpgradeMaintenance = setAuto(auto.MinorUpgradeMaintenance)
	model.PatchUpgradeMaintenance = setAuto(auto.PatchUpgradeMaintenance)

	obj, diags := types.ObjectValueFrom(ctx, data.AutoMaintenances.AttributeTypes(ctx), &model)
	if diags.HasError() {
		return diags
	}
	data.AutoMaintenances = obj

	return nil
}

func (r *oksClusterResource) setOKSStatuses(ctx context.Context, data *ClusterModel, auto oks.Statuses) diag.Diagnostics {
	var model StatusesModel
	model.AvailableUpgrade = to.String(auto.AvailableUpgrade)
	model.CreatedAt = to.RFC3339(auto.CreatedAt)
	model.DeletedAt = to.RFC3339(auto.DeletedAt)
	model.Status = to.String(auto.Status)
	model.UpdatedAt = to.RFC3339(auto.UpdatedAt)

	obj, diags := types.ObjectValueFrom(ctx, data.Statuses.AttributeTypes(ctx), model)
	if diags.HasError() {
		return diags
	}
	data.Statuses = obj

	return nil
}

func (r *oksClusterResource) setOKSClusterState(ctx context.Context, data ClusterModel, timeout time.Duration) (ClusterModel, error) {
	resp, err := r.waitForClusterState(ctx, data.Id.ValueString(), []string{"pending", "deploying", "updating"}, []string{"ready", "deleting"}, timeout)
	if err != nil {
		return data, err
	}
	cluster := resp.Cluster

	adminWhiteList, diags := types.SetValueFrom(ctx, types.StringType, cluster.AdminWhitelist)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert AdminWhitelist into a Set. Error: %v", diags.Errors())
	}
	cpSubregions, diags := types.SetValueFrom(ctx, types.StringType, cluster.CpSubregions)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert CPSubregions into a Set. Error: %v", diags.Errors())
	}
	diags = r.setOKSAdmissionFlags(ctx, &data, cluster.AdmissionFlags)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert AdmissionFlags into the Schema Model. Error: %v", diags.Errors())
	}
	diags = r.setOKSAutoMaintenances(ctx, &data, cluster.AutoMaintenances)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert AutoMaintenances into the Schema Model. Error: %v", diags.Errors())
	}
	diags = r.setOKSStatuses(ctx, &data, cluster.Statuses)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert Statuses into the Schema Model. Error: %v", diags.Errors())
	}

	data.AdminLbu = to.Bool(cluster.AdminLbu)
	data.AdminWhitelist = adminWhiteList
	data.CidrPods = to.String(cluster.CidrPods)
	data.CidrService = to.String(cluster.CidrService)
	data.ClusterDns = to.String(cluster.ClusterDns)
	data.Cni = to.String(cluster.Cni)
	data.ControlPlanes = to.String(cluster.ControlPlanes)
	data.CpMultiAz = to.Bool(cluster.CpMultiAz)
	data.CpSubregions = cpSubregions
	data.Description = to.String(cluster.Description)
	data.DisableApiTermination = to.Bool(cluster.DisableApiTermination)
	data.ExpectedControlPlanes = to.String(cluster.ExpectedControlPlanes)
	data.ExpectedVersion = to.String(cluster.ExpectedVersion)
	data.Id = to.String(cluster.Id)
	data.Name = to.String(cluster.Name)
	data.ProjectId = to.String(cluster.ProjectId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Version = to.String(cluster.Version)

	tags, diags := flattenOKSTags(ctx, cluster.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("%v", diags.Errors())
	}
	data.Tags = tags

	return data, nil
}
