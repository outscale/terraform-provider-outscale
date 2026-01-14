package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource                = &resourceUser{}
	_ resource.ResourceWithConfigure   = &resourceUser{}
	_ resource.ResourceWithImportState = &resourceUser{}
	_ resource.ResourceWithModifyPlan  = &resourceUser{}
)

type UserModel struct {
	CreationDate         types.String   `tfsdk:"creation_date"`
	LastModificationDate types.String   `tfsdk:"last_modification_date"`
	Path                 types.String   `tfsdk:"path"`
	UserEmail            types.String   `tfsdk:"user_email"`
	UserId               types.String   `tfsdk:"user_id"`
	UserName             types.String   `tfsdk:"user_name"`
	Policy               types.Set      `tfsdk:"policy"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
	Id                   types.String   `tfsdk:"id"`
}

type resourceUser struct {
	Client *osc.APIClient
}

func NewResourceUser() resource.Resource {
	return &resourceUser{}
}

func (r *resourceUser) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceUser) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceUser) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *resourceUser) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceUser) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"creation_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modification_date": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("/"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						EIMPathRegexp,
						ErrResourceInvalidEIMPath.Error(),
					),
				},
			},
			"user_email": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"user_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				Required: true,
			},
			"policy": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_orn": schema.StringAttribute{
							Required: true,
						},
						"default_version_id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"policy_name": schema.StringAttribute{
							Computed: true,
						},
						"policy_id": schema.StringAttribute{
							Computed: true,
						},
						"creation_date": schema.StringAttribute{
							Computed: true,
						},
						"last_modification_date": schema.StringAttribute{
							Computed: true,
						},
					},
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

func (r *resourceUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

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
	user := createResp.GetPolicy()

	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(user.Orn)

	stateData, err := r.setPolicyState(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Policy state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserModel
	diag := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.setPolicyState(ctx, data)
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

func (r *resourceUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	deleteTimeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

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

func (r *resourceUser) deleteVersions(ctx context.Context, to time.Duration, orn string) error {
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

func (r *resourceUser) unlinkEntities(ctx context.Context, to time.Duration, orn string) error {
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

func (r *resourceUser) getDocument(ctx context.Context, to time.Duration, orn, versionId string) (string, error) {
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
		return "", fmt.Errorf("cannot find user version: %v", versionId)
	}

	return *resp.GetPolicyVersion().Body, err
}

func (r *resourceUser) setPolicyState(ctx context.Context, data UserModel) (UserModel, error) {
	readTimeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diag.HasError() {
		return data, fmt.Errorf("unable to parse 'user' read timeout value: %v", diag.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	req := osc.NewReadPolicyRequest(data.Id.ValueString())

	var resp osc.ReadPolicyResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
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
	data.RequestId = to.String(resp.ResponseContext.GetRequestId())
	if resp.Policy == nil {
		return data, ErrResourceEmpty
	}

	user := resp.GetPolicy()
	document, err := r.getDocument(ctx, readTimeout, user.GetOrn(), user.GetPolicyDefaultVersionId())
	if err != nil {
		return data, err
	}

	data.CreationDate = to.String(user.CreationDate)
	data.Description = to.String(user.Description)
	data.IsLinkable = to.Bool(user.IsLinkable)
	data.LastModificationDate = to.String(user.LastModificationDate)
	data.Orn = to.String(user.Orn)
	data.Path = to.String(user.Path)
	data.PolicyDefaultVersionId = to.String(user.PolicyDefaultVersionId)
	data.PolicyName = to.String(user.PolicyName)
	data.ResourcesCount = to.Int32(user.ResourcesCount)
	data.Document = to.String(document)
	data.Id = to.String(user.Orn)
	data.PolicyId = to.String(user.PolicyId)

	return data, nil
}
