package oks

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
)

var (
	_ resource.Resource                = &oksClusterResource{}
	_ resource.ResourceWithConfigure   = &oksClusterResource{}
	_ resource.ResourceWithImportState = &oksClusterResource{}
	_ resource.ResourceWithModifyPlan  = &oksClusterResource{}
)

type ClusterModel struct {
	AdminLbu              types.Bool     `tfsdk:"admin_lbu"`
	AdminWhitelist        types.Set      `tfsdk:"admin_whitelist"`
	AdmissionFlags        types.Object   `tfsdk:"admission_flags"`
	AutoMaintenances      types.Object   `tfsdk:"auto_maintenances"`
	MaintenanceWindow     types.Object   `tfsdk:"maintenance_window"`
	CidrPods              types.String   `tfsdk:"cidr_pods"`
	CidrService           types.String   `tfsdk:"cidr_service"`
	ClusterDns            types.String   `tfsdk:"cluster_dns"`
	Cni                   types.String   `tfsdk:"cni"`
	ControlPlanes         types.String   `tfsdk:"control_planes"`
	CpMultiAz             types.Bool     `tfsdk:"cp_multi_az"`
	CpSubregions          types.Set      `tfsdk:"cp_subregions"`
	Description           types.String   `tfsdk:"description"`
	DisableApiTermination types.Bool     `tfsdk:"disable_api_termination"`
	Id                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	ProjectId             types.String   `tfsdk:"project_id"`
	Quirks                types.Set      `tfsdk:"quirks"`
	Statuses              types.Object   `tfsdk:"statuses"`
	Version               types.String   `tfsdk:"version"`
	Kubeconfig            types.String   `tfsdk:"kubeconfig"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
	RequestId             types.String   `tfsdk:"request_id"`
	OKSTagsModel
}

type AutoMaintenancesModel struct {
	MinorUpgradeMaintenanceActual types.Object `tfsdk:"minor_upgrade_maintenance_actual"`
	PatchUpgradeMaintenanceActual types.Object `tfsdk:"patch_upgrade_maintenance_actual"`
	MinorUpgradeMaintenance       types.Object `tfsdk:"minor_upgrade_maintenance"`
	PatchUpgradeMaintenance       types.Object `tfsdk:"patch_upgrade_maintenance"`
}

type MaintenanceWindowModel struct {
	DurationHours types.Int64  `tfsdk:"duration_hours"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	StartHour     types.Int64  `tfsdk:"start_hour"`
	Tz            types.String `tfsdk:"tz"`
	WeekDay       types.String `tfsdk:"week_day"`
}

type MaintenanceModel struct {
	DurationHours types.Int64  `tfsdk:"duration_hours"`
	StartHour     types.Int64  `tfsdk:"start_hour"`
	Tz            types.String `tfsdk:"tz"`
	WeekDay       types.String `tfsdk:"week_day"`
}

var (
	maintenanceWindowAttrTypes = fwhelpers.GetAttrTypes(MaintenanceWindowModel{})
	maintenanceAttrTypes       = fwhelpers.GetAttrTypes(MaintenanceModel{})
)

type AdmissionFlagsModel struct {
	AppliedAdmissionPlugins       types.Set `tfsdk:"applied_admission_plugins"`
	DisableAdmissionPluginsActual types.Set `tfsdk:"disable_admission_plugins_actual"`
	EnableAdmissionPluginsActual  types.Set `tfsdk:"enable_admission_plugins_actual"`
	DisableAdmissionPlugins       types.Set `tfsdk:"disable_admission_plugins"`
	EnableAdmissionPlugins        types.Set `tfsdk:"enable_admission_plugins"`
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
	client, ok := req.ProviderData.(client.OutscaleClient)
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
	attr := []string{"admission_flags.disable_admission_plugins", "admission_flags.enable_admission_plugins", "auto_maintenances.minor_upgrade_maintenance", "auto_maintenances.patch_upgrade_maintenance"}
	resp.Diagnostics.AddWarning("Resource needs an apply", fmt.Sprintf("%q attributes are optional and Terraform cannot verify that the values are applied after Import. If at least one of the values is set in the configuration, an apply is necessary to udpate the state.", attr))
}

func (r *oksClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_cluster"
}

