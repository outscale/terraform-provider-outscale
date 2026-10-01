package oos

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource              = &bucketVersioningResource{}
	_ resource.ResourceWithConfigure = &bucketVersioningResource{}
)

const (
	bucketVersioningErrCreate = "Unable to create Bucket Versioning"
	bucketVersioningErrUpdate = "Unable to update Bucket Versioning"
	bucketVersioningErrDelete = "Unable to delete Bucket Versioning"
)

type bucketVersioningModel struct {
	Bucket   types.String   `tfsdk:"bucket"`
	Status   types.String   `tfsdk:"status"`
	Id       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type bucketVersioningResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucketVersioning() resource.Resource {
	return &bucketVersioningResource{}
}

func (r *bucketVersioningResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketVersioningResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket_versioning"
}

func (r *bucketVersioningResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"status": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(lo.Map(s3types.BucketVersioningStatus("").Values(), func(status s3types.BucketVersioningStatus, _ int) string {
						return string(status)
					})...),
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

func (r *bucketVersioningResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketVersioningModel
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
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	if err := r.put(ctx, timeout, data); err != nil {
		resp.Diagnostics.AddError(bucketVersioningErrCreate, err.Error())
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

func (r *bucketVersioningResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketVersioningModel
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
	case IsNotFound(err):
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketVersioningResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketVersioningModel
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
	defer r.mu.UnlockKey(bucket) //nolint: errcheck

	if err := r.put(ctx, timeout, data); err != nil {
		resp.Diagnostics.AddError(bucketVersioningErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketVersioningResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketVersioningModel
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

	data.Status = to.String(s3types.BucketVersioningStatusSuspended)

	if err := r.put(ctx, timeout, data); err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(bucketVersioningErrDelete, err.Error())
	}
}

func (r *bucketVersioningResource) put(ctx context.Context, timeout time.Duration, data bucketVersioningModel) error {
	input := &s3.PutBucketVersioningInput{
		Bucket: data.Bucket.ValueStringPointer(),
		VersioningConfiguration: &s3types.VersioningConfiguration{
			Status: s3types.BucketVersioningStatus(data.Status.ValueString()),
		},
	}

	_, err := r.Client.PutBucketVersioning(ctx, input, oos.WithRetryTimeout(timeout))

	return err
}

func (r *bucketVersioningResource) read(ctx context.Context, timeout time.Duration, data bucketVersioningModel) (bucketVersioningModel, error) {
	input := &s3.GetBucketVersioningInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	output, err := r.Client.GetBucketVersioning(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if output == nil {
		return data, ErrEmptyResponse
	}

	data.Status = to.String(output.Status)
	data.Id = to.String(data.Bucket.ValueString())

	return data, nil
}
