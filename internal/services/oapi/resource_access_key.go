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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/internal/fwvalidators"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.AccessKeyId = types.StringValue(accessKeyId)
	data.Id = types.StringValue(accessKeyId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
					fwvalidators.DateValidator(),
				},
				PlanModifiers: []planmodifier.String{
					fwmodifyplan.CheckExpirationDate(),
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateAccessKeyRequest{}
	if expirDate := data.ExpirationDate.ValueString(); expirDate != "" {
		createReq.SetExpirationDate(data.ExpirationDate.ValueString())
	}
	if useName := data.UserName.ValueString(); useName != "" {
		createReq.SetUserName(useName)
	}

	var createResp oscgo.CreateAccessKeyResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.AccessKeyApi.CreateAccessKey(ctx).CreateAccessKeyRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create access_key resource",
			err.Error(),
		)
		return
	}

	accessKey := createResp.GetAccessKey()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.Id = types.StringValue(accessKey.GetAccessKeyId())
	data.SecretKey = types.StringValue(accessKey.GetSecretKey())

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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceAccessKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := setAccessKeyState(ctx, r, &data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := oscgo.UpdateAccessKeyRequest{
		AccessKeyId: stateData.AccessKeyId.ValueString(),
		State:       planData.State.ValueString(),
	}
	if expireDate := planData.ExpirationDate.ValueString(); expireDate != "" {
		updateReq.SetExpirationDate(expireDate)
		stateData.ExpirationDate = planData.ExpirationDate
	}
	if useName := stateData.UserName.ValueString(); useName != "" {
		updateReq.SetUserName(useName)
	}
	var updateResp oscgo.UpdateAccessKeyResponse
	err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
		resp, httpResp, err := r.Client.AccessKeyApi.UpdateAccessKey(ctx).UpdateAccessKeyRequest(updateReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		updateResp = resp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update access_key resource",
			err.Error(),
		)
		return
	}

	accKey := updateResp.GetAccessKey()
	if !accKey.HasExpirationDate() {
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

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceAccessKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	to, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteAccessKeyRequest{
		AccessKeyId: data.AccessKeyId.ValueString(),
	}
	if userName := data.UserName.ValueString(); userName != "" {
		delReq.SetUserName(userName)
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

	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		_, httpResp, err := r.Client.AccessKeyApi.DeleteAccessKey(ctx).DeleteAccessKeyRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete access_key",
			err.Error(),
		)
		return
	}
}

func setAccessKeyState(ctx context.Context, r *resourceAccessKey, data *AccessKeyModel) error {
	accessKeyFilters := oscgo.FiltersAccessKeys{
		AccessKeyIds: &[]string{data.Id.ValueString()},
	}
	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse 'access_key' read timeout value. Error: %v: ", diags.Errors())
	}

	readReq := oscgo.ReadAccessKeysRequest{
		Filters: &accessKeyFilters,
	}
	if !data.UserName.IsNull() {
		readReq.SetUserName(data.UserName.ValueString())
	}

	var readResp oscgo.ReadAccessKeysResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.AccessKeyApi.ReadAccessKeys(ctx).ReadAccessKeysRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return utils.GetErrorResponse(err)
	}
	if len(readResp.GetAccessKeys()) == 0 {
		return errors.New("Empty")
	}

	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	acckey := readResp.GetAccessKeys()[0]

	datesEqual := false
	if stateExpirDate := data.ExpirationDate.ValueString(); stateExpirDate != "" {
		if stateExpirDate != "" && acckey.GetExpirationDate() != "" {
			stateDate, err := iso8601.Parse([]byte(stateExpirDate))
			if err != nil {
				return err
			}
			remoteDate, err := iso8601.Parse([]byte(acckey.GetExpirationDate()))
			if err != nil {
				return err
			}
			datesEqual = stateDate.Equal(remoteDate)
		}
	}
	if !datesEqual {
		data.ExpirationDate = types.StringValue(acckey.GetExpirationDate())
	}

	data.AccessKeyId = types.StringValue(acckey.GetAccessKeyId())
	data.State = types.StringValue(acckey.GetState())
	data.CreationDate = types.StringValue(acckey.GetCreationDate())
	data.LastModificationDate = types.StringValue(acckey.GetLastModificationDate())

	return nil
}

func inactiveAccessKey(ctx context.Context, to time.Duration, r *resourceAccessKey, data AccessKeyModel) error {
	req := oscgo.UpdateAccessKeyRequest{
		AccessKeyId: data.Id.ValueString(),
		State:       "INACTIVE",
	}
	if data.UserName.ValueString() != "" {
		req.SetUserName(data.UserName.ValueString())
	}
	if expireDate := data.ExpirationDate.ValueString(); expireDate != "" {
		req.SetExpirationDate(expireDate)
	}

	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		_, httpResp, err := r.Client.AccessKeyApi.UpdateAccessKey(ctx).UpdateAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return utils.GetErrorResponse(err)
	}
	return nil
}
