package outscale

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var (
	_ resource.Resource                = &resourceSubnet{}
	_ resource.ResourceWithConfigure   = &resourceSubnet{}
	_ resource.ResourceWithImportState = &resourceSubnet{}
	_ resource.ResourceWithModifyPlan  = &resourceSubnet{}
)

type SubnetModel struct {
	AvailableIpsCount   types.Int32   `tfsdk:"available_ips_count"`
	IpRange             types.String  `tfsdk:"ip_range"`
	MapPublicIpOnLaunch types.Bool    `tfsdk:"map_public_ip_on_launch"`
	NetId               types.String  `tfsdk:"net_id"`
	State               types.String  `tfsdk:"state"`
	SubnetId            types.String  `tfsdk:"subnet_id"`
	SubregionName       types.String  `tfsdk:"subregion_name"`
	Tags                []ResourceTag `tfsdk:"tags"`

	RequestId types.String   `tfsdk:"request_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
	Id        types.String   `tfsdk:"id"`
}

type resourceSubnet struct {
	Client *oscgo.APIClient
}

func NewResourceSubnet() resource.Resource {
	return &resourceSubnet{}
}

func (r *resourceSubnet) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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
	data.SubnetId = types.StringValue(subnedId)
	data.Id = types.StringValue(subnedId)
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
			"tags": TagsSchema(),
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
					fwmodifyplan.ForceNewFramework(),
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

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateSubnetRequest{
		IpRange: data.IpRange.ValueString(),
		NetId:   data.NetId.ValueString(),
	}

	if !data.SubregionName.IsUnknown() && !data.SubregionName.IsNull() {
		createReq.SetSubregionName(data.SubregionName.ValueString())
	}

	var createResp oscgo.CreateSubnetResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.SubnetApi.CreateSubnet(ctx).CreateSubnetRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create subnet resource.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(*createResp.ResponseContext.RequestId)
	subnet := createResp.GetSubnet()

	if len(data.Tags) > 0 {
		err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(data.Tags), subnet.GetSubnetId())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add Tags on outscale_subnet resource.",
				err.Error(),
			)
			return
		}
	}
	readReq := oscgo.ReadSubnetsRequest{Filters: &oscgo.FiltersSubnet{SubnetIds: &[]string{subnet.GetSubnetId()}}}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: func() (any, string, error) {
			var resp oscgo.ReadSubnetsResponse

			err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
				rp, httpResp, err := r.Client.SubnetApi.ReadSubnets(ctx).ReadSubnetsRequest(readReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				resp = rp
				return nil
			})
			if err != nil {
				return resp, "failed", nil
			}
			subnet := resp.GetSubnets()[0]

			return resp, subnet.GetState(), nil
		},
		Timeout:    createTimeout,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Subnet to be ready.",
			err.Error(),
		)
		return
	}

	data.SubnetId = types.StringValue(subnet.GetSubnetId())
	data.Id = types.StringValue(subnet.GetSubnetId())
	if data.MapPublicIpOnLaunch.ValueBool() != subnet.GetMapPublicIpOnLaunch() {
		updateReq := oscgo.UpdateSubnetRequest{
			SubnetId:            data.SubnetId.ValueString(),
			MapPublicIpOnLaunch: data.MapPublicIpOnLaunch.ValueBool(),
		}
		err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.SubnetApi.UpdateSubnet(ctx).UpdateSubnetRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update MapPublicIpOnLaunch.",
				err.Error(),
			)
			return
		}
	}

	data, err = setSubnetState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set subnet state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSubnet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setSubnetState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSubnet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		tagsPlan, tagsState []ResourceTag
		resourceId          types.String
		err                 error
	)

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tags"), &tagsPlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tags"), &tagsState)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("subnet_id"), &resourceId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !reflect.DeepEqual(tagsPlan, tagsState) {
		toRemove, toCreate := diffOSCAPITags(tagsToOSCResourceTag(tagsPlan), tagsToOSCResourceTag(tagsState))
		err := updateFrameworkTags(ctx, r.Client, toCreate, toRemove, resourceId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Tags on subnet resource.",
				err.Error(),
			)
			return
		}
	}

	var stateData, planData SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := planData.Timeouts.Update(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := oscgo.UpdateSubnetRequest{
		SubnetId:            stateData.SubnetId.ValueString(),
		MapPublicIpOnLaunch: planData.MapPublicIpOnLaunch.ValueBool(),
	}
	err = retry.RetryContext(ctx, updateTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.SubnetApi.UpdateSubnet(ctx).UpdateSubnetRequest(updateReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update subnet resource.",
			err.Error(),
		)
		return
	}

	planData.SubnetId = resourceId
	data, err := setSubnetState(ctx, r, planData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set subnet state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSubnet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteSubnetRequest{
		SubnetId: data.SubnetId.ValueString(),
	}

	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.SubnetApi.DeleteSubnet(ctx).DeleteSubnetRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete subnet.",
			err.Error(),
		)
		return
	}
}

func setSubnetState(ctx context.Context, r *resourceSubnet, data SubnetModel) (SubnetModel, error) {
	subnetFilters := oscgo.FiltersSubnet{
		SubnetIds: &[]string{data.SubnetId.ValueString()},
	}
	readReq := oscgo.ReadSubnetsRequest{
		Filters: &subnetFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'subnet' read timeout value. Error: %v: ", diags.Errors())
	}

	var readResp oscgo.ReadSubnetsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.SubnetApi.ReadSubnets(ctx).ReadSubnetsRequest(readReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}
	if len(readResp.GetSubnets()) == 0 {
		return data, errors.New("Empty")
	}

	subnet := readResp.GetSubnets()[0]
	data.Tags = getTagsFromApiResponse(subnet.GetTags())

	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	data.AvailableIpsCount = types.Int32Value(subnet.GetAvailableIpsCount())
	data.IpRange = types.StringValue(subnet.GetIpRange())
	data.MapPublicIpOnLaunch = types.BoolValue(subnet.GetMapPublicIpOnLaunch())
	data.NetId = types.StringValue(subnet.GetNetId())
	data.State = types.StringValue(subnet.GetState())
	data.SubnetId = types.StringValue(subnet.GetSubnetId())
	data.SubregionName = types.StringValue(subnet.GetSubregionName())
	data.Id = types.StringValue(subnet.GetSubnetId())
	return data, nil
}
