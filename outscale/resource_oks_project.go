package outscale

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	sdkv3_oks "github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/fwvalidators"
	"github.com/outscale/terraform-provider-outscale/utils"
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
	Tags                  types.Map         `tfsdk:"tags"`
	RequestId             types.String      `tfsdk:"request_id"`
	Timeouts              timeouts.Value    `tfsdk:"timeouts"`
}

type oksProjectResource struct {
	Client *sdkv3_oks.Client
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
			"tags": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *oksProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := sdkv3_oks.ProjectInput{
		Cidr:   data.Cidr.ValueString(),
		Name:   data.Name.ValueString(),
		Region: data.Region.ValueString(),
	}

	if !data.Description.IsUnknown() {
		input.Description = data.Description.ValueStringPointer()
	}
	if !data.DisableApiTermination.IsUnknown() {
		input.DisableApiTermination = data.DisableApiTermination.ValueBoolPointer()
	}
	if !data.Quirks.IsUnknown() && !data.Quirks.IsNull() {
		var quirks []string
		resp.Diagnostics.Append(data.Quirks.ElementsAs(ctx, quirks, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Quirks = &quirks
	}

	// TODO: make a generic function for oks tags
	var tags map[string]string
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
	if resp.Diagnostics.HasError() {
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
	data.RequestId = types.StringValue(createResp.ResponseContext.RequestId)
	data.Id = types.StringValue(createResp.Project.Id)

	to, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oksProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	to, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.setOKSProjectState(ctx, data, to)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Project state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oksProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state ProjectModel
		update      sdkv3_oks.ProjectUpdate
	)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Description.Equal(state.Description) {
		update.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DisableApiTermination.Equal(state.DisableApiTermination) {
		update.DisableApiTermination = plan.DisableApiTermination.ValueBoolPointer()
	}
	if !plan.Quirks.Equal(state.Quirks) {
		var quirks []string
		resp.Diagnostics.Append(plan.Quirks.ElementsAs(ctx, quirks, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		update.Quirks = &quirks
	}
	if !plan.Tags.Equal(state.Tags) {
		var tags map[string]string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		update.Tags = &tags
	}
	updateResp, err := r.Client.UpdateProject(ctx, state.Id.ValueString(), update)
	if err != nil {
		return
	}
	state.RequestId = types.StringValue(updateResp.ResponseContext.RequestId)

	to, diags := state.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oksProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
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

	to, diags := data.Timeouts.Update(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err = r.waitForProjectState(ctx, data.Id.ValueString(), []string{"deleting"}, []string{}, to)
	if err != nil {
		return
	}
}

func (r *oksProjectResource) waitForProjectState(ctx context.Context, id string, pending []string, target []string, timeout time.Duration) (*sdkv3_oks.ProjectResponse, error) {
	resp, err := utils.WaitForResource[sdkv3_oks.ProjectResponse](ctx, &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			resp, err := r.Client.GetProject(ctx, id)
			if err != nil {
				return resp, "", errors.Join(fmt.Errorf("Error waiting for Project (%s) to become ready.",
					id), err)
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
	resp, err := r.waitForProjectState(ctx, data.Id.ValueString(), []string{"pending", "updating"}, []string{"ready"}, timeout)
	if err != nil {
		return data, err
	}
	project := resp.Project

	data.Cidr = types.StringValue(project.Cidr)
	data.CreatedAt = timetypes.NewRFC3339TimeValue(project.CreatedAt)
	data.Description = types.StringPointerValue(project.Description)
	data.DisableApiTermination = types.BoolPointerValue(project.DisableApiTermination)
	data.Name = types.StringValue(project.Name)
	data.Region = types.StringValue(project.Region)
	data.Status = types.StringValue(project.Status)
	data.UpdatedAt = timetypes.NewRFC3339TimeValue(project.UpdatedAt)
	data.Id = types.StringValue(project.Id)
	data.RequestId = types.StringValue(resp.ResponseContext.RequestId)

	tags, diags := types.MapValueFrom(ctx, types.StringType, project.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Tags to the schema Set. Error: %v: ", diags.Errors())
	}
	data.Tags = tags

	return data, nil
}
