package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
)

var (
	_ resource.Resource               = &fgpuLinkResource{}
	_ resource.ResourceWithConfigure  = &fgpuLinkResource{}
	_ resource.ResourceWithModifyPlan = &fgpuLinkResource{}
)

const (
	fgpuLinktimeout = 5 * time.Minute

	fgpuLinkErrLink             = "Unable to link fGPU"
	fgpuLinkErrUnlink           = "Unable to unlink fGPU"
	fgpuLinkErrShutdownBehavior = "Unable to change VM shutdown behavior"
)

type fgpuLinkModel struct {
	FlexibleGpuIds types.Set      `tfsdk:"flexible_gpu_ids"`
	VmId           types.String   `tfsdk:"vm_id"`
	RequestId      types.String   `tfsdk:"request_id"`
	Id             types.String   `tfsdk:"id"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

type fgpuLinkResource struct {
	Client *osc.Client
}

func NewResourcefGPULink() resource.Resource {
	return &fgpuLinkResource{}
}

func (r *fgpuLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
}

func (r *fgpuLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	vmId := req.ID
	if vmId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import flexible_gpu_link identifier. Got: %v", req.ID),
		)
		return
	}

	var data fgpuLinkModel
	var timeouts timeouts.Value
	data.Id = to.String(id.UniqueId())
	data.VmId = to.String(vmId)
	data.FlexibleGpuIds = types.SetNull(types.StringType)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *fgpuLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flexible_gpu_link"
}

func (r *fgpuLinkResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// destroy: no-op
	if req.Plan.Raw.IsNull() {
		return
	}

	var planData fgpuLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.State.Raw.IsNull() {
		// When vm_id comes from a resource being created in the same apply, the restart warning is not useful
		if !fwhelpers.IsSet(planData.VmId) || planData.VmId.ValueString() == "" || !fwhelpers.IsSet(planData.Id) {
			return
		}

		resp.Diagnostics.AddWarning(
			"Flexible GPU link creation restarts the VM",
			fmt.Sprintf("Creating this resource links the Flexible GPUs to VM %q and the provider will stop then restart the VM during apply.", planData.VmId.ValueString()),
		)
		return
	}

	var stateData, configData fgpuLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duplicateMessage := "\n\n(This warning is expected to appear twice: during `terraform plan` and again after a successful `terraform apply`. The VMs should be running once apply completes.)"

	if !planData.VmId.Equal(stateData.VmId) {
		detail := "Replacing this resource unlinks the Flexible GPUs from the current VM and links them to the new VM. The provider will stop and then restart the affected VMs during apply."
		if fwhelpers.IsSet(stateData.VmId) && fwhelpers.IsSet(planData.VmId) {
			detail = fmt.Sprintf("Replacing this resource unlinks the Flexible GPUs from VM %q and links them to VM %q. The provider will stop and then restart the affected VMs during apply.", stateData.VmId.ValueString(), planData.VmId.ValueString())
		}
		detail += duplicateMessage

		resp.Diagnostics.AddWarning(
			"Flexible GPU link replacement restarts VMs",
			detail,
		)

		return
	}

	if !planData.FlexibleGpuIds.Equal(stateData.FlexibleGpuIds) {
		detail := "Updating this resource changes the Flexible GPUs linked to the VM, and the provider will stop and then restart the VM during apply."
		if fwhelpers.IsSet(planData.VmId) {
			detail = fmt.Sprintf("Updating this resource changes the Flexible GPUs linked to VM %q, and the provider will stop and then restart the VM during apply.", planData.VmId.ValueString())
		}
		detail += duplicateMessage

		resp.Diagnostics.AddWarning(
			"Flexible GPU link update restarts the VM",
			detail,
		)
	}
}

func (r *fgpuLinkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"flexible_gpu_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"vm_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
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

func (r *fgpuLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data fgpuLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Create(ctx, fgpuLinktimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	fgpuIds, diag := to.Slice[string](ctx, data.FlexibleGpuIds)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	for _, fgpuId := range fgpuIds {
		createReq := osc.LinkFlexibleGpuRequest{
			FlexibleGpuId: fgpuId,
			VmId:          data.VmId.ValueString(),
		}

		_, err := r.Client.LinkFlexibleGpu(ctx, createReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(
				fgpuLinkErrLink,
				err.Error(),
			)
			return
		}
		// The API response does not contain enough information to set the state directly, which would cause an error.
		// The next read will fill the state
	}

	data.Id = to.String(id.UniqueId())
	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			errSetTerraformState,
			err.Error(),
		)
		return
	}
	diag = resp.State.Set(ctx, &stateData)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	err = r.changeShutdownBehavior(ctx, data.VmId.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fgpuLinkErrShutdownBehavior,
			err.Error(),
		)
		return
	}
}

func (r *fgpuLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data fgpuLinkModel
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
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *fgpuLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData fgpuLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateIds, diag := to.Slice[string](ctx, stateData.FlexibleGpuIds)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	planIds, diag := to.Slice[string](ctx, planData.FlexibleGpuIds)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	toUnlink, toLink := lo.Difference(stateIds, planIds)

	for _, id := range toUnlink {
		req := osc.UnlinkFlexibleGpuRequest{
			FlexibleGpuId: id,
		}
		_, err := r.Client.UnlinkFlexibleGpu(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(
				fgpuLinkErrUnlink,
				err.Error(),
			)
			return
		}
	}
	for _, id := range toLink {
		req := osc.LinkFlexibleGpuRequest{
			FlexibleGpuId: id,
			VmId:          stateData.VmId.ValueString(),
		}
		_, err := r.Client.LinkFlexibleGpu(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(
				fgpuLinkErrLink,
				err.Error(),
			)
			return
		}
	}

	err := r.changeShutdownBehavior(ctx, planData.VmId.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fgpuLinkErrShutdownBehavior,
			err.Error(),
		)
		return
	}

	data, err := r.read(ctx, timeout, planData)
	if err != nil {
		resp.Diagnostics.AddError(errSetTerraformState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *fgpuLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data fgpuLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ids, diag := to.Slice[string](ctx, data.FlexibleGpuIds)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	for _, id := range ids {
		req := osc.UnlinkFlexibleGpuRequest{
			FlexibleGpuId: id,
		}
		_, err := r.Client.UnlinkFlexibleGpu(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(
				fgpuLinkErrUnlink,
				err.Error(),
			)
			return
		}
	}

	err := r.changeShutdownBehavior(ctx, data.VmId.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fgpuLinkErrShutdownBehavior,
			err.Error(),
		)
	}
}

func (r *fgpuLinkResource) read(ctx context.Context, timeout time.Duration, data fgpuLinkModel) (fgpuLinkModel, error) {
	readReq := osc.ReadFlexibleGpusRequest{
		Filters: &osc.FiltersFlexibleGpu{
			VmIds: &[]string{
				data.VmId.ValueString(),
			},
		},
	}

	readResp, err := r.Client.ReadFlexibleGpus(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if readResp.FlexibleGpus == nil || len(*readResp.FlexibleGpus) == 0 {
		return data, ErrResourceEmpty
	}

	gpuIds := lo.Map(*readResp.FlexibleGpus, func(gpu osc.FlexibleGpu, _ int) string {
		return gpu.FlexibleGpuId
	})
	idsSet, diag := to.Set(ctx, gpuIds)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.FlexibleGpuIds = idsSet
	data.RequestId = to.String(readResp.ResponseContext.RequestId)

	return data, nil
}

func (r *fgpuLinkResource) changeShutdownBehavior(ctx context.Context, vmId string, timeout time.Duration) error {
	resp, err := r.Client.ReadVms(ctx, osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmIds: &[]string{vmId},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return fmt.Errorf("error reading the vm %w", err)
	}

	if len(ptr.From(resp.Vms)) == 0 {
		return fmt.Errorf("error reading the vm %s err %w ", vmId, err)
	}
	vm := (*resp.Vms)[0]

	shutdownBehOpt := vm.VmInitiatedShutdownBehavior
	if shutdownBehOpt != "stop" {
		sbOpts := osc.UpdateVmRequest{VmId: vm.VmId}
		sbOpts.VmInitiatedShutdownBehavior = new("stop")
		if err := updateVmAttr(ctx, r.Client, timeout, sbOpts); err != nil {
			return err
		}
	}

	err = stopVM(ctx, r.Client, timeout, vmId)
	if err != nil {
		return err
	}

	if shutdownBehOpt != "stop" {
		sbReq := osc.UpdateVmRequest{VmId: vmId}
		sbReq.VmInitiatedShutdownBehavior = new(shutdownBehOpt)
		if err = updateVmAttr(ctx, r.Client, timeout, sbReq); err != nil {
			return err
		}
	}

	err = startVM(ctx, r.Client, timeout, vmId)
	if err != nil {
		return err
	}

	return nil
}
