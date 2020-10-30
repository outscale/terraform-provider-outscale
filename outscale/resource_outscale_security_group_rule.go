package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openlyinc/pointy"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceOutscaleOAPIOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIOutboundRuleCreate,
		Read:   resourceOutscaleOAPIOutboundRuleRead,
		Delete: resourceOutscaleOAPIOutboundRuleDelete,
		Schema: map[string]*schema.Schema{
			"flow": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Inbound", "Outbound"}, false),
			},
			"from_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"rules", "security_group_name_to_link"},
			},
			"ip_range": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rules": getRulesSchema(false),
			"security_group_account_id_to_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_name_to_link": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_protocol", "rules"},
			},
			"to_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inbound_rules":  getRulesSchema(true),
			"outbound_rules": getRulesSchema(true),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateSecurityGroupRuleRequest{
		Flow:            d.Get("flow").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
		Rules:           expandRules(d),
	}

	if v, ok := d.GetOkExists("from_port_range"); ok {
		req.FromPortRange = pointy.Int32(cast.ToInt32(v))
	}
	if v, ok := d.GetOkExists("to_port_range"); ok {
		req.ToPortRange = pointy.Int32(cast.ToInt32(v))
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		req.IpProtocol = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("ip_range"); ok {
		req.IpRange = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("security_group_account_id_to_link"); ok {
		req.SecurityGroupAccountIdToLink = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("security_group_name_to_link"); ok {
		req.SecurityGroupNameToLink = pointy.String(v.(string))
	}

	var err error
	var resp oscgo.CreateSecurityGroupRuleResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SecurityGroupRuleApi.CreateSecurityGroupRule(context.Background()).CreateSecurityGroupRuleRequest(req).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf(
			"Error authorizing security group rule type: %s", utils.GetErrorResponse(err))
	}

	d.SetId(*resp.GetSecurityGroup().SecurityGroupId)

	return resourceOutscaleOAPIOutboundRuleRead(d, meta)
}

func resourceOutscaleOAPIOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	sg, resp, err := readSecurityGroups(conn, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("inbound_rules", flattenRules(sg.GetInboundRules())); err != nil {
		return fmt.Errorf("error setting `inbound_rules` for Outscale Security Group Rule(%s): %s", d.Id(), err)
	}
	if err := d.Set("outbound_rules", flattenRules(sg.GetOutboundRules())); err != nil {
		return fmt.Errorf("error setting `outbound_rules` for Outscale Security Group Rule(%s): %s", d.Id(), err)
	}
	if err := d.Set("security_group_name", sg.GetSecurityGroupName()); err != nil {
		return fmt.Errorf("error setting `security_group_name` for Outscale Security Group Rule(%s): %s", d.Id(), err)
	}
	if err := d.Set("net_id", sg.GetNetId()); err != nil {
		return fmt.Errorf("error setting `net_id` for Outscale Security Group Rule(%s): %s", d.Id(), err)
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return fmt.Errorf("error setting `request_id` for Outscale Security Group Rule(%s): %s", d.Id(), err)
	}
	return nil
}

func resourceOutscaleOAPIOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteSecurityGroupRuleRequest{
		Flow:            d.Get("flow").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
		Rules:           expandRules(d),
	}

	if v, ok := d.GetOkExists("from_port_range"); ok {
		req.FromPortRange = pointy.Int32(cast.ToInt32(v))
	}
	if v, ok := d.GetOkExists("to_port_range"); ok {
		req.ToPortRange = pointy.Int32(cast.ToInt32(v))
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		req.IpProtocol = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("ip_range"); ok {
		req.IpRange = pointy.String(v.(string))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err := conn.SecurityGroupRuleApi.DeleteSecurityGroupRule(context.Background()).DeleteSecurityGroupRuleRequest(req).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error revoking security group %s rules: %s", d.Id(), err)
	}

	return nil
}

func expandRules(d *schema.ResourceData) *[]oscgo.SecurityGroupRule {
	if len(d.Get("rules").([]interface{})) > 0 {
		rules := make([]oscgo.SecurityGroupRule, len(d.Get("rules").([]interface{})))

		for i, rule := range d.Get("rules").([]interface{}) {
			r := rule.(map[string]interface{})

			rules[i] = oscgo.SecurityGroupRule{
				SecurityGroupsMembers: expandSecurityGroupsMembers(r["security_groups_members"].([]interface{})),
			}

			if ipRanges := expandStringValueListPointer(r["ip_ranges"].([]interface{})); len(*ipRanges) > 0 {
				rules[i].IpRanges = expandStringValueListPointer(r["ip_ranges"].([]interface{}))
			}
			if serviceIDs := expandStringValueListPointer(r["service_ids"].([]interface{})); len(*serviceIDs) > 0 {
				rules[i].ServiceIds = expandStringValueListPointer(r["service_ids"].([]interface{}))
			}
			if v, ok := r["from_port_range"]; ok {
				rules[i].FromPortRange = pointy.Int32(cast.ToInt32(v))
			}
			if v, ok := r["ip_protocol"]; ok && v != "" {
				rules[i].IpProtocol = pointy.String(cast.ToString(v))
			}
			if v, ok := r["to_port_range"]; ok {
				rules[i].ToPortRange = pointy.Int32(cast.ToInt32(v))
			}
		}
		return &rules
	}
	return nil
}

func flattenRules(securityGroupsRules []oscgo.SecurityGroupRule) []map[string]interface{} {
	sgrs := make([]map[string]interface{}, len(securityGroupsRules))

	for i, s := range securityGroupsRules {
		sgrs[i] = map[string]interface{}{
			"from_port_range":         s.GetFromPortRange(),
			"ip_protocol":             s.GetIpProtocol(),
			"ip_ranges":               s.GetIpRanges(),
			"service_ids":             s.GetServiceIds(),
			"to_port_range":           s.GetToPortRange(),
			"security_groups_members": flattenSecurityGroupsMembers(s.GetSecurityGroupsMembers()),
		}
	}
	return sgrs
}

func flattenSecurityGroupsMembers(securityGroupMembers []oscgo.SecurityGroupsMember) []map[string]interface{} {
	sgms := make([]map[string]interface{}, len(securityGroupMembers))

	for i, s := range securityGroupMembers {
		sgms[i] = map[string]interface{}{
			"account_id":          s.GetAccountId(),
			"security_group_id":   s.GetSecurityGroupId(),
			"security_group_name": s.GetSecurityGroupName(),
		}
	}
	return sgms
}

func expandSecurityGroupsMembers(gps []interface{}) *[]oscgo.SecurityGroupsMember {
	groups := make([]oscgo.SecurityGroupsMember, len(gps))

	for i, group := range gps {
		g := group.(map[string]interface{})
		groups[i] = oscgo.SecurityGroupsMember{}

		if v, ok := g["account_id"]; ok && v != "" {
			groups[i].AccountId = pointy.String(cast.ToString(v))
		}
		if v, ok := g["security_group_id"]; ok && v != "" {
			groups[i].SecurityGroupId = pointy.String(cast.ToString(v))
		}
		if v, ok := g["security_group_name"]; ok && v != "" {
			groups[i].SecurityGroupName = pointy.String(cast.ToString(v))
		}
	}
	return &groups
}

func getRulesSchema(isForAttr bool) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Computed:      isForAttr,
		ForceNew:      !isForAttr,
		ConflictsWith: []string{"ip_protocol", "security_group_name_to_link"},

		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port_range": {
					Type:     schema.TypeInt,
					Optional: !isForAttr,
					ForceNew: !isForAttr,
					Computed: isForAttr,
				},
				"ip_protocol": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: !isForAttr,
					ForceNew: !isForAttr,
				},
				"ip_ranges": {
					Type:     schema.TypeList,
					Optional: !isForAttr,
					ForceNew: !isForAttr,
					Computed: isForAttr,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"service_ids": {
					Type:     schema.TypeList,
					Optional: !isForAttr,
					ForceNew: !isForAttr,
					Computed: isForAttr,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"to_port_range": {
					Type:     schema.TypeInt,
					Optional: !isForAttr,
					ForceNew: !isForAttr,
					Computed: isForAttr,
				},
				"security_groups_members": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_id": {
								Type:     schema.TypeString,
								Optional: !isForAttr,
								ForceNew: !isForAttr,
								Computed: isForAttr,
							},
							"security_group_id": {
								Type:     schema.TypeString,
								Optional: !isForAttr,
								ForceNew: !isForAttr,
								Computed: isForAttr,
							},
							"security_group_name": {
								Type:     schema.TypeString,
								Optional: !isForAttr,
								ForceNew: !isForAttr,
								Computed: isForAttr,
							},
						},
					},
				},
			},
		},
	}
}
