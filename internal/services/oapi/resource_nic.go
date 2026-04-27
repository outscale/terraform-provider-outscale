package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/framework/stateconf"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

var (
	_ resource.Resource                   = &nicResource{}
	_ resource.ResourceWithConfigure      = &nicResource{}
	_ resource.ResourceWithImportState    = &nicResource{}
	_ resource.ResourceWithModifyPlan     = &nicResource{}
	_ resource.ResourceWithValidateConfig = &nicResource{}
)

const (
	nicErrCreate = "Unable to create NIC"
	nicErrRead   = "Unable to read NIC"
	nicErrUpdate = "Unable to update NIC"
	nicErrDelete = "Unable to delete NIC"
	nicErrState  = "Unable to set NIC state"
	nicErrDetach = "Unable to detach NIC"
)

type nicModel struct {
	Description         types.String   `tfsdk:"description"`
	SecurityGroupIds    types.Set      `tfsdk:"security_group_ids"`
	SubnetId            types.String   `tfsdk:"subnet_id"`
	LinkPublicIp        types.Set      `tfsdk:"link_public_ip"`
	LinkNic             types.Set      `tfsdk:"link_nic"`
	SubregionName       types.String   `tfsdk:"subregion_name"`
	SecurityGroups      types.List     `tfsdk:"security_groups"`
	MacAddress          types.String   `tfsdk:"mac_address"`
	NicID               types.String   `tfsdk:"nic_id"`
	AccountID           types.String   `tfsdk:"account_id"`
	PrivateDNSName      types.String   `tfsdk:"private_dns_name"`
	PrivateIp           types.String   `tfsdk:"private_ip"`
	PrivateIps          types.Set      `tfsdk:"private_ips"`
	RequesterManaged    types.Bool     `tfsdk:"requester_managed"`
	IsSourceDestChecked types.Bool     `tfsdk:"is_source_dest_checked"`
	State               types.String   `tfsdk:"state"`
	NetId               types.String   `tfsdk:"net_id"`
	Id                  types.String   `tfsdk:"id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	RequestId           types.String   `tfsdk:"request_id"`
	TagsModel
}

type nicLinkPublicIPModel struct {
	PublicIPId        types.String `tfsdk:"public_ip_id"`
	LinkPublicIPId    types.String `tfsdk:"link_public_ip_id"`
	PublicIPAccountID types.String `tfsdk:"public_ip_account_id"`
	PublicDNSName     types.String `tfsdk:"public_dns_name"`
	PublicIP          types.String `tfsdk:"public_ip"`
}

var linkPublicIPAttrTypes = fwhelpers.GetAttrTypes(nicLinkPublicIPModel{})

type nicLinkNicModel struct {
	LinkNicID          types.String `tfsdk:"link_nic_id"`
	DeleteOnVMDeletion types.String `tfsdk:"delete_on_vm_deletion"`
	DeviceNumber       types.Int64  `tfsdk:"device_number"`
	VMID               types.String `tfsdk:"vm_id"`
	VMAccountID        types.String `tfsdk:"vm_account_id"`
	State              types.String `tfsdk:"state"`
}

type nicSecurityGroupModel struct {
	SecurityGroupID   types.String `tfsdk:"security_group_id"`
	SecurityGroupName types.String `tfsdk:"security_group_name"`
}

type nicPrivateIpModel struct {
	LinkPublicIp   types.Set    `tfsdk:"link_public_ip"`
	PrivateDNSName types.String `tfsdk:"private_dns_name"`
	PrivateIP      types.String `tfsdk:"private_ip"`
	IsPrimary      types.Bool   `tfsdk:"is_primary"`
}

var nicPrivateIpAttrTypes = map[string]attr.Type{
	"link_public_ip":   types.SetType{ElemType: to.Object(linkPublicIPAttrTypes)},
	"private_dns_name": types.StringType,
	"private_ip":       types.StringType,
	"is_primary":       types.BoolType,
}

type nicResource struct {
	Client *osc.Client
}

func NewResourceNic() resource.Resource {
	return &nicResource{}
}

func (r *nicResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *nicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nic"
}

func (r *nicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *nicResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return // destroy: we do nothing
	}

	if req.State.Raw.IsNull() {
		// create: the SetNestedBlockMerge plan modifier on private_ips already
		// injects an unknown set when the block is omitted
		return
	}

	var stateData, planData, configData nicModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	descriptionUnchanged := configData.Description.IsNull() || configData.Description.IsUnknown() || configData.Description.Equal(stateData.Description)
	tagsUnchanged := configData.Tags.IsNull() || configData.Tags.IsUnknown() || len(configData.Tags.Elements()) == 0 || configData.Tags.Equal(stateData.Tags)

	noInputChange := planData.SubnetId.Equal(stateData.SubnetId) &&
		planData.SecurityGroupIds.Equal(stateData.SecurityGroupIds) &&
		planData.PrivateIps.Equal(stateData.PrivateIps) &&
		descriptionUnchanged &&
		tagsUnchanged

	if noInputChange {
		// When no input has changed, copy computed
		// fields from state into the plan so the framework does not mark
		// them as unknown and produce a perpetual diff
		planData.Description = stateData.Description
		planData.IsSourceDestChecked = stateData.IsSourceDestChecked
		planData.LinkNic = stateData.LinkNic
		planData.LinkPublicIp = stateData.LinkPublicIp
		planData.PrivateDNSName = stateData.PrivateDNSName
		planData.PrivateIp = stateData.PrivateIp
		planData.RequestId = stateData.RequestId
		planData.RequesterManaged = stateData.RequesterManaged
		planData.SecurityGroupIds = stateData.SecurityGroupIds
		planData.SecurityGroups = stateData.SecurityGroups
		planData.State = stateData.State
		planData.Tags = stateData.Tags
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *nicResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configPrivateIps types.Set
	diag := req.Config.GetAttribute(ctx, path.Root("private_ips"), &configPrivateIps)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	if !fwhelpers.IsSet(configPrivateIps) {
		return
	}
	if len(configPrivateIps.Elements()) == 0 {
		return
	}

	privateIpsModel, diags := to.Slice[nicPrivateIpModel](ctx, configPrivateIps)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}
	primaryCount := lo.CountBy(privateIpsModel, func(m nicPrivateIpModel) bool {
		return fwhelpers.IsSet(m.IsPrimary) && m.IsPrimary.ValueBool()
	})

	switch {
	case primaryCount == 0:
		resp.Diagnostics.AddAttributeError(
			path.Root("private_ips"),
			"Invalid private_ips configuration",
			"At least one private_ips block must set is_primary = true.",
		)
	case primaryCount > 1:
		resp.Diagnostics.AddAttributeError(
			path.Root("private_ips"),
			"Invalid private_ips configuration",
			"Only one private_ips block can set is_primary = true.",
		)
	}
}

func (r *nicResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"tags": TagsSchemaFW(),
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"private_ips": schema.SetNestedBlock{
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIf(
						nicPrivateIPsRequiresReplace,
						"Replace only when changing requested primary private IP",
						"Replace only when changing requested primary private IP",
					),
					modifyplans.SetNestedBlockMerge(modifyplans.SetNestedBlockMergeOptions{
						ObjectType: to.Object(nicPrivateIpAttrTypes),
						Fields:     []string{"is_primary", "private_ip"},
					}),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"link_public_ip": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"public_ip_id": schema.StringAttribute{
										Computed: true,
									},
									"link_public_ip_id": schema.StringAttribute{
										Computed: true,
									},
									"public_ip_account_id": schema.StringAttribute{
										Computed: true,
									},
									"public_dns_name": schema.StringAttribute{
										Computed: true,
									},
									"public_ip": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"private_dns_name": schema.StringAttribute{
							Computed: true,
						},
						"private_ip": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"is_primary": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"security_group_ids": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"link_public_ip": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"public_ip_id": schema.StringAttribute{
							Computed: true,
						},
						"link_public_ip_id": schema.StringAttribute{
							Computed: true,
						},
						"public_ip_account_id": schema.StringAttribute{
							Computed: true,
						},
						"public_dns_name": schema.StringAttribute{
							Computed: true,
						},
						"public_ip": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"link_nic": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"link_nic_id": schema.StringAttribute{
							Computed: true,
						},
						"delete_on_vm_deletion": schema.StringAttribute{
							Computed: true,
						},
						"device_number": schema.Int64Attribute{
							Computed: true,
						},
						"vm_id": schema.StringAttribute{
							Computed: true,
						},
						"vm_account_id": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"subregion_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"security_group_id": schema.StringAttribute{
							Computed: true,
						},
						"security_group_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"mac_address": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nic_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_dns_name": schema.StringAttribute{
				Computed: true,
			},
			// private_ip was present in the SDKv2 schema but never documented or used
			// we keep it to not break the state
			"private_ip": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"requester_managed": schema.BoolAttribute{
				Computed: true,
			},
			"is_source_dest_checked": schema.BoolAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
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

func (r *nicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data nicModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateNicRequest{
		SubnetId: data.SubnetId.ValueString(),
	}
	if fwhelpers.IsSet(data.Description) {
		createReq.Description = data.Description.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.SecurityGroupIds) {
		ids, diag := to.Slice[string](ctx, data.SecurityGroupIds)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.SecurityGroupIds = &ids
	}
	if fwhelpers.IsSet(data.PrivateIps) {
		ips, diag := r.expandPrivateIPLight(ctx, data.PrivateIps)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		createReq.PrivateIps = &ips
	}

	createResp, err := r.Client.CreateNic(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(nicErrCreate, err.Error())
		return
	}

	nicID := createResp.Nic.NicId
	data.Id = to.String(nicID)
	data.NicID = to.String(nicID)

	diag := createOAPITagsFW(ctx, r.Client, timeout, data.Tags, nicID)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(nicErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data nicModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		resp.Diagnostics.AddError(nicErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *nicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData nicModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	statePrivateIPs, diag := privateIPsStringsFromSet(ctx, stateData.PrivateIps)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	planPrivateIPs, diag := privateIPsStringsFromSet(ctx, planData.PrivateIps)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	removed := lo.Without(statePrivateIPs, planPrivateIPs...)
	created := lo.Without(planPrivateIPs, statePrivateIPs...)

	if len(removed) > 0 {
		_, err := r.Client.UnlinkPrivateIps(ctx, osc.UnlinkPrivateIpsRequest{
			NicId:      stateData.Id.ValueString(),
			PrivateIps: removed,
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(nicErrUpdate, fmt.Sprintf("failure to unassign private ips: %s", err))
			return
		}
	}

	if len(created) > 0 {
		_, err := r.Client.LinkPrivateIps(ctx, osc.LinkPrivateIpsRequest{
			NicId:      stateData.Id.ValueString(),
			PrivateIps: &created,
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(nicErrUpdate, fmt.Sprintf("failure to assign private ips: %s", err))
			return
		}
	}

	updateReq := osc.UpdateNicRequest{
		NicId: stateData.Id.ValueString(),
	}
	var doUpdate bool

	if !stateData.SecurityGroupIds.Equal(planData.SecurityGroupIds) {
		sgIDs, diag := to.Slice[string](ctx, planData.SecurityGroupIds)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		updateReq.SecurityGroupIds = &sgIDs
		doUpdate = true
	}

	if fwhelpers.HasChange(planData.Description, stateData.Description) && planData.Description.ValueString() != "" {
		updateReq.Description = planData.Description.ValueStringPointer()
		doUpdate = true
	}

	if doUpdate {
		_, err := r.Client.UpdateNic(ctx, updateReq, options.WithRetryTimeout(timeout))
		if err != nil {
			resp.Diagnostics.AddError(nicErrUpdate, err.Error())
			return
		}
	}

	if diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.Id.ValueString()); fwhelpers.CheckDiags(resp, diag) {
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(nicErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *nicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data nicModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	if err := r.detachNic(ctx, data.Id.ValueString(), timeout); err != nil {
		resp.Diagnostics.AddError(nicErrDetach, err.Error())
		return
	}

	_, err := r.Client.DeleteNic(ctx, osc.DeleteNicRequest{
		NicId: data.Id.ValueString(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(nicErrDelete, err.Error())
	}
}

func (r *nicResource) read(ctx context.Context, timeout time.Duration, data nicModel) (nicModel, error) {
	resp, err := r.Client.ReadNics(ctx, osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{
			NicIds: &[]string{data.Id.ValueString()},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.Nics == nil || len(*resp.Nics) == 0 {
		return data, ErrResourceEmpty
	}

	nic := (*resp.Nics)[0]

	tags, diag := flattenOAPITagsFW(ctx, nic.Tags)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	securityGroupsModel, securityGroupsIdsModel := r.flattenSecurityGroups(nic.SecurityGroups)
	securityGroupIDs, diag := to.Set(ctx, securityGroupsIdsModel)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	securityGroups, diag := to.ListObject(ctx, securityGroupsModel, to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	linkPublicIP, diag := to.SetObject(ctx, r.flattenLinkPublicIP(nic.LinkPublicIp))
	if diag.HasError() {
		return data, from.Diag(diag)
	}
	linkNic, diag := to.SetObject(ctx, r.flattenLinkNic(nic.LinkNic))
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	ipsModel, err := r.flattenPrivateIPs(ctx, nic.PrivateIps)
	if err != nil {
		return data, err
	}
	privateIps, diag := to.SetFromAttrType(ctx, ipsModel, to.Object(nicPrivateIpAttrTypes), to.ZeroValueAsEmpty)
	if diag.HasError() {
		return data, from.Diag(diag)
	}

	data.Id = to.String(nic.NicId)
	data.Description = to.String(nic.Description)
	data.SubnetId = to.String(nic.SubnetId)
	data.LinkPublicIp = linkPublicIP
	data.LinkNic = linkNic
	data.SubregionName = to.String(nic.SubregionName)
	data.SecurityGroups = securityGroups
	data.SecurityGroupIds = securityGroupIDs
	data.MacAddress = to.String(nic.MacAddress)
	data.NicID = to.String(nic.NicId)
	data.AccountID = to.String(nic.AccountId)
	data.PrivateDNSName = to.String(nic.PrivateDnsName)
	data.PrivateIps = privateIps
	data.IsSourceDestChecked = to.Bool(nic.IsSourceDestChecked)
	data.State = to.String(string(nic.State))
	data.Tags = tags
	data.NetId = to.String(nic.NetId)
	data.RequestId = to.String(resp.ResponseContext.RequestId)
	data.RequesterManaged = types.BoolNull()

	// Always stored as an empty string in the SDKv2
	data.PrivateIp = to.String("")

	return data, nil
}

func (r *nicResource) flattenSecurityGroups(ips []osc.SecurityGroupLight) ([]nicSecurityGroupModel, []string) {
	sgIDs := make([]string, 0, len(ips))

	sgModels := lo.Map(ips, func(pip osc.SecurityGroupLight, _ int) nicSecurityGroupModel {
		sgIDs = append(sgIDs, pip.SecurityGroupId)
		return nicSecurityGroupModel{
			SecurityGroupID:   to.String(pip.SecurityGroupId),
			SecurityGroupName: to.String(pip.SecurityGroupName),
		}
	})

	return sgModels, sgIDs
}

func (r *nicResource) flattenLinkPublicIP(ip *osc.LinkPublicIp) []nicLinkPublicIPModel {
	if ip == nil {
		return nil
	}

	return []nicLinkPublicIPModel{{
		PublicIPId:        to.String(ip.PublicIpId),
		LinkPublicIPId:    to.String(ip.LinkPublicIpId),
		PublicIPAccountID: to.String(ip.PublicIpAccountId),
		PublicDNSName:     to.String(ip.PublicDnsName),
		PublicIP:          to.String(ip.PublicIp),
	}}
}

func (r *nicResource) flattenLinkNic(nic *osc.LinkNic) []nicLinkNicModel {
	if nic == nil {
		return nil
	}

	return []nicLinkNicModel{{
		LinkNicID:          to.String(nic.LinkNicId),
		DeleteOnVMDeletion: to.String(cast.ToString(nic.DeleteOnVmDeletion)),
		DeviceNumber:       to.Int64(nic.DeviceNumber),
		VMID:               to.String(nic.VmId),
		VMAccountID:        to.String(nic.VmAccountId),
		State:              to.String(nic.State),
	}}
}

func (r *nicResource) flattenPrivateIPs(ctx context.Context, ips []osc.PrivateIp) ([]nicPrivateIpModel, error) {
	return lo.MapErr(ips, func(ip osc.PrivateIp, _ int) (nicPrivateIpModel, error) {
		set, diag := to.SetObject(ctx, r.flattenLinkPublicIP(ip.LinkPublicIp))
		if diag.HasError() {
			return nicPrivateIpModel{}, from.Diag(diag)
		}

		return nicPrivateIpModel{
			LinkPublicIp:   set,
			PrivateDNSName: to.String(ip.PrivateDnsName),
			PrivateIP:      to.String(ip.PrivateIp),
			IsPrimary:      to.Bool(ip.IsPrimary),
		}, nil
	})
}

func (r *nicResource) detachNic(ctx context.Context, nicID string, timeout time.Duration) error {
	stateConf := &stateconf.StateChangeConf[osc.LinkNicState]{
		Pending: stateconf.States(osc.LinkNicStateAttaching, osc.LinkNicStateDetaching),
		Target:  stateconf.States(osc.LinkNicStateAttached, osc.LinkNicStateDetached, osc.LinkNicState("failed")),
		Timeout: timeout,
		Refresh: r.nicLinkStateRefresh(nicID, timeout),
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for nic (%s) to be detached: %w", nicID, err)
	}

	readResp, ok := resp.(*osc.ReadNicsResponse)
	if !ok || readResp == nil || readResp.Nics == nil || len(*readResp.Nics) == 0 {
		return fmt.Errorf("nic (%s) not found", nicID)
	}

	linkNic := ptr.From((*readResp.Nics)[0].LinkNic)
	if linkNic != (osc.LinkNic{}) {
		log.Printf("[DEBUG] Waiting for NIC (%s) to become detached", nicID)
		_, err := r.Client.UnlinkNic(ctx, osc.UnlinkNicRequest{LinkNicId: linkNic.LinkNicId}, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *nicResource) nicLinkStateRefresh(nicID string, timeout time.Duration) stateconf.StateRefreshFunc[osc.LinkNicState] {
	return func(ctx context.Context) (any, osc.LinkNicState, error) {
		resp, err := r.Client.ReadNics(ctx, osc.ReadNicsRequest{
			Filters: &osc.FiltersNic{NicIds: &[]string{nicID}},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, osc.LinkNicState("failed"), err
		}
		if resp.Nics == nil || len(*resp.Nics) == 0 {
			return nil, osc.LinkNicState("failed"), fmt.Errorf("nic (%s) not found", nicID)
		}

		linkNic := ptr.From((*resp.Nics)[0].LinkNic)
		if linkNic == (osc.LinkNic{}) {
			return resp, osc.LinkNicStateDetached, nil
		}

		return resp, linkNic.State, nil
	}
}

func (r *nicResource) expandPrivateIPLight(ctx context.Context, set types.Set) ([]osc.PrivateIpLight, diag.Diagnostics) {
	models, diag := to.Slice[nicPrivateIpModel](ctx, set)
	if diag.HasError() {
		return nil, diag
	}

	return lo.FlatMap(models, func(m nicPrivateIpModel, _ int) []osc.PrivateIpLight {
		if !fwhelpers.IsSet(m.PrivateIP) {
			return nil
		}
		return []osc.PrivateIpLight{{
			PrivateIp: m.PrivateIP.ValueString(),
			IsPrimary: m.IsPrimary.ValueBool(),
		}}
	}), nil
}

func nicPrivateIPsRequiresReplace(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifier.RequiresReplaceIfFuncResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	requiresReplace, diags := func() (bool, diag.Diagnostics) {
		statePrimary, hasStatePrimary, diags := primaryPrivateIPFromSet(ctx, req.StateValue)
		if diags.HasError() || !hasStatePrimary {
			return false, diags
		}
		configModels, configDiags := to.Slice[nicPrivateIpModel](ctx, req.ConfigValue)
		diags.Append(configDiags...)
		if diags.HasError() {
			return false, diags
		}

		for _, privateIP := range configModels {
			if !fwhelpers.IsSet(privateIP.IsPrimary) || !privateIP.IsPrimary.ValueBool() {
				continue
			}
			if !fwhelpers.IsSet(privateIP.PrivateIP) {
				continue
			}
			if privateIP.PrivateIP.ValueString() != statePrimary {
				return true, diags
			}
		}
		return false, diags
	}()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplace
}

func primaryPrivateIPFromSet(ctx context.Context, set types.Set) (string, bool, diag.Diagnostics) {
	models, diags := to.Slice[nicPrivateIpModel](ctx, set)
	if diags.HasError() {
		return "", false, diags
	}

	for _, privateIP := range models {
		if !fwhelpers.IsSet(privateIP.IsPrimary) || !privateIP.IsPrimary.ValueBool() {
			continue
		}
		if !fwhelpers.IsSet(privateIP.PrivateIP) {
			return "", false, diags
		}
		return privateIP.PrivateIP.ValueString(), true, diags
	}

	return "", false, diags
}

func privateIPsStringsFromSet(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	models, diag := to.Slice[nicPrivateIpModel](ctx, set)
	if diag.HasError() {
		return nil, diag
	}

	privateIPs := lo.FlatMap(models, func(m nicPrivateIpModel, _ int) []string {
		if !fwhelpers.IsSet(m.PrivateIP) {
			return nil
		}
		return []string{m.PrivateIP.ValueString()}
	})

	return privateIPs, nil
}

// Used in SDKv2 datasource
func flattenLinkPublicIp(linkIp *osc.LinkPublicIp) []map[string]any {
	return []map[string]any{{
		"public_ip_id":         linkIp.PublicIpId,
		"link_public_ip_id":    linkIp.LinkPublicIpId,
		"public_ip_account_id": linkIp.PublicIpAccountId,
		"public_dns_name":      linkIp.PublicDnsName,
		"public_ip":            linkIp.PublicIp,
	}}
}

func flattenLinkNic(linkNic *osc.LinkNic) []map[string]any {
	return []map[string]any{{
		"link_nic_id":           linkNic.LinkNicId,
		"delete_on_vm_deletion": strconv.FormatBool(linkNic.DeleteOnVmDeletion),
		"device_number":         linkNic.DeviceNumber,
		"vm_id":                 linkNic.VmId,
		"vm_account_id":         linkNic.VmAccountId,
		"state":                 linkNic.State,
	}}
}