func (r *oksClusterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state, config ClusterModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Config.Get(ctx, &config)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if fwhelpers.IsSet(plan.Version) && !plan.Version.Equal(state.Version) {
		resp.Diagnostics.AddWarning("Cluster Upgrade.",
			"Changing the Cluster version will call an Upgrade on the Cluster.")
	}
	if fwhelpers.IsSet(plan.ControlPlanes) && !plan.ControlPlanes.Equal(state.ControlPlanes) {
		resp.Diagnostics.AddWarning("Cluster upgrade.", "Changing the Cluster control planes will call an Upgrade on the Cluster.")
	}
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
					"disable_admission_plugins_actual": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"enable_admission_plugins_actual": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"disable_admission_plugins": schema.SetAttribute{
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"enable_admission_plugins": schema.SetAttribute{
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"auto_maintenances": schema.SingleNestedAttribute{
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "Configure \"maintenance_window\" instead. This attribute will be removed in the next major version of the provider.",
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("maintenance_window")),
				},
				Attributes: map[string]schema.Attribute{
					"minor_upgrade_maintenance_actual": schema.SingleNestedAttribute{
						Computed:   true,
						Attributes: maintenancesWindowSchemaActual,
					},
					"patch_upgrade_maintenance_actual": schema.SingleNestedAttribute{
						Computed:   true,
						Attributes: maintenancesWindowSchemaActual,
					},
					"minor_upgrade_maintenance": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: maintenancesWindowSchema,
					},
					"patch_upgrade_maintenance": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: maintenancesWindowSchema,
					},
				},
			},
			"maintenance_window": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("auto_maintenances")),
				},
				Attributes: map[string]schema.Attribute{
					"duration_hours": schema.Int64Attribute{
						Required: true,
					},
					"start_hour": schema.Int64Attribute{
						Required: true,
					},
					"week_day": schema.StringAttribute{
						Required: true,
					},
					"tz": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
				},
			},
			"cidr_pods": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validatorstring.IsCIDR(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr_service": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validatorstring.IsCIDR(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_dns": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					validatorstring.IsIP(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cni": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					validatorstring.IsIP(),
				},
			},
			"control_planes": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"cp_multi_az": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"cp_subregions": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"disable_api_termination": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 40),
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]$"),
						"Unique cluster name per project, must start with a letter and contain only lowercase letters, numbers, or hyphens.",
					),
				},
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
			"kubeconfig": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"tags": OKSTagsSchema(),
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

var maintenancesWindowSchema = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Optional: true,
	},
	"duration_hours": schema.Int64Attribute{
		Optional: true,
	},
	"start_hour": schema.Int64Attribute{
		Optional: true,
	},
	"week_day": schema.StringAttribute{
		Optional: true,
	},
	"tz": schema.StringAttribute{
		Optional: true,
	},
}

var maintenancesWindowSchemaActual = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Computed: true,
	},
	"duration_hours": schema.Int64Attribute{
		Computed: true,
	},
	"start_hour": schema.Int64Attribute{
		Computed: true,
	},
	"week_day": schema.StringAttribute{
		Computed: true,
	},
	"tz": schema.StringAttribute{
		Computed: true,
	},
}

func (r *oksClusterResource) expandAutoMaintenances(ctx context.Context, obj basetypes.ObjectValue) (auto oks.MaintenanceWindow, _ diag.Diagnostics) {
	if fwhelpers.IsSet(obj) {
		window, diags := to.Model[MaintenanceWindowModel](ctx, obj)
		if diags.HasError() {
			return auto, diags
		}

		if fwhelpers.IsSet(window.DurationHours) {
			hours := int(window.DurationHours.ValueInt64())
			auto.DurationHours = &hours
		}
		if fwhelpers.IsSet(window.Enabled) {
			auto.Enabled = window.Enabled.ValueBoolPointer()
		}
		if fwhelpers.IsSet(window.StartHour) {
			hour := int(window.StartHour.ValueInt64())
			auto.StartHour = &hour
		}
		if fwhelpers.IsSet(window.Tz) {
			auto.Tz = window.Tz.ValueStringPointer()
		}
		if fwhelpers.IsSet(window.WeekDay) {
			auto.WeekDay = (*oks.MaintenanceWindowWeekDay)(window.WeekDay.ValueStringPointer())
		}
	}
	return
}

