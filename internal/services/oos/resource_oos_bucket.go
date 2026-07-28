package oos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"k8s.io/utils/keymutex"
)

var (
	_ resource.Resource                   = &bucketResource{}
	_ resource.ResourceWithConfigure      = &bucketResource{}
	_ resource.ResourceWithModifyPlan     = &bucketResource{}
	_ resource.ResourceWithValidateConfig = &bucketResource{}
)

const (
	bucketErrCreate           = "Unable to create Bucket"
	bucketObjectLockErrCreate = "Unable to put Bucket Object Lock Configuration"
	bucketErrUpdate           = "Unable to update Bucket"
	bucketErrDelete           = "Unable to delete Bucket"
	bucketErrForceDelete      = "Unable to force delete Bucket"
)

type bucketModel struct {
	Name         types.String   `tfsdk:"name"`
	LocationName types.String   `tfsdk:"location_name"`
	ObjectLock   types.Object   `tfsdk:"object_lock"`
	ACL          types.String   `tfsdk:"acl"`
	Grant        types.Object   `tfsdk:"grant"`
	Permissions  types.Object   `tfsdk:"permissions"`
	ForceDelete  types.Bool     `tfsdk:"force_delete"`
	CreationDate types.String   `tfsdk:"creation_date"`
	Id           types.String   `tfsdk:"id"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

type bucketObjectLockModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	DefaultRetention types.Object `tfsdk:"default_retention"`
}

func (bucketObjectLockModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "object_lock")
}

type bucketDefaultRetentionModel struct {
	Mode  types.String `tfsdk:"mode"`
	Days  types.Int32  `tfsdk:"days"`
	Years types.Int32  `tfsdk:"years"`
}

func (bucketDefaultRetentionModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "object_lock", "default_retention")
}

type bucketGrantModel struct {
	FullControl types.Object `tfsdk:"full_control"`
	Read        types.Object `tfsdk:"read"`
	Write       types.Object `tfsdk:"write"`
	ReadACP     types.Object `tfsdk:"read_acp"`
	WriteACP    types.Object `tfsdk:"write_acp"`
}

func (bucketGrantModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "grant")
}

type bucketGrantPermissionModel struct {
	Ids            types.Set `tfsdk:"ids"`
	EmailAddresses types.Set `tfsdk:"email_addresses"`
}

type bucketPermissionsModel struct {
	Grants types.Set    `tfsdk:"grants"`
	Owner  types.Object `tfsdk:"owner"`
}

func (bucketPermissionsModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "permissions")
}

type bucketGrantsModel struct {
	Permission types.String `tfsdk:"permission"`
	Grantee    types.Object `tfsdk:"grantee"`
}

func (bucketGrantsModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "permissions", "grants")
}

type bucketGranteeModel struct {
	Type         types.String `tfsdk:"type"`
	DisplayName  types.String `tfsdk:"display_name"`
	EmailAddress types.String `tfsdk:"email_address"`
	Id           types.String `tfsdk:"id"`
	Uri          types.String `tfsdk:"uri"`
}

func (bucketGranteeModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "permissions", "grants", "grantee")
}

type bucketOwnerModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func (bucketOwnerModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketSchemaType, "permissions", "owner")
}

type bucketResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucket() resource.Resource {
	return &bucketResource{}
}

func (r *bucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *oos.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Client = client.OOS
	r.mu = client.KeyMutex
}

func (r *bucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket"
}

func (r *bucketResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("acl"),
			path.MatchRoot("grant"),
		),
	}
}

func (r *bucketResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var objectLock types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("object_lock"), &objectLock)...)
	if resp.Diagnostics.HasError() || !fwhelpers.IsSet(objectLock) {
		return
	}

	model, diags := to.Model[bucketObjectLockModel](ctx, objectLock)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if fwhelpers.IsSet(model.DefaultRetention) && !model.Enabled.IsUnknown() && (model.Enabled.IsNull() || !model.Enabled.ValueBool()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("object_lock").AtName("default_retention"),
			"Invalid Attribute Combination",
			"Attribute \"default_retention\" can only be configured when \"enabled\" is true.",
		)
	}
}

func (r *bucketResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var stateForceDelete types.Bool
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("force_delete"), &stateForceDelete)...)
	if resp.Diagnostics.HasError() || stateForceDelete.IsUnknown() || stateForceDelete.ValueBool() {
		return
	}

	var planForceDelete types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("force_delete"), &planForceDelete)...)
	if resp.Diagnostics.HasError() || planForceDelete.IsUnknown() || !planForceDelete.ValueBool() {
		return
	}

	resp.Diagnostics.AddAttributeWarning(
		path.Root("force_delete"),
		"Apply \"force_delete\" before destroying the Bucket",
		"Terraform destruction uses the \"force_delete\" value stored in state. Run `terraform apply` to persist \"force_delete = true\" before running the destroy command.",
	)
}

func bucketSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
				},
			},
			"object_lock": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
							if !req.ConfigValue.IsNull() || !fwhelpers.IsSet(req.StateValue) {
								return
							}

							enabled, ok := req.StateValue.Attributes()["enabled"].(types.Bool)
							if ok && fwhelpers.IsSet(enabled) && enabled.ValueBool() {
								resp.RequiresReplace = true
							}
						},
						"Removing an enabled object_lock configuration requires replacing the bucket.",
						"Removing an enabled object_lock configuration requires replacing the bucket.",
					),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						Default:  booldefault.StaticBool(false),
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"default_retention": schema.SingleNestedAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplaceIf(
								func(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
									if fwhelpers.IsSet(req.StateValue) && req.PlanValue.IsNull() {
										resp.RequiresReplace = true
									}
								},
								"Removing the default_retention configuration requires replacing the bucket.",
								"Removing the default_retention configuration requires replacing the bucket.",
							),
						},
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Computed: true,
								Optional: true,
								Default:  stringdefault.StaticString("COMPLIANCE"),
								Validators: []validator.String{
									stringvalidator.OneOf("COMPLIANCE"),
								},
							},
							"days": schema.Int32Attribute{
								Optional: true,
								Validators: []validator.Int32{
									int32validator.ExactlyOneOf(
										path.MatchRoot("object_lock").AtName("default_retention").AtName("days"),
										path.MatchRoot("object_lock").AtName("default_retention").AtName("years"),
									),
								},
							},
							"years": schema.Int32Attribute{
								Optional: true,
							},
						},
					},
				},
			},
			"permissions": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"grants": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"permission": schema.StringAttribute{
									Computed: true,
								},
								"grantee": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Computed: true,
										},
										"display_name": schema.StringAttribute{
											Computed: true,
										},
										"email_address": schema.StringAttribute{
											Computed: true,
										},
										"id": schema.StringAttribute{
											Computed: true,
										},
										"uri": schema.StringAttribute{
											Computed: true,
										},
									},
								},
							},
						},
					},
					"owner": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"display_name": schema.StringAttribute{
								Computed: true,
							},
							"id": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
			"acl": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(lo.Map(s3types.BucketCannedACL("").Values(), func(acl s3types.BucketCannedACL, _ int) string {
						return string(acl)
					})...),
				},
			},
			"grant": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"full_control": schema.SingleNestedAttribute{
						Optional: true,
						Validators: []validator.Object{
							objectvalidator.AtLeastOneOf(
								path.MatchRoot("grant").AtName("read"),
								path.MatchRoot("grant").AtName("read_acp"),
								path.MatchRoot("grant").AtName("write"),
								path.MatchRoot("grant").AtName("write_acp"),
							),
						},
						Attributes: bucketGrantPermissionAttributes("full_control"),
					},
					"read": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: bucketGrantPermissionAttributes("read"),
					},
					"read_acp": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: bucketGrantPermissionAttributes("read_acp"),
					},
					"write": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: bucketGrantPermissionAttributes("write"),
					},
					"write_acp": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: bucketGrantPermissionAttributes("write_acp"),
					},
				},
			},
			"force_delete": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"location_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
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

func bucketGrantPermissionAttributes(permission string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ids": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.AtLeastOneOf(path.MatchRoot("grant").AtName(permission).AtName("email_addresses")),
			},
		},
		"email_addresses": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
	}
}

var bucketSchemaType = bucketSchema(context.Background()).Type()

func (r *bucketResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketSchema(ctx)
}

func (r *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	if !fwhelpers.IsSet(data.Name) {
		data.Name = to.String(id.PrefixedUniqueId("bucket-"))
	}

	inputs := &s3.CreateBucketInput{
		Bucket: data.Name.ValueStringPointer(),
	}

	if fwhelpers.IsSet(data.ACL) {
		inputs.ACL = s3types.BucketCannedACL(data.ACL.ValueString())
	}

	if fwhelpers.IsSet(data.Grant) {
		grantInput, diag := r.expandBucketGrantInput(ctx, data.Grant)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		inputs.GrantFullControl = grantInput.GrantFullControl
		inputs.GrantRead = grantInput.GrantRead
		inputs.GrantReadACP = grantInput.GrantReadACP
		inputs.GrantWrite = grantInput.GrantWrite
		inputs.GrantWriteACP = grantInput.GrantWriteACP
	}

	inputsLock := s3.PutObjectLockConfigurationInput{}

	if fwhelpers.IsSet(data.ObjectLock) {
		lockModel, diag := to.Model[bucketObjectLockModel](ctx, data.ObjectLock)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		inputs.ObjectLockEnabledForBucket = lockModel.Enabled.ValueBoolPointer()

		if fwhelpers.IsSet(lockModel.DefaultRetention) {
			defaultRetentionModel, diag := to.Model[bucketDefaultRetentionModel](ctx, lockModel.DefaultRetention)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}

			inputsLock.Bucket = data.Name.ValueStringPointer()
			inputsLock.ObjectLockConfiguration = &s3types.ObjectLockConfiguration{
				Rule: &s3types.ObjectLockRule{
					DefaultRetention: &s3types.DefaultRetention{
						Mode: s3types.ObjectLockRetentionModeCompliance,
					},
				},
			}

			if fwhelpers.IsSet(lockModel.Enabled) && lockModel.Enabled.ValueBool() {
				inputsLock.ObjectLockConfiguration.ObjectLockEnabled = s3types.ObjectLockEnabledEnabled
			}
			if fwhelpers.IsSet(defaultRetentionModel.Days) {
				inputsLock.ObjectLockConfiguration.Rule.DefaultRetention.Days = defaultRetentionModel.Days.ValueInt32Pointer()
			}
			if fwhelpers.IsSet(defaultRetentionModel.Years) {
				inputsLock.ObjectLockConfiguration.Rule.DefaultRetention.Years = defaultRetentionModel.Years.ValueInt32Pointer()
			}
		}
	}

	bucket := data.Name.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket)

	_, err := r.Client.CreateBucket(ctx, inputs, oos.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(bucketErrCreate, err.Error())
		return
	}
	data.Id = to.String(data.Name.ValueString())

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (s3.PutObjectLockConfigurationInput{}) != inputsLock {
		_, err = r.Client.PutObjectLockConfiguration(ctx, &inputsLock, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(bucketObjectLockErrCreate, err.Error())
			return
		}
	}

	stateData, err = r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	stateData, err := r.read(ctx, timeout, data)
	switch {
	case IsNotFound(err), errors.Is(err, ErrResourceEmpty):
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	var objectLockInput s3.PutObjectLockConfigurationInput
	if fwhelpers.HasChange(plan.ObjectLock, state.ObjectLock) {
		planLock, diag := to.Model[bucketObjectLockModel](ctx, plan.ObjectLock)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		stateLock, diag := to.Model[bucketObjectLockModel](ctx, state.ObjectLock)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		if fwhelpers.HasChange(planLock.DefaultRetention, stateLock.DefaultRetention) {
			defaultRetentionModel, diag := to.Model[bucketDefaultRetentionModel](ctx, planLock.DefaultRetention)
			if fwhelpers.CheckDiags(resp, diag) {
				return
			}

			objectLockInput = s3.PutObjectLockConfigurationInput{
				Bucket: plan.Id.ValueStringPointer(),
				ObjectLockConfiguration: &s3types.ObjectLockConfiguration{
					ObjectLockEnabled: s3types.ObjectLockEnabledEnabled,
					Rule: &s3types.ObjectLockRule{
						DefaultRetention: &s3types.DefaultRetention{
							Mode: s3types.ObjectLockRetentionModeCompliance,
						},
					},
				},
			}

			if fwhelpers.IsSet(defaultRetentionModel.Days) {
				objectLockInput.ObjectLockConfiguration.Rule.DefaultRetention.Days = defaultRetentionModel.Days.ValueInt32Pointer()
			}

			if fwhelpers.IsSet(defaultRetentionModel.Years) {
				objectLockInput.ObjectLockConfiguration.Rule.DefaultRetention.Years = defaultRetentionModel.Years.ValueInt32Pointer()
			}
		}
	}

	var aclInput s3.PutBucketAclInput

	switch {
	case fwhelpers.HasChange(plan.ACL, state.ACL):
		aclInput.Bucket = plan.Id.ValueStringPointer()
		aclInput.ACL = s3types.BucketCannedACL(plan.ACL.ValueString())

	case fwhelpers.HasChange(plan.Grant, state.Grant):
		grantInput, diag := r.expandBucketGrantInput(ctx, plan.Grant)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		grantInput.Bucket = plan.Id.ValueStringPointer()
		aclInput = grantInput

	case (!plan.ACL.Equal(state.ACL) || !plan.Grant.Equal(state.Grant)) && plan.ACL.IsNull() && plan.Grant.IsNull():
		// When both ACL and Grant are set to null, revert to default ACL Permission
		aclInput.Bucket = plan.Id.ValueStringPointer()
		aclInput.ACL = s3types.BucketCannedACLPrivate
	}

	bucket := plan.Id.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket)

	if objectLockInput != (s3.PutObjectLockConfigurationInput{}) {
		_, err := r.Client.PutObjectLockConfiguration(ctx, &objectLockInput, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(bucketErrUpdate, err.Error())
			return
		}
	}

	if aclInput != (s3.PutBucketAclInput{}) {
		_, err := r.Client.PutBucketAcl(ctx, &aclInput, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(bucketErrUpdate, err.Error())
			return
		}
	}

	newData, err := r.read(ctx, timeout, plan)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *bucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	bucket := data.Id.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket)

	forceDelete := fwhelpers.IsSet(data.ForceDelete) && data.ForceDelete.ValueBool()
	if forceDelete {
		err := r.cleanObjects(ctx, timeout, data)
		if err != nil {
			resp.Diagnostics.AddError(bucketErrForceDelete, err.Error())
			return
		}
	}

	input := &s3.DeleteBucketInput{
		Bucket: data.Id.ValueStringPointer(),
	}
	_, err := r.Client.DeleteBucket(ctx, input, oos.WithRetryTimeout(timeout))
	switch {
	case IsNotFound(err):
	case IsErrorCode(err, "BucketNotEmpty") && !forceDelete:
		resp.Diagnostics.AddError(
			bucketErrDelete,
			err.Error()+"\n\nThe bucket still contains objects, object versions, or delete markers. If you intend to permanently delete the bucket and all of its contents, set `force_delete = true` on this resource before destroying it.",
		)
	case err != nil:
		resp.Diagnostics.AddError(bucketErrDelete, err.Error())
	}
}

func (r *bucketResource) cleanObjects(ctx context.Context, timeout time.Duration, data bucketModel) error {
	bucket := data.Id.ValueStringPointer()
	paginator := s3.NewListObjectVersionsPaginator(r.Client, &s3.ListObjectVersionsInput{
		Bucket: bucket,
	})
	var versions []s3types.ObjectVersion
	var objectsIds []s3types.ObjectIdentifier

	for paginator.HasMorePages() {
		objects, err := paginator.NextPage(ctx, oos.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}

		versions = append(versions, objects.Versions...)
		for _, object := range objects.Versions {
			objectsIds = append(objectsIds, s3types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			})
		}
		for _, marker := range objects.DeleteMarkers {
			objectsIds = append(objectsIds, s3types.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			})
		}
	}
	if len(objectsIds) == 0 {
		return nil
	}

	now := time.Now()
	var retentionErr error
	var retentionErrMu sync.Mutex
	var group errgroup.Group

	// Verify that no object has a blocking retention policy before starting deletion
	for _, version := range versions {
		group.Go(func() error {
			output, err := r.Client.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{
				Bucket:    bucket,
				Key:       version.Key,
				VersionId: version.VersionId,
			}, oos.WithRetryTimeout(timeout))
			if err != nil {
				// If the Object Lock Configuration is not found, ignore the error
				if IsErrorMessage(err, "InvalidRequest", "Bucket is missing Object Lock Configuration") {
					return nil
				}

				retentionErrMu.Lock()
				retentionErr = errors.Join(retentionErr, fmt.Errorf(
					"check object %q version %q retention: %w",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					err,
				))
				retentionErrMu.Unlock()
				return nil
			}
			if output.Retention == nil {
				return nil
			}

			var objectErr error
			switch {
			case output.Retention.RetainUntilDate == nil:
				objectErr = fmt.Errorf(
					"object %q version %q has %s retention without an expiration date",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					output.Retention.Mode,
				)
			case now.Before(*output.Retention.RetainUntilDate):
				objectErr = fmt.Errorf(
					"object %q version %q is protected by %s retention until %s",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					output.Retention.Mode,
					output.Retention.RetainUntilDate.Format(time.RFC3339),
				)
			}
			if objectErr == nil {
				return nil
			}

			retentionErrMu.Lock()
			retentionErr = errors.Join(retentionErr, objectErr)
			retentionErrMu.Unlock()
			return nil
		})
	}
	_ = group.Wait()
	if retentionErr != nil {
		return retentionErr
	}

	batchSize := 1000
	for start := 0; start < len(objectsIds); start += batchSize {
		end := min(start+batchSize, len(objectsIds))
		output, err := r.Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: bucket,
			Delete: &s3types.Delete{
				Objects: objectsIds[start:end],
			},
		}, oos.WithRetryTimeout(timeout))
		if err != nil {
			return fmt.Errorf("delete bucket objects: %w", err)
		}

		deleteErrors := lo.Map(output.Errors, func(objectErr s3types.Error, _ int) error {
			return fmt.Errorf(
				"delete object %q version %q: %s: %s",
				ptr.From(objectErr.Key),
				ptr.From(objectErr.VersionId),
				ptr.From(objectErr.Code),
				ptr.From(objectErr.Message),
			)
		})
		if err := errors.Join(deleteErrors...); err != nil {
			return err
		}
	}

	return nil
}

func (r *bucketResource) read(ctx context.Context, timeout time.Duration, data bucketModel) (bucketModel, error) {
	resp, err := r.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: data.Id.ValueStringPointer(),
	}, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	data.LocationName = to.String(resp.BucketLocationName)

	respList, err := r.Client.ListBuckets(ctx, &s3.ListBucketsInput{
		Prefix: new(data.Id.ValueString()),
	}, oos.WithRetryTimeout(timeout))

	bucket, ok := lo.Find(respList.Buckets, func(b s3types.Bucket) bool {
		return ptr.From(b.Name) == data.Id.ValueString()
	})
	if !ok {
		return data, ErrResourceEmpty
	}
	data.CreationDate = to.String(bucket.CreationDate.String())

	input := &s3.GetObjectLockConfigurationInput{
		Bucket: data.Id.ValueStringPointer(),
	}

	respLock, err := r.Client.GetObjectLockConfiguration(ctx, input, oos.WithRetryTimeout(timeout))
	switch {
	case IsNotFound(err):
	case err != nil:
		return data, err
	}

	objLock, diag := r.flattenObjectLock(ctx, respLock)
	if diag != nil {
		return data, diag
	}
	data.ObjectLock = objLock

	respAcl, err := r.Client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: data.Id.ValueStringPointer(),
	}, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	permissions, diag := r.flattenPermissions(ctx, respAcl)
	if diag != nil {
		return data, diag
	}
	data.Permissions = permissions

	return data, nil
}

func (r *bucketResource) flattenObjectLock(ctx context.Context, lockOutput *s3.GetObjectLockConfigurationOutput) (types.Object, error) {
	var obj types.Object
	model := bucketObjectLockModel{
		DefaultRetention: to.ObjectNull[bucketDefaultRetentionModel](),
	}

	if lockOutput == nil || lockOutput.ObjectLockConfiguration == nil {
		model.Enabled = to.Bool(false)
	} else {
		model.Enabled = to.Bool(true)

		if lockOutput.ObjectLockConfiguration.Rule != nil && lockOutput.ObjectLockConfiguration.Rule.DefaultRetention != nil {
			retentionModel := bucketDefaultRetentionModel{
				Days:  to.Int32(lockOutput.ObjectLockConfiguration.Rule.DefaultRetention.Days),
				Years: to.Int32(lockOutput.ObjectLockConfiguration.Rule.DefaultRetention.Years),
				Mode:  to.String(lockOutput.ObjectLockConfiguration.Rule.DefaultRetention.Mode),
			}

			retention, diag := to.Object(ctx, retentionModel)
			if diag.HasError() {
				return obj, from.Diag(diag)
			}
			model.DefaultRetention = retention
		}
	}

	obj, diag := to.Object(ctx, model)
	if diag.HasError() {
		return obj, from.Diag(diag)
	}

	return obj, nil
}

func (r *bucketResource) flattenPermissions(ctx context.Context, acl *s3.GetBucketAclOutput) (types.Object, error) {
	if acl == nil {
		return to.ObjectNull[bucketPermissionsModel](), nil
	}

	var obj types.Object
	model := bucketPermissionsModel{
		Owner: to.ObjectNull[bucketOwnerModel](),
	}

	grantObjects, err := lo.MapErr(acl.Grants, func(grant s3types.Grant, _ int) (types.Object, error) {
		var grantObject types.Object
		grantModel := bucketGrantsModel{
			Permission: to.String(grant.Permission),
			Grantee:    to.ObjectNull[bucketGranteeModel](),
		}

		if grant.Grantee != nil {
			granteeModel := bucketGranteeModel{
				Type:         to.String(grant.Grantee.Type),
				DisplayName:  to.String(grant.Grantee.DisplayName),
				EmailAddress: to.String(grant.Grantee.EmailAddress),
				Id:           to.String(grant.Grantee.ID),
				Uri:          to.String(grant.Grantee.URI),
			}

			granteeObj, diag := to.Object(ctx, granteeModel)
			if diag.HasError() {
				return grantObject, from.Diag(diag)
			}
			grantModel.Grantee = granteeObj
		}

		grantObject, diag := to.Object(ctx, grantModel)
		if diag.HasError() {
			return grantObject, from.Diag(diag)
		}

		return grantObject, nil
	})
	if err != nil {
		return obj, err
	}

	uniqueGrants := make([]attr.Value, 0, len(grantObjects))
	for _, grant := range grantObjects {
		if !lo.ContainsBy(uniqueGrants, grant.Equal) {
			uniqueGrants = append(uniqueGrants, grant)
		}
	}

	grants, diag := to.SetFromAttrType(ctx, uniqueGrants, to.ObjType(bucketGrantsModel{}.AttributeTypes()))
	if diag.HasError() {
		return obj, from.Diag(diag)
	}
	model.Grants = grants

	if acl.Owner != nil {
		ownerModel := bucketOwnerModel{
			DisplayName: to.String(acl.Owner.DisplayName),
			Id:          to.String(acl.Owner.ID),
		}

		ownerObj, diag := to.Object(ctx, ownerModel)
		if diag.HasError() {
			return obj, from.Diag(diag)
		}
		model.Owner = ownerObj
	}

	obj, diag = to.Object(ctx, model)
	if diag.HasError() {
		return obj, from.Diag(diag)
	}

	return obj, nil
}

func (r *bucketResource) expandBucketGrantInput(ctx context.Context, grant types.Object) (s3.PutBucketAclInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	var input s3.PutBucketAclInput

	grantModel, modelDiags := to.Model[bucketGrantModel](ctx, grant)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return input, diags
	}

	format := func(value types.Object) *string {
		if !fwhelpers.IsSet(value) {
			return nil
		}

		permission, diag := to.Model[bucketGrantPermissionModel](ctx, value)
		diags.Append(diag...)
		if diag.HasError() {
			return nil
		}

		var grantees []string
		if fwhelpers.IsSet(permission.Ids) {
			ids, diag := to.Slice[string](ctx, permission.Ids)
			diags.Append(diag...)

			grantees = append(grantees, lo.Map(ids, func(id string, _ int) string {
				return "id=" + id
			})...)
		}
		if fwhelpers.IsSet(permission.EmailAddresses) {
			emails, diag := to.Slice[string](ctx, permission.EmailAddresses)
			diags.Append(diag...)

			grantees = append(grantees, lo.Map(emails, func(email string, _ int) string {
				return "emailaddress=" + email
			})...)
		}
		if diags.HasError() || len(grantees) == 0 {
			return nil
		}

		return new(strings.Join(grantees, ","))
	}

	input.GrantFullControl = format(grantModel.FullControl)
	input.GrantRead = format(grantModel.Read)
	input.GrantReadACP = format(grantModel.ReadACP)
	input.GrantWrite = format(grantModel.Write)
	input.GrantWriteACP = format(grantModel.WriteACP)

	return input, diags
}
