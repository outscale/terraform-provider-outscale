package outscale

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var (
	_ resource.Resource              = &netResource{}
	_ resource.ResourceWithConfigure = &netResource{}
)

type NetModel struct {
	DhcpOptionsSetId types.String   `tfsdk:"dhcp_options_set_id"`
	IpRange          types.String   `tfsdk:"ip_range"`
	NetId            types.String   `tfsdk:"net_id"`
	State            types.String   `tfsdk:"state"`
	Tenancy          types.String   `tfsdk:"tenancy"`
	Tags             []ResourceTag  `tfsdk:"tags"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	RequestId        types.String   `tfsdk:"request_id"`
	Id               types.String   `tfsdk:"id"`
}

type ResourceTag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type netResource struct {
	Client *oscgo.APIClient
}

func NewResourceNet() resource.Resource {
	return &netResource{}
}

func (r *netResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *netResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	data.NetId = types.StringValue(net_id)
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

func (r *netResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net"
}
func (r *netResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
					fwmodifyplan.ForceNewFramework(),
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

func (r *netResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oscgo.CreateNetRequest{
		IpRange: data.IpRange.ValueString(),
	}

	if utils.IsSet(data.Tenancy) {
		createReq.SetTenancy(data.Tenancy.ValueString())
	}
	var createResp oscgo.CreateNetResponse

	err = retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.NetApi.CreateNet(ctx).CreateNetRequest(createReq).Execute()
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

	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	net := createResp.GetNet()
	if len(data.Tags) > 0 {
		err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(data.Tags), net.GetNetId())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add Tags on outscale_net resource",
				"Error: "+utils.GetErrorResponse(err).Error(),
			)
			return
		}

	}
	data.NetId = types.StringValue(net.GetNetId())
	data, err = setNetState(ctx, r, data)

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

func (r *netResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err = setNetState(ctx, r, data)
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

func (r *netResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		tagsPlan, tagsState []ResourceTag
		resourceId          types.String
		err                 error
	)

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tags"), &tagsPlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tags"), &tagsState)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("net_id"), &resourceId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !reflect.DeepEqual(tagsPlan, tagsState) {
		toRemove, toCreate := diffOSCAPITags(tagsToOSCResourceTag(tagsPlan), tagsToOSCResourceTag(tagsState))
		err := updateFrameworkTags(ctx, r.Client, toCreate, toRemove, resourceId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Tags on net resource",
				"Error: "+utils.GetErrorResponse(err).Error(),
			)
			return
		}
	}
	var data NetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err = setNetState(ctx, r, data)
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

func (r *netResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delReq := oscgo.DeleteNetRequest{
		NetId: data.NetId.ValueString(),
	}
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.NetApi.DeleteNet(ctx).DeleteNetRequest(delReq).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete net",
			"Error: "+err.Error(),
		)
		return
	}
}

func TagsSchema() *schema.SetNestedBlock {
	return &schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Required: true,
				},
				"value": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func tagsToOSCResourceTag(tags []ResourceTag) []oscgo.ResourceTag {
	result := make([]oscgo.ResourceTag, 0, len(tags))
	for _, tag := range tags {
		rTag := oscgo.NewResourceTag(tag.Key.ValueString(), tag.Value.ValueString())
		result = append(result, *rTag)
	}
	return result
}

func setNetState(ctx context.Context, r *netResource, data NetModel) (NetModel, error) {
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
	data.Id = types.StringValue(net.GetNetId())
	data.Tags = getTagsFromApiResponse(net.GetTags())
	data.NetId = types.StringValue(net.GetNetId())
	data.DhcpOptionsSetId = types.StringValue(net.GetDhcpOptionsSetId())
	data.IpRange = types.StringValue(net.GetIpRange())
	data.State = types.StringValue(net.GetState())
	data.Tenancy = types.StringValue(net.GetTenancy())
	return data, nil
}