func (r *oksClusterResource) expandMaintenances(ctx context.Context, obj basetypes.ObjectValue) (maintenance oks.Maintenance, _ diag.Diagnostics) {
	if fwhelpers.IsSet(obj) {
		maintenanceModel, diags := to.Model[MaintenanceModel](ctx, obj)
		if diags.HasError() {
			return maintenance, diags
		}

		if fwhelpers.IsSet(maintenanceModel.DurationHours) {
			maintenance.DurationHours = int(maintenanceModel.DurationHours.ValueInt64())
		}
		if fwhelpers.IsSet(maintenanceModel.StartHour) {
			maintenance.StartHour = int(maintenanceModel.StartHour.ValueInt64())
		}
		if fwhelpers.IsSet(maintenanceModel.Tz) {
			maintenance.Tz = maintenanceModel.Tz.ValueStringPointer()
		}
		if fwhelpers.IsSet(maintenanceModel.WeekDay) {
			maintenance.WeekDay = (oks.MaintenanceWeekDay)(maintenanceModel.WeekDay.ValueString())
		}
	}
	return
}

func (r *oksClusterResource) expandAdmissionFlags(ctx context.Context, data AdmissionFlagsModel) (oks.AdmissionFlagsInput, diag.Diagnostics) {
	var (
		admissionFlags oks.AdmissionFlagsInput
		diags          diag.Diagnostics
	)

	if fwhelpers.IsSet(data.DisableAdmissionPlugins) {
		disablePlugins, diag := to.Slice[string](ctx, data.DisableAdmissionPlugins)
		diags.Append(diag...)
		admissionFlags.DisableAdmissionPlugins = &disablePlugins
	}
	if fwhelpers.IsSet(data.EnableAdmissionPlugins) {
		enablePlugins, diag := to.Slice[string](ctx, data.EnableAdmissionPlugins)
		diags.Append(diag...)
		admissionFlags.EnableAdmissionPlugins = &enablePlugins
	}

	return admissionFlags, diags
}

func (r *oksClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterModel

	diag := req.Plan.Get(ctx, &plan)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	input := oks.ClusterInput{
		Name:      plan.Name.ValueString(),
		ProjectId: plan.ProjectId.ValueString(),
		Version:   plan.Version.ValueString(),
	}

	whitelist, diag := to.Slice[string](ctx, plan.AdminWhitelist)
	resp.Diagnostics.Append(diag...)
	input.AdminWhitelist = whitelist

	if fwhelpers.IsSet(plan.Description) {
		input.Description = plan.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(plan.CpMultiAz) {
		input.CpMultiAz = plan.CpMultiAz.ValueBoolPointer()
	}
	if fwhelpers.IsSet(plan.CpSubregions) {
		sub, diag := to.Slice[string](ctx, plan.CpSubregions)
		resp.Diagnostics.Append(diag...)
		input.CpSubregions = &sub
	}
	if fwhelpers.IsSet(plan.AdminLbu) {
		input.AdminLbu = plan.AdminLbu.ValueBoolPointer()
	}
	if fwhelpers.IsSet(plan.CidrPods) {
		input.CidrPods = plan.CidrPods.ValueString()
	}
	if fwhelpers.IsSet(plan.CidrService) {
		input.CidrService = plan.CidrService.ValueString()
	}
	if fwhelpers.IsSet(plan.ClusterDns) {
		input.ClusterDns = plan.ClusterDns.ValueStringPointer()
	}
	if fwhelpers.IsSet(plan.ControlPlanes) {
		input.ControlPlanes = plan.ControlPlanes.ValueStringPointer()
	}
	if fwhelpers.IsSet(plan.Quirks) {
		quirks, diag := to.Slice[string](ctx, plan.Quirks)
		resp.Diagnostics.Append(diag...)
		input.Quirks = &quirks
	}
	if fwhelpers.IsSet(plan.DisableApiTermination) {
		input.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}

	tags, diag := expandOKSTags(ctx, plan.OKSTagsModel)
	resp.Diagnostics.Append(diag...)
	input.Tags = &tags

	if fwhelpers.IsSet(plan.AutoMaintenances) {
		autoInput := oks.AutoMaintenances{}
		auto, diag := to.Model[AutoMaintenancesModel](ctx, plan.AutoMaintenances)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		minor, diag := r.expandAutoMaintenances(ctx, auto.MinorUpgradeMaintenance)
		resp.Diagnostics.Append(diag...)
		patch, diag := r.expandAutoMaintenances(ctx, auto.PatchUpgradeMaintenance)
		resp.Diagnostics.Append(diag...)

		autoInput.MinorUpgradeMaintenance = minor
		autoInput.PatchUpgradeMaintenance = patch
		input.AutoMaintenances = &autoInput
	}
	if fwhelpers.IsSet(plan.MaintenanceWindow) {
		maintenance, diag := r.expandMaintenances(ctx, plan.MaintenanceWindow)
		resp.Diagnostics.Append(diag...)

		input.MaintenanceWindow = &maintenance
	}
	if fwhelpers.IsSet(plan.AdmissionFlags) {
		admissionModel, diag := to.Model[AdmissionFlagsModel](ctx, plan.AdmissionFlags)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		admissionFlags, diag := r.expandAdmissionFlags(ctx, admissionModel)
		resp.Diagnostics.Append(diag...)
		input.AdmissionFlags = &admissionFlags
	}

	if resp.Diagnostics.HasError() {
		return
	}
	createResp, err := r.Client.CreateCluster(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Cluster",
			"Error: "+err.Error(),
		)
		return
	}
	plan.RequestId = to.String(createResp.ResponseContext.RequestId)
	plan.Id = to.String(createResp.Cluster.Id)

	to, diag := plan.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	data, err := r.read(ctx, plan, to)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diag = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diag...)
}

