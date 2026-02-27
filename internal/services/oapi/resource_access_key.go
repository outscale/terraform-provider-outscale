package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
)

var (
	_ resource.Resource               = &resourceAccessKey{}
	_ resource.ResourceWithConfigure  = &resourceAccessKey{}
	_ resource.ResourceWithModifyPlan = &resourceAccessKey{}
)

type AccessKeyModel struct {
	AccessKeyId          types.String   `tfsdk:"access_key_id"`
	SecretKey            types.String   `tfsdk:"secret_key"`
	UserName             types.String   `tfsdk:"user_name"`
	State                types.String   `tfsdk:"state"`
	CreationDate         types.String   `tfsdk:"creation_date"`
	ExpirationDate       types.String   `tfsdk:"expiration_date"`
	LastModificationDate types.String   `tfsdk:"last_modification_date"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
	RequestId            types.String   `tfsdk:"request_id"`
	Id                   types.String   `tfsdk:"id"`
}

type resourceAccessKey struct {
	Client *osc.Client
}

func NewResourceAccessKey() resource.Resource {
	return &resourceAccessKey{}
}

func (r *resourceAccessKey) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
}

func (r *resourceAccessKey) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource. "+
				"And users will not be able to use this credentials.",
		)
	}
}

func (r *resourceAccessKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	accessKeyId := req.ID
	if accessKeyId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import access_key_resource identifier Got: %v", req.ID),
		)
		return
	}

	var data AccessKeyModel
	var timeouts timeouts.Value
	data.AccessKeyId = to.String(accessKeyId)
	data.Id = to.String(accessKeyId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceAccessKey) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_key"
}

func (r *resourceAccessKey) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"access_key_id": schema.StringAttribute{
				Computed: true,
			},
			"secret_key": schema.StringAttribute{
				Computed: true,
			},
			"user_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ACTIVE"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ACTIVE", "INACTIVE"}...),
				},
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
			},
			"expiration_date": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					validatorstring.DateValidator(),
				},
				PlanModifiers: []planmodifier.String{
					modifyplans.CheckExpirationDate(),
				},
			},
			"last_modification_date": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *resourceAccessKey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AccessKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateAccessKeyRequest{}
	if expireDate := data.ExpirationDate.ValueString(); expireDate != "" {
		time, err := to.ISO8601(expireDate)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse expiration date",
				err.Error(),
			)
			return
		}
		createReq.ExpirationDate = &time
	}
	if useName := data.UserName.ValueString(); useName != "" {
		createReq.UserName = new(useName)
	}

	createResp, err := r.Client.CreateAccessKey(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create access_key resource",
			err.Error(),
		)
		return
	}

	accessKey := ptr.From(createResp.AccessKey)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(accessKey.AccessKeyId)
	data.SecretKey = to.String(accessKey.SecretKey)

	if data.State.ValueString() != "ACTIVE" {
		if err := inactiveAccessKey(ctx, createTimeout, r, data); err != nil {
			resp.Diagnostics.AddError(
				"Unable to update access_key state",
				err.Error(),
			)
			return
		}
	}

	err = setAccessKeyState(ctx, r, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set access_key state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceAccessKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setAccessKeyState(ctx, r, &data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set access_key API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceAccessKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData AccessKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateAccessKeyRequest{
		AccessKeyId: stateData.AccessKeyId.ValueString(),
		State:       planData.State.ValueString(),
	}
	if expireDate := planData.ExpirationDate.ValueString(); expireDate != "" {
		time, err := to.ISO8601(expireDate)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse expiration date",
				err.Error(),
			)
			return
		}
		updateReq.ExpirationDate = &time
		stateData.ExpirationDate = planData.ExpirationDate
	}
	if useName := stateData.UserName.ValueString(); useName != "" {
		updateReq.UserName = new(useName)
	}
	updateResp, err := r.Client.UpdateAccessKey(ctx, updateReq, options.WithRetryTimeout(updateTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update access_key resource",
			err.Error(),
		)
		return
	}

	accKey := ptr.From(updateResp.AccessKey)
	if accKey.ExpirationDate == nil {
		stateData.ExpirationDate = types.StringNull()
	}
	err = setAccessKeyState(ctx, r, &stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update access_key API response values.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceAccessKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	to, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteAccessKeyRequest{
		AccessKeyId: data.AccessKeyId.ValueString(),
	}
	if userName := data.UserName.ValueString(); userName != "" {
		delReq.UserName = new(userName)
		if data.State.ValueString() != "INACTIVE" {
			err := inactiveAccessKey(ctx, to, r, data)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to INACTIVE access_key state",
					err.Error(),
				)
				return
			}
		}
	}

	_, err := r.Client.DeleteAccessKey(ctx, delReq, options.WithRetryTimeout(to))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete access_key",
			err.Error(),
		)
	}
}

func setAccessKeyState(ctx context.Context, r *resourceAccessKey, data *AccessKeyModel) error {
	accessKeyFilters := osc.FiltersAccessKeys{
		AccessKeyIds: &[]string{data.Id.ValueString()},
	}
	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'access_key' read timeout value: %v", diags.Errors())
	}

	readReq := osc.ReadAccessKeysRequest{
		Filters: &accessKeyFilters,
	}
	if !data.UserName.IsNull() {
		readReq.UserName = data.UserName.ValueStringPointer()
	}

	readResp, err := r.Client.ReadAccessKeys(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return err
	}
	if readResp.AccessKeys == nil || len(*readResp.AccessKeys) == 0 {
		return ErrResourceEmpty
	}

	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	acckey := (*readResp.AccessKeys)[0]

	datesEqual := false
	if stateExpirDate := data.ExpirationDate.ValueString(); stateExpirDate != "" {
		if from.ISO8601(acckey.ExpirationDate) != "" {
			stateDate, err := to.ISO8601(stateExpirDate)
			if err != nil {
				return err
			}
			datesEqual = stateDate.Equal(*acckey.ExpirationDate)
		}
	}
	if !datesEqual {
		data.ExpirationDate = to.String(from.ISO8601(acckey.ExpirationDate))
	}

	data.AccessKeyId = to.String(acckey.AccessKeyId)
	data.State = to.String(ptr.From(acckey.State))
	data.CreationDate = to.String(from.ISO8601(acckey.CreationDate))
	data.LastModificationDate = to.String(from.ISO8601(acckey.LastModificationDate))

	return nil
}

func inactiveAccessKey(ctx context.Context, timeout time.Duration, r *resourceAccessKey, data AccessKeyModel) error {
	req := osc.UpdateAccessKeyRequest{
		AccessKeyId: data.Id.ValueString(),
		State:       "INACTIVE",
	}
	if data.UserName.ValueString() != "" {
		req.UserName = data.UserName.ValueStringPointer()
	}
	if expireDate := data.ExpirationDate.ValueString(); expireDate != "" {
		time, err := to.ISO8601(expireDate)
		if err != nil {
			return err
		}
		req.ExpirationDate = &time
	}

	_, err := r.Client.UpdateAccessKey(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return err
	}
	return nil
}
