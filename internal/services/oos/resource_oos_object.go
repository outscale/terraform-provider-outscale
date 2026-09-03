package oos

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                     = &objectResource{}
	_ resource.ResourceWithConfigure        = &objectResource{}
	_ resource.ResourceWithConfigValidators = &objectResource{}
	_ resource.ResourceWithModifyPlan       = &objectResource{}
)

const (
	objectErrCreate      = "Unable to create Object"
	objectErrUpdate      = "Unable to update Object"
	objectErrDelete      = "Unable to delete Object"
	objectErrForceDelete = "Unable to force delete Object"
	objectErrPutTags     = "Unable to put tags on Object"
)

type objectModel struct {
	Bucket             types.String      `tfsdk:"bucket"`
	Key                types.String      `tfsdk:"key"`
	Content            types.String      `tfsdk:"content"`
	ContentB64         types.String      `tfsdk:"content_base64"`
	Source             types.String      `tfsdk:"source"`
	ACL                types.String      `tfsdk:"acl"`
	Grant              types.Object      `tfsdk:"grant"`
	Permissions        types.Object      `tfsdk:"permissions"`
	CacheControl       types.String      `tfsdk:"cache_control"`
	ContentDisposition types.String      `tfsdk:"content_disposition"`
	ContentEncoding    types.String      `tfsdk:"content_encoding"`
	ContentLanguage    types.String      `tfsdk:"content_language"`
	ContentType        types.String      `tfsdk:"content_type"`
	Expires            types.String      `tfsdk:"expires"`
	Metadata           types.Map         `tfsdk:"metadata"`
	EncryptionType     types.String      `tfsdk:"encryption_type"`
	Tags               types.Map         `tfsdk:"tags"`
	ForceDelete        types.Bool        `tfsdk:"force_delete"`
	ContentLength      types.Int64       `tfsdk:"content_length"`
	Expiration         types.String      `tfsdk:"expiration"`
	ETag               types.String      `tfsdk:"etag"`
	LastModified       timetypes.RFC3339 `tfsdk:"last_modified"`
	VersionId          types.String      `tfsdk:"version_id"`
	IsLatest           types.Bool        `tfsdk:"is_latest"`
	Id                 types.String      `tfsdk:"id"`
	Timeouts           timeouts.Value    `tfsdk:"timeouts"`
}

type objectResource struct {
	resourceCommon
}

func NewResourceObject() resource.Resource {
	return &objectResource{}
}

func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_object"
}

func (r *objectResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || r.Client == nil {
		return
	}

	var bucket, key types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("bucket"), &bucket)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("key"), &key)...)
	if resp.Diagnostics.HasError() || !fwhelpers.IsSet(bucket) || !fwhelpers.IsSet(key) {
		return
	}

	input := s3.HeadObjectInput{
		Bucket: bucket.ValueStringPointer(),
		Key:    key.ValueStringPointer(),
	}
	output, err := r.Client.HeadObject(ctx, &input, oos.WithRetryTimeout(ReadDefaultTimeout))
	if err != nil || output == nil {
		return
	}

	inputVersioning := s3.GetBucketVersioningInput{
		Bucket: bucket.ValueStringPointer(),
	}
	versioning, err := r.Client.GetBucketVersioning(ctx, &inputVersioning, oos.WithRetryTimeout(ReadDefaultTimeout))
	if err != nil || versioning == nil || versioning.Status != s3types.BucketVersioningStatusEnabled {
		return
	}

	resp.Diagnostics.AddAttributeWarning(
		path.Root("key"),
		"Existing Object Will Be Overwritten",
		fmt.Sprintf(
			"Object %q already exists in bucket %q, where versioning is not enabled. Applying this plan will overwrite the existing object with PutObject, and its current contents may not be recoverable. You may enable bucket versioning or use a different key before applying.",
			key.ValueString(),
			bucket.ValueString(),
		),
	)
}

func (r *objectResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("acl"),
			path.MatchRoot("grant"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("content"),
			path.MatchRoot("content_base64"),
			path.MatchRoot("source"),
		),
	}
}

