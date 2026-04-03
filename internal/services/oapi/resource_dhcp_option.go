package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                     = &dhcpOptionResource{}
	_ resource.ResourceWithConfigure        = &dhcpOptionResource{}
	_ resource.ResourceWithImportState      = &dhcpOptionResource{}
	_ resource.ResourceWithConfigValidators = &dhcpOptionResource{}
)

const (
	dhcpErrCreate  = "Unable to create DHCP Option"
	dhcpErrRead    = "Unable to read DHCP Option"
	dhcpErrDelete  = "Unable to delete DHCP Option"
	dhcpErrState   = "Unable to set DHCP Option state"
	dhcpErrDetach  = "Unable to detach DHCP Option from Nets"
	dhcpErrGetNets = "Unable to get attached Nets of DHCP Option"
)

type dhcpOptionModel struct {
	DomainName        types.String   `tfsdk:"domain_name"`
	DomainNameServers types.List     `tfsdk:"domain_name_servers"`
	LogServers        types.List     `tfsdk:"log_servers"`
	NtpServers        types.List     `tfsdk:"ntp_servers"`
	Default           types.Bool     `tfsdk:"default"`
	DhcpOptionsSetId  types.String   `tfsdk:"dhcp_options_set_id"`
	Id                types.String   `tfsdk:"id"`
	RequestId         types.String   `tfsdk:"request_id"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	TagsModel
}

type dhcpOptionResource struct {
	Client *osc.Client
}

func NewResourceDhcpOption() resource.Resource {
	return &dhcpOptionResource{}
}

func (r *dhcpOptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dhcpOptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_option"
}

func (r *dhcpOptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id := req.ID
	if id == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import DHCP option identifier. Got: %v", req.ID),
		)
		return
	}

	var data dhcpOptionModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(id)
	data.DhcpOptionsSetId = to.String(id)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal
	data.Tags = TagsNull()
	data.DomainNameServers = types.ListNull(types.StringType)
	data.LogServers = types.ListNull(types.StringType)
	data.NtpServers = types.ListNull(types.StringType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dhcpOptionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("domain_name"),
			path.MatchRoot("domain_name_servers"),
			path.MatchRoot("log_servers"),
			path.MatchRoot("ntp_servers"),
		),
	}
}

func (r *dhcpOptionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_name_servers": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
			},
			"log_servers": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
			},
			"ntp_servers": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
			},
			"default": schema.BoolAttribute{
				Computed: true,
			},
			"dhcp_options_set_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *dhcpOptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dhcpOptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateDhcpOptionsRequest{}

	if fwhelpers.IsSet(data.DomainName) {
		createReq.DomainName = data.DomainName.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.DomainNameServers) {
		servers, diag := to.Slice[string](ctx, data.DomainNameServers)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.DomainNameServers = &servers
	}
	if fwhelpers.IsSet(data.LogServers) {
		servers, diag := to.Slice[string](ctx, data.LogServers)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.LogServers = &servers
	}
	if fwhelpers.IsSet(data.NtpServers) {
		servers, diag := to.Slice[string](ctx, data.NtpServers)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.NtpServers = &servers
	}

	createResp, err := r.Client.CreateDhcpOptions(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(dhcpErrCreate, err.Error())
		return
	}

	dhcpId := ptr.From(createResp.DhcpOptionsSet.DhcpOptionsSetId)
	data.Id = to.String(dhcpId)
	data.DhcpOptionsSetId = to.String(dhcpId)

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, dhcpId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(dhcpErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *dhcpOptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dhcpOptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(dhcpErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *dhcpOptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData dhcpOptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.Id.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(dhcpErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *dhcpOptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dhcpOptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	// If not default DHCP Option, detach from any Nets before deleting
	if !data.Default.ValueBool() {
		nets, err := r.Client.ReadNets(ctx, osc.ReadNetsRequest{
			Filters: &osc.FiltersNet{
				DhcpOptionsSetIds: &[]string{data.Id.ValueString()},
			},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(dhcpErrGetNets, err.Error())
			return
		}

		for _, net := range ptr.From(nets.Nets) {
			_, err := r.Client.UpdateNet(ctx, osc.UpdateNetRequest{
				DhcpOptionsSetId: "default",
				NetId:            net.NetId,
			}, options.WithRetryTimeout(timeout))
			if err != nil {
				resp.Diagnostics.AddError(dhcpErrDetach, err.Error())
				return
			}
		}
	}

	_, err := r.Client.DeleteDhcpOptions(ctx, osc.DeleteDhcpOptionsRequest{
		DhcpOptionsSetId: data.Id.ValueString(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(dhcpErrDelete, err.Error())
	}
}

func (r *dhcpOptionResource) read(ctx context.Context, timeout time.Duration, data dhcpOptionModel) (dhcpOptionModel, error) {
	readReq := osc.ReadDhcpOptionsRequest{
		Filters: &osc.FiltersDhcpOptions{
			DhcpOptionsSetIds: &[]string{data.Id.ValueString()},
		},
	}

	resp, err := r.Client.ReadDhcpOptions(ctx, readReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if len(ptr.From(resp.DhcpOptionsSets)) == 0 {
		return data, ErrResourceEmpty
	}

	dhcp := (*resp.DhcpOptionsSets)[0]

	tags, diag := flattenOAPITagsFW(ctx, ptr.From(dhcp.Tags))
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	domainNameServers, diag := to.List(ctx, ptr.From(dhcp.DomainNameServers))
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	logServers, diag := to.List(ctx, ptr.From(dhcp.LogServers))
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	ntpServers, diag := to.List(ctx, ptr.From(dhcp.NtpServers))
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Tags = tags
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.Id = to.String(dhcp.DhcpOptionsSetId)
	data.DhcpOptionsSetId = to.String(dhcp.DhcpOptionsSetId)
	data.DomainName = to.String(ptr.From(dhcp.DomainName))
	data.DomainNameServers = domainNameServers
	data.LogServers = logServers
	data.NtpServers = ntpServers
	data.Default = to.Bool(dhcp.Default)

	return data, nil
}
