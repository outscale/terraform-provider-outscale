package outscale

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/fwmodifyplan"
	"github.com/outscale/terraform-provider-outscale/utils"
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
	Tags                      []ResourceTag  `tfsdk:"tags"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
	RequestId                 types.String   `tfsdk:"request_id"`
	Id                        types.String   `tfsdk:"id"`
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

	for _, sgRule := range sgRules {
		ipRanges, diags := types.SetValueFrom(ctx, types.StringType, sgRule.GetIpRanges())
		if diags.HasError() {
			return sgRulesModels, diags
		}
		serviceIds, diags := types.SetValueFrom(ctx, types.StringType, sgRule.GetServiceIds())
		if diags.HasError() {
			return sgRulesModels, diags
		}
		sgMembers, diags := types.ListValueFrom(ctx, securityGroupsMemberModelAttrTypes, flattenSecurityGroupsMembers(sgRule.GetSecurityGroupsMembers()))
		if diags.HasError() {
			return sgRulesModels, diags
		}
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
	return sgRulesModels, nil
}

func expandSecurityGroupRules(ctx context.Context, sgRulesModels []SecurityGroupRulesModel) ([]oscgo.SecurityGroupRule, diag.Diagnostics) {
	sgRules := []oscgo.SecurityGroupRule{}

	for _, sgRuleModel := range sgRulesModels {
		var ipRanges, serviceIds []string
		diags := sgRuleModel.IpRanges.ElementsAs(ctx, &ipRanges, false)
		if diags.HasError() {
			return sgRules, diags
		}
		diags = sgRuleModel.ServiceIds.ElementsAs(ctx, &serviceIds, false)
		if diags.HasError() {
			return sgRules, diags
		}
		var sgMembers []SecurityGroupsMembersModel
		diags = sgRuleModel.SecurityGroupsMembers.ElementsAs(ctx, &sgMembers, false)
		if diags.HasError() {
			return sgRules, diags
		}

		var rule oscgo.SecurityGroupRule
		if len(ipRanges) > 0 {
			rule.SetIpRanges(ipRanges)
		}
		if len(serviceIds) > 0 {
			rule.SetServiceIds(serviceIds)
		}
		if !sgRuleModel.FromPortRange.IsUnknown() && !sgRuleModel.FromPortRange.IsNull() {
			rule.SetFromPortRange(sgRuleModel.FromPortRange.ValueInt32())
		}
		if !sgRuleModel.IpProtocol.IsUnknown() && !sgRuleModel.IpProtocol.IsNull() {
			rule.SetIpProtocol(sgRuleModel.IpProtocol.ValueString())
		}
		if !sgRuleModel.SecurityGroupsMembers.IsUnknown() && !sgRuleModel.SecurityGroupsMembers.IsNull() {
			rule.SetSecurityGroupsMembers(expandSecurityGroupsMembers(sgMembers))
		}
		if !sgRuleModel.ToPortRange.IsUnknown() && !sgRuleModel.ToPortRange.IsNull() {
			rule.SetToPortRange(sgRuleModel.ToPortRange.ValueInt32())
		}

		sgRules = append(sgRules, rule)
	}

	return sgRules, nil
}

func expandSecurityGroupsMembers(sgMembersModels []SecurityGroupsMembersModel) []oscgo.SecurityGroupsMember {
	sgMembers := []oscgo.SecurityGroupsMember{}

	for _, sgMemberModel := range sgMembersModels {
		var member oscgo.SecurityGroupsMember
		if !sgMemberModel.AccountId.IsUnknown() && !sgMemberModel.AccountId.IsNull() {
			member.SetAccountId(sgMemberModel.AccountId.ValueString())
		}
		if !sgMemberModel.SecurityGroupId.IsUnknown() && !sgMemberModel.SecurityGroupId.IsNull() {
			member.SetSecurityGroupId(sgMemberModel.SecurityGroupId.ValueString())
		}
		if !sgMemberModel.SecurityGroupName.IsUnknown() && !sgMemberModel.SecurityGroupName.IsNull() {
			member.SetSecurityGroupName(sgMemberModel.SecurityGroupName.ValueString())
		}

		sgMembers = append(sgMembers, member)
	}

	return sgMembers
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
	client, ok := req.ProviderData.(OutscaleClient_fw)
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
	data.SecurityGroupId = types.StringValue(securityGroupId)
	data.Id = types.StringValue(securityGroupId)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeouts
	data.InboundRules = types.ListNull(securityGroupRulesModelAttrTypes)
	data.OutboundRules = types.ListNull(securityGroupRulesModelAttrTypes)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSecurityGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group"
}

func (r *resourceSecurityGroup) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
}

func (r *resourceSecurityGroup) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"inbound_rules": SecurityGroupIntOutRulesSchema(),
			"net_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					fwmodifyplan.ForceNewFramework(),
				},
			},
			"outbound_rules": SecurityGroupIntOutRulesSchema(),
			"security_group_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"remove_default_outbound_rule": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(utils.RemoveDefaultOutboundRuleDefaultValue),
				PlanModifiers: []planmodifier.Bool{
					// TODO: ForceNewFramework ?
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
			},
		},
	}
}

func SecurityGroupIntOutRulesSchema() *schema.ListAttribute {
	return &schema.ListAttribute{
		Computed:    true,
		ElementType: securityGroupRulesModelAttrTypes,
	}
}

func (r *resourceSecurityGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecurityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.CreateDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

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

		deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
		defer cancel()

		err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
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

	if len(data.Tags) > 0 {
		err = createFrameworkTags(ctx, r.Client, tagsToOSCResourceTag(data.Tags), sg.GetSecurityGroupId())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add Tags on the Security Group resource.",
				err.Error(),
			)
			return
		}
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSecurityGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SecurityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := setSecurityGroupState(ctx, r, data)
	if err != nil {
		if err.Error() == "Empty" {
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSecurityGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData SecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !reflect.DeepEqual(planData.Tags, stateData.Tags) {
		toRemove, toCreate := diffOSCAPITags(tagsToOSCResourceTag(planData.Tags), tagsToOSCResourceTag(stateData.Tags))
		err := updateFrameworkTags(ctx, r.Client, toCreate, toRemove, stateData.SecurityGroupId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update Tags on Security Group resource.",
				err.Error(),
			)
			return
		}
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
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSecurityGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SecurityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DeleteDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	sgId := data.Id.ValueString()
	delReq := oscgo.DeleteSecurityGroupRequest{
		SecurityGroupId: &sgId,
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
			"Unable to delete Security Group.",
			err.Error(),
		)
		return
	}
}

func setSecurityGroupState(ctx context.Context, r *resourceSecurityGroup, data SecurityGroupModel) (SecurityGroupModel, error) {
	readReq := oscgo.ReadSecurityGroupsRequest{
		Filters: &oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &[]string{data.Id.ValueString()},
		},
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.ReadDefaultTimeout)
	if diags.HasError() {
		return data, fmt.Errorf("unable to parse 'security_group' read timeout value. Error: %v: ", diags.Errors())
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
		return data, errors.New("Empty")
	}

	securityGroup := readResp.GetSecurityGroups()[0]
	data.Tags = getTagsFromApiResponse(securityGroup.GetTags())

	inboundRulesModels, diags := flattenSecurityGroupRules(ctx, securityGroup.GetInboundRules())
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Inbound Rules to the model. Error: %v: ", diags.Errors())
	}
	inboundRules, diags := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, inboundRulesModels)
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Inbound Rules to the schema List. Error: %v: ", diags.Errors())
	}
	outboundRulesModels, diags := flattenSecurityGroupRules(ctx, securityGroup.GetOutboundRules())
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Outbound Rules to the model. Error: %v: ", diags.Errors())
	}
	outboundRules, diags := types.ListValueFrom(ctx, securityGroupRulesModelAttrTypes, outboundRulesModels)
	if diags.HasError() {
		return data, fmt.Errorf("Unable to convert Outbound Rules to the schema List. Error: %v: ", diags.Errors())
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
