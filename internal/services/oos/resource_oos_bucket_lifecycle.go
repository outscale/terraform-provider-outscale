package oos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorset"
	"github.com/samber/lo"
	"k8s.io/utils/keymutex"
)

var (
	_ resource.Resource              = &bucketLifecycleResource{}
	_ resource.ResourceWithConfigure = &bucketLifecycleResource{}
)

const (
	bucketLifecycleErrCreate = "Unable to create Bucket Lifecycle configuration"
	bucketLifecycleErrUpdate = "Unable to update Bucket Lifecycle configuration"
	bucketLifecycleErrDelete = "Unable to delete Bucket Lifecycle configuration"
)

type bucketLifecycleModel struct {
	Bucket   types.String   `tfsdk:"bucket"`
	Rules    types.Set      `tfsdk:"rule"`
	Id       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type bucketLifecycleRuleModel struct {
	Id                             types.String `tfsdk:"id"`
	Status                         types.String `tfsdk:"status"`
	Expiration                     types.Object `tfsdk:"expiration"`
	NoncurrentVersionExpiration    types.Object `tfsdk:"noncurrent_version_expiration"`
	AbortIncompleteMultipartUpload types.Object `tfsdk:"abort_incomplete_multipart_upload"`
	Filter                         types.Object `tfsdk:"filter"`
}

func (bucketLifecycleRuleModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule")
}

type bucketLifecycleExpirationModel struct {
	Date                      fwtypes.ISO8601Value `tfsdk:"date"`
	Days                      types.Int32          `tfsdk:"days"`
	ExpiredObjectDeleteMarker types.Bool           `tfsdk:"expired_object_delete_marker"`
}

func (bucketLifecycleExpirationModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule", "expiration")
}

type bucketLifecycleNoncurrentVersionExpirationModel struct {
	NoncurrentDays types.Int32 `tfsdk:"noncurrent_days"`
}

func (bucketLifecycleNoncurrentVersionExpirationModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule", "noncurrent_version_expiration")
}

type bucketLifecycleAbortIncompleteMultipartUploadModel struct {
	DaysAfterInitiation types.Int32 `tfsdk:"days_after_initiation"`
}

func (bucketLifecycleAbortIncompleteMultipartUploadModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule", "abort_incomplete_multipart_upload")
}

type bucketLifecycleFilterModel struct {
	Prefix types.String `tfsdk:"prefix"`
	And    types.Object `tfsdk:"and"`
}

func (bucketLifecycleFilterModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule", "filter")
}

type bucketLifecycleFilterAndModel struct {
	Prefix types.String `tfsdk:"prefix"`
	Tags   types.Map    `tfsdk:"tags"`
}

func (bucketLifecycleFilterAndModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(bucketLifecycleSchemaType, "rule", "filter", "and")
}

type bucketLifecycleResource struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

func NewResourceBucketLifecycle() resource.Resource {
	return &bucketLifecycleResource{}
}

func (r *bucketLifecycleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketLifecycleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_bucket_lifecycle"
}

func bucketLifecycleSchema(ctx context.Context) schema.Schema {
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
					validatorset.IsRequired(),
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(255),
							},
						},
						"status": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(lo.Map(s3types.ExpirationStatus("").Values(), func(v s3types.ExpirationStatus, _ int) string {
									return string(v)
								})...),
							},
						},
						"expiration": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"date": schema.StringAttribute{
									CustomType: fwtypes.ISO8601Type{},
									Optional:   true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"days": schema.Int32Attribute{
									Optional: true,
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int32Attribute{
									Required: true,
								},
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int32Attribute{
									Optional: true,
								},
							},
						},
						"filter": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									Optional: true,
								},
								"and": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"prefix": schema.StringAttribute{
											Required: true,
										},
										"tags": schema.MapAttribute{
											Required:    true,
											ElementType: types.StringType,
											Validators: []validator.Map{
												mapvalidator.SizeAtLeast(1),
											},
										},
									},
								},
							},
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

var bucketLifecycleSchemaType = bucketLifecycleSchema(context.Background()).Type()

func (r *bucketLifecycleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketLifecycleSchema(ctx)
}

