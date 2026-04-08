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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &natServiceResource{}
	_ resource.ResourceWithConfigure   = &natServiceResource{}
	_ resource.ResourceWithImportState = &natServiceResource{}
)

const (
	natSvcErrCreate       = "Unable to create NAT Service"
	natSvcErrRead         = "Unable to read NAT Service"
	natSvcErrDelete       = "Unable to delete NAT Service"
	natSvcErrState        = "Unable to set NAT Service state"
	natSvcErrNotAvailable = "Unable to wait for NAT Service to be available"
)

type natServiceModel struct {
	Id           types.String   `tfsdk:"id"`
	PublicIpId   types.String   `tfsdk:"public_ip_id"`
	SubnetId     types.String   `tfsdk:"subnet_id"`
	NatServiceId types.String   `tfsdk:"nat_service_id"`
	NetId        types.String   `tfsdk:"net_id"`
	PublicIps    types.List     `tfsdk:"public_ips"`
	State        types.String   `tfsdk:"state"`
	RequestId    types.String   `tfsdk:"request_id"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type natServicePublicIpModel struct {
	PublicIpId types.String `tfsdk:"public_ip_id"`
	PublicIp   types.String `tfsdk:"public_ip"`
}

var natSvcPublicIpAttrTypes = fwhelpers.GetAttrTypes(natServicePublicIpModel{})

type natServiceResource struct {
	Client *osc.Client
}

func NewResourceNatService() resource.Resource {
	return &natServiceResource{}
}

func (r *natServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *natServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nat_service"
}

func (r *natServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import NAT service identifier. Got: %v", req.ID),
		)
		return
	}

	var data natServiceModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.NatServiceId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal
	data.Tags = TagsNull()
	data.PublicIps = types.ListNull(types.ObjectType{AttrTypes: natSvcPublicIpAttrTypes})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *natServiceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nat_service_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"net_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_ips": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: natSvcPublicIpAttrTypes,
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *natServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data natServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateNatServiceRequest{
		PublicIpId: data.PublicIpId.ValueString(),
		SubnetId:   data.SubnetId.ValueString(),
	}

	// When creating a NAT service, the API may return a 9045 code if the subnet is not public yet
	respAny, err := oapihelpers.RetryOnCodes(ctx, []string{"9045"}, func() (any, error) {
		return r.Client.CreateNatService(ctx, createReq, options.WithRetryTimeout(timeout))
	}, timeout)
	if err != nil {
		resp.Diagnostics.AddError(natSvcErrCreate, err.Error())
		return
	}
	createResp := respAny.(*osc.CreateNatServiceResponse)

	natServiceId := createResp.NatService.NatServiceId
	data.Id = to.String(natServiceId)
	data.NatServiceId = to.String(natServiceId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{string(osc.NatServiceStatePending)},
		Target:  []string{string(osc.NatServiceStateAvailable)},
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(ctx, timeout, natServiceId),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(natSvcErrNotAvailable, err.Error())
		return
	}

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, natServiceId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(natSvcErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *natServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data natServiceModel
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
		resp.Diagnostics.AddError(natSvcErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *natServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData natServiceModel
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
		resp.Diagnostics.AddError(natSvcErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *natServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data natServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	_, err := r.Client.DeleteNatService(ctx, osc.DeleteNatServiceRequest{
		NatServiceId: data.Id.ValueString(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(natSvcErrDelete, err.Error())
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{string(osc.NatServiceStateDeleting)},
		Target:  []string{string(osc.NatServiceStateDeleted), string(osc.NatServiceStateAvailable)},
		Timeout: timeout,
		Refresh: r.stateRefreshFunc(ctx, timeout, data.Id.ValueString()),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(natSvcErrDelete, err.Error())
	}
}

func (r *natServiceResource) read(ctx context.Context, timeout time.Duration, data natServiceModel) (natServiceModel, error) {
	readReq := osc.ReadNatServicesRequest{
		Filters: &osc.FiltersNatService{NatServiceIds: &[]string{data.Id.ValueString()}},
	}

	resp, err := r.Client.ReadNatServices(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.NatServices == nil || len(*resp.NatServices) == 0 {
		return data, ErrResourceEmpty
	}

	natService := (*resp.NatServices)[0]
	if natService.State == osc.NatServiceStateDeleted {
		return data, ErrResourceEmpty
	}

	tags, diag := flattenOAPITagsFW(ctx, natService.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	publicIps, diag := to.ListObject(ctx, r.flattenPublicIpModel(natService.PublicIps))
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Tags = tags
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Id = to.String(natService.NatServiceId)
	data.NatServiceId = to.String(natService.NatServiceId)
	data.NetId = to.String(natService.NetId)
	data.SubnetId = to.String(natService.SubnetId)
	data.State = to.String(string(natService.State))
	data.PublicIps = publicIps

	// public_ip_id is a create-only input not returned directly by the API.
	// During import it is unknown, we recover it from public_ips[0].
	// In the case where the NAT has multiple linked public_ips, the resource will be recreated
	if data.PublicIpId.IsNull() || data.PublicIpId.IsUnknown() {
		if len(natService.PublicIps) > 0 {
			data.PublicIpId = to.String(natService.PublicIps[0].PublicIpId)
		}
	}

	return data, nil
}

func (r *natServiceResource) flattenPublicIpModel(ips []osc.PublicIpLight) []natServicePublicIpModel {
	return lo.Map(ips, func(pip osc.PublicIpLight, _ int) natServicePublicIpModel {
		return natServicePublicIpModel{
			PublicIpId: to.String(pip.PublicIpId),
			PublicIp:   to.String(pip.PublicIp),
		}
	})
}

func (r *natServiceResource) stateRefreshFunc(ctx context.Context, timeout time.Duration, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		req := osc.ReadNatServicesRequest{
			Filters: &osc.FiltersNatService{NatServiceIds: &[]string{id}},
		}
		resp, err := r.Client.ReadNatServices(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "", err
		}
		if resp.NatServices == nil || len(*resp.NatServices) == 0 {
			return nil, "", fmt.Errorf("nat service %s not found", id)
		}

		natService := (*resp.NatServices)[0]
		return resp, string(natService.State), nil
	}
}
