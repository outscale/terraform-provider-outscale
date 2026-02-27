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
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceUserGroup{}
	_ resource.ResourceWithConfigure   = &resourceUserGroup{}
	_ resource.ResourceWithImportState = &resourceUserGroup{}
)

type UserGroupModel struct {
	CreationDate         types.String   `tfsdk:"creation_date"`
	LastModificationDate types.String   `tfsdk:"last_modification_date"`
	UserGroupName        types.String   `tfsdk:"user_group_name"`
	Orn                  types.String   `tfsdk:"orn"`
	UserGroupId          types.String   `tfsdk:"user_group_id"`
	Path                 types.String   `tfsdk:"path"`
	Policies             types.Set      `tfsdk:"policy"`
	Users                types.Set      `tfsdk:"user"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
	Id                   types.String   `tfsdk:"id"`
}

type UserGroupUserModel struct {
	UserName             types.String `tfsdk:"user_name"`
	Path                 types.String `tfsdk:"path"`
	UserId               types.String `tfsdk:"user_id"`
	CreationDate         types.String `tfsdk:"creation_date"`
	LastModificationDate types.String `tfsdk:"last_modification_date"`
}

type resourceUserGroup struct {
	UserCommon
}

func NewResourceUserGroup() resource.Resource {
	return &resourceUserGroup{}
}

func (r *resourceUserGroup) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceUserGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceUserGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *resourceUserGroup) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"user": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"user_name": schema.StringAttribute{
							Required: true,
						},
						"path": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									EIMPathRegexp,
									ErrResourceInvalidEIMPath.Error(),
								),
								stringvalidator.LengthBetween(1, 512),
							},
						},
						"user_id": schema.StringAttribute{
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
			"orn": schema.StringAttribute{
				Computed: true,
			},
			"user_group_name": schema.StringAttribute{
				Required: true,
			},
			"user_group_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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

func (r *resourceUserGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreateUserGroupRequest{
		UserGroupName: data.UserGroupName.ValueString(),
		Path:          new(data.Path.ValueString()),
	}

	createResp, err := r.Client.CreateUserGroup(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create User Group",
			err.Error(),
		)
		return
	}

	data.Id = to.String(createResp.UserGroup.UserGroupId)

	if fwhelpers.IsSet(data.Users) {
		users, diag := to.Slice[UserGroupUserModel](ctx, data.Users)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		diag = r.addUsers(ctx, createTimeout, data.UserGroupName.ValueString(), data.Path.ValueString(), users)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}
	if fwhelpers.IsSet(data.Policies) {
		policies, diag := to.Slice[UserPolicyModel](ctx, data.Policies)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		diag = r.linkPolicies(ctx, createTimeout, data.UserGroupName.ValueString(), policies)
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

func (r *resourceUserGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserGroupModel
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

func (r *resourceUserGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData UserGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	updateReq := osc.UpdateUserGroupRequest{
		UserGroupName: stateData.UserGroupName.ValueString(),
	}
	updateGroup := false

	if fwhelpers.HasChange(planData.UserGroupName, stateData.UserGroupName) {
		updateReq.NewUserGroupName = planData.UserGroupName.ValueStringPointer()
		updateGroup = true
	}
	if fwhelpers.HasChange(planData.Path, stateData.Path) {
		updateReq.NewPath = planData.Path.ValueStringPointer()
		updateGroup = true
	}
	if fwhelpers.HasChange(planData.Policies, stateData.Policies) {
		diag := r.updatePolicies(ctx, timeout, planData, stateData)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}
	if fwhelpers.HasChange(planData.Users, stateData.Users) {
		diag := r.updateUsers(ctx, timeout, planData, stateData)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	if updateGroup {
		_, err := r.Client.UpdateUserGroup(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update User Group",
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
			"Unable to set User Group API response values",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *resourceUserGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	deleteReq := osc.DeleteUserGroupRequest{
		UserGroupName: data.UserGroupName.ValueString(),
		Force:         new(true),
	}

	_, err := r.Client.DeleteUserGroup(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete User Group",
			err.Error(),
		)
	}
}

func (r *resourceUserGroup) read(ctx context.Context, timeout time.Duration, data UserGroupModel) (UserGroupModel, error) {
	req := osc.ReadUserGroupsRequest{
		Filters: &osc.FiltersUserGroup{
			UserGroupIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadUserGroups(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.UserGroups == nil || len(*resp.UserGroups) == 0 {
		return data, ErrResourceEmpty
	}

	userGroup := (*resp.UserGroups)[0]
	reqUserGroup := osc.ReadUserGroupRequest{
		UserGroupName: *userGroup.Name,
	}

	respUsers, err := r.Client.ReadUserGroup(ctx, reqUserGroup, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	reqLink := osc.ReadManagedPoliciesLinkedToUserGroupRequest{
		UserGroupName: *userGroup.Name,
	}

	respLink, err := r.Client.ReadManagedPoliciesLinkedToUserGroup(ctx, reqLink, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	data.UserGroupName = to.String(ptr.From(userGroup.Name))
	data.Path = to.String(ptr.From(userGroup.Path))
	data.UserGroupId = to.String(userGroup.UserGroupId)
	data.Id = to.String(userGroup.UserGroupId)
	data.Orn = to.String(ptr.From(userGroup.Orn))
	data.CreationDate = to.String(from.ISO8601(ptr.From(userGroup.CreationDate)))
	data.LastModificationDate = to.String(from.ISO8601(userGroup.LastModificationDate))

	policies, err := r.flattenPolicies(ctx, timeout, ptr.From(respLink.Policies))
	if err != nil {
		return data, err
	}
	policiesSet, diag := to.SetObject(ctx, policies)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert policies to a set: %v", diag.Errors())
	}
	data.Policies = policiesSet

	users := r.flattenUsers(ptr.From(respUsers.Users))
	usersSet, diag := to.SetObject(ctx, users)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert users to a set: %v", diag.Errors())
	}
	data.Users = usersSet

	return data, nil
}

func (r *resourceUserGroup) addUsers(ctx context.Context, timeout time.Duration, groupName, groupPath string, users []UserGroupUserModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, user := range users {
		req := osc.AddUserToUserGroupRequest{
			UserGroupName: groupName,
			UserGroupPath: &groupPath,
		}
		if fwhelpers.IsSet(user.UserName) {
			req.UserName = user.UserName.ValueString()
		}

		if fwhelpers.IsSet(user.Path) {
			req.UserPath = user.Path.ValueStringPointer()
		}

		_, err := r.Client.AddUserToUserGroup(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			diags.AddError(
				"Unable to add user to User Group",
				err.Error(),
			)
			return diags
		}
	}

	return diags
}

func (r *resourceUserGroup) removeUsers(ctx context.Context, timeout time.Duration, groupName, groupPath string, users []UserGroupUserModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, user := range users {
		req := osc.RemoveUserFromUserGroupRequest{
			UserGroupName: groupName,
			UserGroupPath: &groupPath,
		}
		if fwhelpers.IsSet(user.UserName) {
			req.UserName = user.UserName.ValueString()
		}

		if fwhelpers.IsSet(user.Path) {
			req.UserPath = user.Path.ValueStringPointer()
		}

		_, err := r.Client.RemoveUserFromUserGroup(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			// This case happens when a user linked to a group changes its username
			// Terraform detects a change and need to remove and readd the user to the group
			// Trying to remove the user based on the state name fails since the username changed
			errCode := oapihelpers.GetError(err)
			if errCode.Code == "5098" {
				continue
			}
			diags.AddError(
				"Unable to remove user from User Group",
				err.Error(),
			)
			return diags
		}
	}

	return diags
}

func (r *resourceUserGroup) updatePolicies(ctx context.Context, timeout time.Duration, planData, stateData UserGroupModel) diag.Diagnostics {
	var diags diag.Diagnostics

	statePolicies, diag := to.Slice[UserPolicyModel](ctx, stateData.Policies)
	diags.Append(diag...)

	planPolicies, diag := to.Slice[UserPolicyModel](ctx, planData.Policies)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	toRemove, toAdd := lo.Difference(statePolicies, planPolicies)
	diags.Append(r.unlinkPolicies(ctx, timeout, stateData.UserGroupName.ValueString(), toRemove)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(r.linkPolicies(ctx, timeout, stateData.UserGroupName.ValueString(), toAdd)...)

	return diags
}

func (r *resourceUserGroup) updateUsers(ctx context.Context, timeout time.Duration, planData, stateData UserGroupModel) diag.Diagnostics {
	var diags diag.Diagnostics
	stateUsers, diag := to.Slice[UserGroupUserModel](ctx, stateData.Users)
	diags.Append(diag...)

	planUsers, diag := to.Slice[UserGroupUserModel](ctx, planData.Users)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	toRemove, toAdd := lo.Difference(stateUsers, planUsers)
	diags.Append(r.removeUsers(ctx, timeout, stateData.UserGroupName.ValueString(), stateData.Path.ValueString(), toRemove)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(r.addUsers(ctx, timeout, stateData.UserGroupName.ValueString(), stateData.Path.ValueString(), toAdd)...)

	return diags
}

func (r *resourceUserGroup) unlinkPolicies(ctx context.Context, timeout time.Duration, groupname string, policies []UserPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, policy := range policies {
		req := osc.UnlinkManagedPolicyFromUserGroupRequest{
			UserGroupName: groupname,
			PolicyOrn:     policy.PolicyOrn.ValueString(),
		}

		_, err := r.Client.UnlinkManagedPolicyFromUserGroup(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			diags.AddError(
				"Unable to unlink policy from User Group",
				err.Error(),
			)
			return diags
		}
	}

	return diags
}

func (r *resourceUserGroup) linkPolicies(ctx context.Context, timeout time.Duration, usergroupname string, policies []UserPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, policy := range policies {
		req := osc.LinkManagedPolicyToUserGroupRequest{
			UserGroupName: usergroupname,
			PolicyOrn:     policy.PolicyOrn.ValueString(),
		}

		_, err := r.Client.LinkManagedPolicyToUserGroup(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			diags.AddError(
				"Unable to link policy to User Group",
				err.Error(),
			)
			return diags
		}

		if fwhelpers.IsSet(policy.DefaultVersionId) {
			err := r.setDefaultPolicyVersion(ctx, timeout, policy.PolicyOrn.ValueString(), policy.DefaultVersionId.ValueString())
			if err != nil {
				diags.AddError(
					"Unable to set default policy version for policy linked to User Group",
					err.Error(),
				)
				return diags
			}
		}
	}

	return diags
}

func (r *resourceUserGroup) flattenUsers(users []osc.User) []UserGroupUserModel {
	return lo.Map(users, func(user osc.User, _ int) UserGroupUserModel {
		return UserGroupUserModel{
			UserId:               to.String(user.UserId),
			UserName:             to.String(user.UserName),
			Path:                 to.String(user.Path),
			CreationDate:         to.String(from.ISO8601(user.CreationDate)),
			LastModificationDate: to.String(from.ISO8601(user.LastModificationDate)),
		}
	})
}