func (r *oksClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterModel
	diag := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	to, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	data, err := r.read(ctx, data, to)
	if err != nil {
		if oks.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diag = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diag...)
}

func (r *oksClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state ClusterModel
		updateReq   oks.ClusterUpdate
		doUpgrade   bool
	)
	diag := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diag...)

	diag = req.State.Get(ctx, &state)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	statuses, diag := to.Model[StatusesModel](ctx, state.Statuses)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	if statuses.Status.ValueString() == "deleting" {
		resp.Diagnostics.AddError(
			"Unable to Update Cluster",
			"Cluster is currently being deleted and cannot be updated. The resource will be removed from the state once deleted.",
		)
		return
	}

	if fwhelpers.IsSet(plan.AdminWhitelist) && !plan.AdminWhitelist.Equal(state.AdminWhitelist) {
		admins, diag := to.Slice[string](ctx, plan.AdminWhitelist)
		resp.Diagnostics.Append(diag...)
		updateReq.AdminWhitelist = &admins
	}
	if !plan.AdmissionFlags.Equal(state.AdmissionFlags) {
		updateReq.AdmissionFlags = &oks.AdmissionFlagsInput{}

		if fwhelpers.IsSet(plan.AdmissionFlags) {
			planAdmission, diag := to.Model[AdmissionFlagsModel](ctx, plan.AdmissionFlags)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}
			stateAdmission, diag := to.Model[AdmissionFlagsModel](ctx, state.AdmissionFlags)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}

			if planAdmission.DisableAdmissionPlugins.IsNull() && !planAdmission.DisableAdmissionPlugins.Equal(stateAdmission.DisableAdmissionPlugins) {
				updateReq.AdmissionFlags.DisableAdmissionPlugins = &[]string{}
			} else {
				disablePlugins, diag := to.Slice[string](ctx, planAdmission.DisableAdmissionPlugins)
				resp.Diagnostics.Append(diag...)
				updateReq.AdmissionFlags.DisableAdmissionPlugins = &disablePlugins
			}
			if planAdmission.EnableAdmissionPlugins.IsNull() && !planAdmission.EnableAdmissionPlugins.Equal(stateAdmission.EnableAdmissionPlugins) {
				updateReq.AdmissionFlags.EnableAdmissionPlugins = &[]string{}
			} else {
				enablePlugins, diag := to.Slice[string](ctx, planAdmission.EnableAdmissionPlugins)
				resp.Diagnostics.Append(diag...)
				updateReq.AdmissionFlags.EnableAdmissionPlugins = &enablePlugins
			}
		} else {
			// Send empty slice to reset to default
			updateReq.AdmissionFlags.DisableAdmissionPlugins = &[]string{}
			updateReq.AdmissionFlags.EnableAdmissionPlugins = &[]string{}
		}
		// Set new optional value to state, since Read does not write it
		state.AdmissionFlags = plan.AdmissionFlags
	}
	if !plan.AutoMaintenances.Equal(state.AutoMaintenances) {
		updateReq.AutoMaintenances = &oks.AutoMaintenances{}

		if fwhelpers.IsSet(plan.AutoMaintenances) {
			stateAuto, diag := to.Model[AutoMaintenancesModel](ctx, state.AutoMaintenances)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}
			planAuto, diag := to.Model[AutoMaintenancesModel](ctx, plan.AutoMaintenances)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}

			if planAuto.MinorUpgradeMaintenance.IsNull() && !planAuto.MinorUpgradeMaintenance.Equal(stateAuto.MinorUpgradeMaintenance) {
				updateReq.AutoMaintenances.MinorUpgradeMaintenance = oks.MaintenanceWindow{}
			} else {
				minor, diag := r.expandAutoMaintenances(ctx, planAuto.MinorUpgradeMaintenance)
				resp.Diagnostics.Append(diag...)
				updateReq.AutoMaintenances.MinorUpgradeMaintenance = minor
			}
			if planAuto.PatchUpgradeMaintenance.IsNull() && !planAuto.PatchUpgradeMaintenance.Equal(stateAuto.PatchUpgradeMaintenance) {
				updateReq.AutoMaintenances.PatchUpgradeMaintenance = oks.MaintenanceWindow{}
			} else {
				patch, diag := r.expandAutoMaintenances(ctx, planAuto.PatchUpgradeMaintenance)
				resp.Diagnostics.Append(diag...)
				updateReq.AutoMaintenances.PatchUpgradeMaintenance = patch
			}
		} else {
			// Send empty struct to reset to default value
			updateReq.AutoMaintenances.MinorUpgradeMaintenance = oks.MaintenanceWindow{}
			updateReq.AutoMaintenances.PatchUpgradeMaintenance = oks.MaintenanceWindow{}
		}
		// Set new optional value to state, since Read does not write it
		state.AutoMaintenances = plan.AutoMaintenances
	}
	if fwhelpers.IsSet(plan.MaintenanceWindow) && !plan.MaintenanceWindow.Equal(state.MaintenanceWindow) {
		maintenance, diag := r.expandMaintenances(ctx, plan.MaintenanceWindow)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.MaintenanceWindow = &maintenance
	}
	if fwhelpers.IsSet(plan.ControlPlanes) && !plan.ControlPlanes.Equal(state.ControlPlanes) {
		updateReq.ControlPlanes = plan.ControlPlanes.ValueStringPointer()
		doUpgrade = true
	}
	if fwhelpers.IsSet(plan.Description) && !plan.Description.Equal(state.Description) {
		updateReq.Description = plan.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(plan.DisableApiTermination) && !plan.DisableApiTermination.Equal(state.DisableApiTermination) {
		updateReq.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}
	if fwhelpers.IsSet(plan.Quirks) && !plan.Quirks.Equal(state.Quirks) {
		quirks, diag := to.Slice[string](ctx, plan.Quirks)
		resp.Diagnostics.Append(diag...)
		updateReq.Quirks = &quirks
	}
	if fwhelpers.IsSet(plan.Version) && !plan.Version.Equal(state.Version) {
		updateReq.Version = plan.Version.ValueStringPointer()
		doUpgrade = true
	}
	tags, diag := cmpOKSTags(ctx, plan.OKSTagsModel, state.OKSTagsModel)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	if tags != nil {
		updateReq.Tags = &tags
	}

	updateResp, err := r.Client.UpdateCluster(ctx, state.Id.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Cluster",
			"Error: "+err.Error(),
		)
		return
	}
	state.RequestId = to.String(updateResp.ResponseContext.RequestId)

	timeout, diag := state.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	if doUpgrade {
		_, err := r.waitForClusterState(ctx, state.Id.ValueString(), []string{"pending", "updating"}, []string{"ready"}, timeout)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to wait for Cluster to complete updating.", "Error: "+err.Error(),
			)
			return
		}
		upgradeResp, err := r.Client.UpgradeCluster(ctx, state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Upgrade Cluster",
				"Error: "+err.Error(),
			)
			return
		}
		state.RequestId = to.String(upgradeResp.ResponseContext.RequestId)

		timeout, diag = state.Timeouts.Update(ctx, UpgradeDefaultTimeout)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		tflog.Info(ctx, fmt.Sprintf("Cluster upgrading. Timeout is set at %s.", timeout.String()))
	}

	data, err := r.read(ctx, state, timeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Cluster state",
			"Error: "+err.Error(),
		)
		return
	}
	diag = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diag...)
}

