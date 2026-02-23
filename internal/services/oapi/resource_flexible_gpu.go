package oapi

import (
	"context"
	"errors"
	"fmt"

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
	_ resource.Resource              = &fgpuResource{}
	_ resource.ResourceWithConfigure = &fgpuResource{}
)

type GpuModel struct {
	DeleteOnVmDeletion types.Bool     `tfsdk:"delete_on_vm_deletion"`
	SubregionName      types.String   `tfsdk:"subregion_name"`
	FlexibleGpuId      types.String   `tfsdk:"flexible_gpu_id"`
	Generation         types.String   `tfsdk:"generation"`
	ModelName          types.String   `tfsdk:"model_name"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	VmId               types.String   `tfsdk:"vm_id"`
	State              types.String   `tfsdk:"state"`

	Id types.String `tfsdk:"id"`
}

type fgpuResource struct {
	Client *osc.Client
}

func NewResourcefGPU() resource.Resource {
	return &fgpuResource{}
}

func (r *fgpuResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *fgpuResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	flexible_gpu_id := req.ID
	if flexible_gpu_id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import net_resource identifier Got: %v", req.ID),
		)
		return
	}

	var data GpuModel
	var timeouts timeouts.Value
	data.FlexibleGpuId = to.String(flexible_gpu_id)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *fgpuResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flexible_gpu"
}

func (r *fgpuResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"subregion_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_on_vm_deletion": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"flexible_gpu_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vm_id": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
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

func (r *fgpuResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GpuModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := osc.CreateFlexibleGpuRequest{
		ModelName:     data.ModelName.ValueString(),
		SubregionName: data.SubregionName.ValueString(),
	}
	if !data.DeleteOnVmDeletion.IsNull() {
		createReq.DeleteOnVmDeletion = data.DeleteOnVmDeletion.ValueBoolPointer()
	}
	if !data.Generation.IsNull() {
		createReq.Generation = data.Generation.ValueStringPointer()
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createResp, err := r.Client.CreateFlexibleGpu(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource Net",
			"Error: "+err.Error(),
		)
		return
	}
	fGpu := ptr.From(createResp.FlexibleGpu)

	data.FlexibleGpuId = to.String(fGpu.FlexibleGpuId)
	data, err = read(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set net state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *fgpuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GpuModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err = read(ctx, r, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set net state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *fgpuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		planData   GpuModel
		resourceId types.String
		err        error
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("flexible_gpu_id"), &resourceId)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateFlexibleGpuRequest{
		FlexibleGpuId:      resourceId.ValueString(),
		DeleteOnVmDeletion: planData.DeleteOnVmDeletion.ValueBoolPointer(),
	}
	_, err = r.Client.UpdateFlexibleGpu(ctx, updateReq, options.WithRetryTimeout(updateTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Flexible GPU resource",
			"Error: "+err.Error(),
		)
		return
	}

	planData.FlexibleGpuId = resourceId
	data, err := read(ctx, r, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Flexible GPU state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *fgpuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GpuModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteFlexibleGpuRequest{
		FlexibleGpuId: data.FlexibleGpuId.ValueString(),
	}
	_, err := r.Client.DeleteFlexibleGpu(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Flexible GPU resource",
			"Error: "+err.Error(),
		)
	}
}

func read(ctx context.Context, r *fgpuResource, data GpuModel) (GpuModel, error) {
	netFilters := osc.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{data.FlexibleGpuId.ValueString()},
	}
	readReq := osc.ReadFlexibleGpusRequest{
		Filters: &netFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'flexible_gpu' read timeout value: %v", diags.Errors())
	}
	readResp, err := r.Client.ReadFlexibleGpus(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}

	if readResp.FlexibleGpus == nil || len(*readResp.FlexibleGpus) == 0 {
		return data, ErrResourceEmpty
	}
	fgpu := (*readResp.FlexibleGpus)[0]
	data.DeleteOnVmDeletion = to.Bool(fgpu.DeleteOnVmDeletion)
	data.FlexibleGpuId = to.String(fgpu.FlexibleGpuId)
	data.SubregionName = to.String(fgpu.SubregionName)
	data.Generation = to.String(fgpu.Generation)
	data.ModelName = to.String(fgpu.ModelName)
	data.Id = to.String(fgpu.FlexibleGpuId)
	data.State = to.String(*fgpu.State)
	data.VmId = to.String(ptr.From(fgpu.VmId))

	return data, nil
}
