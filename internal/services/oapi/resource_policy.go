package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
)

var (
	_ resource.Resource                = &resourcePolicy{}
	_ resource.ResourceWithConfigure   = &resourcePolicy{}
	_ resource.ResourceWithImportState = &resourcePolicy{}
	_ resource.ResourceWithModifyPlan  = &resourcePolicy{}
)

type PolicyModel struct {
	CreationDate           types.String   `tfsdk:"creation_date"`
	Description            types.String   `tfsdk:"description"`
	IsLinkable             types.Bool     `tfsdk:"is_linkable"`
	LastModificationDate   types.String   `tfsdk:"last_modification_date"`
	Orn                    types.String   `tfsdk:"orn"`
	Path                   types.String   `tfsdk:"path"`
	PolicyDefaultVersionId types.String   `tfsdk:"policy_default_version_id"`
	PolicyId               types.String   `tfsdk:"policy_id"`
	PolicyName             types.String   `tfsdk:"policy_name"`
	ResourcesCount         types.Int32    `tfsdk:"resources_count"`
	Document               types.String   `tfsdk:"document"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
	Id                     types.String   `tfsdk:"id"`
}

type resourcePolicy struct {
	Client *osc.Client
}

func NewResourcePolicy() resource.Resource {
	return &resourcePolicy{}
}

func (r *resourcePolicy) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourcePolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourcePolicy) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *resourcePolicy) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourcePolicy) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"creation_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_linkable": schema.BoolAttribute{
				Computed: true,
			},
			"last_modification_date": schema.StringAttribute{
				Computed: true,
			},
			"orn": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"policy_default_version_id": schema.StringAttribute{
				Computed: true,
			},
			"policy_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resources_count": schema.Int32Attribute{
				Computed: true,
			},
			"document": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (r *resourcePolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreatePolicyRequest{
		PolicyName: data.PolicyName.ValueString(),
		Document:   data.Document.ValueString(),
	}

	if fwhelpers.IsSet(data.Path) {
		createReq.Path = data.Path.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.Description) {
		createReq.Description = data.Description.ValueStringPointer()
	}

	createResp, err := r.Client.CreatePolicy(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Policy",
			err.Error(),
		)
		return
	}
	policy := ptr.From(createResp.Policy)

	data.Id = to.String(policy.Orn)

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Policy state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourcePolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	readTimeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, readTimeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Policy API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourcePolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourcePolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	deleteTimeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	err := r.unlinkEntities(ctx, deleteTimeout, data.Orn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unlink entities from Policy",
			err.Error(),
		)
		return
	}

	err = r.deleteVersions(ctx, deleteTimeout, data.Orn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete versions linked to Policy",
			err.Error(),
		)
		return
	}

	delReq := osc.DeletePolicyRequest{
		PolicyOrn: data.Orn.ValueString(),
	}

	_, err = r.Client.DeletePolicy(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Policy",
			err.Error(),
		)
	}
}

func (r *resourcePolicy) deleteVersions(ctx context.Context, to time.Duration, orn string) error {
	req := osc.ReadPolicyVersionsRequest{
		PolicyOrn: orn,
	}

	versions, err := r.Client.ReadPolicyVersions(ctx, req, options.WithRetryTimeout(to))
	if err != nil || versions.PolicyVersions == nil {
		return err
	}
	if len(*versions.PolicyVersions) > 1 {
		for _, version := range *versions.PolicyVersions {
			if *version.DefaultVersion {
				continue
			}
			delReq := osc.DeletePolicyVersionRequest{
				PolicyOrn: orn,
				VersionId: *version.VersionId,
			}
			_, err := r.Client.DeletePolicyVersion(ctx, delReq, options.WithRetryTimeout(to))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *resourcePolicy) unlinkEntities(ctx context.Context, to time.Duration, orn string) error {
	req := osc.ReadEntitiesLinkedToPolicyRequest{PolicyOrn: orn}

	var users, groups []osc.MinimalPolicy
	resp, err := r.Client.ReadEntitiesLinkedToPolicy(ctx, req, options.WithRetryTimeout(to))
	if err != nil {
		return err
	}
	if resp.PolicyEntities != nil {
		if resp.PolicyEntities.Users != nil {
			users = *resp.PolicyEntities.Users
		}
		if resp.PolicyEntities.Groups != nil {
			groups = *resp.PolicyEntities.Groups
		}
	}
	if len(users) > 0 {
		req := osc.UnlinkPolicyRequest{
			PolicyOrn: orn,
		}

		for _, user := range users {
			if user.Name == nil {
				continue
			}
			req.UserName = *user.Name
			_, err := r.Client.UnlinkPolicy(ctx, req, options.WithRetryTimeout(to))
			if err != nil {
				return err
			}
		}
	}
	if len(groups) > 0 {
		req := osc.UnlinkManagedPolicyFromUserGroupRequest{
			PolicyOrn: orn,
		}

		for _, group := range groups {
			if group.Name == nil {
				continue
			}
			req.UserGroupName = *group.Name
			_, err := r.Client.UnlinkManagedPolicyFromUserGroup(ctx, req, options.WithRetryTimeout(to))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *resourcePolicy) getDocument(ctx context.Context, to time.Duration, orn, versionId string) (string, error) {
	req := osc.ReadPolicyVersionRequest{
		PolicyOrn: orn,
		VersionId: versionId,
	}

	resp, err := r.Client.ReadPolicyVersion(ctx, req, options.WithRetryTimeout(to))
	if err != nil {
		return "", err
	}
	if resp.PolicyVersion == nil {
		return "", fmt.Errorf("cannot find policy version: %v", versionId)
	}

	return ptr.From(resp.PolicyVersion.Body), err
}

func (r *resourcePolicy) read(ctx context.Context, timeout time.Duration, data PolicyModel) (PolicyModel, error) {
	req := osc.ReadPolicyRequest{
		PolicyOrn: data.Id.ValueString(),
	}

	resp, err := r.Client.ReadPolicy(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Policy == nil {
		return data, ErrResourceEmpty
	}

	policy := ptr.From(resp.Policy)
	document, err := r.getDocument(ctx, timeout, *policy.Orn, "v1")
	if err != nil {
		return data, err
	}

	data.CreationDate = to.String(from.ISO8601(policy.CreationDate))
	data.Description = to.String(ptr.From(policy.Description))
	data.IsLinkable = to.Bool(ptr.From(policy.IsLinkable))
	data.LastModificationDate = to.String(from.ISO8601(policy.LastModificationDate))
	data.Orn = to.String(*policy.Orn)
	data.Path = to.String(ptr.From(policy.Path))
	data.PolicyDefaultVersionId = to.String(ptr.From(policy.PolicyDefaultVersionId))
	data.PolicyName = to.String(ptr.From(policy.PolicyName))
	data.ResourcesCount = to.Int32(int32(ptr.From(policy.ResourcesCount)))
	data.Document = to.String(document)
	data.Id = to.String(policy.Orn)
	data.PolicyId = to.String(policy.PolicyId)

	return data, nil
}
