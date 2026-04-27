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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
)

var (
	_ resource.Resource              = &nicPrivateIpResource{}
	_ resource.ResourceWithConfigure = &nicPrivateIpResource{}
)

const (
	nicPrivateIpErrCreate = "Unable to link Private Ips to the NIC"
	nicPrivateIpErrRead   = "Unable to read related NIC"
	nicPrivateIpErrDelete = "Unable to unlink Private Ips from the NIC"
)

type NicPrivateIpsModel struct {
	AllowRelink             types.Bool     `tfsdk:"allow_relink"`
	SecondaryPrivateIpCount types.Int64    `tfsdk:"secondary_private_ip_count"`
	NicId                   types.String   `tfsdk:"nic_id"`
	PrivateIps              types.List     `tfsdk:"private_ips"`
	PrimaryPrivateIp        types.String   `tfsdk:"primary_private_ip"`
	Id                      types.String   `tfsdk:"id"`
	RequestId               types.String   `tfsdk:"request_id"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
}

type nicPrivateIpResource struct {
	Client *osc.Client
}

func NewResourceNicPrivateIp() resource.Resource {
	return &nicPrivateIpResource{}
}

func (r *nicPrivateIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *nicPrivateIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nic_private_ip"
}

func (r *nicPrivateIpResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("private_ips"),
			path.MatchRoot("secondary_private_ip_count"),
		),
	}
}

func (r *nicPrivateIpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"allow_relink": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"secondary_private_ip_count": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"nic_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_ips": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_private_ip": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *nicPrivateIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NicPrivateIpsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	linkReq := osc.LinkPrivateIpsRequest{
		NicId: data.NicId.ValueString(),
	}

	if fwhelpers.IsSet(data.AllowRelink) {
		linkReq.AllowRelink = data.AllowRelink.ValueBoolPointer()
	}

	if fwhelpers.IsSet(data.SecondaryPrivateIpCount) {
		linkReq.SecondaryPrivateIpCount = new(int(data.SecondaryPrivateIpCount.ValueInt64()))
	}

	if fwhelpers.IsSet(data.PrivateIps) {
		ips, diag := to.Slice[string](ctx, data.PrivateIps)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		linkReq.PrivateIps = new(ips)
	}

	_, err := r.Client.LinkPrivateIps(ctx, linkReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(nicPrivateIpErrCreate, err.Error())
		return
	}

	data.Id = data.NicId

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(nicPrivateIpErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicPrivateIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NicPrivateIpsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	switch {
	case errors.Is(err, ErrResourceEmpty):
		resp.State.RemoveResource(ctx)
		return
	case err != nil:
		resp.Diagnostics.AddError(nicPrivateIpErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicPrivateIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *nicPrivateIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NicPrivateIpsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	privateIps, diag := to.Slice[string](ctx, data.PrivateIps)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	if len(privateIps) == 0 {
		return
	}

	_, err := r.Client.UnlinkPrivateIps(ctx, osc.UnlinkPrivateIpsRequest{
		NicId:      data.Id.ValueString(),
		PrivateIps: privateIps,
	}, options.WithRetryTimeout(timeout))
	switch {
	case osc.HasErrorCode(err, []string{"9108"}):
	// IPs already unlinked from the NIC
	case err != nil:
		resp.Diagnostics.AddError(nicPrivateIpErrDelete, err.Error())
		return
	}
}

func (r *nicPrivateIpResource) read(ctx context.Context, timeout time.Duration, data NicPrivateIpsModel) (NicPrivateIpsModel, error) {
	nicId := data.Id.ValueString()

	req := osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{NicIds: &[]string{nicId}},
	}

	resp, err := r.Client.ReadNics(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Nics == nil || len(*resp.Nics) == 0 {
		return data, ErrResourceEmpty
	}

	nic := (*resp.Nics)[0]

	privateIps, primaryIp := r.flattenPrivateIps(nic.PrivateIps)
	privateIpsList, diag := to.List(ctx, privateIps)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.PrivateIps = privateIpsList
	data.SecondaryPrivateIpCount = to.Int64(len(privateIps))
	data.Id = to.String(nic.NicId)
	data.PrimaryPrivateIp = to.String(primaryIp)
	data.NicId = to.String(nic.NicId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return data, nil
}

func (r *nicPrivateIpResource) flattenPrivateIps(ips []osc.PrivateIp) ([]string, string) {
	var primaryIp string

	return lo.Reduce(ips, func(acc []string, ip osc.PrivateIp, _ int) []string {
		if ip.IsPrimary {
			primaryIp = ip.PrivateIp
			return acc
		}
		return append(acc, ip.PrivateIp)
	}, []string{}), primaryIp
}
