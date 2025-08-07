package outscale

import (
	"context"
	"fmt"
	"time"

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
	DeletedAt             timetypes.RFC3339 `tfsdk:"deleted_at"`
	Description           types.String      `tfsdk:"description"`
	DisableApiTermination types.Bool        `tfsdk:"disable_api_termination"`
	Name                  types.String      `tfsdk:"name"`
	Region                types.String      `tfsdk:"region"`
	Status                types.String      `tfsdk:"status"`
	UpdatedAt             timetypes.RFC3339 `tfsdk:"updated_at"`
	Quirks                types.Set         `tfsdk:"quirks"`
	Id                    types.String      `tfsdk:"id"`
	Tags                  types.Set         `tfsdk:"tags"`
	RequestId             types.String      `tfsdk:"request_id"`
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
	// net_id := req.ID
	// if net_id == "" {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Import Identifier",
	// 		fmt.Sprintf("Expected import net_resource identifier Got: %v", req.ID),
	// 	)
	// 	return
	// }

	// var data ProjectModel
	// var timeouts timeouts.Value
	// data.NetId = types.StringValue(net_id)
	// resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// data.Timeouts = timeouts
	// diags := resp.State.Set(ctx, &data)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

func (r *oksProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_project"
}

func (r *oksProjectResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags": TagsSchema(),
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
			},
			"deleted_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 40),
				},
			},
			"disable_api_termination": schema.BoolAttribute{
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

	createResp, err := r.Client.CreateProject(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource Project",
			"Error: "+err.Error(),
		)
		return
	}

	data.RequestId = types.StringValue(createResp.ResponseContext.RequestId)
	// if len(createResp.Project.Tags) > 0 {
	// 	err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(data.Tags), net.GetNetId())
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable to add Tags on outscale_net resource",
	// 			"Error: "+utils.GetErrorResponse(err).Error(),
	// 		)
	// 		return
	// 	}
	// }

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"ready"},
		Refresh: func() (any, string, error) {
			resp, err := r.Client.GetProject(ctx, createResp.Project.Id)
			if err != nil {
				return resp, "error", err
			}
			return resp, resp.Project.Status, nil
		},
		Timeout:    utils.CreateDefaultTimeout,
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error waiting for Project (%s) to become ready.",
				createResp.Project.Id),
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(createResp.Project.Id)
	data, err = setOKSProjectState(ctx, r, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set OKS state",
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

	// data, err = setNetState(ctx, r, data)
	// if err != nil {
	// 	if err.Error() == "Empty" {
	// 		resp.State.RemoveResource(ctx)
	// 		return
	// 	}
	// 	resp.Diagnostics.AddError(
	// 		"Unable to set net state",
	// 		"Error: "+err.Error(),
	// 	)
	// 	return
	// }
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oksProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// var (
	// 	tagsPlan, tagsState []ResourceTag
	// 	resourceId          types.String
	// 	err                 error
	// )

	// resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tags"), &tagsPlan)...)
	// resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tags"), &tagsState)...)
	// resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("net_id"), &resourceId)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// if !reflect.DeepEqual(tagsPlan, tagsState) {
	// 	toRemove, toCreate := diffOSCAPITags(tagsToOSCResourceTag(tagsPlan), tagsToOSCResourceTag(tagsState))
	// 	err := updateFrameworkTags(ctx, r.Client, toCreate, toRemove, resourceId.ValueString())
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable to update Tags on net resource",
	// 			"Error: "+utils.GetErrorResponse(err).Error(),
	// 		)
	// 		return
	// 	}
	// }
	// var data ProjectModel
	// resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// data, err = setNetState(ctx, r, data)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to set net state",
	// 		"Error: "+err.Error(),
	// 	)
	// 	return
	// // }
	// resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

func (r *oksProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// var data ProjectModel

	// resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// delReq := oscgo.DeleteNetRequest{
	// 	NetId: data.NetId.ValueString(),
	// }
	// err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
	// 	_, httpResp, err := r.Client.NetApi.DeleteNet(ctx).DeleteNetRequest(delReq).Execute()

	// 	if err != nil {
	// 		return utils.CheckThrottling(httpResp, err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to Delete net",
	// 		"Error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func setOKSProjectState(ctx context.Context, r *oksProjectResource, data ProjectModel) (ProjectModel, error) {
	resp, err := r.Client.GetProject(ctx, data.Id.ValueString())
	if err != nil {
		return data, err
	}
	project := resp.Project

	data.Cidr = types.StringValue(project.Cidr)
	data.CreatedAt = timetypes.NewRFC3339TimeValue(project.CreatedAt)
	data.DeletedAt = timetypes.NewRFC3339TimePointerValue(project.DeletedAt)
	data.Description = types.StringPointerValue(project.Description)
	data.DisableApiTermination = types.BoolPointerValue(project.DisableApiTermination)
	data.Name = types.StringValue(project.Name)
	data.Region = types.StringValue(project.Region)
	data.Status = types.StringValue(project.Status)
	data.UpdatedAt = timetypes.NewRFC3339TimeValue(project.UpdatedAt)
	data.Id = types.StringValue(project.Id)
	// data.Tags =
	data.RequestId = types.StringValue(resp.ResponseContext.RequestId)

	return data, nil
}
