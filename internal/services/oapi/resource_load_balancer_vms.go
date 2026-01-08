package oapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

var (
	_ resource.Resource               = &resourceLbuVms{}
	_ resource.ResourceWithConfigure  = &resourceLbuVms{}
	_ resource.ResourceWithModifyPlan = &resourceLbuVms{}
)

type LbuBackendVmsModel struct {
	LoadBalancerName types.String   `tfsdk:"load_balancer_name"`
	BackendVmIds     types.Set      `tfsdk:"backend_vm_ids"`
	BackendIps       types.Set      `tfsdk:"backend_ips"`
	RequestId        types.String   `tfsdk:"request_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id               types.String   `tfsdk:"id"`
}

type resourceLbuVms struct {
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *oscgo.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
}

func (r *resourceLbuVms) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data LbuBackendVmsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.BackendIps.IsNull() && data.BackendVmIds.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Attribute Configuration",
			"You need to specify at least the 'backend_ips' or the 'backend_vm_ids' parameter.",
		)
		return
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
		return
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
	var data LbuBackendVmsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listVmsIds, listVmsIps, diags := getSlicesLbuBackendVms(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	lbuName := data.LoadBalancerName.ValueString()
	createReq := oscgo.NewLinkLoadBalancerBackendMachinesRequest(lbuName)

	if len(listVmsIds) > 0 {
		createReq.SetBackendVmIds(listVmsIds)
	}
	if len(listVmsIps) > 0 {
		createReq.SetBackendIps(listVmsIps)
	}

	var createResp oscgo.LinkLoadBalancerBackendMachinesResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.LoadBalancerApi.LinkLoadBalancerBackendMachines(ctx).LinkLoadBalancerBackendMachinesRequest(*createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create LBU Backend_vms resource",
			err.Error(),
		)
		return

	}
	data.RequestId = types.StringValue(*createResp.ResponseContext.RequestId)
	data.Id = types.StringValue(lbuName)

	err = setLbuBackendState(ctx, r, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set LBU Backend_vms state",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLbuVms) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LbuBackendVmsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := setLbuBackendState(ctx, r, &data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set LBU Backend_vms API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLbuVms) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataPlan, dataState LbuBackendVmsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateTimeout, diags := dataPlan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	linkReq, unLinkReq, err := buildFwUpdateBackendsRequest(ctx, dataState.LoadBalancerName.ValueString(), &dataState, &dataPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update LBU backends.",
			err.Error(),
		)
		return
	}
	if unLinkReq.HasBackendVmIds() || unLinkReq.HasBackendIps() {
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.
				UnlinkLoadBalancerBackendMachines(ctx).UnlinkLoadBalancerBackendMachinesRequest(*unLinkReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to remove LBU backend vms.",
				err.Error(),
			)
			return
		}
	}

	if linkReq.HasBackendVmIds() || linkReq.HasBackendIps() {
		err := retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.
				LinkLoadBalancerBackendMachines(ctx).LinkLoadBalancerBackendMachinesRequest(*linkReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add LBU backend vms.",
				err.Error(),
			)
			return
		}
	}
	err = setLbuBackendState(ctx, r, &dataPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set LBU backend vms state after updating.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &dataPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLbuVms) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LbuBackendVmsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listVmsIds, listVmsIps, diags := getSlicesLbuBackendVms(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lbuName := data.LoadBalancerName.ValueString()
	unLinkReq := oscgo.NewUnlinkLoadBalancerBackendMachinesRequest(lbuName)
	if len(listVmsIds) > 0 {
		unLinkReq.SetBackendVmIds(listVmsIds)
	}
	if len(listVmsIps) > 0 {
		unLinkReq.SetBackendIps(listVmsIps)
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.LoadBalancerApi.
			UnlinkLoadBalancerBackendMachines(ctx).UnlinkLoadBalancerBackendMachinesRequest(*unLinkReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete LBU backend vms.",
			err.Error(),
		)
		return
	}
}

func setLbuBackendState(ctx context.Context, r *resourceLbuVms, data *LbuBackendVmsModel) error {
	lbuFilters := oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{data.LoadBalancerName.ValueString()},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return fmt.Errorf("unable to parse read timeout value: %v", diags.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	readReq := oscgo.ReadLoadBalancersRequest{
		Filters: &lbuFilters,
	}
	var readResp oscgo.ReadLoadBalancersResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.LoadBalancerApi.ReadLoadBalancers(context.Background()).ReadLoadBalancersRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if len(readResp.GetLoadBalancers()) == 0 {
		return ErrResourceEmpty
	}

	data.RequestId = types.StringValue(*readResp.ResponseContext.RequestId)
	lbu := readResp.GetLoadBalancers()[0]
	if !data.BackendVmIds.IsUnknown() && !data.BackendVmIds.IsNull() {
		data.BackendVmIds, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetBackendVmIds())
		if diags.HasError() {
			return fmt.Errorf("unable to set lbu backend_vm_ips: %v", diags.Errors())
		}
	}
	if !data.BackendIps.IsUnknown() && !data.BackendIps.IsNull() {
		data.BackendIps, diags = types.SetValueFrom(ctx, types.StringType, lbu.GetBackendIps())
		if diags.HasError() {
			return fmt.Errorf("unable to set lbu backend_ips: %v", diags.Errors())
		}
	}
	data.LoadBalancerName = types.StringValue(lbu.GetLoadBalancerName())
	data.Id = types.StringValue(lbu.GetLoadBalancerName())
	return nil
}

func getSlicesLbuBackendVms(ctx context.Context, data *LbuBackendVmsModel) ([]string, []string, diag.Diagnostics) {
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

func buildFwUpdateBackendsRequest(ctx context.Context, lbuName string, stateData, planData *LbuBackendVmsModel) (*oscgo.LinkLoadBalancerBackendMachinesRequest, *oscgo.UnlinkLoadBalancerBackendMachinesRequest, error) {
	linkReq := oscgo.NewLinkLoadBalancerBackendMachinesRequest(lbuName)
	unLinkReq := oscgo.NewUnlinkLoadBalancerBackendMachinesRequest(lbuName)
	var (
		ipsToAdd, ipsToRemove, vmIdsToAdd, vmIdsToRemove []string
		diags                                            diag.Diagnostics
	)

	if !reflect.DeepEqual(planData.BackendIps, stateData.BackendIps) {
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
		linkReq.SetBackendIps(ipsToAdd)
	}
	if len(vmIdsToAdd) > 0 {
		linkReq.SetBackendVmIds(vmIdsToAdd)
	}
	if len(ipsToRemove) > 0 {
		unLinkReq.SetBackendIps(ipsToRemove)
	}
	if len(vmIdsToRemove) > 0 {
		unLinkReq.SetBackendVmIds(vmIdsToRemove)
	}
	return linkReq, unLinkReq, nil
}
