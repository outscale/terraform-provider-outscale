package outscale

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
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
	ModeName           types.String   `tfsdk:"model_name"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	VmId               types.String   `tfsdk:"vm_id"`
	State              types.String   `tfsdk:"state"`

	Id types.String `tfsdk:"id"`
}

type fgpuResource struct {
	Client *oscgo.APIClient
}

func NewResourcefGPU() resource.Resource {
	return &fgpuResource{}
}

func (r *fgpuResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.FlexibleGpuId = types.StringValue(flexible_gpu_id)
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

	createReq := oscgo.NewCreateFlexibleGpuRequest(data.ModeName.ValueString(), data.SubregionName.ValueString())
	if !data.DeleteOnVmDeletion.IsNull() {
		createReq.SetDeleteOnVmDeletion(data.DeleteOnVmDeletion.ValueBool())
	}
	if !data.Generation.IsNull() {
		createReq.SetGeneration(data.Generation.ValueString())
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createResp oscgo.CreateFlexibleGpuResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.FlexibleGpuApi.CreateFlexibleGpu(ctx).CreateFlexibleGpuRequest(*createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource Net",
			"Error: "+utils.GetErrorResponse(err).Error(),
		)
		return
	}

	fGpu := createResp.GetFlexibleGpu()

	data.FlexibleGpuId = types.StringValue(fGpu.GetFlexibleGpuId())
	data, err = setFlexibleGpuState(ctx, r, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set net state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *fgpuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GpuModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err = setFlexibleGpuState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
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
	updateTimeout, diags := planData.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := oscgo.UpdateFlexibleGpuRequest{
		FlexibleGpuId:      resourceId.ValueString(),
		DeleteOnVmDeletion: planData.DeleteOnVmDeletion.ValueBoolPointer(),
	}
	err = retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.FlexibleGpuApi.UpdateFlexibleGpu(ctx).UpdateFlexibleGpuRequest(updateReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Flexible GPU resource",
			"Error: "+utils.GetErrorResponseToString(err),
		)
		return
	}

	planData.FlexibleGpuId = resourceId
	data, err := setFlexibleGpuState(ctx, r, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Flexible GPU state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *fgpuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GpuModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteFlexibleGpuRequest{
		FlexibleGpuId: data.FlexibleGpuId.ValueString(),
	}
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.FlexibleGpuApi.DeleteFlexibleGpu(ctx).DeleteFlexibleGpuRequest(delReq).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Flexible GPU resource",
			"Error: "+err.Error(),
		)
		return
	}
}

func setFlexibleGpuState(ctx context.Context, r *fgpuResource, data GpuModel) (GpuModel, error) {
	netFilters := oscgo.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{data.FlexibleGpuId.ValueString()},
	}
	readReq := oscgo.ReadFlexibleGpusRequest{
		Filters: &netFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'flexible_gpu' read timeout value. Error: %v: ", diags.Errors())
	}
	var readResp oscgo.ReadFlexibleGpusResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.FlexibleGpuApi.ReadFlexibleGpus(ctx).ReadFlexibleGpusRequest(readReq).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}

	if len(readResp.GetFlexibleGpus()) == 0 {
		return data, errors.New("Empty")
	}
	fgpu := readResp.GetFlexibleGpus()[0]
	data.DeleteOnVmDeletion = types.BoolValue(fgpu.GetDeleteOnVmDeletion())
	data.FlexibleGpuId = types.StringValue(fgpu.GetFlexibleGpuId())
	data.SubregionName = types.StringValue(fgpu.GetSubregionName())
	data.Generation = types.StringValue(fgpu.GetGeneration())
	data.ModeName = types.StringValue(fgpu.GetModelName())
	data.Id = types.StringValue(fgpu.GetFlexibleGpuId())
	data.State = types.StringValue(fgpu.GetState())
	data.VmId = types.StringValue(fgpu.GetVmId())
	return data, nil
}
