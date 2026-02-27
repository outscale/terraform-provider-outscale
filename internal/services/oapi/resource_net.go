package oapi

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
)

var (
	_ resource.Resource              = &resourceNet{}
	_ resource.ResourceWithConfigure = &resourceNet{}
)

type NetModel struct {
	DhcpOptionsSetId types.String   `tfsdk:"dhcp_options_set_id"`
	IpRange          types.String   `tfsdk:"ip_range"`
	NetId            types.String   `tfsdk:"net_id"`
	State            types.String   `tfsdk:"state"`
	Tenancy          types.String   `tfsdk:"tenancy"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	RequestId        types.String   `tfsdk:"request_id"`
	Id               types.String   `tfsdk:"id"`
	TagsModel
}

type resourceNet struct {
	Client *osc.Client
}

func NewResourceNet() resource.Resource {
	return &resourceNet{}
}

func (r *resourceNet) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceNet) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	net_id := req.ID
	if net_id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import net_resource identifier Got: %v", req.ID),
		)
		return
	}

	var data NetModel
	var timeouts timeouts.Value
	data.NetId = to.String(net_id)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.Tags = TagsNull()

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceNet) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net"
}

func (r *resourceNet) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"dhcp_options_set_id": schema.StringAttribute{
				Computed: true,
			},
			"ip_range": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"net_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"tenancy": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
				Validators: []validator.String{stringvalidator.NoneOf("")},
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

func (r *resourceNet) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, valideIpRange, err := net.ParseCIDR(data.IpRange.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse ip_range value: "+data.IpRange.ValueString(),
			"Error: "+err.Error(),
		)
		return
	}
	if data.IpRange.ValueString() != valideIpRange.String() {
		resp.Diagnostics.AddError(
			"Invalide net ip_range value: "+data.IpRange.ValueString(),
			"Error: ip_range value should be: "+valideIpRange.String(),
		)
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateNetRequest{
		IpRange: data.IpRange.ValueString(),
	}

	if fwhelpers.IsSet(data.Tenancy) {
		createReq.Tenancy = data.Tenancy.ValueStringPointer()
	}

	createResp, err := r.Client.CreateNet(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource Net",
			"Error: "+err.Error(),
		)
		return
	}

	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	net := ptr.From(createResp.Net)

	diag := createOAPITagsFW(ctx, r.Client, createTimeout, data.Tags, net.NetId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	data.NetId = to.String(net.NetId)
	data, err = setNetState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set net state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err = setNetState(ctx, r, data)
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

func (r *resourceNet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData NetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.NetId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := setNetState(ctx, r, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set net state",
			"Error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceNet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	delReq := osc.DeleteNetRequest{
		NetId: data.NetId.ValueString(),
	}
	_, err := r.Client.DeleteNet(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete net",
			"Error: "+err.Error(),
		)
	}
}

func setNetState(ctx context.Context, r *resourceNet, data NetModel) (NetModel, error) {
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
	tags, diag := flattenOAPITagsFW(ctx, net.Tags)
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
