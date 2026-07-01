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
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &virtualGatewayLinkResource{}
	_ resource.ResourceWithConfigure   = &virtualGatewayLinkResource{}
	_ resource.ResourceWithImportState = &virtualGatewayLinkResource{}
)

const (
	virtualGatewayLinkErrCreate = "Unable to link Virtual Gateway"
	virtualGatewayLinkErrDelete = "Unable to unlink Virtual Gateway"
)

type virtualGatewayLinkModel struct {
	NetId                    types.String   `tfsdk:"net_id"`
	VirtualGatewayId         types.String   `tfsdk:"virtual_gateway_id"`
	NetToVirtualGatewayLinks types.List     `tfsdk:"net_to_virtual_gateway_links"`
	RequestId                types.String   `tfsdk:"request_id"`
	Id                       types.String   `tfsdk:"id"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}

type virtualGatewayCommon struct {
	Client *osc.Client
}

type virtualGatewayLinkResource struct {
	virtualGatewayCommon
}

func NewResourceVirtualGatewayLink() resource.Resource {
	return &virtualGatewayLinkResource{}
}

func (r *virtualGatewayLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *virtualGatewayLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_gateway_link"
}

func (r *virtualGatewayLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import virtual gateway link identifier. Got: %v", req.ID),
		)
		return
	}

	var data virtualGatewayLinkModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.VirtualGatewayId = to.String(id)
	data.NetToVirtualGatewayLinks = types.ListNull(to.ObjType(virtualGatewayNetLinkAttrTypes))

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *virtualGatewayLinkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtual_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"net_to_virtual_gateway_links": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": schema.StringAttribute{
							Computed: true,
						},
						"net_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *virtualGatewayLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data virtualGatewayLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	linkReq := osc.LinkVirtualGatewayRequest{
		NetId:            data.NetId.ValueString(),
		VirtualGatewayId: data.VirtualGatewayId.ValueString(),
	}
	linkResp, err := r.Client.LinkVirtualGateway(ctx, linkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(virtualGatewayLinkErrCreate, err.Error())
		return
	}

	data.RequestId = to.String(linkResp.ResponseContext.RequestId)
	data.Id = to.String(data.VirtualGatewayId.ValueString())

	stateConf := &stateconf.StateChangeConf[osc.NetToVirtualGatewayLinkState]{
		Pending: stateconf.States(osc.NetToVirtualGatewayLinkStateDetached, osc.NetToVirtualGatewayLinkStateAttaching),
		Target:  stateconf.States(osc.NetToVirtualGatewayLinkStateAttached),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.NetId.ValueString(), data.VirtualGatewayId.ValueString()),
	}
	vgwAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	// We set the last read response to the state
	stateData, err := r.flatten(ctx, vgwAny.(osc.VirtualGateway), data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *virtualGatewayLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data virtualGatewayLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	data, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *virtualGatewayLinkResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *virtualGatewayLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data virtualGatewayLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.UnlinkVirtualGatewayRequest{
		VirtualGatewayId: data.VirtualGatewayId.ValueString(),
		NetId:            data.NetId.ValueString(),
	}
	_, err := r.Client.UnlinkVirtualGateway(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(virtualGatewayLinkErrDelete, err.Error())
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.NetToVirtualGatewayLinkState]{
		Pending: stateconf.States(osc.NetToVirtualGatewayLinkStateAttached, osc.NetToVirtualGatewayLinkStateDetaching),
		Target:  stateconf.States(osc.NetToVirtualGatewayLinkStateDetached),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.NetId.ValueString(), data.VirtualGatewayId.ValueString()),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	switch {
	case errors.Is(err, ErrResourceEmpty):
	case err != nil:
		resp.Diagnostics.AddError(virtualGatewayLinkErrDelete, err.Error())
	}
}

func (r *virtualGatewayLinkResource) read(ctx context.Context, timeout time.Duration, data virtualGatewayLinkModel) (virtualGatewayLinkModel, error) {
	readReq := osc.ReadVirtualGatewaysRequest{
		Filters: &osc.FiltersVirtualGateway{
			VirtualGatewayIds: &[]string{data.VirtualGatewayId.ValueString()},
		},
	}

	readResp, err := r.Client.ReadVirtualGateways(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if len(ptr.From(readResp.VirtualGateways)) == 0 {
		return data, ErrResourceEmpty
	}

	vgw := (*readResp.VirtualGateways)[0]
	if vgw.State == osc.VirtualGatewayStateDeleted {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	return r.flatten(ctx, vgw, data)
}

func (r *virtualGatewayLinkResource) flatten(ctx context.Context, virtualGateway osc.VirtualGateway, data virtualGatewayLinkModel) (virtualGatewayLinkModel, error) {
	vgwLink, err := r.findLink(virtualGateway)
	if err != nil {
		return data, err
	}

	links, diag := to.ListObject(ctx, r.flattenLinks(virtualGateway.NetToVirtualGatewayLinks), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag.Errors())
	}

	data.Id = to.String(virtualGateway.VirtualGatewayId)
	data.VirtualGatewayId = to.String(virtualGateway.VirtualGatewayId)
	data.NetId = to.String(ptr.From(vgwLink.NetId))
	data.NetToVirtualGatewayLinks = links

	return data, nil
}

func (r *virtualGatewayLinkResource) refreshFunc(netId, virtualGatewayID string) stateconf.StateRefreshFunc[osc.NetToVirtualGatewayLinkState] {
	return func(ctx context.Context) (any, osc.NetToVirtualGatewayLinkState, error) {
		readReq := osc.ReadVirtualGatewaysRequest{Filters: &osc.FiltersVirtualGateway{
			LinkNetIds:        &[]string{netId},
			VirtualGatewayIds: &[]string{virtualGatewayID},
		}}

		readResp, err := r.Client.ReadVirtualGateways(ctx, readReq)
		if err != nil {
			return nil, "", err
		}
		if len(ptr.From(readResp.VirtualGateways)) == 0 {
			return nil, "", ErrResourceEmpty
		}

		vgw := (*readResp.VirtualGateways)[0]
		vgwLink, err := r.findLink(vgw)
		if err != nil {
			return nil, "", err
		}

		return vgw, ptr.From(vgwLink.State), nil
	}
}

func (r *virtualGatewayLinkResource) findLink(virtualGateway osc.VirtualGateway) (*osc.NetToVirtualGatewayLink, error) {
	vgwLink, ok := lo.Find(virtualGateway.NetToVirtualGatewayLinks, func(link osc.NetToVirtualGatewayLink) bool {
		return ptr.From(link.State) == osc.NetToVirtualGatewayLinkStateAttached
	})
	if !ok {
		return nil, ErrResourceEmpty
	}

	return &vgwLink, nil
}
