package oapi

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &resourceSecurityGroup{}
	_ resource.ResourceWithConfigure   = &resourceSecurityGroup{}
	_ resource.ResourceWithImportState = &resourceSecurityGroup{}
	_ resource.ResourceWithModifyPlan  = &resourceSecurityGroup{}
)

type SecurityGroupModel struct {
	AccountId                 types.String   `tfsdk:"account_id"`
	Description               types.String   `tfsdk:"description"`
	InboundRules              types.List     `tfsdk:"inbound_rules"`
	NetId                     types.String   `tfsdk:"net_id"`
	OutboundRules             types.List     `tfsdk:"outbound_rules"`
	SecurityGroupId           types.String   `tfsdk:"security_group_id"`
	SecurityGroupName         types.String   `tfsdk:"security_group_name"`
	RemoveDefaultOutboundRule types.Bool     `tfsdk:"remove_default_outbound_rule"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
	RequestId                 types.String   `tfsdk:"request_id"`
	Id                        types.String   `tfsdk:"id"`
	TagsModel
}

type resourceSecurityGroup struct {
	Client *osc.Client
}

func NewResourceSecurityGroup() resource.Resource {
	return &resourceSecurityGroup{}
}

func (r *resourceSecurityGroup) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceSecurityGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	securityGroupId := req.ID

	if securityGroupId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Security Group identifier, got: %v", req.ID),
		)
		return
	}

	var data SecurityGroupModel
	var timeouts timeouts.Value
	diag := resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	data.Timeouts = timeouts

	data.SecurityGroupId = to.String(securityGroupId)
	data.Id = to.String(securityGroupId)
	data.InboundRules = types.ListNull(securityGroupRulesModelAttrTypes)
	data.OutboundRules = types.ListNull(securityGroupRulesModelAttrTypes)
	data.Tags = TagsNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSecurityGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group"
}

func (r *resourceSecurityGroup) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Do not verify plan if this is a new resource
	if req.State.Raw.IsNull() {
		return
	}

	var stateData SecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isDestroy := req.Plan.Raw.IsNull()
	isReplace := false
	sameName := false

	if !isDestroy {
		var planData SecurityGroupModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// A replacement happens if any "RequiresReplace" attribute is changed
		isReplace = !planData.Description.Equal(stateData.Description) ||
			!planData.NetId.Equal(stateData.NetId) ||
			!planData.RemoveDefaultOutboundRule.Equal(stateData.RemoveDefaultOutboundRule) ||
			!planData.SecurityGroupName.Equal(stateData.SecurityGroupName)
		sameName = planData.SecurityGroupName.Equal(stateData.SecurityGroupName)
	}
	if !isDestroy && !isReplace {
		return
	}

	to, diag := stateData.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	vms, nics, lbus, diag := r.getLinkedResources(ctx, to, stateData.Id.ValueString())
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	resources := []string{}
	for _, vm := range vms {
		if len(vm.SecurityGroups) == 1 {
			resources = append(resources, vm.VmId)
		}
	}
	for _, nic := range nics {
		if len(nic.SecurityGroups) == 1 {
			resources = append(resources, nic.NicId)
		}
	}
	targetSG := []string{stateData.Id.ValueString()}
	for _, lbu := range lbus {
		if lbu.SecurityGroups != nil && slices.Equal(lbu.SecurityGroups, targetSG) {
			resources = append(resources, lbu.LoadBalancerName)
		}
	}

	if len(resources) > 0 {
		errDetailMsg := `The resource is being destroyed and this operation will fail because other resources depend on it as their only security group.`
		if isReplace {
			errDetailMsg += `

Add this lifecycle block to create the new Security Group before destroying the old one (prevents dependent resources from having zero security groups):

resource "outscale_security_group" "foo" {
  # ... configuration ...

  lifecycle {
    create_before_destroy = true
  }
}`
		}

		resp.Diagnostics.AddWarning(
			fmt.Sprintf("The Security Group is the unique Security Group of the following resources: %v.", resources),
			errDetailMsg,
		)

		if isReplace && sameName {
			resp.Diagnostics.AddWarning(
				"Security Group name conflict.",
				`When using "create_before_destroy", change the "security_group_name" to a different value. Both security groups exist temporarily during replacement, and the API requires unique names.`,
			)
		}
	}
}

func (r *resourceSecurityGroup) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"tags": TagsSchemaFW(),
		},
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"inbound_rules": schema.ListAttribute{
				Computed:    true,
				ElementType: securityGroupRulesModelAttrTypes,
			},
			"net_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
			},
			"outbound_rules": schema.ListAttribute{
				Computed:    true,
				ElementType: securityGroupRulesModelAttrTypes,
			},
			"security_group_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					modifyplans.ForceNewFramework(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"remove_default_outbound_rule": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(RemoveDefaultOutboundRuleDefaultValue),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(
						path.MatchRoot("net_id"),
					),
				},
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

func (r *resourceSecurityGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	if !fwhelpers.IsSet(data.SecurityGroupName) {
		data.SecurityGroupName = to.String(id.UniqueId())
	}

	createReq := osc.CreateSecurityGroupRequest{
		Description:       data.Description.ValueString(),
		SecurityGroupName: data.SecurityGroupName.ValueString(),
	}

	if !data.NetId.IsUnknown() && !data.NetId.IsNull() {
		createReq.NetId = data.NetId.ValueStringPointer()
	}

	createResp, err := r.Client.CreateSecurityGroup(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Security Group.",
			err.Error(),
		)
		return
	}
	sg := ptr.From(createResp.SecurityGroup)
	data.RequestId = to.String(createResp.ResponseContext.RequestId)
	data.Id = to.String(sg.SecurityGroupId)
	data.SecurityGroupId = to.String(sg.SecurityGroupId)

	if data.RemoveDefaultOutboundRule.ValueBool() {
		ipRange := "0.0.0.0/0"
		ipProtocol := "-1"
		emptySGReq := osc.DeleteSecurityGroupRuleRequest{
			Flow:            "Outbound",
			SecurityGroupId: sg.SecurityGroupId,
			IpRange:         &ipRange,
			IpProtocol:      &ipProtocol,
		}

		_, err := r.Client.DeleteSecurityGroupRule(ctx, emptySGReq, options.WithRetryTimeout(createTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to empty the Security Group rules.",
				err.Error(),
			)
			return
		}
	}

	diag = createOAPITagsFW(ctx, r.Client, createTimeout, data.Tags, sg.SecurityGroupId)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := setSecurityGroupState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Security Group state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceSecurityGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SecurityGroupModel
	diag := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := setSecurityGroupState(ctx, r, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Security Group API response values.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSecurityGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData SecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	diag := updateOAPITagsFW(ctx, r.Client, timeout, stateData.Tags, planData.Tags, stateData.SecurityGroupId.ValueString())
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := setSecurityGroupState(ctx, r, stateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Security Group state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSecurityGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	deleteTimeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	sgId := data.Id.ValueString()
	delReq := osc.DeleteSecurityGroupRequest{
		SecurityGroupId: &sgId,
	}

	sgLinked := false
	_, err := r.Client.DeleteSecurityGroup(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		oscErr := oapihelpers.GetError(err)
		if oscErr.Code == "9085" {
			sgLinked = true
		} else {
			resp.Diagnostics.AddError(
				"Unable to delete Security Group.",
				err.Error(),
			)
			return
		}
	}

	if sgLinked {
		diag = r.unlink(ctx, deleteTimeout, data)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		// // Retry on 409 as API can take time to see a security group as not in use anymore
		// err := oapihelpers.RetryOnCode(ctx, "9085", func() (resp interface{}, err error) {
		// 	return r.Client.DeleteSecurityGroup(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
		// }, deleteTimeout)
		err := oapihelpers.RetryOnCodes(ctx, []string{"9085"}, func() (resp any, err error) {
			return r.Client.DeleteSecurityGroup(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
		}, deleteTimeout)
		// err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		// 	_, err := r.Client.DeleteSecurityGroup(ctx, delReq, options.WithRetryTimeout(deleteTimeout))
		// 	if err != nil {
		// 		if osc.IsConflict(err) {
		// 			return retry.RetryableError(err)
		// 		}
		// 		return retry.NonRetryableError(err)
		// 	}
		// 	return nil
		// })
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to delete Security Group after unlinking from resources.",
				err.Error(),
			)
			return
		}
	}
}

func (r *resourceSecurityGroup) getLinkedResources(ctx context.Context, to time.Duration, sgId string) ([]osc.Vm, []osc.Nic, []osc.LoadBalancer, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get VMs linked to Security Group
	readVmReq := osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			SecurityGroupIds: &[]string{sgId},
		},
	}

	vms, err := r.Client.ReadVms(ctx, readVmReq, options.WithRetryTimeout(to))
	if err != nil {
		diags.AddError(
			"Unable to find VMs linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	// Get NICs linked to Security Group
	readNicsReq := osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{
			SecurityGroupIds: &[]string{sgId},
		},
	}

	nics, err := r.Client.ReadNics(ctx, readNicsReq, options.WithRetryTimeout(to))
	if err != nil {
		diags.AddError(
			"Unable to find NICs linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	// Get LBUs linked to Security Group
	lbus, err := r.Client.ReadLoadBalancers(ctx, osc.ReadLoadBalancersRequest{}, options.WithRetryTimeout(to))
	if err != nil {
		diags.AddError(
			"Unable to find Load Balancers linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	return *vms.Vms, *nics.Nics, *lbus.LoadBalancers, nil
}

func (r *resourceSecurityGroup) unlink(ctx context.Context, to time.Duration, data SecurityGroupModel) (diags diag.Diagnostics) {
	vms, nics, lbus, diag := r.getLinkedResources(ctx, to, data.Id.ValueString())
	if diag.HasError() {
		diags = append(diags, diag...)
		return diags
	}

	// Check if the Security Group is the unique Security Group of any VM
	vmsUniqueSG := lo.FilterMap(vms, func(vm osc.Vm, _ int) (string, bool) {
		return vm.VmId, (len(vm.SecurityGroups) == 1)
	})
	// Check if the Security Group is the unique Security Group of any NIC
	nicsUniqueSG := lo.FilterMap(nics, func(nic osc.Nic, _ int) (string, bool) {
		return nic.NicId, (len(nic.SecurityGroups) == 1)
	})
	targetSG := []string{data.Id.ValueString()}
	// Check if the Security Group is the unique Security Group of any LBU
	lbusUniqueSG := lo.FilterMap(lbus, func(lbu osc.LoadBalancer, _ int) (string, bool) {
		return lbu.LoadBalancerName, (lbu.SecurityGroups != nil && slices.Equal(lbu.SecurityGroups, targetSG))
	})

	if len(vmsUniqueSG) > 0 {
		diags.AddError(
			fmt.Sprintf("The Security Group is the unique Security Group of the following VMs: %v.", vmsUniqueSG),
			"The Security Group cannot be deleted and needs to be removed from the VMs first.",
		)
	}
	if len(nicsUniqueSG) > 0 {
		diags.AddError(
			fmt.Sprintf("The Security Group is the unique Security Group of the following NICs: %v.", nicsUniqueSG),
			"The Security Group cannot be deleted and needs to be removed from the NICs first.",
		)
	}
	if len(lbusUniqueSG) > 0 {
		diags.AddError(
			fmt.Sprintf("The Security Group is the unique Security Group of the following LBUs: %v.", lbusUniqueSG),
			"The Security Group cannot be deleted and needs to be removed from the LBUs first.",
		)
	}

	if diags.HasError() {
		return
	}

	for _, vm := range vms {
		// Get the Security Group IDs of the VM without the current Security Group
		sgIds := lo.FilterMap(vm.SecurityGroups, func(sg osc.SecurityGroupLight, _ int) (string, bool) {
			return sg.SecurityGroupId, (sg.SecurityGroupId != data.Id.ValueString())
		})
		updateReq := osc.UpdateVmRequest{
			VmId:             vm.VmId,
			SecurityGroupIds: sgIds,
		}

		_, err := r.Client.UpdateVm(ctx, updateReq, options.WithRetryTimeout(to))
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from VM (%s)", data.Id.ValueString(), vm.VmId),
				err.Error(),
			)
			return
		}
	}

	for _, nic := range nics {
		// Get the Security Group IDs of the NIC without the current Security Group
		sgIds := lo.FilterMap(nic.SecurityGroups, func(sg osc.SecurityGroupLight, _ int) (string, bool) {
			return sg.SecurityGroupId, (sg.SecurityGroupId != data.Id.ValueString())
		})
		updateReq := osc.UpdateNicRequest{
			NicId:            nic.NicId,
			SecurityGroupIds: &sgIds,
		}

		_, err := r.Client.UpdateNic(ctx, updateReq, options.WithRetryTimeout(to))
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from NIC (%s)", data.Id.ValueString(), nic.NicId),
				err.Error(),
			)
			return
		}
	}

	for _, lbu := range lbus {
		if lbu.SecurityGroups == nil {
			continue
		}
		// Get the Security Group IDs of the LBU without the current Security Group
		sgIds := lo.FilterMap(lbu.SecurityGroups, func(sg string, _ int) (string, bool) {
			return sg, (sg != data.Id.ValueString())
		})
		updateReq := osc.UpdateLoadBalancerRequest{
			LoadBalancerName: lbu.LoadBalancerName,
			SecurityGroups:   &sgIds,
		}

		// Retry on 6031 has it can happen if the LBU is not in an available state
		err := oapihelpers.RetryOnCodes(ctx, []string{"6031"}, func() (resp any, err error) {
			return r.Client.UpdateLoadBalancer(ctx, updateReq, options.WithRetryTimeout(to))
		}, to)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from Load Balancer (%s)", data.Id.ValueString(), lbu.LoadBalancerName),
				err.Error(),
			)
			return
		}
	}

	return
}

func setSecurityGroupState(ctx context.Context, r *resourceSecurityGroup, data SecurityGroupModel) (SecurityGroupModel, error) {
	readReq := osc.ReadSecurityGroupsRequest{
		Filters: &osc.FiltersSecurityGroup{
			SecurityGroupIds: &[]string{data.Id.ValueString()},
		},
	}

	readTimeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diag.HasError() {
		return data, fmt.Errorf("unable to parse 'security_group' read timeout value: %v", diag.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	readResp, err := r.Client.ReadSecurityGroups(ctx, readReq, options.WithRetryTimeout(readTimeout))
	if err != nil {
		return data, err
	}
	data.RequestId = to.String(readResp.ResponseContext.RequestId)
	if readResp.SecurityGroups == nil || len(*readResp.SecurityGroups) == 0 {
		return data, ErrResourceEmpty
	}

	securityGroup := (*readResp.SecurityGroups)[0]
	tags, diags := flattenOAPITagsFW(ctx, securityGroup.Tags)
	if diags.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	inboundRulesModels, diag := flattenSecurityGroupRules(ctx, securityGroup.InboundRules)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert inbound rules to the model: %v", diag.Errors())
	}
	inboundRules, diag := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, inboundRulesModels)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert inbound rules to the schema list: %v", diag.Errors())
	}
	outboundRulesModels, diag := flattenSecurityGroupRules(ctx, securityGroup.OutboundRules)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert outbound rules to the model: %v", diag.Errors())
	}
	outboundRules, diag := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, outboundRulesModels)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert outbound rules to the schema list: %v", diag.Errors())
	}

	data.AccountId = to.String(securityGroup.AccountId)
	data.Description = to.String(securityGroup.Description)
	data.InboundRules = inboundRules
	data.NetId = to.String(ptr.From(securityGroup.NetId))
	data.OutboundRules = outboundRules
	data.SecurityGroupId = to.String(securityGroup.SecurityGroupId)
	data.Id = to.String(securityGroup.SecurityGroupId)
	data.SecurityGroupName = to.String(securityGroup.SecurityGroupName)

	return data, nil
}

func flattenSecurityGroupsMembers(sgMembers []osc.SecurityGroupsMember) []SecurityGroupsMembersModel {
	sgMembersModels := []SecurityGroupsMembersModel{}

	for _, sgMember := range sgMembers {
		member := SecurityGroupsMembersModel{
			AccountId:         to.String(sgMember.AccountId),
			SecurityGroupId:   to.String(sgMember.SecurityGroupId),
			SecurityGroupName: to.String(sgMember.SecurityGroupName),
		}
		sgMembersModels = append(sgMembersModels, member)
	}
	return sgMembersModels
}

func flattenSecurityGroupRules(ctx context.Context, sgRules []osc.SecurityGroupRule) ([]SecurityGroupRulesModel, diag.Diagnostics) {
	sgRulesModels := []SecurityGroupRulesModel{}
	diags := diag.Diagnostics{}

	for _, sgRule := range sgRules {
		ipRanges, diag := types.ListValueFrom(ctx, types.StringType, sgRule.IpRanges)
		diags.Append(diag...)
		serviceIds, diag := types.ListValueFrom(ctx, types.StringType, sgRule.ServiceIds)
		diags.Append(diag...)
		sgMembers, diag := types.ListValueFrom(ctx, securityGroupsMemberModelAttrTypes, flattenSecurityGroupsMembers(sgRule.SecurityGroupsMembers))
		diags.Append(diag...)

		rule := SecurityGroupRulesModel{
			FromPortRange:         types.Int32Value(int32(sgRule.FromPortRange)),
			IpProtocol:            to.String(sgRule.IpProtocol),
			IpRanges:              ipRanges,
			SecurityGroupsMembers: sgMembers,
			ServiceIds:            serviceIds,
			ToPortRange:           types.Int32Value(int32(sgRule.ToPortRange)),
		}
		sgRulesModels = append(sgRulesModels, rule)
	}
	if diags.HasError() {
		return nil, diags
	}

	return sgRulesModels, nil
}

func expandSecurityGroupRules(ctx context.Context, sgRulesModels []SecurityGroupRulesModel) ([]osc.SecurityGroupRule, diag.Diagnostics) {
	sgRules := []osc.SecurityGroupRule{}
	diags := diag.Diagnostics{}

	for _, sgRuleModel := range sgRulesModels {
		rule := osc.SecurityGroupRule{}

		if fwhelpers.IsSet(sgRuleModel.IpRanges) && len(sgRuleModel.IpRanges.Elements()) > 0 {
			ipRanges, diag := to.Slice[string](ctx, sgRuleModel.IpRanges)
			diags.Append(diag...)
			rule.IpRanges = ipRanges
		}
		if fwhelpers.IsSet(sgRuleModel.ServiceIds) && len(sgRuleModel.ServiceIds.Elements()) > 0 {
			serviceIds, diag := to.Slice[string](ctx, sgRuleModel.ServiceIds)
			diags.Append(diag...)
			rule.ServiceIds = serviceIds
		}
		if fwhelpers.IsSet(sgRuleModel.FromPortRange) {
			rule.FromPortRange = int(sgRuleModel.FromPortRange.ValueInt32())
		}
		if fwhelpers.IsSet(sgRuleModel.IpProtocol) {
			rule.IpProtocol = sgRuleModel.IpProtocol.ValueString()
		}
		if fwhelpers.IsSet(sgRuleModel.SecurityGroupsMembers) {
			sgMembers, diag := to.Slice[SecurityGroupsMembersModel](ctx, sgRuleModel.SecurityGroupsMembers)
			diags.Append(diag...)
			if diags.HasError() {
				return nil, diags
			}

			rule.SecurityGroupsMembers = expandSecurityGroupsMembers(sgMembers)
		}
		if fwhelpers.IsSet(sgRuleModel.ToPortRange) {
			rule.ToPortRange = int(sgRuleModel.ToPortRange.ValueInt32())
		}
		if diags.HasError() {
			return nil, diags
		}

		sgRules = append(sgRules, rule)
	}

	return sgRules, nil
}

func expandSecurityGroupsMembers(sgMembersModels []SecurityGroupsMembersModel) []osc.SecurityGroupsMember {
	sgMembers := []osc.SecurityGroupsMember{}

	for _, sgMemberModel := range sgMembersModels {
		member := osc.SecurityGroupsMember{}

		if !sgMemberModel.AccountId.IsUnknown() && !sgMemberModel.AccountId.IsNull() && sgMemberModel.AccountId.ValueString() != "" {
			member.AccountId = sgMemberModel.AccountId.ValueStringPointer()
		}
		if !sgMemberModel.SecurityGroupId.IsUnknown() && !sgMemberModel.SecurityGroupId.IsNull() && sgMemberModel.SecurityGroupId.ValueString() != "" {
			member.SecurityGroupId = sgMemberModel.SecurityGroupId.ValueString()
		}
		if !sgMemberModel.SecurityGroupName.IsUnknown() && !sgMemberModel.SecurityGroupName.IsNull() && sgMemberModel.SecurityGroupName.ValueString() != "" {
			member.SecurityGroupName = sgMemberModel.SecurityGroupName.ValueStringPointer()
		}

		sgMembers = append(sgMembers, member)
	}

	return sgMembers
}
