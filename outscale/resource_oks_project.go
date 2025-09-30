package outscale

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/fwvalidators"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/to"
)

var (
	_ resource.Resource              = &oksProjectResource{}
	_ resource.ResourceWithConfigure = &oksProjectResource{}
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
					fwvalidators.IsCIDR(),
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
	if utils.CheckDiags(resp, diags) {
		return
	}

	input := oks.ProjectInput{
		Cidr:   data.Cidr.ValueString(),
		Name:   data.Name.ValueString(),
		Region: data.Region.ValueString(),
	}

	if utils.IsSet(data.Description) {
		input.Description = data.Description.ValueStringPointer()
	}
	if utils.IsSet(data.DisableApiTermination) {
		input.DisableApiTermination = data.DisableApiTermination.ValueBoolPointer()
	}
	if utils.IsSet(data.Quirks) {
		quirks, diags := to.Slice[string](ctx, data.Quirks)
		resp.Diagnostics.Append(diags...)
		input.Quirks = &quirks
	}

	tags, diags := expandOKSTags(ctx, data.OKSTagsModel)
	if utils.CheckDiags(resp, diags) {
		return
	}
	input.Tags = &tags

	createResp, err := r.Client.CreateProject(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource Project",
			"Error: "+err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(createResp.Project.Id)

	to, diags := data.Timeouts.Create(ctx, utils.CreateOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}

	data, err = r.setOKSProjectState(ctx, data, to)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Project state",
			"Error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *oksProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectModel

	diags := req.State.Get(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}

	to, diags := data.Timeouts.Read(ctx, utils.ReadOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	data, err := r.setOKSProjectState(ctx, data, to)
	if err != nil {
		if code := oks.StatusCodeHelper(err); code != nil && *code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Project state",
			"Error: "+err.Error(),
		)
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
	if utils.CheckDiags(resp, diags) {
		return
	}

	if state.Status.ValueString() == "deleting" {
		resp.Diagnostics.AddError(
			"Unable to Update Project",
			"Project is currently being deleted and cannot be updated. The resource will be removed from the state once deleted.",
		)
		return
	}

	if utils.IsSet(plan.Description) && !plan.Description.Equal(state.Description) {
		update.Description = plan.Description.ValueStringPointer()
	}
	if utils.IsSet(plan.DisableApiTermination) && !plan.DisableApiTermination.Equal(state.DisableApiTermination) {
		update.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}
	if utils.IsSet(plan.Quirks) && !plan.Quirks.Equal(state.Quirks) {
		quirks, diags := to.Slice[string](ctx, plan.Quirks)
		resp.Diagnostics.Append(diags...)
		update.Quirks = &quirks
	}
	tags, diags := cmpOKSTags(ctx, plan.OKSTagsModel, state.OKSTagsModel)
	if utils.CheckDiags(resp, diags) {
		return
	}
	if tags != nil {
		update.Tags = &tags
	}

	updateResp, err := r.Client.UpdateProject(ctx, state.Id.ValueString(), update)
	if err != nil {
		return
	}
	state.RequestId = to.String(updateResp.ResponseContext.RequestId)

	to, diags := state.Timeouts.Update(ctx, utils.UpdateOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	data, err := r.setOKSProjectState(ctx, state, to)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Project state",
			"Error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *oksProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

	diags := req.State.Get(ctx, &data)
	if utils.CheckDiags(resp, diags) {
		return
	}
	_, err := r.Client.DeleteProject(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Resource Project",
			"Error: "+err.Error(),
		)
		return
	}

	to, diags := data.Timeouts.Update(ctx, utils.DeleteOKSDefaultTimeout)
	if utils.CheckDiags(resp, diags) {
		return
	}
	_, err = r.waitForProjectState(ctx, data.Id.ValueString(), []string{"deleting"}, []string{}, to)
	if err != nil {
		if code := oks.StatusCodeHelper(err); code != nil && *code != 404 {
			resp.Diagnostics.AddError("Unable to wait for Project complete deletion.", "Error: "+err.Error())
		}
	}
}

func (r *oksProjectResource) waitForProjectState(ctx context.Context, id string, pending []string, target []string, timeout time.Duration) (*oks.ProjectResponse, error) {
	resp, err := utils.WaitForResource[oks.ProjectResponse](ctx, &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			resp, err := r.Client.GetProject(ctx, id)
			if err != nil {
				return resp, "", err
			}
			return resp, resp.Project.Status, nil
		},
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *oksProjectResource) setOKSProjectState(ctx context.Context, data ProjectModel, timeout time.Duration) (ProjectModel, error) {
	resp, err := r.waitForProjectState(ctx, data.Id.ValueString(), []string{"pending", "updating"}, []string{"ready", "deleting"}, timeout)
	if err != nil {
		return data, err
	}
	project := resp.Project

	data.Cidr = to.String(project.Cidr)
	data.CreatedAt = to.RFC3339(project.CreatedAt)
	data.Description = to.String(project.Description)
	data.DisableApiTermination = to.Bool(project.DisableApiTermination)
	data.Name = to.String(project.Name)
	data.Region = to.String(project.Region)
	data.Status = to.String(project.Status)
	data.UpdatedAt = to.RFC3339(project.UpdatedAt)
	data.Id = to.String(project.Id)
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	tags, diags := flattenOKSTags(ctx, project.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("%v", diags.Errors())
	}
	data.Tags = tags

	return data, nil
}
