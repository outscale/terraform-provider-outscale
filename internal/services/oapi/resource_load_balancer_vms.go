package oapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource               = &resourceLbuVms{}
	_ resource.ResourceWithConfigure  = &resourceLbuVms{}
	_ resource.ResourceWithModifyPlan = &resourceLbuVms{}
)

const (
	lbuVmsErrCreate = "Unable to create Load Balancer backends"
	lbuVmsErrUpdate = "Unable to update Load Balancer backends"
	lbuVmsErrDelete = "Unable to delete Load Balancer backends"
	lbuVmsErrWait   = "Unable to wait for Load Balancer state"
	lbuVmsErrRemove = "Unable to remove Load Balancer backends"
	lbuVmsErrAdd    = "Unable to add Load Balancer backends"
)

type lbuBackendVmsModel struct {
	LoadBalancerName types.String   `tfsdk:"load_balancer_name"`
	BackendVmIds     types.Set      `tfsdk:"backend_vm_ids"`
	BackendIps       types.Set      `tfsdk:"backend_ips"`
	RequestId        types.String   `tfsdk:"request_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id               types.String   `tfsdk:"id"`
}

type resourceLbuVms struct {
	Client *osc.Client
}

func NewResourceLBUVms() resource.Resource {
	return &resourceLbuVms{}
}

func (r *resourceLbuVms) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceLbuVms) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data lbuBackendVmsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.BackendIps.IsNull() && data.BackendVmIds.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Attribute Configuration",
			"You need to specify at least the 'backend_ips' or the 'backend_vm_ids' parameter.",
		)
	}
}

func (r *resourceLbuVms) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will only unlink backend vms from load_balancer.",
		)
		return
	}

	if req.Plan.Raw.IsKnown() && req.State.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource 'outscale_load_balancer_vms' Considerations",
			"You have to apply twice or run 'terraform refesh' after the fist apply to get"+
				" the 'backend_ips' or 'backend_vm_ids' block values in load_balancer resource state",
		)
	}
}

func (r *resourceLbuVms) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_vms"
}

func (r *resourceLbuVms) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"load_balancer_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backend_vm_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"backend_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("backend_vm_ids"),
					}...),
				},
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

func (r *resourceLbuVms) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data lbuBackendVmsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	listVmsIds, listVmsIps, diags := getSlicesLbuBackendVms(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	lbuName := data.LoadBalancerName.ValueString()
	createReq := osc.LinkLoadBalancerBackendMachinesRequest{
		LoadBalancerName: lbuName,
	}

	if len(listVmsIds) > 0 {
		createReq.BackendVmIds = &listVmsIds
	}
	if len(listVmsIps) > 0 {
		createReq.BackendIps = &listVmsIps
	}

	createResp, err := r.Client.LinkLoadBalancerBackendMachines(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(lbuVmsErrCreate, err.Error())
		return
	}
	data.RequestId = to.String(*createResp.ResponseContext.RequestId)
	data.Id = to.String(lbuName)

	// LinkLoadBalancerBackendMachines response does not return the LBU object, a read is required to store the initial state
	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}
	diags = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	lbu, err := waitForLbuActive(ctx, r.Client, lbuName, timeout)
	if err != nil {
		resp.Diagnostics.AddError(lbuVmsErrWait, err.Error())
		return
	}

	// We set the last read response to the state
	stateData, err = r.flatten(ctx, data, *lbu)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceLbuVms) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data lbuBackendVmsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceLbuVms) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataPlan, dataState lbuBackendVmsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	timeout, diags := dataPlan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	linkReq, unLinkReq, err := buildUpdateBackendsRequest(ctx, dataState.LoadBalancerName.ValueString(), &dataState, &dataPlan)
	if err != nil {
		resp.Diagnostics.AddError(lbuVmsErrUpdate, err.Error())
		return
	}
	if (unLinkReq.BackendVmIds != nil && len(*unLinkReq.BackendVmIds) > 0) || (unLinkReq.BackendIps != nil && len(*unLinkReq.BackendIps) > 0) {
		respUpdate, err := r.Client.UnlinkLoadBalancerBackendMachines(ctx, unLinkReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(lbuVmsErrRemove, err.Error())
			return
		}
		dataPlan.RequestId = to.String(respUpdate.ResponseContext.RequestId)
	}

	if (linkReq.BackendVmIds != nil && len(*linkReq.BackendVmIds) > 0) || (linkReq.BackendIps != nil && len(*linkReq.BackendIps) > 0) {
		respUpdate, err := r.Client.LinkLoadBalancerBackendMachines(ctx, linkReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(lbuVmsErrAdd, err.Error())
			return
		}
		dataPlan.RequestId = to.String(respUpdate.ResponseContext.RequestId)
	}
	lbu, err := waitForLbuActive(ctx, r.Client, dataState.LoadBalancerName.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(lbuVmsErrWait, err.Error())
		return
	}

	// We set the last read response to the state
	stateData, err := r.flatten(ctx, dataPlan, *lbu)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceLbuVms) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data lbuBackendVmsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	listVmsIds, listVmsIps, diags := getSlicesLbuBackendVms(ctx, &data)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	lbuName := data.LoadBalancerName.ValueString()
	unLinkReq := osc.UnlinkLoadBalancerBackendMachinesRequest{
		LoadBalancerName: lbuName,
	}
	if len(listVmsIds) > 0 {
		unLinkReq.BackendVmIds = &listVmsIds
	}
	if len(listVmsIps) > 0 {
		unLinkReq.BackendIps = &listVmsIps
	}

	_, err := r.Client.UnlinkLoadBalancerBackendMachines(ctx, unLinkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(lbuVmsErrDelete, err.Error())
	}
}

func (r *resourceLbuVms) read(ctx context.Context, timeout time.Duration, data lbuBackendVmsModel) (lbuBackendVmsModel, error) {
	lbuFilters := osc.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{data.LoadBalancerName.ValueString()},
	}

	readReq := osc.ReadLoadBalancersRequest{
		Filters: &lbuFilters,
	}
	readResp, err := r.Client.ReadLoadBalancers(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if readResp.LoadBalancers == nil || len(*readResp.LoadBalancers) == 0 {
		return data, ErrResourceEmpty
	}

	data.RequestId = to.String(*readResp.ResponseContext.RequestId)
	lbu := (*readResp.LoadBalancers)[0]

	return r.flatten(ctx, data, lbu)
}

func (r *resourceLbuVms) flatten(ctx context.Context, data lbuBackendVmsModel, lbu osc.LoadBalancer) (lbuBackendVmsModel, error) {
	if fwhelpers.IsSet(data.BackendVmIds) {
		vmIds, diag := to.Set(ctx, lbu.BackendVmIds)
		if diag.HasError() {
			return data, from.Diag(diag)
		}
		data.BackendVmIds = vmIds
	}
	if !data.BackendIps.IsUnknown() && !data.BackendIps.IsNull() {
		ips, diag := to.Set(ctx, lbu.BackendIps)
		if diag.HasError() {
			return data, from.Diag(diag)
		}
		data.BackendIps = ips
	}
	data.LoadBalancerName = to.String(lbu.LoadBalancerName)
	data.Id = to.String(lbu.LoadBalancerName)

	return data, nil
}

func getSlicesLbuBackendVms(ctx context.Context, data *lbuBackendVmsModel) ([]string, []string, diag.Diagnostics) {
	listVmsIds := []string{}
	listVmsIps := []string{}
	diags := data.BackendVmIds.ElementsAs(ctx, &listVmsIds, false)
	if diags.HasError() {
		return listVmsIds, listVmsIds, diags
	}
	diags = data.BackendIps.ElementsAs(ctx, &listVmsIps, false)
	if diags.HasError() {
		return listVmsIds, listVmsIds, diags
	}
	return listVmsIds, listVmsIps, diags
}

func buildUpdateBackendsRequest(ctx context.Context, lbuName string, stateData, planData *lbuBackendVmsModel) (osc.LinkLoadBalancerBackendMachinesRequest, osc.UnlinkLoadBalancerBackendMachinesRequest, error) {
	linkReq := osc.LinkLoadBalancerBackendMachinesRequest{
		LoadBalancerName: lbuName,
	}
	unLinkReq := osc.UnlinkLoadBalancerBackendMachinesRequest{
		LoadBalancerName: lbuName,
	}
	var (
		ipsToAdd, ipsToRemove, vmIdsToAdd, vmIdsToRemove []string
		diags                                            diag.Diagnostics
	)

	if !planData.BackendIps.Equal(stateData.BackendIps) {
		ipsToAdd, ipsToRemove, diags = fwhelpers.GetSlicesFromTypesSetForUpdating(ctx, stateData.BackendIps, planData.BackendIps)
		if diags.HasError() {
			return linkReq, unLinkReq, fmt.Errorf("unable to get 'backend_ips' form typeset: %v", diags.Errors())
		}
	}
	if !reflect.DeepEqual(planData.BackendVmIds, stateData.BackendVmIds) {
		vmIdsToAdd, vmIdsToRemove, diags = fwhelpers.GetSlicesFromTypesSetForUpdating(ctx, stateData.BackendVmIds, planData.BackendVmIds)
		if diags.HasError() {
			return linkReq, unLinkReq, fmt.Errorf("unable to get 'backend_vm_ids' form typeset: %v", diags.Errors())
		}
	}
	if len(ipsToAdd) > 0 {
		linkReq.BackendIps = &ipsToAdd
	}
	if len(vmIdsToAdd) > 0 {
		linkReq.BackendVmIds = &vmIdsToAdd
	}
	if len(ipsToRemove) > 0 {
		unLinkReq.BackendIps = &ipsToRemove
	}
	if len(vmIdsToRemove) > 0 {
		unLinkReq.BackendVmIds = &vmIdsToRemove
	}
	return linkReq, unLinkReq, nil
}
