package oapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &publicIpLinkResource{}
	_ resource.ResourceWithConfigure   = &publicIpLinkResource{}
	_ resource.ResourceWithImportState = &publicIpLinkResource{}
)

const (
	publicIpLinkErrLink   = "Unable to link Public IP"
	publicIpLinkErrRead   = "Unable to read Public IP Link"
	publicIpLinkErrUnlink = "Unable to unlink Public IP"
	publicIpLinkErrState  = "Unable to set Public IP Link state"
)

type publicIpLinkModel struct {
	Id             types.String   `tfsdk:"id"`
	PublicIpId     types.String   `tfsdk:"public_ip_id"`
	PublicIp       types.String   `tfsdk:"public_ip"`
	AllowRelink    types.Bool     `tfsdk:"allow_relink"`
	VmId           types.String   `tfsdk:"vm_id"`
	NicId          types.String   `tfsdk:"nic_id"`
	PrivateIp      types.String   `tfsdk:"private_ip"`
	LinkPublicIpId types.String   `tfsdk:"link_public_ip_id"`
	NicAccountId   types.String   `tfsdk:"nic_account_id"`
	RequestId      types.String   `tfsdk:"request_id"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	TagsComputedModel
}

type publicIpLinkResource struct {
	Client *osc.Client
}

func NewResourcePublicIpLink() resource.Resource {
	return &publicIpLinkResource{}
}

func (r *publicIpLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *publicIpLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip_link"
}

func (r *publicIpLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	importId := req.ID
	if importId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import public IP link identifier. Got: %v", req.ID),
		)
		return
	}

	var data publicIpLinkModel
	var timeouts timeouts.Value
	data.Id = to.String(importId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = ComputedTagsNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicIpLinkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
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
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_ip": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_relink": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vm_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nic_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"link_public_ip_id": schema.StringAttribute{
				Computed: true,
			},
			"nic_account_id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (r *publicIpLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data publicIpLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	linkReq := osc.LinkPublicIpRequest{}
	if fwhelpers.IsSet(data.PublicIpId) {
		linkReq.PublicIpId = data.PublicIpId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.AllowRelink) {
		linkReq.AllowRelink = data.AllowRelink.ValueBoolPointer()
	}
	if fwhelpers.IsSet(data.VmId) {
		linkReq.VmId = data.VmId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.NicId) {
		linkReq.NicId = data.NicId.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.PrivateIp) {
		linkReq.PrivateIp = data.PrivateIp.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.PublicIp) {
		linkReq.PublicIp = data.PublicIp.ValueStringPointer()
	}

	linkResp, err := r.Client.LinkPublicIp(ctx, linkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(publicIpLinkErrLink, err.Error())
		return
	}

	if ptr.From(linkResp.LinkPublicIpId) != "" {
		data.Id = to.String(linkResp.LinkPublicIpId)
	} else {
		data.Id = data.PublicIp
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(publicIpLinkErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *publicIpLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data publicIpLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(publicIpLinkErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *publicIpLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *publicIpLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data publicIpLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	unlinkReq := osc.UnlinkPublicIpRequest{
		LinkPublicIpId: data.LinkPublicIpId.ValueStringPointer(),
	}

	_, err := r.Client.UnlinkPublicIp(ctx, unlinkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(publicIpLinkErrUnlink, err.Error())
	}
}

func (r *publicIpLinkResource) read(ctx context.Context, timeout time.Duration, data publicIpLinkModel) (publicIpLinkModel, error) {
	id := data.Id.ValueString()

	var readReq osc.ReadPublicIpsRequest
	if strings.Contains(id, "eipassoc") {
		readReq.Filters = &osc.FiltersPublicIp{
			LinkPublicIpIds: &[]string{id},
		}
	} else {
		readReq.Filters = &osc.FiltersPublicIp{
			PublicIps: &[]string{id},
		}
	}

	resp, err := r.Client.ReadPublicIps(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	if resp.PublicIps == nil || len(*resp.PublicIps) == 0 {
		return data, ErrResourceEmpty
	}

	publicIp := (*resp.PublicIps)[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, publicIp.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("%v", diag.Errors())
	}

	data.Tags = tags
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.PublicIpId = to.String(publicIp.PublicIpId)
	data.PublicIp = to.String(publicIp.PublicIp)
	data.LinkPublicIpId = to.String(ptr.From(publicIp.LinkPublicIpId))
	data.VmId = to.String(ptr.From(publicIp.VmId))
	data.NicId = to.String(ptr.From(publicIp.NicId))
	data.NicAccountId = to.String(ptr.From(publicIp.NicAccountId))
	data.PrivateIp = to.String(ptr.From(publicIp.PrivateIp))

	return data, nil
}
