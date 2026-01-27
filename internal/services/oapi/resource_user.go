package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceUser{}
	_ resource.ResourceWithConfigure   = &resourceUser{}
	_ resource.ResourceWithImportState = &resourceUser{}
)

type UserModel struct {
	CreationDate         types.String   `tfsdk:"creation_date"`
	LastModificationDate types.String   `tfsdk:"last_modification_date"`
	Path                 types.String   `tfsdk:"path"`
	UserEmail            types.String   `tfsdk:"user_email"`
	UserId               types.String   `tfsdk:"user_id"`
	UserName             types.String   `tfsdk:"user_name"`
	Policies             types.Set      `tfsdk:"policy"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
	Id                   types.String   `tfsdk:"id"`
}

type UserPolicyModel struct {
	PolicyOrn            types.String                       `tfsdk:"policy_orn"`
	DefaultVersionId     fwtypes.CaseInsensitiveStringValue `tfsdk:"default_version_id"`
	PolicyName           types.String                       `tfsdk:"policy_name"`
	PolicyId             types.String                       `tfsdk:"policy_id"`
	CreationDate         types.String                       `tfsdk:"creation_date"`
	LastModificationDate types.String                       `tfsdk:"last_modification_date"`
}

type UserCommon struct {
	Client *osc.APIClient
}

type resourceUser struct {
	UserCommon
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

func (r *resourceUser) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"policy": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"policy_orn": schema.StringAttribute{
							Required: true,
						},
						"default_version_id": schema.StringAttribute{
							CustomType: fwtypes.CaseInsensitiveStringType{},
							Computed:   true,
							Optional:   true,
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
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(EIMPathDefaultValue),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						EIMPathRegexp,
						ErrResourceInvalidEIMPath.Error(),
					),
					stringvalidator.LengthBetween(1, 512),
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

	createReq := osc.CreateUserRequest{
		UserName: data.UserName.ValueString(),
		Path:     ptr.To(data.Path.ValueString()),
	}

	if fwhelpers.IsSet(data.UserEmail) {
		createReq.SetUserEmail(data.UserEmail.ValueString())
	}

	var createResp osc.CreateUserResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.UserApi.CreateUser(ctx).CreateUserRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create User",
			err.Error(),
		)
		return
	}
	data.Id = to.String(createResp.GetUser().UserId)

	if fwhelpers.IsSet(data.Policies) {
		policies, diag := to.Slice[UserPolicyModel](ctx, data.Policies)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		diag = r.linkPolicies(ctx, createTimeout, data.UserName.ValueString(), policies)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set User state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *UserCommon) setDefaultPolicyVersion(ctx context.Context, to time.Duration, orn, version string) error {
	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		_, httpResp, err := r.Client.PolicyApi.SetDefaultPolicyVersion(ctx).SetDefaultPolicyVersionRequest(
			osc.SetDefaultPolicyVersionRequest{
				PolicyOrn: orn,
				VersionId: version,
			}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	return err
}

func (r *resourceUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set User API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	updateReq := osc.UpdateUserRequest{
		UserName: stateData.UserName.ValueString(),
	}

	updateUser := false
	if fwhelpers.HasChange(planData.UserName, stateData.UserName) {
		updateReq.SetNewUserName(planData.UserName.ValueString())
		updateUser = true
	}
	if fwhelpers.HasChange(planData.Path, stateData.Path) {
		updateReq.SetNewPath(planData.Path.ValueString())
		updateUser = true
	}
	if fwhelpers.HasChange(planData.UserEmail, stateData.UserEmail) {
		updateReq.SetNewUserEmail(planData.UserEmail.ValueString())
		updateUser = true
	}
	if fwhelpers.HasChange(planData.Policies, stateData.Policies) {
		diag := r.updatePolicies(ctx, timeout, planData, stateData)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	if updateUser {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.UserApi.UpdateUser(ctx).UpdateUserRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update User",
				err.Error(),
			)
		}
	}

	newStateData, err := r.read(ctx, timeout, stateData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set User API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *resourceUser) updatePolicies(ctx context.Context, timeout time.Duration, planData, stateData UserModel) diag.Diagnostics {
	var diags diag.Diagnostics

	statePolicies, diag := to.Slice[UserPolicyModel](ctx, stateData.Policies)
	diags.Append(diag...)

	planPolicies, diag := to.Slice[UserPolicyModel](ctx, planData.Policies)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	toRemove, toAdd := lo.Difference(statePolicies, planPolicies)
	diags.Append(r.unlinkPolicies(ctx, timeout, stateData.UserName.ValueString(), toRemove)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(r.linkPolicies(ctx, timeout, stateData.UserName.ValueString(), toAdd)...)

	return diags
}

func (r *resourceUser) unlinkPolicies(ctx context.Context, timeout time.Duration, username string, policies []UserPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, policy := range policies {
		req := osc.UnlinkPolicyRequest{
			UserName:  username,
			PolicyOrn: policy.PolicyOrn.ValueString(),
		}

		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.PolicyApi.UnlinkPolicy(ctx).UnlinkPolicyRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			diags.AddError(
				"Unable to unlink policy from User",
				err.Error(),
			)
			return diags
		}
	}

	return diags
}

func (r *resourceUser) linkPolicies(ctx context.Context, timeout time.Duration, username string, policies []UserPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, policy := range policies {
		req := osc.LinkPolicyRequest{
			UserName:  username,
			PolicyOrn: policy.PolicyOrn.ValueString(),
		}

		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.PolicyApi.LinkPolicy(ctx).LinkPolicyRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			diags.AddError(
				"Unable to link policy to User",
				err.Error(),
			)
			return diags
		}

		if fwhelpers.IsSet(policy.DefaultVersionId) {
			err := r.setDefaultPolicyVersion(ctx, timeout, policy.PolicyOrn.ValueString(), policy.DefaultVersionId.ValueString())
			if err != nil {
				diags.AddError(
					"Unable to set default policy version for policy linked to User",
					err.Error(),
				)
				return diags
			}
		}
	}

	return diags
}

func (r *resourceUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	if fwhelpers.IsSet(data.Policies) {
		policies, diag := to.Slice[UserPolicyModel](ctx, data.Policies)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		diag = r.unlinkPolicies(ctx, timeout, data.UserName.ValueString(), policies)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	deleteReq := osc.DeleteUserRequest{
		UserName: data.UserName.ValueString(),
	}

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.UserApi.DeleteUser(ctx).DeleteUserRequest(deleteReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete User",
			err.Error(),
		)
	}
}

func (r *resourceUser) read(ctx context.Context, timeout time.Duration, data UserModel) (UserModel, error) {
	req := osc.ReadUsersRequest{
		Filters: &osc.FiltersUsers{
			UserIds: &[]string{data.Id.ValueString()},
		},
	}

	var resp osc.ReadUsersResponse
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.UserApi.ReadUsers(ctx).ReadUsersRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return data, err
	}
	if len(resp.GetUsers()) == 0 {
		return data, ErrResourceEmpty
	}

	user := resp.GetUsers()[0]
	linkReq := osc.NewReadLinkedPoliciesRequest(user.GetUserName())
	var linkResp osc.ReadLinkedPoliciesResponse
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.PolicyApi.ReadLinkedPolicies(ctx).ReadLinkedPoliciesRequest(*linkReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		linkResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}

	policies, err := r.flattenPolicies(ctx, timeout, linkResp.GetPolicies())
	if err != nil {
		return data, err
	}

	policiesSet, diag := to.SetObject(ctx, policies)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert policies to a set: %v", diag.Errors())
	}

	data.Policies = policiesSet
	data.CreationDate = to.String(user.CreationDate)
	data.LastModificationDate = to.String(user.LastModificationDate)
	data.Path = to.String(user.Path)
	data.UserEmail = to.String(user.UserEmail)
	data.UserId = to.String(user.UserId)
	data.UserName = to.String(user.UserName)
	data.Id = to.String(user.UserId)

	return data, nil
}

func (r *UserCommon) flattenPolicies(ctx context.Context, timeout time.Duration, policies []osc.LinkedPolicy) ([]UserPolicyModel, error) {
	flattenedPolicies := make([]UserPolicyModel, 0, len(policies))
	for _, policy := range policies {
		versionID, err := r.getPolicyVersion(ctx, timeout, policy.GetOrn())
		if err != nil {
			return nil, err
		}
		flattenedPolicies = append(flattenedPolicies, UserPolicyModel{
			PolicyOrn:            to.String(policy.Orn),
			DefaultVersionId:     fwtypes.CaseInsensitiveString(versionID),
			PolicyName:           to.String(policy.PolicyName),
			PolicyId:             to.String(policy.PolicyId),
			CreationDate:         to.String(policy.CreationDate),
			LastModificationDate: to.String(policy.LastModificationDate),
		})
	}
	return flattenedPolicies, nil
}

func (r *UserCommon) getPolicyVersion(ctx context.Context, timeout time.Duration, orn string) (string, error) {
	var versionId string
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		resp, httpResp, err := r.Client.PolicyApi.ReadPolicy(ctx).ReadPolicyRequest(
			osc.ReadPolicyRequest{PolicyOrn: orn}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		versionId = ptr.From(resp.GetPolicy().PolicyDefaultVersionId)
		return nil
	})
	if err != nil {
		return versionId, err
	}
	return versionId, err
}
