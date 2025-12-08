package outscale

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var (
	_ resource.Resource                = &resourceNetAttributes{}
	_ resource.ResourceWithConfigure   = &resourceNetAttributes{}
	_ resource.ResourceWithImportState = &resourceNetAttributes{}
	_ resource.ResourceWithModifyPlan  = &resourceNetAttributes{}
)

type NetAttributesModel struct {
	DhcpOptionsSetId types.String   `tfsdk:"dhcp_options_set_id"`
	IpRange          types.String   `tfsdk:"ip_range"`
	NetId            types.String   `tfsdk:"net_id"`
	State            types.String   `tfsdk:"state"`
	Tenancy          types.String   `tfsdk:"tenancy"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	RequestId        types.String   `tfsdk:"request_id"`
	Id               types.String   `tfsdk:"id"`
	TagsComputedModel
}

type resourceNetAttributes struct {
	Client *oscgo.APIClient
}

func NewResourceNetAttributes() resource.Resource {
	return &resourceNetAttributes{}
}

func (r *resourceNetAttributes) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceNetAttributes) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	netId := req.ID

	if netId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Net Peering identifier, got: %v", req.ID),
		)
		return
	}

	var data NetAttributesModel
	var timeouts timeouts.Value
	data.NetId = types.StringValue(netId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = ComputedTagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAttributes) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net_attributes"
}

func (r *resourceNetAttributes) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will remove resource from terraform state but not fully delete. It will be gone when deleting the Net.",
		)
	}
}

func (r *resourceNetAttributes) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"dhcp_options_set_id": schema.StringAttribute{
				Required: true,
			},
			"net_id": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"ip_range": schema.StringAttribute{
				Computed: true,
			},
			"tenancy": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"tags": TagsSchemaComputedFW(),
		},
	}
}

func (r *resourceNetAttributes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetAttributesModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.UpdateNetRequest{
		NetId:            data.NetId.ValueString(),
		DhcpOptionsSetId: data.DhcpOptionsSetId.ValueString(),
	}

	var createResp oscgo.UpdateNetResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetApi.UpdateNet(ctx).UpdateNetRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Net Attributes.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	net := createResp.GetNet()
	data.NetId = types.StringValue(net.GetNetId())
	data.Id = types.StringValue(net.GetNetId())

	data, err = r.setNetAttributesState(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Attributes state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAttributes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetAttributesModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.setNetAttributesState(ctx, data)
	if err != nil {
		if err.Error() == "Empty" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Net Attributes API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAttributes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NetAttributesModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.UpdateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.UpdateNetRequest{
		NetId:            data.NetId.ValueString(),
		DhcpOptionsSetId: data.DhcpOptionsSetId.ValueString(),
	}

	var createResp oscgo.UpdateNetResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetApi.UpdateNet(ctx).UpdateNetRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Net Attributes.",
			err.Error(),
		)
		return
	}
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	net := createResp.GetNet()
	data.NetId = types.StringValue(net.GetNetId())
	data.Id = types.StringValue(net.GetNetId())

	data, err = r.setNetAttributesState(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Attributes state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceNetAttributes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *resourceNetAttributes) setNetAttributesState(ctx context.Context, data NetAttributesModel) (NetAttributesModel, error) {
	netFilters := oscgo.FiltersNet{
		NetIds: &[]string{data.NetId.ValueString()},
	}
	readReq := oscgo.ReadNetsRequest{
		Filters: &netFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net' read timeout value. Error: %v: ", diags.Errors())
	}
	var readResp oscgo.ReadNetsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetApi.ReadNets(ctx).ReadNetsRequest(readReq).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		readResp = rp
		return nil
	})
	if err != nil {
		return data, err
	}

	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	if len(readResp.GetNets()) == 0 {
		return data, errors.New("Empty")
	}

	net := readResp.GetNets()[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, net.GetTags())
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags. Error: %v: ", diags.Errors())
	}
	data.Tags = tags
	data.Id = types.StringValue(net.GetNetId())
	data.NetId = types.StringValue(net.GetNetId())
	data.DhcpOptionsSetId = types.StringValue(net.GetDhcpOptionsSetId())
	data.IpRange = types.StringValue(net.GetIpRange())
	data.State = types.StringValue(net.GetState())
	data.Tenancy = types.StringValue(net.GetTenancy())

	return data, nil
}
