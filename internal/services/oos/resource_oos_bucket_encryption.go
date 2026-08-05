package oos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
	"k8s.io/utils/keymutex"
)

var (
	_ resource.Resource              = &bucketEncryptionResource{}
	_ resource.ResourceWithConfigure = &bucketEncryptionResource{}
)

const (
	bucketEncryptionErrCreate = "Unable to create Bucket Encryption"
	bucketEncryptionErrUpdate = "Unable to update Bucket Encryption"
	bucketEncryptionErrDelete = "Unable to delete Bucket Encryption"
)

type bucketEncryptionModel struct {
	Bucket           types.String   `tfsdk:"bucket"`
	Type             types.String   `tfsdk:"encryption_type"`
	BucketKeyEnabled types.Bool     `tfsdk:"bucket_key_enabled"`
	Id               types.String   `tfsdk:"id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type bucketEncryptionResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucketEncryption() resource.Resource {
	return &bucketEncryptionResource{}
}

func (r *bucketEncryptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketEncryptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket_encryption"
}

func (r *bucketEncryptionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"encryption_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(string(s3types.ServerSideEncryptionAes256)),
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.ServerSideEncryptionAes256), string(s3types.ServerSideEncryptionAwsKms)),
				},
			},
			"bucket_key_enabled": schema.BoolAttribute{
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

func (r *bucketEncryptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketEncryptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	bucket := data.Bucket.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket)

	if err := r.put(ctx, timeout, data); err != nil {
		resp.Diagnostics.AddError(bucketEncryptionErrCreate, err.Error())
		return
	}

	data.Id = to.String(data.Bucket.ValueString())

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketEncryptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketEncryptionModel
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

func (r *bucketEncryptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketEncryptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	bucket := data.Bucket.ValueString()
	r.mu.LockKey(bucket)
	defer r.mu.UnlockKey(bucket)

	if err := r.put(ctx, timeout, data); err != nil {
		resp.Diagnostics.AddError(bucketEncryptionErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketEncryptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketEncryptionModel
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
	defer r.mu.UnlockKey(bucket)

	input := &s3.DeleteBucketEncryptionInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	_, err := r.Client.DeleteBucketEncryption(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(bucketEncryptionErrDelete, err.Error())
	}
}

func (r *bucketEncryptionResource) put(ctx context.Context, timeout time.Duration, data bucketEncryptionModel) error {
	input := &s3.PutBucketEncryptionInput{
		Bucket: data.Bucket.ValueStringPointer(),
		ServerSideEncryptionConfiguration: &s3types.ServerSideEncryptionConfiguration{
			Rules: []s3types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &s3types.ServerSideEncryptionByDefault{
						SSEAlgorithm: s3types.ServerSideEncryption(data.Type.ValueString()),
					},
				},
			},
		},
	}

	_, err := r.Client.PutBucketEncryption(ctx, input, oos.WithRetryTimeout(timeout))
	return err
}

func (r *bucketEncryptionResource) read(ctx context.Context, timeout time.Duration, data bucketEncryptionModel) (bucketEncryptionModel, error) {
	input := &s3.GetBucketEncryptionInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	output, err := r.Client.GetBucketEncryption(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if output.ServerSideEncryptionConfiguration == nil {
		return data, ErrResourceEmpty
	}

	rule, ok := lo.Find(output.ServerSideEncryptionConfiguration.Rules, func(rule s3types.ServerSideEncryptionRule) bool {
		return rule.ApplyServerSideEncryptionByDefault != nil
	})
	if !ok {
		return data, ErrResourceEmpty
	}

	data.Type = to.String(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	data.BucketKeyEnabled = to.Bool(rule.BucketKeyEnabled)
	data.Id = to.String(data.Bucket.ValueString())

	return data, nil
}
