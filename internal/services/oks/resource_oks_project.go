package oks

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
)

var (
	_ resource.Resource              = &oksProjectResource{}
	_ resource.ResourceWithConfigure = &oksProjectResource{}
)

const (
	projectErrCreate = "Unable to create OKS Project"
	projectErrUpdate = "Unable to update OKS Project"
	projectErrDelete = "Unable to delete OKS Project"
	projectErrWait   = "Unable to wait for OKS Project state"
)

type ProjectModel struct {
	Cidr                  types.String      `tfsdk:"cidr"`
	CreatedAt             timetypes.RFC3339 `tfsdk:"created_at"`
	Description           types.String      `tfsdk:"description"`
	DisableApiTermination types.Bool        `tfsdk:"disable_api_termination"`
	Name                  types.String      `tfsdk:"name"`
	Region                types.String      `tfsdk:"region"`
	Status                types.String      `tfsdk:"status"`
	UpdatedAt             timetypes.RFC3339 `tfsdk:"updated_at"`
	Quirks                types.Set         `tfsdk:"quirks"`
	Id                    types.String      `tfsdk:"id"`
	RequestId             types.String      `tfsdk:"request_id"`
	Timeouts              timeouts.Value    `tfsdk:"timeouts"`
	OKSTagsModel
}

type oksProjectResource struct {
	Client *oks.Client
}

func NewResourceProject() resource.Resource {
	return &oksProjectResource{}
}

func (r *oksProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *oksProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *oksProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_project"
}

func (r *oksProjectResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"cidr": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validatorstring.IsCIDR(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 40),
				},
			},
			"disable_api_termination": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 40),
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]$"),
						"Unique name for the project, must start with a letter and contain only lowercase letters, numbers, or hyphens.",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
			"quirks": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
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
			"tags": OKSTagsSchema(),
		},
	}
}

func (r *oksProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectModel
	diags := req.Plan.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	input := oks.ProjectInput{
		Cidr:   data.Cidr.ValueString(),
		Name:   data.Name.ValueString(),
		Region: data.Region.ValueString(),
	}

	if fwhelpers.IsSet(data.Description) {
		input.Description = data.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.DisableApiTermination) {
		input.DisableApiTermination = data.DisableApiTermination.ValueBoolPointer()
	}
	if fwhelpers.IsSet(data.Quirks) {
		quirks, diags := to.Slice[string](ctx, data.Quirks)
		resp.Diagnostics.Append(diags...)
		input.Quirks = &quirks
	}

	tags, diags := expandOKSTags(ctx, data.OKSTagsModel)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	input.Tags = &tags

	createResp, err := r.Client.CreateProject(ctx, input, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(projectErrCreate, err.Error())
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(createResp.Project.Id)

	stateData, err := r.flatten(ctx, data, createResp.Project)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	// A read call is necessary after the flatten as it waits for the project to be ready
	data, err = r.read(ctx, data, timeout)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *oksProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectModel

	diags := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	to, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	data, err := r.read(ctx, data, to)
	if err != nil {
		if oks.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *oksProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state ProjectModel
		update      oks.ProjectUpdate
	)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if state.Status.ValueString() == "deleting" {
		resp.Diagnostics.AddError(
			projectErrUpdate,
			"Project is currently being deleted and cannot be updated. The resource will be removed from the state once deleted.",
		)
		return
	}

	if fwhelpers.IsSet(plan.Description) && !plan.Description.Equal(state.Description) {
		update.Description = plan.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(plan.DisableApiTermination) && !plan.DisableApiTermination.Equal(state.DisableApiTermination) {
		update.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}
	if fwhelpers.IsSet(plan.Quirks) && !plan.Quirks.Equal(state.Quirks) {
		quirks, diags := to.Slice[string](ctx, plan.Quirks)
		resp.Diagnostics.Append(diags...)
		update.Quirks = &quirks
	}
	tags, diags := cmpOKSTags(ctx, plan.OKSTagsModel, state.OKSTagsModel)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	if tags != nil {
		update.Tags = &tags
	}

	if update != (oks.ProjectUpdate{}) {
		updateResp, err := r.Client.UpdateProject(ctx, state.Id.ValueString(), update, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(projectErrUpdate, err.Error())
			return
		}

		state.RequestId = to.String(updateResp.ResponseContext.RequestId)
	}

	state.Timeouts = plan.Timeouts
	state.Quirks = plan.Quirks

	data, err := r.read(ctx, state, timeout)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *oksProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

	diags := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	_, err := r.Client.DeleteProject(ctx, data.Id.ValueString(), options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(projectErrDelete, err.Error())
		return
	}

	_, err = r.waitForProjectState(ctx, data.Id.ValueString(), stateconf.States(oks.ProjectStatusDeleting), []oks.ProjectStatus{}, timeout)
	if err != nil {
		if !oks.IsNotFound(err) {
			resp.Diagnostics.AddError(projectErrWait, err.Error())
		}
	}
}

func (r *oksProjectResource) waitForProjectState(ctx context.Context, id string, pending []oks.ProjectStatus, target []oks.ProjectStatus, timeout time.Duration) (*oks.ProjectResponse, error) {
	conf := stateconf.StateChangeConf[oks.ProjectStatus]{
		Pending: pending,
		Target:  target,
		Timeout: timeout,
		Refresh: func(ctx context.Context) (any, oks.ProjectStatus, error) {
			resp, err := r.Client.GetProject(ctx, id)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Project.Status, nil
		},
	}
	respAny, err := conf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	return respAny.(*oks.ProjectResponse), nil
}

func (r *oksProjectResource) read(ctx context.Context, data ProjectModel, timeout time.Duration) (ProjectModel, error) {
	resp, err := r.waitForProjectState(ctx, data.Id.ValueString(), stateconf.States(oks.ProjectStatusPending, oks.ProjectStatusUpdating), stateconf.States(oks.ProjectStatusReady, oks.ProjectStatusDeleting), timeout)
	if err != nil {
		return data, err
	}

	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return r.flatten(ctx, data, resp.Project)
}

func (r *oksProjectResource) flatten(ctx context.Context, data ProjectModel, project oks.Project) (ProjectModel, error) {
	data.Cidr = to.String(project.Cidr)
	data.CreatedAt = to.RFC3339(project.CreatedAt)
	data.Description = to.String(project.Description)
	data.DisableApiTermination = to.Bool(project.DisableApiTermination)
	data.Name = to.String(project.Name)
	data.Region = to.String(project.Region)
	data.Status = to.String(project.Status)
	data.UpdatedAt = to.RFC3339(project.UpdatedAt)
	data.Id = to.String(project.Id)

	tags, diags := flattenOKSTags(ctx, project.Tags)
	if diags.HasError() {
		return data, from.Diag(diags)
	}
	data.Tags = tags

	return data, nil
}
