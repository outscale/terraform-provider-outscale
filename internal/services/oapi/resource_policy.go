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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
	Client *osc.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
		createReq.SetPath(data.Path.ValueString())
	}
	if fwhelpers.IsSet(data.Description) {
		createReq.SetDescription(data.Description.ValueString())
	}

	var createResp osc.CreatePolicyResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.PolicyApi.CreatePolicy(ctx).CreatePolicyRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Policy",
			err.Error(),
		)
		return
	}
	policy := createResp.GetPolicy()

	data.Id = to.String(policy.Orn)

	stateData, err := r.setPolicyState(ctx, createTimeout, data)
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

	stateData, err := r.setPolicyState(ctx, readTimeout, data)
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

	err = retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.PolicyApi.DeletePolicy(ctx).DeletePolicyRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Policy",
			err.Error(),
		)
		return
	}
}

func (r *resourcePolicy) deleteVersions(ctx context.Context, to time.Duration, orn string) error {
	versions, err := getPolicyVersions(ctx, r.Client, to, orn)
	if err != nil {
		return err
	}
	if len(versions) > 1 {
		for _, version := range versions {
			if version.GetDefaultVersion() {
				continue
			}
			delReq := osc.DeletePolicyVersionRequest{
				PolicyOrn: orn,
				VersionId: version.GetVersionId(),
			}
			err := retry.RetryContext(ctx, to, func() *retry.RetryError {
				_, httpResp, err := r.Client.PolicyApi.DeletePolicyVersion(ctx).DeletePolicyVersionRequest(delReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
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
	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		resp, httpResp, err := r.Client.PolicyApi.ReadEntitiesLinkedToPolicy(ctx).ReadEntitiesLinkedToPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		if resp.HasPolicyEntities() {
			users = *resp.GetPolicyEntities().Users
			groups = *resp.GetPolicyEntities().Groups
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		req := osc.UnlinkPolicyRequest{
			PolicyOrn: orn,
		}

		for _, user := range users {
			req.SetUserName(user.GetName())
			err := retry.RetryContext(ctx, to, func() *retry.RetryError {
				_, httpResp, err := r.Client.PolicyApi.UnlinkPolicy(ctx).UnlinkPolicyRequest(req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
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
			req.SetUserGroupName(group.GetName())
			err := retry.RetryContext(ctx, to, func() *retry.RetryError {
				_, httpResp, err := r.Client.PolicyApi.UnlinkManagedPolicyFromUserGroup(ctx).UnlinkManagedPolicyFromUserGroupRequest(req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *resourcePolicy) getDocument(ctx context.Context, to time.Duration, orn, versionId string) (string, error) {
	req := osc.NewReadPolicyVersionRequest(orn, versionId)

	var resp osc.ReadPolicyVersionResponse
	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		rp, httpResp, err := r.Client.PolicyApi.ReadPolicyVersion(ctx).ReadPolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return "", err
	}
	if _, ok := resp.GetPolicyVersionOk(); !ok {
		return "", fmt.Errorf("cannot find policy version: %v", versionId)
	}

	return *resp.GetPolicyVersion().Body, err
}

func (r *resourcePolicy) setPolicyState(ctx context.Context, timeout time.Duration, data PolicyModel) (PolicyModel, error) {
	req := osc.NewReadPolicyRequest(data.Id.ValueString())

	var resp osc.ReadPolicyResponse
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.PolicyApi.ReadPolicy(ctx).ReadPolicyRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return data, err
	}
	if resp.Policy == nil {
		return data, ErrResourceEmpty
	}

	policy := resp.GetPolicy()
	document, err := r.getDocument(ctx, timeout, policy.GetOrn(), policy.GetPolicyDefaultVersionId())
	if err != nil {
		return data, err
	}

	data.CreationDate = to.String(policy.CreationDate)
	data.Description = to.String(policy.Description)
	data.IsLinkable = to.Bool(policy.IsLinkable)
	data.LastModificationDate = to.String(policy.LastModificationDate)
	data.Orn = to.String(policy.Orn)
	data.Path = to.String(policy.Path)
	data.PolicyDefaultVersionId = to.String(policy.PolicyDefaultVersionId)
	data.PolicyName = to.String(policy.PolicyName)
	data.ResourcesCount = to.Int32(policy.ResourcesCount)
	data.Document = to.String(document)
	data.Id = to.String(policy.Orn)
	data.PolicyId = to.String(policy.PolicyId)

	return data, nil
}
