package oos

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                     = &objectCopyResource{}
	_ resource.ResourceWithConfigure        = &objectCopyResource{}
	_ resource.ResourceWithConfigValidators = &objectCopyResource{}
	_ resource.ResourceWithModifyPlan       = &objectCopyResource{}
)

const (
	objectCopyErrCreate      = "Unable to copy Object"
	objectCopyErrUpdate      = "Unable to update copied Object"
	objectCopyErrDelete      = "Unable to delete copied Object"
	objectCopyErrForceDelete = "Unable to force delete copied Object"
)

type objectCopyModel struct {
	objectCommonModel
	CopySource        types.Object `tfsdk:"copy_source"`
	MetadataDirective types.String `tfsdk:"metadata_directive"`
	TaggingDirective  types.String `tfsdk:"tagging_directive"`
}

type objectCopySourceModel struct {
	Path              types.String `tfsdk:"path"`
	IfMatch           types.String `tfsdk:"if_match"`
	IfModifiedSince   types.String `tfsdk:"if_modified_since"`
	IfNoneMatch       types.String `tfsdk:"if_none_match"`
	IfUnmodifiedSince types.String `tfsdk:"if_unmodified_since"`
	VersionId         types.String `tfsdk:"version_id"`
}

type objectCopyResource struct {
	objectCommon
}

func NewResourceObjectCopy() resource.Resource {
	return &objectCopyResource{}
}

func (r *objectCopyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configuredClient, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *oos.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Client = configuredClient.OOS
	r.mu = configuredClient.KeyMutex
}

func (r *objectCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_object_copy"
}

func (r *objectCopyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.modifyPlan(ctx, req, resp)
}

func (r *objectCopyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("acl"),
			path.MatchRoot("grant"),
		),
	}
}

