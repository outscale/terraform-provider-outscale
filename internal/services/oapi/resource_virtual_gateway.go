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
	_ resource.Resource                = &virtualGatewayResource{}
	_ resource.ResourceWithConfigure   = &virtualGatewayResource{}
	_ resource.ResourceWithImportState = &virtualGatewayResource{}
)

const (
	virtualGatewayErrCreate = "Unable to create Virtual Gateway"
	virtualGatewayErrDelete = "Unable to delete Virtual Gateway"
)

type virtualGatewayModel struct {
	ConnectionType           types.String   `tfsdk:"connection_type"`
	NetToVirtualGatewayLinks types.List     `tfsdk:"net_to_virtual_gateway_links"`
	State                    types.String   `tfsdk:"state"`
	VirtualGatewayId         types.String   `tfsdk:"virtual_gateway_id"`
	RequestId                types.String   `tfsdk:"request_id"`
	Id                       types.String   `tfsdk:"id"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type virtualGatewayNetLinkModel struct {
	State types.String `tfsdk:"state"`
	NetId types.String `tfsdk:"net_id"`
}

var virtualGatewayNetLinkAttrTypes = fwhelpers.GetAttrTypes(virtualGatewayNetLinkModel{})

type virtualGatewayResource struct {
	virtualGatewayCommon
}

func NewResourceVirtualGateway() resource.Resource {
	return &virtualGatewayResource{}
}

func (r *virtualGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *virtualGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_gateway"
}

func (r *virtualGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import virtual gateway identifier. Got: %v", req.ID),
		)
		return
	}

	var data virtualGatewayModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.VirtualGatewayId = to.String(id)
	data.Tags = TagsNull()
	data.NetToVirtualGatewayLinks = types.ListNull(to.ObjType(virtualGatewayNetLinkAttrTypes))

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *virtualGatewayResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"connection_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
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
			"virtual_gateway_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *virtualGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data virtualGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateVirtualGatewayRequest{
		ConnectionType: data.ConnectionType.ValueString(),
	}
	createResp, err := r.Client.CreateVirtualGateway(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(virtualGatewayErrCreate, err.Error())
		return
	}

	vgw := *createResp.VirtualGateway
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(vgw.VirtualGatewayId)

	stateData, err := r.flatten(ctx, data, vgw)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, vgw.VirtualGatewayId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.VirtualGatewayState]{
		Pending: stateconf.States(osc.VirtualGatewayStatePending),
		Target:  stateconf.States(osc.VirtualGatewayStateAvailable),
		Timeout: timeout,
		Refresh: r.refreshFunc(vgw.VirtualGatewayId),
	}
	vgwAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	// We set the last read response to the state
	stateData, err = r.flatten(ctx, data, vgwAny.(osc.VirtualGateway))
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *virtualGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data virtualGatewayModel
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

func (r *virtualGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData virtualGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.Id.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *virtualGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data virtualGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteVirtualGatewayRequest{
		VirtualGatewayId: data.Id.ValueString(),
	}
	_, err := r.Client.DeleteVirtualGateway(ctx, delReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(virtualGatewayErrDelete, err.Error())
		return
	}

	stateConf := &stateconf.StateChangeConf[osc.VirtualGatewayState]{
		Pending: stateconf.States(osc.VirtualGatewayStateDeleting),
		Target:  stateconf.States(osc.VirtualGatewayStateDeleted),
		Timeout: timeout,
		Refresh: r.refreshFunc(data.Id.ValueString()),
	}
	_, err = stateConf.WaitForStateContext(ctx)

	switch {
	case errors.Is(err, ErrResourceEmpty):
	case err != nil:
		resp.Diagnostics.AddError(virtualGatewayErrDelete, err.Error())
	}
}

func (r *virtualGatewayResource) read(ctx context.Context, timeout time.Duration, data virtualGatewayModel) (virtualGatewayModel, error) {
	readReq := osc.ReadVirtualGatewaysRequest{
		Filters: &osc.FiltersVirtualGateway{
			VirtualGatewayIds: &[]string{data.Id.ValueString()},
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

	return r.flatten(ctx, data, vgw)
}

func (r *virtualGatewayResource) flatten(ctx context.Context, data virtualGatewayModel, vgw osc.VirtualGateway) (virtualGatewayModel, error) {
	tags, diag := flattenOAPITagsFW(ctx, vgw.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	links, diag := to.ListFromAttrType(ctx, r.flattenLinks(vgw.NetToVirtualGatewayLinks), to.ObjType(virtualGatewayNetLinkAttrTypes), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Tags = tags
	data.Id = to.String(vgw.VirtualGatewayId)
	data.VirtualGatewayId = to.String(vgw.VirtualGatewayId)
	data.ConnectionType = to.String(vgw.ConnectionType)
	data.State = to.String(vgw.State)
	data.NetToVirtualGatewayLinks = links

	return data, nil
}

func (r *virtualGatewayResource) refreshFunc(id string) stateconf.StateRefreshFunc[osc.VirtualGatewayState] {
	return func(ctx context.Context) (any, osc.VirtualGatewayState, error) {
		readReq := osc.ReadVirtualGatewaysRequest{
			Filters: &osc.FiltersVirtualGateway{
				VirtualGatewayIds: &[]string{id},
			},
		}
		readResp, err := r.Client.ReadVirtualGateways(ctx, readReq)
		if err != nil {
			return nil, "", err
		}
		if len(ptr.From(readResp.VirtualGateways)) == 0 {
			return nil, "", ErrResourceEmpty
		}

		vgw := (*readResp.VirtualGateways)[0]
		return vgw, vgw.State, nil
	}
}

func (r *virtualGatewayCommon) flattenLinks(links []osc.NetToVirtualGatewayLink) []virtualGatewayNetLinkModel {
	return lo.Map(links, func(link osc.NetToVirtualGatewayLink, _ int) virtualGatewayNetLinkModel {
		return virtualGatewayNetLinkModel{
			State: to.String(ptr.From(link.State)),
			NetId: to.String(ptr.From(link.NetId)),
		}
	})
}