func (r *oksClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterModel

	diags := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
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

	to, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	_, err = r.waitForClusterState(ctx, data.Id.ValueString(), []string{"pending", "deleting"}, []string{}, to)
	if err != nil {
		if !oks.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to wait for Cluster complete deletion.", "Error: "+err.Error())
		}
	}
}

func (r *oksClusterResource) waitForClusterState(ctx context.Context, id string, pending []string, target []string, timeout time.Duration) (*oks.ClusterResponse, error) {
	resp, err := fwhelpers.WaitForResource[oks.ClusterResponse](ctx, &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			resp, err := r.Client.GetCluster(ctx, id)
			if err != nil {
				return resp, "", err
			}
			if resp.Cluster.Statuses.Status == nil {
				return resp, "", nil
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

func (r *oksClusterResource) setAdmissionFlags(ctx context.Context, data *ClusterModel, auto oks.AdmissionFlags) diag.Diagnostics {
	var diags diag.Diagnostics

	var model AdmissionFlagsModel
	applied, diag := types.SetValueFrom(ctx, types.StringType, auto.AppliedAdmissionPlugins)
	diags.Append(diag...)

	disable, diag := types.SetValueFrom(ctx, types.StringType, auto.DisableAdmissionPlugins)
	diags.Append(diag...)

	enable, diag := types.SetValueFrom(ctx, types.StringType, auto.EnableAdmissionPlugins)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	model.DisableAdmissionPlugins = types.SetValueMust(types.StringType, []attr.Value{})
	model.EnableAdmissionPlugins = types.SetValueMust(types.StringType, []attr.Value{})
	if fwhelpers.IsSet(data.AdmissionFlags) {
		stateModel, diag := to.Model[AdmissionFlagsModel](ctx, data.AdmissionFlags)
		diags.Append(diag...)
		if diag.HasError() {
			return diags
		}
		if fwhelpers.IsSet(stateModel.EnableAdmissionPlugins) {
			model.EnableAdmissionPlugins = stateModel.EnableAdmissionPlugins
		}
		if fwhelpers.IsSet(stateModel.DisableAdmissionPlugins) {
			model.DisableAdmissionPlugins = stateModel.DisableAdmissionPlugins
		}
	}

	// Set computed values
	model.AppliedAdmissionPlugins = applied
	model.DisableAdmissionPluginsActual = disable
	model.EnableAdmissionPluginsActual = enable

	obj, diag := types.ObjectValueFrom(ctx, data.AdmissionFlags.AttributeTypes(ctx), model)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}
	data.AdmissionFlags = obj

	return diags
}

func (r *oksClusterResource) setAutoMaintenances(ctx context.Context, data *ClusterModel, auto *oks.AutoMaintenances) diag.Diagnostics {
	var diags diag.Diagnostics
	if auto == nil {
		return diags
	}

	setMaintenanceWindow := func(window oks.MaintenanceWindow) (types.Object, diag.Diagnostics) {
		var windowModel MaintenanceWindowModel
		windowModel.DurationHours = to.Int64(window.DurationHours)
		windowModel.Enabled = to.Bool(window.Enabled)
		windowModel.StartHour = to.Int64(window.StartHour)
		windowModel.Tz = to.String(window.Tz)
		windowModel.WeekDay = to.String((*string)(window.WeekDay))

		return types.ObjectValueFrom(ctx, maintenanceWindowAttrTypes, windowModel)
	}
	var planModel AutoMaintenancesModel

	minor, diag := setMaintenanceWindow(auto.MinorUpgradeMaintenance)
	diags.Append(diag...)
	patch, diag := setMaintenanceWindow(auto.PatchUpgradeMaintenance)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	planModel.MinorUpgradeMaintenance = types.ObjectNull(maintenanceWindowAttrTypes)
	planModel.PatchUpgradeMaintenance = types.ObjectNull(maintenanceWindowAttrTypes)
	if fwhelpers.IsSet(data.AutoMaintenances) {
		stateModel, diag := to.Model[AutoMaintenancesModel](ctx, data.AutoMaintenances)
		diags.Append(diag...)
		if diags.HasError() {
			return diags
		}
		if fwhelpers.IsSet(stateModel.MinorUpgradeMaintenance) {
			planModel.MinorUpgradeMaintenance = stateModel.MinorUpgradeMaintenance
		}
		if fwhelpers.IsSet(stateModel.PatchUpgradeMaintenance) {
			planModel.PatchUpgradeMaintenance = stateModel.PatchUpgradeMaintenance
		}
	}

	// Set computed values
	planModel.MinorUpgradeMaintenanceActual = minor
	planModel.PatchUpgradeMaintenanceActual = patch

	obj, diag := types.ObjectValueFrom(ctx, data.AutoMaintenances.AttributeTypes(ctx), &planModel)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}
	data.AutoMaintenances = obj

	return diags
}

func (r *oksClusterResource) setMaintenance(ctx context.Context, data *ClusterModel, maintenance *oks.Maintenance) diag.Diagnostics {
	data.MaintenanceWindow = types.ObjectNull(maintenanceAttrTypes)
	if maintenance == nil {
		return nil
	}

	maintenanceModel := MaintenanceModel{
		DurationHours: to.Int64(maintenance.DurationHours),
		StartHour:     to.Int64(maintenance.StartHour),
		Tz:            to.String(maintenance.Tz),
		WeekDay:       to.String((string)(maintenance.WeekDay)),
	}

	obj, diag := types.ObjectValueFrom(ctx, maintenanceAttrTypes, maintenanceModel)
	if diag.HasError() {
		return diag
	}
	data.MaintenanceWindow = obj

	return nil
}

func (r *oksClusterResource) setStatuses(ctx context.Context, data *ClusterModel, auto oks.Statuses) diag.Diagnostics {
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

func (r *oksClusterResource) read(ctx context.Context, data ClusterModel, timeout time.Duration) (ClusterModel, error) {
	resp, err := r.waitForClusterState(ctx, data.Id.ValueString(), []string{"pending", "deploying", "updating", "upgrading"}, []string{"ready", "deleting"}, timeout)
	if err != nil {
		return data, err
	}
	cluster := resp.Cluster

	adminWhiteList, diags := types.SetValueFrom(ctx, types.StringType, cluster.AdminWhitelist)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert admin_whitelist into a set: %v", diags.Errors())
	}
	cpSubregions, diags := types.SetValueFrom(ctx, types.StringType, cluster.CpSubregions)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert cp_subregions into a set: %v", diags.Errors())
	}
	diags = r.setAdmissionFlags(ctx, &data, cluster.AdmissionFlags)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert admission_flags into the schema model: %v", diags.Errors())
	}
	diags = r.setAutoMaintenances(ctx, &data, cluster.AutoMaintenances)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert auto_maintenances into the schema model: %v", diags.Errors())
	}
	diags = r.setMaintenance(ctx, &data, cluster.MaintenanceWindow)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert maintenance_window into the schema model: %v", diags.Errors())
	}
	diags = r.setStatuses(ctx, &data, cluster.Statuses)
	if diags.HasError() {
		return data, fmt.Errorf("unable to convert statuses into the schema model: %v", diags.Errors())
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
	data.Id = to.String(cluster.Id)
	data.Name = to.String(cluster.Name)
	data.ProjectId = to.String(cluster.ProjectId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Version = to.String(cluster.Version)

	kubeconfigResp, err := r.Client.GetKubeconfig(ctx, cluster.Id, nil)
	if err != nil {
		return data, err
	}
	data.Kubeconfig = to.String(kubeconfigResp.Cluster.Data.Kubeconfig)

	tags, diags := flattenOKSTags(ctx, cluster.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("%v", diags.Errors())
	}
	data.Tags = tags

	return data, nil
}