func (r *objectCopyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"permissions": permissionsAttributes(),
			"acl": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("private", "public-read", "public-read-write", "authenticated-read"),
				},
			},
			"grant": grantAttributes(),
			"cache_control": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("max-age", "max-stale", "min-fresh", "no-cache", "no-store", "no-transform", "only-if-cached", "stale-if-error"),
				},
			},
			"content_disposition": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"content_encoding": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("gzip", "compress", "deflate", "identity", "br"),
				},
			},
			"content_language": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"content_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"copy_source": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(3),
						},
					},
					"if_match": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"if_modified_since": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"if_none_match": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"if_unmodified_since": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"version_id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"expires": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"metadata_directive": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.MetadataDirectiveCopy), string(s3types.MetadataDirectiveReplace)),
				},
			},
			"encryption_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.ServerSideEncryptionAes256)),
				},
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.SizeAtMost(50),
					mapvalidator.KeysAre(stringvalidator.LengthBetween(1, 128)),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtMost(256)),
				},
			},
			"tagging_directive": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.TaggingDirectiveCopy), string(s3types.TaggingDirectiveReplace)),
				},
			},
			"force_delete": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"expiration": schema.StringAttribute{
				Computed: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"last_modified": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_latest": schema.BoolAttribute{
				Computed: true,
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

func (r *objectCopyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data objectCopyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	input := s3.CopyObjectInput{
		Bucket: data.Bucket.ValueStringPointer(),
		Key:    data.Key.ValueStringPointer(),
	}

	copySourceModel, diags := to.Model[objectCopySourceModel](ctx, data.CopySource)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	input.CopySourceIfMatch = copySourceModel.IfMatch.ValueStringPointer()
	input.CopySourceIfNoneMatch = copySourceModel.IfNoneMatch.ValueStringPointer()
	input.CopySource = new(url.QueryEscape(copySourceModel.Path.ValueString()))

	if fwhelpers.IsSet(data.CacheControl) {
		input.CacheControl = data.CacheControl.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.ContentDisposition) {
		input.ContentDisposition = data.ContentDisposition.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.ContentEncoding) {
		input.ContentEncoding = data.ContentEncoding.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.ContentLanguage) {
		input.ContentLanguage = data.ContentLanguage.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.ContentType) {
		input.ContentType = data.ContentType.ValueStringPointer()
	}

	modifiedSince, err := from.RFC3339(copySourceModel.IfModifiedSince)
	if err != nil {
		diags.AddError(objectCopyErrCreate, err.Error())
	}
	input.CopySourceIfModifiedSince = modifiedSince

	unmodifiedSince, err := from.RFC3339(copySourceModel.IfUnmodifiedSince)
	if err != nil {
		diags.AddError(objectCopyErrCreate, err.Error())
	}
	input.CopySourceIfUnmodifiedSince = unmodifiedSince

	expires, err := from.RFC3339(data.Expires)
	if err != nil {
		diags.AddError(objectCopyErrCreate, err.Error())
	}
	input.Expires = expires

	if fwhelpers.IsSet(data.Metadata) {
		metadata, d := from.Map[string](ctx, data.Metadata)
		diags.Append(d...)

		input.Metadata = metadata
	}
	if fwhelpers.IsSet(data.Tags) {
		tags, d := from.Map[string](ctx, data.Tags)
		diags.Append(d...)

		values := url.Values{}
		for key, value := range tags {
			values.Set(key, value)
		}
		encoded := values.Encode()

		input.Tagging = &encoded
	}
	if fwhelpers.IsSet(data.MetadataDirective) {
		input.MetadataDirective = s3types.MetadataDirective(data.MetadataDirective.ValueString())
	}
	if fwhelpers.IsSet(data.TaggingDirective) {
		input.TaggingDirective = s3types.TaggingDirective(data.TaggingDirective.ValueString())
	}
	if fwhelpers.IsSet(data.EncryptionType) {
		input.ServerSideEncryption = s3types.ServerSideEncryption(data.EncryptionType.ValueString())
	}

	before, key, found := strings.CutLast(copySourceModel.Path.ValueString(), "/")
	switch {
	case !found:
		resp.Diagnostics.AddAttributeError(path.Empty().AtName("copy_source").AtName("path"), objectCopyErrCreate, "The path to the source bucket containing the object you want to copy is not in the appropriate format")
	case key == "":
		key = before
	}
	if resp.Diagnostics.HasError() {
		return
	}

	r.mu.LockKey(key)
	defer r.mu.UnlockKey(key) //nolint: errcheck

	output, err := r.Client.CopyObject(ctx, &input, oos.WithRetryTimeout(timeout))
	switch {
	case err != nil:
		resp.Diagnostics.AddError(objectCopyErrCreate, err.Error())
		return
	case output == nil:
		resp.Diagnostics.AddError(objectCopyErrCreate, ErrEmptyResponse.Error())
		return
	}

	data.Id = to.String(data.Key.ValueString())
	data.VersionId = to.String(output.VersionId)
	data.Expiration = to.String(output.Expiration)
	data.EncryptionType = to.StringEnum(output.ServerSideEncryption)
	copySourceModel.VersionId = to.String(output.CopySourceVersionId)

	copySource, d := to.Object(ctx, copySourceModel)
	if fwhelpers.CheckDiags(resp, d) {
		return
	}
	data.CopySource = copySource

	if fwhelpers.IsSet(data.Grant) || fwhelpers.IsSet(data.ACL) {
		aclInput, diags := r.expandGrant(ctx, data.objectCommonModel)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}

		_, err := r.Client.PutObjectAcl(ctx, &aclInput, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(objectCopyErrCreate, err.Error())
			return
		}
	}

	stateData, err := r.read(ctx, timeout, data.objectCommonModel)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	data.objectCommonModel = stateData
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *objectCopyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data objectCopyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data.objectCommonModel)
	switch {
	case IsNotFound(err), errors.Is(err, ErrEmptyResponse):
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	data.objectCommonModel = stateData
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *objectCopyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectCopyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	r.mu.LockKey(state.Key.ValueString())
	defer r.mu.UnlockKey(state.Key.ValueString()) //nolint: errcheck

	if !plan.ACL.Equal(state.ACL) || !plan.Grant.Equal(state.Grant) {
		input, diags := r.expandGrant(ctx, plan.objectCommonModel)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}
		if _, err := r.Client.PutObjectAcl(ctx, &input, oos.WithRetryTimeout(timeout)); err != nil {
			resp.Diagnostics.AddError(objectCopyErrUpdate, err.Error())
			return
		}
	}
	if !plan.Tags.Equal(state.Tags) {
		diags := r.putTags(ctx, timeout, plan.objectCommonModel)
		if fwhelpers.CheckDiags(resp, diags) {
			return
		}
	}

	stateData, err := r.read(ctx, timeout, plan.objectCommonModel)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	plan.objectCommonModel = stateData
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *objectCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data objectCopyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	bucket := data.Bucket.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	if fwhelpers.IsSet(data.ForceDelete) && data.ForceDelete.ValueBool() {
		if err := r.cleanObjects(ctx, timeout, bucket, data.Key.ValueString()); err != nil {
			resp.Diagnostics.AddError(objectErrForceDelete, err.Error())
		}
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}

	_, err := r.Client.DeleteObject(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(objectErrDelete, err.Error())
	}
}
