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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/modifyplans"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
	Client *oscgo.APIClient
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
			fmt.Sprintf("Expected *osc.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSCAPI
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

	data.SecurityGroupId = types.StringValue(securityGroupId)
	data.Id = types.StringValue(securityGroupId)
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
		if vm.SecurityGroups != nil && len(*vm.SecurityGroups) == 1 {
			resources = append(resources, *vm.VmId)
		}
	}
	for _, nic := range nics {
		if nic.SecurityGroups != nil && len(*nic.SecurityGroups) == 1 {
			resources = append(resources, *nic.NicId)
		}
	}
	targetSG := []string{stateData.Id.ValueString()}
	for _, lbu := range lbus {
		if lbu.SecurityGroups != nil && slices.Equal(*lbu.SecurityGroups, targetSG) {
			resources = append(resources, *lbu.LoadBalancerName)
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
		data.SecurityGroupName = types.StringValue(id.UniqueId())
	}

	createReq := oscgo.CreateSecurityGroupRequest{
		Description:       data.Description.ValueString(),
		SecurityGroupName: data.SecurityGroupName.ValueString(),
	}

	if !data.NetId.IsUnknown() && !data.NetId.IsNull() {
		createReq.SetNetId(data.NetId.ValueString())
	}

	var createResp oscgo.CreateSecurityGroupResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.SecurityGroupApi.CreateSecurityGroup(ctx).CreateSecurityGroupRequest(createReq).Execute()
		if err != nil {
			apiErr := oapihelpers.GetError(err)
			// Fail fast when the security group name already exists
			if apiErr.GetCode() == "9008" {
				errBody := utils.GetHttpErrorResponse(httpResp.Body, err)
				return retry.NonRetryableError(errBody)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Security Group.",
			err.Error(),
		)
		return
	}
	sg := createResp.GetSecurityGroup()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.Id = types.StringValue(sg.GetSecurityGroupId())
	data.SecurityGroupId = types.StringValue(sg.GetSecurityGroupId())

	if data.RemoveDefaultOutboundRule.ValueBool() {
		ipRange := "0.0.0.0/0"
		ipProtocol := "-1"
		emptySGReq := oscgo.DeleteSecurityGroupRuleRequest{
			Flow:            "Outbound",
			SecurityGroupId: sg.GetSecurityGroupId(),
			IpRange:         &ipRange,
			IpProtocol:      &ipProtocol,
		}

		err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.SecurityGroupRuleApi.DeleteSecurityGroupRule(ctx).DeleteSecurityGroupRuleRequest(emptySGReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to empty the Security Group rules.",
				err.Error(),
			)
			return
		}
	}

	diag = createOAPITagsFW(ctx, r.Client, data.Tags, sg.GetSecurityGroupId())
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

	diag := updateOAPITagsFW(ctx, r.Client, stateData.Tags, planData.Tags, stateData.SecurityGroupId.ValueString())
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
	delReq := oscgo.DeleteSecurityGroupRequest{
		SecurityGroupId: &sgId,
	}

	sgLinked := false
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.SecurityGroupApi.DeleteSecurityGroup(ctx).DeleteSecurityGroupRequest(delReq).Execute()
		if err != nil {
			oscErr := oapihelpers.GetError(err)
			if oscErr.GetCode() == "9085" {
				sgLinked = true
				return nil
			}
			return utils.CheckThrottling(httpResp, err)

		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Security Group.",
			err.Error(),
		)
		return
	}

	if sgLinked {
		diag = r.unlink(ctx, deleteTimeout, data)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
			_, httpResp, err := r.Client.SecurityGroupApi.DeleteSecurityGroup(ctx).DeleteSecurityGroupRequest(delReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to delete Security Group after unlinking from resources.",
				err.Error(),
			)
			return
		}
	}
}

