package oos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/samber/lo"
	"k8s.io/utils/keymutex"
)

var (
	_ resource.Resource              = &bucketCorsResource{}
	_ resource.ResourceWithConfigure = &bucketCorsResource{}
)

const (
	bucketCorsErrCreate = "Unable to create Bucket CORS configuration"
	bucketCorsErrUpdate = "Unable to update Bucket CORS configuration"
	bucketCorsErrDelete = "Unable to delete Bucket CORS configuration"
)

type bucketCorsModel struct {
	Bucket   types.String   `tfsdk:"bucket"`
	Rules    types.Set      `tfsdk:"rule"`
	Id       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type bucketCorsRulesModel struct {
	AllowedHeaders types.Set   `tfsdk:"allowed_headers"`
	AllowedMethods types.Set   `tfsdk:"allowed_methods"`
	AllowedOrigins types.Set   `tfsdk:"allowed_origins"`
	ExposeHeaders  types.Set   `tfsdk:"expose_headers"`
	MaxAgeSeconds  types.Int32 `tfsdk:"max_age_seconds"`
}

func (bucketCorsRulesModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketCorsSchemaType, "rule")
}

type bucketCorsResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucketCors() resource.Resource {
	return &bucketCorsResource{}
}

func (r *bucketCorsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketCorsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket_cors"
}

func bucketCorsSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"rule": schema.SetNestedBlock{
				Validators: []validator.Set{
					setvalidator.IsRequired(),
					setvalidator.SizeBetween(1, 100),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_headers": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
						"allowed_methods": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(stringvalidator.OneOf("GET", "PUT", "HEAD", "POST", "DELETE")),
							},
						},
						"allowed_origins": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
						"expose_headers": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
						"max_age_seconds": schema.Int32Attribute{
							Required: true,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

var bucketCorsSchemaType = bucketCorsSchema(context.Background()).Type()

func (r *bucketCorsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketCorsSchema(ctx)
}

func (r *bucketCorsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketCorsModel
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
		resp.Diagnostics.AddError(bucketCorsErrCreate, err.Error())
		return
	}

	data.Id = to.String(bucket)

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketCorsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketCorsModel
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

func (r *bucketCorsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketCorsModel
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
		resp.Diagnostics.AddError(bucketCorsErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketCorsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketCorsModel
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

	input := &s3.DeleteBucketCorsInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	_, err := r.Client.DeleteBucketCors(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(bucketCorsErrDelete, err.Error())
	}
}

func (r *bucketCorsResource) put(ctx context.Context, timeout time.Duration, data bucketCorsModel) error {
	rules, err := r.expandRules(ctx, data.Rules)
	if err != nil {
		return err
	}

	input := &s3.PutBucketCorsInput{
		Bucket: data.Bucket.ValueStringPointer(),
		CORSConfiguration: &s3types.CORSConfiguration{
			CORSRules: rules,
		},
	}
	_, err = r.Client.PutBucketCors(ctx, input, oos.WithRetryTimeout(timeout))

	return err
}

func (r *bucketCorsResource) read(ctx context.Context, timeout time.Duration, data bucketCorsModel) (bucketCorsModel, error) {
	input := &s3.GetBucketCorsInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	output, err := r.Client.GetBucketCors(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if len(output.CORSRules) == 0 {
		return data, ErrResourceEmpty
	}

	rules, err := r.flattenRules(ctx, output.CORSRules)
	if err != nil {
		return data, err
	}

	data.Rules = rules
	data.Id = to.String(data.Bucket.ValueString())

	return data, nil
}

func (r *bucketCorsResource) expandRules(ctx context.Context, ruleSet types.Set) ([]s3types.CORSRule, error) {
	models, diags := to.Slice[bucketCorsRulesModel](ctx, ruleSet)
	if diags.HasError() {
		return nil, from.Diag(diags)
	}

	return lo.MapErr(models, func(model bucketCorsRulesModel, _ int) (s3types.CORSRule, error) {
		headers, diags := to.Slice[string](ctx, model.AllowedHeaders)
		methods, d := to.Slice[string](ctx, model.AllowedMethods)
		diags.Append(d...)
		origins, d := to.Slice[string](ctx, model.AllowedOrigins)
		diags.Append(d...)
		exposeHeaders, d := to.Slice[string](ctx, model.ExposeHeaders)
		diags.Append(d...)
		if diags.HasError() {
			return s3types.CORSRule{}, from.Diag(diags)
		}

		return s3types.CORSRule{
			AllowedHeaders: headers,
			AllowedMethods: methods,
			AllowedOrigins: origins,
			ExposeHeaders:  exposeHeaders,
			MaxAgeSeconds:  model.MaxAgeSeconds.ValueInt32Pointer(),
		}, nil
	})
}

func (r *bucketCorsResource) flattenRules(ctx context.Context, rules []s3types.CORSRule) (types.Set, error) {
	models, err := lo.MapErr(rules, func(rule s3types.CORSRule, _ int) (bucketCorsRulesModel, error) {
		headers, diags := to.Set(ctx, rule.AllowedHeaders)
		methods, d := to.Set(ctx, rule.AllowedMethods)
		diags.Append(d...)
		origins, d := to.Set(ctx, rule.AllowedOrigins)
		diags.Append(d...)
		exposeHeaders, d := to.Set(ctx, rule.ExposeHeaders)
		diags.Append(d...)
		if diags.HasError() {
			return bucketCorsRulesModel{}, from.Diag(diags)
		}

		return bucketCorsRulesModel{
			AllowedHeaders: headers,
			AllowedMethods: methods,
			AllowedOrigins: origins,
			ExposeHeaders:  exposeHeaders,
			MaxAgeSeconds:  to.Int32(rule.MaxAgeSeconds),
		}, nil
	})
	if err != nil {
		return types.Set{}, err
	}

	ruleSet, diags := to.SetObject(ctx, models)
	if diags.HasError() {
		return types.Set{}, from.Diag(diags)
	}

	return ruleSet, nil
}
