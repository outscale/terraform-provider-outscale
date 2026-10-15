package oos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource              = &objectRetentionResource{}
	_ resource.ResourceWithConfigure = &objectRetentionResource{}
)

const (
	objectRetentionErrCreate = "Unable to create Object Retention"
	objectRetentionErrUpdate = "Unable to update Object Retention"
)

type objectRetentionModel struct {
	Bucket          types.String      `tfsdk:"bucket"`
	Key             types.String      `tfsdk:"key"`
	VersionId       types.String      `tfsdk:"version_id"`
	Mode            types.String      `tfsdk:"mode"`
	RetainUntilDate timetypes.RFC3339 `tfsdk:"retain_until_date"`
	Id              types.String      `tfsdk:"id"`
	Timeouts        timeouts.Value    `tfsdk:"timeouts"`
}

type objectRetentionResource struct {
	resourceCommon
}

func NewResourceObjectRetention() resource.Resource {
	return &objectRetentionResource{}
}

func (r *objectRetentionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *objectRetentionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_object_retention"
}

func (r *objectRetentionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
			"version_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mode": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(string(s3types.ObjectLockRetentionModeCompliance)),
				Validators: []validator.String{
					stringvalidator.OneOf(string(s3types.ObjectLockRetentionModeCompliance)),
				},
			},
			"retain_until_date": schema.StringAttribute{
				Required:   true,
				CustomType: timetypes.RFC3339Type{},
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

func (r *objectRetentionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data objectRetentionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	objKey := data.Key.ValueString()
	r.mu.LockKey(objKey)
	defer r.mu.UnlockKey(objKey) //nolint: errcheck

	err := r.put(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(objectRetentionErrCreate, err.Error())
		return
	}

	id := data.Bucket.ValueString() + "/" + data.Key.ValueString()
	if fwhelpers.IsSet(data.VersionId) {
		id += "/" + data.VersionId.ValueString()
	}
	data.Id = to.String(id)

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *objectRetentionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data objectRetentionModel
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

func (r *objectRetentionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data objectRetentionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	objKey := data.Key.ValueString()
	r.mu.LockKey(objKey)
	defer r.mu.UnlockKey(objKey) //nolint: errcheck

	err := r.put(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(objectRetentionErrUpdate, err.Error())
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *objectRetentionResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *objectRetentionResource) put(ctx context.Context, timeout time.Duration, data objectRetentionModel) error {
	date, err := from.RFC3339(data.RetainUntilDate)
	if err != nil {
		return err
	}

	input := s3.PutObjectRetentionInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
		Retention: &s3types.ObjectLockRetention{
			Mode:            s3types.ObjectLockRetentionMode(data.Mode.ValueString()),
			RetainUntilDate: date,
		},
	}

	_, err = r.Client.PutObjectRetention(ctx, &input, oos.WithRetryTimeout(timeout))

	return err
}

func (r *objectRetentionResource) read(ctx context.Context, timeout time.Duration, data objectRetentionModel) (objectRetentionModel, error) {
	input := s3.GetObjectRetentionInput{
		Bucket:    data.Bucket.ValueStringPointer(),
		Key:       data.Key.ValueStringPointer(),
		VersionId: data.VersionId.ValueStringPointer(),
	}

	output, err := r.Client.GetObjectRetention(ctx, &input, oos.WithRetryTimeout(timeout))
	switch {
	case err != nil:
		return data, err
	case output == nil:
		return data, ErrEmptyResponse
	case output.Retention == nil || output.Retention.RetainUntilDate == nil:
		return data, ErrResourceNotFound
	}

	data.Mode = to.String(output.Retention.Mode)
	data.RetainUntilDate = to.RFC3339(output.Retention.RetainUntilDate)

	return data, nil
}