func (r *resourceSecurityGroup) getLinkedResources(ctx context.Context, to time.Duration, sgId string) ([]oscgo.Vm, []oscgo.Nic, []oscgo.LoadBalancer, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get VMs linked to Security Group
	readVmReq := oscgo.ReadVmsRequest{
		Filters: &oscgo.FiltersVm{
			SecurityGroupIds: &[]string{sgId},
		},
	}

	var vms []oscgo.Vm
	err := retry.RetryContext(ctx, to, func() *retry.RetryError {
		resp, httpResp, err := r.Client.VmApi.ReadVms(ctx).ReadVmsRequest(readVmReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		vms = *resp.Vms
		return nil
	})
	if err != nil {
		diags.AddError(
			"Unable to find VMs linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	// Get NICs linked to Security Group
	readNicsReq := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{
			SecurityGroupIds: &[]string{sgId},
		},
	}

	var nics []oscgo.Nic
	err = retry.RetryContext(ctx, to, func() *retry.RetryError {
		resp, httpResp, err := r.Client.NicApi.ReadNics(ctx).ReadNicsRequest(readNicsReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		nics = *resp.Nics
		return nil
	})
	if err != nil {
		diags.AddError(
			"Unable to find NICs linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	// Get LBUs linked to Security Group
	var lbus []oscgo.LoadBalancer
	err = retry.RetryContext(ctx, to, func() *retry.RetryError {
		resp, httpResp, err := r.Client.LoadBalancerApi.ReadLoadBalancers(ctx).ReadLoadBalancersRequest(oscgo.ReadLoadBalancersRequest{}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		lbus = *resp.LoadBalancers
		return nil
	})
	if err != nil {
		diags.AddError(
			"Unable to find Load Balancers linked to Security Group.",
			err.Error(),
		)
		return nil, nil, nil, diags
	}

	return vms, nics, lbus, nil
}

func (r *resourceSecurityGroup) unlink(ctx context.Context, to time.Duration, data SecurityGroupModel) (diags diag.Diagnostics) {
	vms, nics, lbus, diag := r.getLinkedResources(ctx, to, data.Id.ValueString())
	if diag.HasError() {
		diags = append(diags, diag...)
		return diags
	}

	// Check if the Security Group is the unique Security Group of any VM
	vmsUniqueSG := lo.FilterMap(vms, func(vm oscgo.Vm, _ int) (string, bool) {
		return *vm.VmId, (vm.SecurityGroups != nil && len(*vm.SecurityGroups) == 1)
	})
	// Check if the Security Group is the unique Security Group of any NIC
	nicsUniqueSG := lo.FilterMap(nics, func(nic oscgo.Nic, _ int) (string, bool) {
		return *nic.NicId, (nic.SecurityGroups != nil && len(*nic.SecurityGroups) == 1)
	})
	targetSG := []string{data.Id.ValueString()}
	// Check if the Security Group is the unique Security Group of any LBU
	lbusUniqueSG := lo.FilterMap(lbus, func(lbu oscgo.LoadBalancer, _ int) (string, bool) {
		return *lbu.LoadBalancerName, (lbu.SecurityGroups != nil && slices.Equal(*lbu.SecurityGroups, targetSG))
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
		sgIds := lo.FilterMap(*vm.SecurityGroups, func(sg oscgo.SecurityGroupLight, _ int) (string, bool) {
			return *sg.SecurityGroupId, (*sg.SecurityGroupId != data.Id.ValueString())
		})
		updateReq := oscgo.UpdateVmRequest{
			VmId:             *vm.VmId,
			SecurityGroupIds: &sgIds,
		}

		err := retry.RetryContext(ctx, to, func() *retry.RetryError {
			_, httpResp, err := r.Client.VmApi.UpdateVm(ctx).UpdateVmRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from VM (%s)", data.Id.ValueString(), *vm.VmId),
				err.Error(),
			)
			return
		}
	}

	for _, nic := range nics {
		// Get the Security Group IDs of the NIC without the current Security Group
		sgIds := lo.FilterMap(*nic.SecurityGroups, func(sg oscgo.SecurityGroupLight, _ int) (string, bool) {
			return *sg.SecurityGroupId, (*sg.SecurityGroupId != data.Id.ValueString())
		})
		updateReq := oscgo.UpdateNicRequest{
			NicId:            *nic.NicId,
			SecurityGroupIds: &sgIds,
		}

		err := retry.RetryContext(ctx, to, func() *retry.RetryError {
			_, httpResp, err := r.Client.NicApi.UpdateNic(ctx).UpdateNicRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from NIC (%s)", data.Id.ValueString(), *nic.NicId),
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
		sgIds := lo.FilterMap(*lbu.SecurityGroups, func(sg string, _ int) (string, bool) {
			return sg, (sg != data.Id.ValueString())
		})
		updateReq := oscgo.UpdateLoadBalancerRequest{
			LoadBalancerName: *lbu.LoadBalancerName,
			SecurityGroups:   &sgIds,
		}

		err := retry.RetryContext(ctx, to, func() *retry.RetryError {
			_, httpResp, err := r.Client.LoadBalancerApi.UpdateLoadBalancer(ctx).UpdateLoadBalancerRequest(updateReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Unable to remove Security Group (%s) from Load Balancer (%s)", data.Id.ValueString(), *lbu.LoadBalancerName),
				err.Error(),
			)
			return
		}
	}

	return
}

func setSecurityGroupState(ctx context.Context, r *resourceSecurityGroup, data SecurityGroupModel) (SecurityGroupModel, error) {
	readReq := oscgo.ReadSecurityGroupsRequest{
		Filters: &oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &[]string{data.Id.ValueString()},
		},
	}

	readTimeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diag.HasError() {
		return data, fmt.Errorf("unable to parse 'security_group' read timeout value: %v", diag.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var readResp oscgo.ReadSecurityGroupsResponse
	err := retry.RetryContext(ctx, readTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.SecurityGroupApi.ReadSecurityGroups(ctx).ReadSecurityGroupsRequest(readReq).Execute()
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
	if len(readResp.GetSecurityGroups()) == 0 {
		return data, ErrResourceEmpty
	}

	securityGroup := readResp.GetSecurityGroups()[0]
	tags, diags := flattenOAPITagsFW(ctx, securityGroup.GetTags())
	if diags.HasError() {
		return data, fmt.Errorf("unable to flatten tags: %v", diags.Errors())
	}
	data.Tags = tags

	inboundRulesModels, diag := flattenSecurityGroupRules(ctx, securityGroup.GetInboundRules())
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert inbound rules to the model: %v", diag.Errors())
	}
	inboundRules, diag := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, inboundRulesModels)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert inbound rules to the schema list: %v", diag.Errors())
	}
	outboundRulesModels, diag := flattenSecurityGroupRules(ctx, securityGroup.GetOutboundRules())
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert outbound rules to the model: %v", diag.Errors())
	}
	outboundRules, diag := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, outboundRulesModels)
	if diag.HasError() {
		return data, fmt.Errorf("unable to convert outbound rules to the schema list: %v", diag.Errors())
	}

	data.AccountId = types.StringValue(securityGroup.GetAccountId())
	data.Description = types.StringValue(securityGroup.GetDescription())
	data.InboundRules = inboundRules
	data.NetId = types.StringValue(securityGroup.GetNetId())
	data.OutboundRules = outboundRules
	data.SecurityGroupId = types.StringValue(securityGroup.GetSecurityGroupId())
	data.Id = types.StringValue(securityGroup.GetSecurityGroupId())
	data.SecurityGroupName = types.StringValue(securityGroup.GetSecurityGroupName())

	return data, nil
}

func flattenSecurityGroupsMembers(sgMembers []oscgo.SecurityGroupsMember) []SecurityGroupsMembersModel {
	sgMembersModels := []SecurityGroupsMembersModel{}

	for _, sgMember := range sgMembers {
		member := SecurityGroupsMembersModel{
			AccountId:         types.StringValue(sgMember.GetAccountId()),
			SecurityGroupId:   types.StringValue(sgMember.GetSecurityGroupId()),
			SecurityGroupName: types.StringValue(sgMember.GetSecurityGroupName()),
		}
		sgMembersModels = append(sgMembersModels, member)
	}
	return sgMembersModels
}

func flattenSecurityGroupRules(ctx context.Context, sgRules []oscgo.SecurityGroupRule) ([]SecurityGroupRulesModel, diag.Diagnostics) {
	sgRulesModels := []SecurityGroupRulesModel{}
	diags := diag.Diagnostics{}

	for _, sgRule := range sgRules {
		ipRanges, diag := types.ListValueFrom(ctx, types.StringType, sgRule.GetIpRanges())
		diags.Append(diag...)
		serviceIds, diag := types.ListValueFrom(ctx, types.StringType, sgRule.GetServiceIds())
		diags.Append(diag...)
		sgMembers, diag := types.ListValueFrom(ctx, securityGroupsMemberModelAttrTypes, flattenSecurityGroupsMembers(sgRule.GetSecurityGroupsMembers()))
		diags.Append(diag...)

		rule := SecurityGroupRulesModel{
			FromPortRange:         types.Int32Value(sgRule.GetFromPortRange()),
			IpProtocol:            types.StringValue(sgRule.GetIpProtocol()),
			IpRanges:              ipRanges,
			SecurityGroupsMembers: sgMembers,
			ServiceIds:            serviceIds,
			ToPortRange:           types.Int32Value(sgRule.GetToPortRange()),
		}
		sgRulesModels = append(sgRulesModels, rule)
	}
	if diags.HasError() {
		return nil, diags
	}

	return sgRulesModels, nil
}

func expandSecurityGroupRules(ctx context.Context, sgRulesModels []SecurityGroupRulesModel) ([]oscgo.SecurityGroupRule, diag.Diagnostics) {
	sgRules := []oscgo.SecurityGroupRule{}
	diags := diag.Diagnostics{}

	for _, sgRuleModel := range sgRulesModels {
		rule := oscgo.SecurityGroupRule{}

		if fwhelpers.IsSet(sgRuleModel.IpRanges) && len(sgRuleModel.IpRanges.Elements()) > 0 {
			ipRanges, diag := to.Slice[string](ctx, sgRuleModel.IpRanges)
			diags.Append(diag...)
			rule.SetIpRanges(ipRanges)
		}
		if fwhelpers.IsSet(sgRuleModel.ServiceIds) && len(sgRuleModel.ServiceIds.Elements()) > 0 {
			serviceIds, diag := to.Slice[string](ctx, sgRuleModel.ServiceIds)
			diags.Append(diag...)
			rule.SetServiceIds(serviceIds)
		}
		if fwhelpers.IsSet(sgRuleModel.FromPortRange) {
			rule.SetFromPortRange(sgRuleModel.FromPortRange.ValueInt32())
		}
		if fwhelpers.IsSet(sgRuleModel.IpProtocol) {
			rule.SetIpProtocol(sgRuleModel.IpProtocol.ValueString())
		}
		if fwhelpers.IsSet(sgRuleModel.SecurityGroupsMembers) {
			sgMembers, diag := to.Slice[SecurityGroupsMembersModel](ctx, sgRuleModel.SecurityGroupsMembers)
			diags.Append(diag...)
			if diags.HasError() {
				return nil, diags
			}

			rule.SetSecurityGroupsMembers(expandSecurityGroupsMembers(sgMembers))
		}
		if fwhelpers.IsSet(sgRuleModel.ToPortRange) {
			rule.SetToPortRange(sgRuleModel.ToPortRange.ValueInt32())
		}
		if diags.HasError() {
			return nil, diags
		}

		sgRules = append(sgRules, rule)
	}

	return sgRules, nil
}

func expandSecurityGroupsMembers(sgMembersModels []SecurityGroupsMembersModel) []oscgo.SecurityGroupsMember {
	sgMembers := []oscgo.SecurityGroupsMember{}

	for _, sgMemberModel := range sgMembersModels {
		member := oscgo.SecurityGroupsMember{}

		if !sgMemberModel.AccountId.IsUnknown() && !sgMemberModel.AccountId.IsNull() && sgMemberModel.AccountId.ValueString() != "" {
			member.SetAccountId(sgMemberModel.AccountId.ValueString())
		}
		if !sgMemberModel.SecurityGroupId.IsUnknown() && !sgMemberModel.SecurityGroupId.IsNull() && sgMemberModel.SecurityGroupId.ValueString() != "" {
			member.SetSecurityGroupId(sgMemberModel.SecurityGroupId.ValueString())
		}
		if !sgMemberModel.SecurityGroupName.IsUnknown() && !sgMemberModel.SecurityGroupName.IsNull() && sgMemberModel.SecurityGroupName.ValueString() != "" {
			member.SetSecurityGroupName(sgMemberModel.SecurityGroupName.ValueString())
		}

		sgMembers = append(sgMembers, member)
	}

	return sgMembers
}