func (r *objectResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"content": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"content_base64": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"source": schema.StringAttribute{
				Optional: true,
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
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("max-age", "max-stale", "min-fresh", "no-cache", "no-store", "no-transform", "only-if-cached", "stale-if-error"),
				},
			},
			"content_disposition": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_encoding": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("gzip", "compress", "deflate", "identity", "br"),
				},
			},
			"content_language": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"expires": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"encryption_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.ServerSideEncryptionAes256)),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.SizeAtMost(50),
					mapvalidator.KeysAre(
						stringvalidator.LengthBetween(1, 128),
					),
					mapvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(256),
					),
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

func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data objectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	var body io.ReadSeeker

	switch {
	case fwhelpers.IsSet(data.Content):
		body = strings.NewReader(data.Content.ValueString())
	case fwhelpers.IsSet(data.ContentB64):
		decoded, err := base64.StdEncoding.DecodeString(data.ContentB64.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(objectErrCreate, err.Error())
		}
		body = bytes.NewReader(decoded)
	case fwhelpers.IsSet(data.Source):
		path := data.Source.ValueString()

		if strings.HasPrefix(path, "~/") {
			dirname, _ := os.UserHomeDir()
			path = filepath.Join(dirname, path[2:])
		}

		file, err := os.Open(path) //nolint: gosec
		if err != nil {
			resp.Diagnostics.AddError(objectErrCreate, err.Error())
		}

		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				resp.Diagnostics.AddWarning("Unable to close source file at path: "+path, err.Error())
			}
		}()
	}

	bucket := data.Bucket.ValueString()
	input := &s3.PutObjectInput{
		Body:   body,
		Bucket: data.Bucket.ValueStringPointer(),
		Key:    data.Key.ValueStringPointer(),
	}
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

	if fwhelpers.IsSet(data.Expires) {
		expires, err := time.Parse(time.RFC3339, data.Expires.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(objectErrCreate, err.Error())
		}

		input.Expires = &expires
	}
	if fwhelpers.IsSet(data.Metadata) {
		diag := data.Metadata.ElementsAs(ctx, &input.Metadata, false)
		resp.Diagnostics.Append(diag...)
	}
	if fwhelpers.IsSet(data.EncryptionType) {
		input.ServerSideEncryption = s3types.ServerSideEncryption(data.EncryptionType.ValueString())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	uploader := oos.NewUploader(r.Client)

	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	output, err := uploader.Upload(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(objectErrCreate, err.Error())
		return
	}
	if output == nil {
		resp.Diagnostics.AddError(objectErrCreate, ErrEmptyResponse.Error())
		return
	}

	data.VersionId = to.String(output.VersionID)
	data.Expiration = to.String(output.Expiration)

	// PutObjectInput does not have the same parameters as PutObjectAcl, so we need to make two separate calls...
	if fwhelpers.IsSet(data.Grant) || fwhelpers.IsSet(data.ACL) {
		input, diag := r.expandGrant(ctx, data)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		_, err := r.Client.PutObjectAcl(ctx, &input, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(objectErrCreate, err.Error())
			return
		}
	}

	if fwhelpers.IsSet(data.Tags) {
		diag := r.putTags(ctx, timeout, data)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	data.Id = to.String(data.Key.ValueString())

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data objectModel
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

func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	bucket := plan.Bucket.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	if !plan.ACL.Equal(state.ACL) || !plan.Grant.Equal(state.Grant) {
		input, diag := r.expandGrant(ctx, plan)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		_, err := r.Client.PutObjectAcl(ctx, &input, oos.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(objectErrUpdate, err.Error())
			return
		}
	}

	if !plan.Tags.Equal(state.Tags) {
		diag := r.putTags(ctx, timeout, plan)
		if fwhelpers.CheckDiags(resp, diag) {
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

func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data objectModel
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

func (r *objectResource) putTags(ctx context.Context, timeout time.Duration, data objectModel) diag.Diagnostics {
	var diags diag.Diagnostics

	tags, diag := from.Map[string](ctx, data.Tags)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	if len(tags) == 0 {
		input := s3.DeleteObjectTaggingInput{
			Bucket:    data.Bucket.ValueStringPointer(),
			Key:       data.Key.ValueStringPointer(),
			VersionId: data.VersionId.ValueStringPointer(),
		}

		_, err := r.Client.DeleteObjectTagging(ctx, &input, oos.WithRetryTimeout(timeout))
		if err != nil && !IsNotFound(err) {
			diags.AddError(objectErrPutTags, err.Error())
		}

		return diags
	}

	tagSet := lo.MapToSlice(tags, func(key string, value string) s3types.Tag {
		return s3types.Tag{
			Key:   &key,
			Value: &value,
		}
	})
	input := s3.PutObjectTaggingInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
		Tagging: &s3types.Tagging{
			TagSet: tagSet,
		},
	}
	_, err := r.Client.PutObjectTagging(ctx, &input, oos.WithRetryTimeout(timeout))
	if err != nil {
		diags.AddError(objectErrPutTags, err.Error())
	}

	return diags
}

func (r *objectResource) expandGrant(ctx context.Context, data objectModel) (s3.PutObjectAclInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	input := s3.PutObjectAclInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}

	switch {
	case fwhelpers.IsSet(data.ACL):
		input.ACL = s3types.ObjectCannedACL(data.ACL.ValueString())
	case fwhelpers.IsSet(data.Grant):
		headers, diag := r.expandGrantHeaders(ctx, data.Grant)
		diags.Append(diag...)

		input.GrantFullControl = headers.FullControl
		input.GrantRead = headers.Read
		input.GrantReadACP = headers.ReadACP
		input.GrantWrite = headers.Write
		input.GrantWriteACP = headers.WriteACP
	default:
		input.ACL = s3types.ObjectCannedACLPrivate
	}

	return input, diags
}

func (r *objectResource) read(ctx context.Context, timeout time.Duration, data objectModel) (objectModel, error) {
	input := s3.HeadObjectInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}
	output, err := r.Client.HeadObject(ctx, &input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if output == nil {
		return data, ErrEmptyResponse
	}

	metadata, diags := to.Map(ctx, output.Metadata)
	if diags.HasError() {
		return data, from.Diag(diags)
	}
	data.Metadata = metadata

	data.CacheControl = to.String(output.CacheControl)
	data.ContentDisposition = to.String(output.ContentDisposition)
	data.ContentEncoding = to.String(output.ContentEncoding)
	data.ContentType = to.String(output.ContentType)
	data.ContentLength = to.Int64(output.ContentLength)
	data.ETag = to.String(output.ETag)
	data.LastModified = to.RFC3339(output.LastModified)
	data.VersionId = to.String(output.VersionId)
	data.Id = to.String(data.Key.ValueString())
	data.Expiration = to.String(output.Expiration)
	data.EncryptionType = to.StringEnum(output.ServerSideEncryption)
	// OOS API does not return the value of ContentLanguage, in case it does, we set it so Terraform can detect changes
	if output.ContentLanguage != nil {
		data.ContentLanguage = to.String(output.ContentLanguage)
	}

	inputAcl := s3.GetObjectAclInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}
	acl, err := r.Client.GetObjectAcl(ctx, &inputAcl, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if acl == nil {
		return data, ErrEmptyResponse
	}

	permissions, err := r.flattenPermissions(ctx, acl.Owner, acl.Grants)
	if err != nil {
		return data, err
	}
	data.Permissions = permissions

	inputTags := s3.GetObjectTaggingInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}
	tagging, err := r.Client.GetObjectTagging(ctx, &inputTags, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	tags := lo.SliceToMap(tagging.TagSet, func(tag s3types.Tag) (string, string) {
		return lo.FromPtr(tag.Key), lo.FromPtr(tag.Value)
	})
	tagsMap, diags := to.Map(ctx, tags)
	if diags.HasError() {
		return data, from.Diag(diags)
	}
	data.Tags = tagsMap

	versionID := data.VersionId.ValueString()
	if data.VersionId.IsNull() {
		// ListObjectVersions returns an unversioned object with a "null" string as a versionId
		versionID = "null"
	}

	paginator := s3.NewListObjectVersionsPaginator(r.Client, &s3.ListObjectVersionsInput{
		Bucket: data.Bucket.ValueStringPointer(),
		Prefix: data.Key.ValueStringPointer(),
	})
	for paginator.HasMorePages() {
		versions, err := paginator.NextPage(ctx, oos.WithRetryTimeout(timeout))
		if err != nil {
			return data, err
		}

		version, found := lo.Find(versions.Versions, func(v s3types.ObjectVersion) bool {
			return ptr.From(v.Key) == data.Key.ValueString() && ptr.From(v.VersionId) == versionID
		})
		if found {
			data.IsLatest = to.Bool(version.IsLatest)
			break
		}
	}

	return data, nil
}