func (r *bucketLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bucketLifecycleModel
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
		resp.Diagnostics.AddError(bucketLifecycleErrCreate, err.Error())
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

func (r *bucketLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bucketLifecycleModel
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

func (r *bucketLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketLifecycleModel
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
		resp.Diagnostics.AddError(bucketLifecycleErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *bucketLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bucketLifecycleModel
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

	input := &s3.DeleteBucketLifecycleInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	_, err := r.Client.DeleteBucketLifecycle(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil && !IsNotFound(err) {
		resp.Diagnostics.AddError(bucketLifecycleErrDelete, err.Error())
	}
}

func (r *bucketLifecycleResource) put(ctx context.Context, timeout time.Duration, data bucketLifecycleModel) error {
	rules, err := r.expandRules(ctx, data.Rules)
	if err != nil {
		return err
	}

	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: data.Bucket.ValueStringPointer(),
		LifecycleConfiguration: &s3types.BucketLifecycleConfiguration{
			Rules: rules,
		},
	}
	_, err = r.Client.PutBucketLifecycleConfiguration(ctx, input, oos.WithRetryTimeout(timeout))

	return err
}

func (r *bucketLifecycleResource) read(ctx context.Context, timeout time.Duration, data bucketLifecycleModel) (bucketLifecycleModel, error) {
	input := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: data.Bucket.ValueStringPointer(),
	}
	output, err := r.Client.GetBucketLifecycleConfiguration(ctx, input, oos.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if len(output.Rules) == 0 {
		return data, ErrResourceEmpty
	}

	rules, err := r.flattenRules(ctx, output.Rules)
	if err != nil {
		return data, err
	}

	data.Rules = rules
	data.Id = to.String(data.Bucket.ValueString())

	return data, nil
}

func (r *bucketLifecycleResource) expandRules(ctx context.Context, ruleSet types.Set) ([]s3types.LifecycleRule, error) {
	models, diags := to.Slice[bucketLifecycleRuleModel](ctx, ruleSet)
	if diags.HasError() {
		return nil, from.Diag(diags)
	}

	return lo.MapErr(models, func(model bucketLifecycleRuleModel, _ int) (s3types.LifecycleRule, error) {
		rule := s3types.LifecycleRule{
			Status: s3types.ExpirationStatus(model.Status.ValueString()),
		}
		if fwhelpers.IsSet(model.Id) {
			rule.ID = model.Id.ValueStringPointer()
		}

		if fwhelpers.IsSet(model.Expiration) {
			expiration, d := to.Model[bucketLifecycleExpirationModel](ctx, model.Expiration)
			diags.Append(d...)

			rule.Expiration = &s3types.LifecycleExpiration{}
			if fwhelpers.IsSet(expiration.Date) {
				date, d := expiration.Date.ValueISO8601Time()
				diags.Append(d...)

				rule.Expiration.Date = &date.Time
			}
			if fwhelpers.IsSet(expiration.Days) {
				rule.Expiration.Days = expiration.Days.ValueInt32Pointer()
			}
			if fwhelpers.IsSet(expiration.ExpiredObjectDeleteMarker) {
				rule.Expiration.ExpiredObjectDeleteMarker = expiration.ExpiredObjectDeleteMarker.ValueBoolPointer()
			}
		}

		if fwhelpers.IsSet(model.NoncurrentVersionExpiration) {
			expiration, d := to.Model[bucketLifecycleNoncurrentVersionExpirationModel](ctx, model.NoncurrentVersionExpiration)
			diags.Append(d...)

			rule.NoncurrentVersionExpiration = &s3types.NoncurrentVersionExpiration{
				NoncurrentDays: expiration.NoncurrentDays.ValueInt32Pointer(),
			}
		}

		if fwhelpers.IsSet(model.AbortIncompleteMultipartUpload) {
			upload, d := to.Model[bucketLifecycleAbortIncompleteMultipartUploadModel](ctx, model.AbortIncompleteMultipartUpload)
			diags.Append(d...)

			rule.AbortIncompleteMultipartUpload = &s3types.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: upload.DaysAfterInitiation.ValueInt32Pointer(),
			}
		}

		if fwhelpers.IsSet(model.Filter) {
			filter, d := to.Model[bucketLifecycleFilterModel](ctx, model.Filter)
			diags.Append(d...)

			rule.Filter = &s3types.LifecycleRuleFilter{}
			if fwhelpers.IsSet(filter.Prefix) {
				rule.Filter.Prefix = filter.Prefix.ValueStringPointer()
			}
			if fwhelpers.IsSet(filter.And) {
				and, d := to.Model[bucketLifecycleFilterAndModel](ctx, filter.And)
				diags.Append(d...)

				tagsMap, d := from.Map[string](ctx, and.Tags)
				diags.Append(d...)
				tags := lo.MapToSlice(tagsMap, func(key string, value string) s3types.Tag {
					return s3types.Tag{
						Key:   &key,
						Value: &value,
					}
				})

				rule.Filter.And = &s3types.LifecycleRuleAndOperator{
					Prefix: and.Prefix.ValueStringPointer(),
					Tags:   tags,
				}
			}
		}
		if diags.HasError() {
			return rule, from.Diag(diags)
		}

		return rule, nil
	})
}

func (r *bucketLifecycleResource) flattenRules(ctx context.Context, rules []s3types.LifecycleRule) (types.Set, error) {
	models, err := lo.MapErr(rules, func(rule s3types.LifecycleRule, _ int) (bucketLifecycleRuleModel, error) {
		model := bucketLifecycleRuleModel{
			Id:                             to.String(rule.ID),
			Status:                         to.String(rule.Status),
			Expiration:                     to.ObjectNull[bucketLifecycleExpirationModel](),
			NoncurrentVersionExpiration:    to.ObjectNull[bucketLifecycleNoncurrentVersionExpirationModel](),
			AbortIncompleteMultipartUpload: to.ObjectNull[bucketLifecycleAbortIncompleteMultipartUploadModel](),
			Filter:                         to.ObjectNull[bucketLifecycleFilterModel](),
		}

		if rule.Expiration != nil {
			expirationModel := bucketLifecycleExpirationModel{
				Date:                      fwtypes.NewISO8601Null(),
				Days:                      to.Int32(rule.Expiration.Days),
				ExpiredObjectDeleteMarker: to.Bool(rule.Expiration.ExpiredObjectDeleteMarker),
			}

			if rule.Expiration.Date != nil {
				time, err := to.ISO8601(rule.Expiration.Date)
				if err != nil {
					return model, err
				}
				expirationModel.Date = to.ISO8601Value(time)
			}

			expiration, d := to.Object(ctx, expirationModel)
			if d.HasError() {
				return model, from.Diag(d)
			}

			model.Expiration = expiration
		}

		if rule.NoncurrentVersionExpiration != nil {
			expiration, d := to.Object(ctx, bucketLifecycleNoncurrentVersionExpirationModel{
				NoncurrentDays: to.Int32(rule.NoncurrentVersionExpiration.NoncurrentDays),
			})
			if d.HasError() {
				return model, from.Diag(d)
			}

			model.NoncurrentVersionExpiration = expiration
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			upload, d := to.Object(ctx, bucketLifecycleAbortIncompleteMultipartUploadModel{
				DaysAfterInitiation: to.Int32(rule.AbortIncompleteMultipartUpload.DaysAfterInitiation),
			})
			if d.HasError() {
				return model, from.Diag(d)
			}

			model.AbortIncompleteMultipartUpload = upload
		}

		filterModel := bucketLifecycleFilterModel{
			And: to.ObjectNull[bucketLifecycleFilterAndModel](),
		}
		if rule.Filter != nil {
			filterModel.Prefix = to.String(rule.Filter.Prefix)
			if rule.Filter.And != nil {
				tags := lo.FilterSliceToMap(rule.Filter.And.Tags, func(t s3types.Tag) (string, string, bool) {
					return ptr.From(t.Key), ptr.From(t.Value), t.Key != nil
				})

				tagsMap, d := to.Map(ctx, tags)
				if d.HasError() {
					return model, from.Diag(d)
				}

				and, d := to.Object(ctx, bucketLifecycleFilterAndModel{
					Prefix: to.String(rule.Filter.And.Prefix),
					Tags:   tagsMap,
				})
				if d.HasError() {
					return model, from.Diag(d)
				}

				filterModel.And = and
			}
		}

		filter, d := to.Object(ctx, filterModel)
		if d.HasError() {
			return model, from.Diag(d)
		}
		model.Filter = filter

		return model, nil
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
