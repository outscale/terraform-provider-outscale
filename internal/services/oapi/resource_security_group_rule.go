package oapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorint32"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorlist"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

var (
	_ resource.Resource                   = &resourceSecurityGroupRule{}
	_ resource.ResourceWithConfigure      = &resourceSecurityGroupRule{}
	_ resource.ResourceWithImportState    = &resourceSecurityGroupRule{}
	_ resource.ResourceWithModifyPlan     = &resourceSecurityGroupRule{}
	_ resource.ResourceWithValidateConfig = &resourceRoute{}
)

type SecurityGroupRuleModel struct {
	Flow                         types.String   `tfsdk:"flow"`
	FromPortRange                types.Int32    `tfsdk:"from_port_range"`
	IpProtocol                   types.String   `tfsdk:"ip_protocol"`
	IpRange                      types.String   `tfsdk:"ip_range"`
	Rules                        types.List     `tfsdk:"rules"`
	SecurityGroupAccountIdToLink types.String   `tfsdk:"security_group_account_id_to_link"`
	SecurityGroupId              types.String   `tfsdk:"security_group_id"`
	SecurityGroupNameToLink      types.String   `tfsdk:"security_group_name_to_link"`
	ToPortRange                  types.Int32    `tfsdk:"to_port_range"`
	SecurityGroupName            types.String   `tfsdk:"security_group_name"`
	NetId                        types.String   `tfsdk:"net_id"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
	RequestId                    types.String   `tfsdk:"request_id"`
	Id                           types.String   `tfsdk:"id"`
}

type SecurityGroupRulesModel struct {
	FromPortRange         types.Int32  `tfsdk:"from_port_range"`
	IpProtocol            types.String `tfsdk:"ip_protocol"`
	IpRanges              types.List   `tfsdk:"ip_ranges"`
	SecurityGroupsMembers types.List   `tfsdk:"security_groups_members"`
	ServiceIds            types.List   `tfsdk:"service_ids"`
	ToPortRange           types.Int32  `tfsdk:"to_port_range"`
}

type SecurityGroupsMembersModel struct {
	AccountId         types.String `tfsdk:"account_id"`
	SecurityGroupId   types.String `tfsdk:"security_group_id"`
	SecurityGroupName types.String `tfsdk:"security_group_name"`
}

var securityGroupsMemberModelAttrTypes = types.ObjectType{AttrTypes: fwhelpers.GetAttrTypes(SecurityGroupsMembersModel{})}

var securityGroupRulesModelAttrTypes = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"from_port_range":         types.Int32Type,
		"ip_protocol":             types.StringType,
		"ip_ranges":               types.ListType{ElemType: types.StringType},
		"security_groups_members": types.ListType{ElemType: securityGroupsMemberModelAttrTypes},
		"service_ids":             types.ListType{ElemType: types.StringType},
		"to_port_range":           types.Int32Type,
	},
}

type resourceSecurityGroupRule struct {
	Client *oscgo.APIClient
}

func NewResourceSecurityGroupRule() resource.Resource {
	return &resourceSecurityGroupRule{}
}

func (r *resourceSecurityGroupRule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceSecurityGroupRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	parts := strings.SplitN(req.ID, "_", 6)
	if len(parts) != 6 || req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import Security Group Rule identifier in the format {security_group_id}_{flow}_{ip_protocol}_{from_port}_{to_port}_{ip_range}, got: %v", req.ID),
		)
		return
	}
	securityGroupId := parts[0]
	flow := parts[1]
	ipProtocol := parts[2]
	fromPort := parts[3]
	toPort := parts[4]
	ipRange := parts[5]

	filters := &oscgo.FiltersSecurityGroup{
		SecurityGroupIds: &[]string{securityGroupId},
	}

	if strings.EqualFold(flow, "Inbound") {
		filters.InboundRuleProtocols = &[]string{ipProtocol}
		filters.InboundRuleFromPortRanges = &[]int32{cast.ToInt32(fromPort)}
		filters.InboundRuleToPortRanges = &[]int32{cast.ToInt32(toPort)}
		filters.InboundRuleIpRanges = &[]string{ipRange}
	} else if strings.EqualFold(flow, "Outbound") {
		filters.OutboundRuleProtocols = &[]string{ipProtocol}
		filters.OutboundRuleFromPortRanges = &[]int32{cast.ToInt32(fromPort)}
		filters.OutboundRuleToPortRanges = &[]int32{cast.ToInt32(toPort)}
		filters.OutboundRuleIpRanges = &[]string{ipRange}
	}

	var data SecurityGroupRuleModel
	readResp, err := r.readSecurityGroupsWithFilters(ctx, data, filters)
	if err != nil || len(readResp.GetSecurityGroups()) != 1 {
		resp.Diagnostics.AddError(
			"Unable to find the Security Group Rule with the requested attributes",
			fmt.Sprintf("Expected import Security Group Rule identifier in the format {security_group_id}_{flow}_{ip_protocol}_{from_port}_{to_port}_{ip_range}, got: %v", req.ID),
		)
		return
	}
	securityGroup := readResp.GetSecurityGroups()[0]

	var timeouts timeouts.Value
	diag := resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeouts)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	data.Timeouts = timeouts
	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	data.Id = types.StringValue(securityGroup.GetSecurityGroupId())
	data.SecurityGroupId = types.StringValue(securityGroup.GetSecurityGroupId())
	data.Flow = types.StringValue(strings.ToUpper(flow[:1]) + strings.ToLower(flow[1:]))
	data.IpProtocol = types.StringValue(ipProtocol)
	data.FromPortRange = types.Int32Value(cast.ToInt32(fromPort))
	data.ToPortRange = types.Int32Value(cast.ToInt32(toPort))
	data.IpRange = types.StringValue(ipRange)
	data.Rules = types.ListNull(securityGroupRulesModelAttrTypes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSecurityGroupRule) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group_rule"
}

func (r *resourceSecurityGroupRule) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() && req.Plan.Raw.IsFullyKnown() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will fully destroy this resource.",
		)
	}
	if req.State.Raw.IsNull() {
		check := func(ctx context.Context, config tfsdk.Config, diags *diag.Diagnostics, parentPath path.Path) {
			var protocol types.String
			var fromPort, toPort types.Int32
			protocolPath := parentPath.AtName("ip_protocol")
			fromPath := parentPath.AtName("from_port_range")
			toPath := parentPath.AtName("to_port_range")

			diags.Append(config.GetAttribute(ctx, protocolPath, &protocol)...)
			diags.Append(config.GetAttribute(ctx, fromPath, &fromPort)...)
			diags.Append(config.GetAttribute(ctx, toPath, &toPort)...)
			if diags.HasError() {
				return
			}

			if fwhelpers.IsSet(protocol) && protocol.ValueString() == "-1" &&
				(fwhelpers.IsSet(fromPort) || fwhelpers.IsSet(toPort)) {
				diags.AddWarning(
					"IP Protocol Configuration",
					"Attribute `ip_protocol` set to -1 targets all protocols, `from_port_range` and `to_port_range` attributes will be ignored. \nDocumentation: https://docs.outscale.com/en/userguide/About-Security-Group-Rules.html",
				)
			}
		}
		check(ctx, req.Config, &resp.Diagnostics, path.Empty())

		var rules types.List
		diag := req.Config.GetAttribute(ctx, path.Root("rules"), &rules)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
		for i := range rules.Elements() {
			check(ctx, req.Config, &resp.Diagnostics, path.Root("rules").AtListIndex(i))
		}
	}
}

func (r *resourceSecurityGroupRule) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	configAttrCount := 0
	attrs := []string{
		"rules",
		"ip_range",
		"from_port_range",
		"to_port_range",
		"ip_protocol",
		"security_group_account_id_to_link",
		"security_group_name_to_link",
	}

	for _, attrName := range attrs {
		var val attr.Value
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(attrName), &val)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !val.IsNull() {
			configAttrCount++
		}
	}
	if configAttrCount < 1 {
		resp.Diagnostics.AddError(
			"Attribute Configuration",
			fmt.Sprintf("At least one of %v should be set.", attrs),
		)
	}
}

func ifAllProtocolInt32(ctx context.Context, req validator.Int32Request) bool {
	var protocol types.String
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("ip_protocol"), &protocol)
	if diags.HasError() {
		return false
	}
	return protocol.IsNull() || protocol.IsUnknown() || protocol.ValueString() != "-1"
}

func ifAllProtocolList(ctx context.Context, req validator.ListRequest) bool {
	var protocol types.String
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("ip_protocol"), &protocol)
	if diags.HasError() {
		return false
	}
	return protocol.IsNull() || protocol.IsUnknown() || protocol.ValueString() != "-1"
}

func (r *resourceSecurityGroupRule) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"rules": schema.ListNestedBlock{
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"security_groups_members": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									// Setting a default value is necessary to replicate the behaviour of nested strings of the SDKv2
									// The Framework handle them as null value, why leads to the resource recreation if not using default
									"account_id": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Default:  stringdefault.StaticString(""),
									},
									"security_group_id": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Default:  stringdefault.StaticString(""),
									},
									"security_group_name": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Default:  stringdefault.StaticString(""),
									},
								},
							},
						},
					},
					Attributes: map[string]schema.Attribute{
						"from_port_range": schema.Int32Attribute{
							Optional: true,
							Validators: []validator.Int32{
								validatorint32.AlsoRequiresIf(path.MatchRelative().AtParent().AtName("to_port_range"), ifAllProtocolInt32),
							},
						},
						"ip_protocol": schema.StringAttribute{
							Optional: true,
						},
						"ip_ranges": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Optional:    true,
							Validators: []validator.List{
								validatorlist.AlsoRequiresIf(path.MatchRelative().AtParent().AtName("to_port_range"), ifAllProtocolList),
								validatorlist.AlsoRequiresIf(path.MatchRelative().AtParent().AtName("from_port_range"), ifAllProtocolList),
							},
							Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
						},

						"service_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Optional:    true,
							Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
						},
						"to_port_range": schema.Int32Attribute{
							Optional: true,
							Validators: []validator.Int32{
								validatorint32.AlsoRequiresIf(path.MatchRelative().AtParent().AtName("from_port_range"), ifAllProtocolInt32),
							},
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"flow": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Inbound", "Outbound",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"from_port_range": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("security_group_account_id_to_link"),
						path.MatchRoot("security_group_name_to_link"),
					),
					validatorint32.AlsoRequiresIf(path.MatchRoot("to_port_range"), ifAllProtocolInt32),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"ip_protocol": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("security_group_name_to_link"),
						path.MatchRoot("security_group_account_id_to_link"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_range": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("security_group_name_to_link"),
						path.MatchRoot("security_group_account_id_to_link"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_account_id_to_link": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("ip_protocol"),
						path.MatchRoot("ip_range"),
						path.MatchRoot("from_port_range"),
						path.MatchRoot("to_port_range"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_name_to_link": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("ip_protocol"),
						path.MatchRoot("ip_range"),
						path.MatchRoot("from_port_range"),
						path.MatchRoot("to_port_range"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"to_port_range": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.ConflictsWith(
						path.MatchRoot("rules"),
						path.MatchRoot("security_group_account_id_to_link"),
						path.MatchRoot("security_group_name_to_link"),
					),
					validatorint32.AlsoRequiresIf(path.MatchRoot("from_port_range"), ifAllProtocolInt32),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"security_group_name": schema.StringAttribute{
				Computed: true,
			},
			"net_id": schema.StringAttribute{
				Computed: true,
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

func (r *resourceSecurityGroupRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecurityGroupRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	createReq := oscgo.CreateSecurityGroupRuleRequest{
		Flow:            data.Flow.ValueString(),
		SecurityGroupId: data.SecurityGroupId.ValueString(),
	}

	if fwhelpers.IsSet(data.FromPortRange) {
		createReq.SetFromPortRange(data.FromPortRange.ValueInt32())
	}
	if fwhelpers.IsSet(data.IpProtocol) {
		createReq.SetIpProtocol(data.IpProtocol.ValueString())
	}
	if fwhelpers.IsSet(data.IpRange) {
		createReq.SetIpRange(data.IpRange.ValueString())
	}
	if fwhelpers.IsSet(data.Rules) {
		model, diag := to.Slice[SecurityGroupRulesModel](ctx, data.Rules)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		rules, diag := expandSecurityGroupRules(ctx, model)
		resp.Diagnostics.Append(diag...)
		createReq.SetRules(rules)
	}
	if !data.SecurityGroupAccountIdToLink.IsUnknown() && !data.SecurityGroupAccountIdToLink.IsNull() {
		createReq.SetSecurityGroupAccountIdToLink(data.SecurityGroupAccountIdToLink.ValueString())
	}
	if !data.SecurityGroupNameToLink.IsUnknown() && !data.SecurityGroupNameToLink.IsNull() {
		createReq.SetSecurityGroupNameToLink(data.SecurityGroupNameToLink.ValueString())
	}
	if !data.ToPortRange.IsUnknown() && !data.ToPortRange.IsNull() {
		createReq.SetToPortRange(data.ToPortRange.ValueInt32())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var createResp oscgo.CreateSecurityGroupRuleResponse
	err := retry.RetryContext(ctx, createTimeout, func() *retry.RetryError {
		rp, httpResp, err := r.Client.SecurityGroupRuleApi.CreateSecurityGroupRule(ctx).CreateSecurityGroupRuleRequest(createReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		createResp = rp
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Security Group Rule.",
			err.Error(),
		)
		return
	}
	sg := createResp.GetSecurityGroup()
	data.RequestId = types.StringValue(createResp.ResponseContext.GetRequestId())
	data.Id = types.StringValue(sg.GetSecurityGroupId())

	stateData, err := setSecurityGroupRuleState(ctx, r, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Security Group Rule state.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourceSecurityGroupRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SecurityGroupRuleModel

	diag := req.State.Get(ctx, &data)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data, err := setSecurityGroupRuleState(ctx, r, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Security Group Rule API response values.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceSecurityGroupRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceSecurityGroupRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SecurityGroupRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	deleteTimeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	delReq := oscgo.DeleteSecurityGroupRuleRequest{
		Flow:            data.Flow.ValueString(),
		SecurityGroupId: data.SecurityGroupId.ValueString(),
	}
	if !data.FromPortRange.IsUnknown() && !data.FromPortRange.IsNull() {
		delReq.SetFromPortRange(data.FromPortRange.ValueInt32())
	}
	if !data.IpProtocol.IsUnknown() && !data.IpProtocol.IsNull() {
		delReq.SetIpProtocol(data.IpProtocol.ValueString())
	}
	if !data.IpRange.IsUnknown() && !data.IpRange.IsNull() {
		delReq.SetIpRange(data.IpRange.ValueString())
	}
	if fwhelpers.IsSet(data.Rules) {
		model, diag := to.Slice[SecurityGroupRulesModel](ctx, data.Rules)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}

		rules, diag := expandSecurityGroupRules(ctx, model)
		resp.Diagnostics.Append(diag...)
		delReq.SetRules(rules)
	}
	if !data.SecurityGroupAccountIdToLink.IsUnknown() && !data.SecurityGroupAccountIdToLink.IsNull() {
		delReq.SetSecurityGroupAccountIdToUnlink(data.SecurityGroupAccountIdToLink.ValueString())
	}
	if !data.SecurityGroupNameToLink.IsUnknown() && !data.SecurityGroupNameToLink.IsNull() {
		delReq.SetSecurityGroupNameToUnlink(data.SecurityGroupNameToLink.ValueString())
	}
	if !data.ToPortRange.IsUnknown() && !data.ToPortRange.IsNull() {
		delReq.SetToPortRange(data.ToPortRange.ValueInt32())
	}

	if resp.Diagnostics.HasError() {
		return
	}
	err := retry.RetryContext(ctx, deleteTimeout, func() *retry.RetryError {
		_, httpResp, err := r.Client.SecurityGroupRuleApi.DeleteSecurityGroupRule(ctx).DeleteSecurityGroupRuleRequest(delReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Security Group Rule.",
			err.Error(),
		)
		return
	}
}

func (r *resourceSecurityGroupRule) readSecurityGroupsWithFilters(ctx context.Context, data SecurityGroupRuleModel, filter *oscgo.FiltersSecurityGroup) (*oscgo.ReadSecurityGroupsResponse, error) {
	readTimeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if diag.HasError() {
		return nil, fmt.Errorf("unable to parse 'security_group_rule' read timeout value: %v", diag.Errors())
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	readReq := oscgo.ReadSecurityGroupsRequest{
		Filters: filter,
	}
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
		return nil, err
	}

	return &readResp, nil
}

func setSecurityGroupRuleState(ctx context.Context, r *resourceSecurityGroupRule, data SecurityGroupRuleModel) (SecurityGroupRuleModel, error) {
	filters := &oscgo.FiltersSecurityGroup{SecurityGroupIds: &[]string{data.Id.ValueString()}}
	readResp, err := r.readSecurityGroupsWithFilters(ctx, data, filters)
	if err != nil {
		return data, err
	}
	data.RequestId = types.StringValue(readResp.ResponseContext.GetRequestId())
	if len(readResp.GetSecurityGroups()) == 0 {
		return data, ErrResourceEmpty
	}
	securityGroup := readResp.GetSecurityGroups()[0]

	data.SecurityGroupName = types.StringValue(securityGroup.GetSecurityGroupName())
	data.NetId = types.StringValue(securityGroup.GetNetId())

	return data, nil
}
