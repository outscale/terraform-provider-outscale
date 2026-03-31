package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                = &publicIpResource{}
	_ resource.ResourceWithConfigure   = &publicIpResource{}
	_ resource.ResourceWithImportState = &publicIpResource{}
)

const (
	publicIpErrCreate = "Unable to create Public IP"
	publicIpErrRead   = "Unable to read Public IP"
	publicIpErrDelete = "Unable to delete Public IP"
	publicIpErrUnlink = "Unable to unlink Public IP from VM or NIC"
	publicIpErrState  = "Unable to set Public IP state"
)

type publicIpModel struct {
	Id             types.String   `tfsdk:"id"`
	PublicIpId     types.String   `tfsdk:"public_ip_id"`
	PublicIp       types.String   `tfsdk:"public_ip"`
	LinkPublicIpId types.String   `tfsdk:"link_public_ip_id"`
	VmId           types.String   `tfsdk:"vm_id"`
	NicId          types.String   `tfsdk:"nic_id"`
	NicAccountId   types.String   `tfsdk:"nic_account_id"`
	PrivateIp      types.String   `tfsdk:"private_ip"`
	RequestId      types.String   `tfsdk:"request_id"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type publicIpResource struct {
	Client *osc.Client
}

func NewResourcePublicIp() resource.Resource {
	return &publicIpResource{}
}

func (r *publicIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *publicIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (r *publicIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	publicIpId := req.ID
	if publicIpId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import public IP identifier. Got: %v", req.ID),
		)
		return
	}

	var data publicIpModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(publicIpId)
	data.PublicIpId = to.String(publicIpId)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal
	data.Tags = TagsNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicIpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags": TagsSchemaFW(),
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_ip_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_ip": schema.StringAttribute{
				Computed: true,
			},
			"link_public_ip_id": schema.StringAttribute{
				Computed: true,
			},
			"vm_id": schema.StringAttribute{
				Computed: true,
			},
			"nic_id": schema.StringAttribute{
				Computed: true,
			},
			"nic_account_id": schema.StringAttribute{
				Computed: true,
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *publicIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data publicIpModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createResp, err := r.Client.CreatePublicIp(ctx, osc.CreatePublicIpRequest{}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(publicIpErrCreate, err.Error())
		return
	}

	data.Id = to.String(createResp.PublicIp.PublicIpId)
	data.PublicIpId = to.String(createResp.PublicIp.PublicIpId)

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, createResp.PublicIp.PublicIpId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(publicIpErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *publicIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data publicIpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(publicIpErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *publicIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData publicIpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.PublicIpId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(publicIpErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *publicIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data publicIpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	// Get current state to check if the IP is linked
	currentData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			return
		}
		resp.Diagnostics.AddError(publicIpErrRead, err.Error())
		return
	}

	// Unlink if associated to a VM or NIC
	if currentData.LinkPublicIpId.ValueString() != "" || currentData.VmId.ValueString() != "" {
		unlinkReq := osc.UnlinkPublicIpRequest{}
		switch {
		case currentData.LinkPublicIpId.ValueString() != "":
			unlinkReq.LinkPublicIpId = currentData.LinkPublicIpId.ValueStringPointer()
		case currentData.VmId.ValueString() != "":
			unlinkReq.PublicIp = currentData.PublicIp.ValueStringPointer()
		}

		_, err := r.Client.UnlinkPublicIp(ctx, unlinkReq, options.WithRetryTimeout(timeout))
		if err != nil {
			// Do not fail if IP is already unlinked
			if osc.HasErrorCode(err, []string{"5080", "5026"}) {
				return
			}
			resp.Diagnostics.AddError(publicIpErrUnlink, err.Error())
			return
		}
	}

	delReq := osc.DeletePublicIpRequest{
		PublicIpId: data.Id.ValueStringPointer(),
	}
	_, err = r.Client.DeletePublicIp(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(publicIpErrDelete, err.Error())
	}
}

func (r *publicIpResource) read(ctx context.Context, timeout time.Duration, data publicIpModel) (publicIpModel, error) {
	readReq := osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{
			PublicIpIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadPublicIps(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	if resp.PublicIps == nil || len(*resp.PublicIps) == 0 {
		return data, ErrResourceEmpty
	}

	publicIp := (*resp.PublicIps)[0]

	tags, diag := flattenOAPITagsFW(ctx, publicIp.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diag.Errors())
	}

	data.Tags = tags
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Id = to.String(publicIp.PublicIpId)
	data.PublicIpId = to.String(publicIp.PublicIpId)
	data.PublicIp = to.String(publicIp.PublicIp)
	data.LinkPublicIpId = to.String(ptr.From(publicIp.LinkPublicIpId))
	data.VmId = to.String(ptr.From(publicIp.VmId))
	data.NicId = to.String(ptr.From(publicIp.NicId))
	data.NicAccountId = to.String(ptr.From(publicIp.NicAccountId))
	data.PrivateIp = to.String(ptr.From(publicIp.PrivateIp))

	return data, nil
}
