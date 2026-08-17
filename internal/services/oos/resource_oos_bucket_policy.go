package oos

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
	"k8s.io/utils/keymutex"
)

var (
	_ resource.Resource              = &bucketPolicyResource{}
	_ resource.ResourceWithConfigure = &bucketPolicyResource{}
)

const (
	bucketPolicyErrCreate = "Unable to create Bucket Policy"
	bucketPolicyErrUpdate = "Unable to update Bucket Policy"
	bucketPolicyErrDelete = "Unable to delete Bucket Policy"
)

type bucketPolicyModel struct {
	Bucket   types.String   `tfsdk:"bucket"`
	Policy   types.String   `tfsdk:"policy"`
	Id       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type bucketPolicyResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucketPolicy() resource.Resource {
	return &bucketPolicyResource{}
}

func (r *bucketPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket_policy"
}

func (r *bucketPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"policy": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validatorstring.IsJSON(),
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

func (r *bucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketPolicyModel
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
		resp.Diagnostics.AddError(bucketPolicyErrCreate, err.Error())
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

func (r *bucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketPolicyModel
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

func (r *bucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketPolicyModel
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
		resp.Diagnostics.AddError(bucketPolicyErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketPolicyModel
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

	input := &s3.DeleteBucketPolicyInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	_, err := r.Client.DeleteBucketPolicy(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(bucketPolicyErrDelete, err.Error())
	}
}

func (r *bucketPolicyResource) put(ctx context.Context, timeout time.Duration, data bucketPolicyModel) error {
	input := &s3.PutBucketPolicyInput{
		Bucket: data.Bucket.ValueStringPointer(),
		Policy: data.Policy.ValueStringPointer(),
	}

	_, err := r.Client.PutBucketPolicy(ctx, input, oos.WithRetryTimeout(timeout))

	return err
}

func (r *bucketPolicyResource) read(ctx context.Context, timeout time.Duration, data bucketPolicyModel) (bucketPolicyModel, error) {
	input := &s3.GetBucketPolicyInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	output, err := r.Client.GetBucketPolicy(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if output == nil {
		return data, ErrEmptyResponse
	}
	if output.Policy == nil {
		return data, ErrResourceNotFound
	}

	policy := *output.Policy
	data.Policy = to.String(policy)
	data.Id = to.String(data.Bucket.ValueString())

	return data, nil
}
