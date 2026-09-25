package oos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

type bucketResource struct {
	resourceCommon
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
			"permissions": permissionsAttributes(),
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
			"grant": grantAttributes(),
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
		headers, d := r.expandGrantHeaders(ctx, data.Grant)
		if fwhelpers.CheckDiags(resp, d) {
			return
		}

		inputs.GrantFullControl = headers.FullControl
		inputs.GrantRead = headers.Read
		inputs.GrantReadACP = headers.ReadACP
		inputs.GrantWrite = headers.Write
		inputs.GrantWriteACP = headers.WriteACP
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
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

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
	case IsNotFound(err), errors.Is(err, ErrResourceNotFound):
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
		headers, d := r.expandGrantHeaders(ctx, plan.Grant)
		if fwhelpers.CheckDiags(resp, d) {
			return
		}

		aclInput.Bucket = plan.Id.ValueStringPointer()
		aclInput.GrantFullControl = headers.FullControl
		aclInput.GrantRead = headers.Read
		aclInput.GrantReadACP = headers.ReadACP
		aclInput.GrantWrite = headers.Write
		aclInput.GrantWriteACP = headers.WriteACP

	case (!plan.ACL.Equal(state.ACL) || !plan.Grant.Equal(state.Grant)) && plan.ACL.IsNull() && plan.Grant.IsNull():
		// When both ACL and Grant are set to null, revert to default ACL Permission
		aclInput.Bucket = plan.Id.ValueStringPointer()
		aclInput.ACL = s3types.BucketCannedACLPrivate
	}

	bucket := plan.Id.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

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
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	forceDelete := fwhelpers.IsSet(data.ForceDelete) && data.ForceDelete.ValueBool()
	if forceDelete {
		err := r.cleanObjects(ctx, timeout, bucket)
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

func (r *bucketResource) read(ctx context.Context, timeout time.Duration, data bucketModel) (bucketModel, error) {
	resp, err := r.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: data.Id.ValueStringPointer(),
	}, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp == nil {
		return data, ErrEmptyResponse
	}

	data.LocationName = to.String(resp.BucketLocationName)

	respList, err := r.Client.ListBuckets(ctx, &s3.ListBucketsInput{
		Prefix: new(data.Id.ValueString()),
	}, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if respList == nil {
		return data, ErrEmptyResponse
	}

	bucket, ok := lo.Find(respList.Buckets, func(b s3types.Bucket) bool {
		return ptr.From(b.Name) == data.Id.ValueString()
	})
	if !ok {
		return data, ErrResourceNotFound
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
	if respAcl == nil {
		return data, ErrEmptyResponse
	}

	permissions, diag := r.flattenPermissions(ctx, respAcl.Owner, respAcl.Grants)
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
