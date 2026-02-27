package oapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
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
	Client *osc.Client
}

func NewResourceNetAttributes() resource.Resource {
	return &resourceNetAttributes{}
}

func (r *resourceNetAttributes) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	data.NetId = to.String(netId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = ComputedTagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.UpdateNetRequest{
		NetId:            data.NetId.ValueString(),
		DhcpOptionsSetId: data.DhcpOptionsSetId.ValueString(),
	}

	createResp, err := r.Client.UpdateNet(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Net Attributes.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	net := ptr.From(createResp.Net)
	data.NetId = to.String(net.NetId)
	data.Id = to.String(net.NetId)

	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Attributes state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetAttributes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetAttributesModel

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
			"Unable to set Net Attributes API response values.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetAttributes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NetAttributesModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.UpdateNetRequest{
		NetId:            data.NetId.ValueString(),
		DhcpOptionsSetId: data.DhcpOptionsSetId.ValueString(),
	}

	createResp, err := r.Client.UpdateNet(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Net Attributes.",
			err.Error(),
		)
		return
	}
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	net := ptr.From(createResp.Net)
	data.NetId = to.String(net.NetId)
	data.Id = to.String(net.NetId)

	data, err = r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Net Attributes state.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNetAttributes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *resourceNetAttributes) read(ctx context.Context, data NetAttributesModel) (NetAttributesModel, error) {
	netFilters := osc.FiltersNet{
		NetIds: &[]string{data.NetId.ValueString()},
	}
	readReq := osc.ReadNetsRequest{
		Filters: &netFilters,
	}

	readTimeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'net' read timeout value: %v", diags.Errors())
	}
	readResp, err := r.Client.ReadNets(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}

	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	if readResp.Nets == nil || len(*readResp.Nets) == 0 {
		return data, ErrResourceEmpty
	}

	net := (*readResp.Nets)[0]

	tags, diag := flattenOAPIComputedTagsFW(ctx, net.Tags)
	if diag.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags
	data.Id = to.String(net.NetId)
	data.NetId = to.String(net.NetId)
	data.DhcpOptionsSetId = to.String(net.DhcpOptionsSetId)
	data.IpRange = to.String(net.IpRange)
	data.State = to.String(net.State)
	data.Tenancy = to.String(net.Tenancy)

	return data, nil
}
