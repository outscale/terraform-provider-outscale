package oapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
)

var (
	_ resource.Resource                = &resourceSubnet{}
	_ resource.ResourceWithConfigure   = &resourceSubnet{}
	_ resource.ResourceWithImportState = &resourceSubnet{}
	_ resource.ResourceWithModifyPlan  = &resourceSubnet{}
)

type SubnetModel struct {
	AvailableIpsCount   types.Int32    `tfsdk:"available_ips_count"`
	IpRange             types.String   `tfsdk:"ip_range"`
	MapPublicIpOnLaunch types.Bool     `tfsdk:"map_public_ip_on_launch"`
	NetId               types.String   `tfsdk:"net_id"`
	State               types.String   `tfsdk:"state"`
	SubnetId            types.String   `tfsdk:"subnet_id"`
	SubregionName       types.String   `tfsdk:"subregion_name"`
	RequestId           types.String   `tfsdk:"request_id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	Id                  types.String   `tfsdk:"id"`
	TagsModel
}

type resourceSubnet struct {
	Client *osc.Client
}

func NewResourceSubnet() resource.Resource {
	return &resourceSubnet{}
}

func (r *resourceSubnet) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceSubnet) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	subnedId := req.ID

	if subnedId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import subnet identifier Got: %v", req.ID),
		)
		return
	}

	var data SubnetModel
	var timeouts timeouts.Value
	data.SubnetId = to.String(subnedId)
	data.Id = to.String(subnedId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceSubnet) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (r *resourceSubnet) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		// Return warning diagnostic to practitioners.
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceSubnet) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"net_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_range": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subregion_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
			},
			"available_ips_count": schema.Int32Attribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"subnet_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"map_public_ip_on_launch": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
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

func (r *resourceSubnet) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SubnetModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateSubnetRequest{
		IpRange: data.IpRange.ValueString(),
		NetId:   data.NetId.ValueString(),
	}

	if !data.SubregionName.IsUnknown() && !data.SubregionName.IsNull() {
		createReq.SubregionName = data.SubregionName.ValueStringPointer()
	}

	createResp, err := r.Client.CreateSubnet(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create subnet resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(*createResp.ResponseContext.RequestId)
	subnet := ptr.From(createResp.Subnet)

	diag := createOAPITagsFW(ctx, r.Client, createTimeout, data.Tags, subnet.SubnetId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	readReq := osc.ReadSubnetsRequest{Filters: &osc.FiltersSubnet{SubnetIds: &[]string{subnet.SubnetId}}}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Timeout: createTimeout,
		Refresh: func() (any, string, error) {
			resp, err := r.Client.ReadSubnets(ctx, readReq, options.WithRetryTimeout(createTimeout))
			if err != nil || resp.Subnets == nil {
				return resp, "failed", nil
			}
			subnet := (*resp.Subnets)[0]

			return resp, string(subnet.State), nil
		},
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Subnet to be ready.",
			err.Error(),
		)
		return
	}

	data.SubnetId = to.String(subnet.SubnetId)
	data.Id = to.String(subnet.SubnetId)
	if data.MapPublicIpOnLaunch.ValueBool() != subnet.MapPublicIpOnLaunch {
		updateReq := osc.UpdateSubnetRequest{
			SubnetId:            data.SubnetId.ValueString(),
			MapPublicIpOnLaunch: data.MapPublicIpOnLaunch.ValueBool(),
		}
		_, err := r.Client.UpdateSubnet(ctx, updateReq, options.WithRetryTimeout(createTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update MapPublicIpOnLaunch.",
				err.Error(),
			)
			return
		}
	}

	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set subnet state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSubnet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.read(ctx, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set subnet API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSubnet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData SubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.SubnetId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	updateTimeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateSubnetRequest{
		SubnetId:            stateData.SubnetId.ValueString(),
		MapPublicIpOnLaunch: planData.MapPublicIpOnLaunch.ValueBool(),
	}
	_, err := r.Client.UpdateSubnet(ctx, updateReq, options.WithRetryTimeout(updateTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update subnet resource",
			err.Error(),
		)
		return
	}

	data, err := r.read(ctx, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set subnet state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *resourceSubnet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteSubnetRequest{
		SubnetId: data.SubnetId.ValueString(),
	}

	// Retry on 409 (subnet in-use) as API can take time to see a subnet as not in use anymore
	err := oapihelpers.RetryOnCodes(ctx, []string{"9095"}, func() (resp any, err error) {
		return r.Client.DeleteSubnet(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	}, deleteTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete subnet.",
			err.Error(),
		)
		return
	}
}

func (r *resourceSubnet) read(ctx context.Context, data SubnetModel) (SubnetModel, error) {
	subnetFilters := osc.FiltersSubnet{
		SubnetIds: &[]string{data.SubnetId.ValueString()},
	}
	readReq := osc.ReadSubnetsRequest{
		Filters: &subnetFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'subnet' read timeout value: %v", diags.Errors())
	}

	readResp, err := r.Client.ReadSubnets(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	if readResp.Subnets == nil || len(*readResp.Subnets) == 0 {
		return data, ErrResourceEmpty
	}

	subnet := (*readResp.Subnets)[0]
	tags, diags := flattenOAPITagsFW(ctx, subnet.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	data.AvailableIpsCount = types.Int32Value(int32(subnet.AvailableIpsCount))
	data.IpRange = to.String(subnet.IpRange)
	data.MapPublicIpOnLaunch = to.Bool(subnet.MapPublicIpOnLaunch)
	data.NetId = to.String(subnet.NetId)
	data.State = to.String(subnet.State)
	data.SubnetId = to.String(subnet.SubnetId)
	data.SubregionName = to.String(subnet.SubregionName)
	data.Id = to.String(subnet.SubnetId)
	return data, nil
}
