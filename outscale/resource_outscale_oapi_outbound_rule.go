package outscale

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func resourceOutscaleOAPIOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIOutboundRuleCreate,
		Read:   resourceOutscaleOAPIOutboundRuleRead,
		Delete: resourceOutscaleOAPIOutboundRuleDelete,

		Schema: map[string]*schema.Schema{
			"ip_range": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"from_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"firewall_rules_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_firewall_rules_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_firewall_rules_set_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"to_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"inbound_rule": getIPOAPIPermissionsSchema(),
		},
	}
}

func getIPOAPIPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port_range": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"groups": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Set:      schema.HashString,
				},
				"to_port_range": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"ip_protocol": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"ip_ranges": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						// ValidateFunc: validateCIDRNetworkAddress,
					},
				},
				"prefix_list_ids": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func resourceOutscaleOAPIOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandOAPIIPPermEgress(d, sg)
	if err != nil {
		return err
	}

	ippems := d.Get("inbound_rule").([]interface{})

	if err := validateOAPISecurityGroupRule(ippems); err != nil {
		return err
	}

	ruleType := "egress"
	isVPC := sg.NetId != ""

	var autherr error
	log.Printf("[DEBUG] Authorizing security group %s %s rule: %#v", sgID, "Egress", perms)

	req := oapi.CreateSecurityGroupRuleRequest{
		SecurityGroupId: sg.SecurityGroupId,
		Rules:           perms,
	}

	autherr = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.POST_CreateSecurityGroupRule(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if autherr != nil {
		if strings.Contains(fmt.Sprint(autherr), "InvalidPermission.Duplicate") {
			return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sgID, "InvalidPermission.Duplicate")
		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			ruleType, autherr)
	}

	id := ipOAPIPermissionIDHash(sgID, ruleType, perms)
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sgID, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []oapi.SecurityGroupRule
		rules = sg.OutboundRules
		rule := findOAPIRuleMatch(perms, rules, isVPC)

		if rule == nil {
			log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				ruleType, id, sgID)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		return nil
	})

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			ruleType, id, sgID)
	}

	d.SetId(id)
	return nil
}

func resourceOutscaleOAPIOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)
	sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)
	if _, notFound := err.(securityGroupNotFound); notFound {
		// The security group containing this rule no longer exists.
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	isVPC := sg.NetId != ""

	var rule *oapi.SecurityGroupRule
	var rules []oapi.SecurityGroupRule
	ruleType := "egress"
	rules = sg.OutboundRules

	p, err := expandOAPIIPPermEgress(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		log.Printf("[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, sg.SecurityGroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findOAPIRuleMatch(p, rules, isVPC)

	if rule == nil {
		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sgID)
		d.SetId("")
		return nil
	}

	if ips, err := setOAPIFromIPPerm(d, sg, p); err != nil {
		return d.Set("inbound_rule", ips)
	}
	return nil
}

func resourceOutscaleOAPIOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	sgID := d.Get("firewall_rules_set_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findOAPIResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandOAPIIPPermEgress(d, sg)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
		sgID, "egress", perms)
	req := oapi.DeleteSecurityGroupRuleRequest{
		SecurityGroupId: sg.SecurityGroupId,
		Rules:           perms,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.POST_DeleteSecurityGroupRule(req)

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
			"Error revoking security group %s rules: %s",
			sgID, err)
	}

	d.SetId("")

	return nil
}

// #################################

func findOAPIResourceSecurityGroup(conn *oapi.Client, id string) (*oapi.SecurityGroup, *string, error) {
	req := oapi.ReadSecurityGroupsRequest{
		Filters: oapi.FiltersSecurityGroup{
			SecurityGroupIds: []string{id},
		},
	}

	var err error
	var resp *oapi.POST_ReadSecurityGroupsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadSecurityGroups(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err, ok := err.(awserr.Error); ok && err.Code() == "InvalidGroup.NotFound" {
		return nil, nil, oapiSecurityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, oapiSecurityGroupNotFound{id, nil}
	}
	if len(resp.OK.SecurityGroups) != 1 {
		return nil, nil, oapiSecurityGroupNotFound{id, resp.OK.SecurityGroups}
	}

	return &resp.OK.SecurityGroups[0], &resp.OK.ResponseContext.RequestId, nil
}

func expandOAPIIPPermEgress(d *schema.ResourceData, sg *oapi.SecurityGroup) ([]oapi.SecurityGroupRule, error) {

	ippems := d.Get("inbound_rule").([]interface{})
	perms := make([]oapi.SecurityGroupRule, len(ippems))

	return expandOAPIIPPerm(d, sg, perms, ippems)
}

func expandOAPIIPPerm(d *schema.ResourceData, sg *oapi.SecurityGroup, perms []oapi.SecurityGroupRule, ippems []interface{}) ([]oapi.SecurityGroupRule, error) {

	for k, ip := range ippems {
		perm := oapi.SecurityGroupRule{}
		v := ip.(map[string]interface{})

		perm.FromPortRange = int64(v["from_port_range"].(int))
		perm.ToPortRange = int64(v["to_port_range"].(int))
		protocol := protocolForValue(v["ip_protocol"].(string))
		perm.IpProtocol = protocol

		groups := make(map[string]bool)
		if raw, ok := d.GetOk("destination_firewall_rules_set_account_id"); ok {
			groups[raw.(string)] = true
		}

		if v, ok := d.GetOk("self"); ok && v.(bool) {
			if sg.NetId != "" {
				groups[sg.SecurityGroupId] = true
			} else {
				groups[sg.SecurityGroupName] = true
			}
		}

		if len(groups) > 0 {
			perm.SecurityGroupsMembers = make([]oapi.SecurityGroupsMember, len(groups))
			// build string list of group name/ids
			var gl []string
			for k := range groups {
				gl = append(gl, k)
			}

			for i, name := range gl {
				ownerID, id := "", name
				if items := strings.Split(id, "/"); len(items) > 1 {
					ownerID, id = items[0], items[1]
				}

				perm.SecurityGroupsMembers[i] = oapi.SecurityGroupsMember{
					SecurityGroupId: id,
					AccountId:       ownerID,
				}

				if sg.NetId == "" {
					perm.SecurityGroupsMembers[i].SecurityGroupId = ""
					perm.SecurityGroupsMembers[i].SecurityGroupName = id
					perm.SecurityGroupsMembers[i].AccountId = ""
				}
			}
		}

		if raw, ok := v["ip_ranges"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.IpRanges = make([]string, len(list))
				for i, v := range list {
					cidrIP, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in cidr_blocks - consider using the compact function")
					}
					perm.IpRanges[i] = cidrIP
				}
			}
		}

		if raw, ok := v["prefix_list_ids"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.PrefixListIds = make([]string, len(list))
				for i, v := range list {
					prefixListID, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in prefix_list_ids - consider using the compact function")
					}
					perm.PrefixListIds[i] = prefixListID
				}
			}
		}

		perms[k] = perm
	}
	return perms, nil
}

func validateOAPISecurityGroupRule(ippems []interface{}) error {

	for _, value := range ippems {
		v := value.(map[string]interface{})

		_, blocksOk := v["ip_ranges"]
		_, sourceOk := v["destination_firewall_rules_set_account_id"]
		_, selfOk := v["self"]
		_, prefixOk := v["prefix_list_ids"]
		if !blocksOk && !sourceOk && !selfOk && !prefixOk {
			return fmt.Errorf(
				"One of ['cidr_blocks', 'self', 'source_security_group_id', 'prefix_list_ids'] must be set to create an AWS Security Group Rule")
		}
	}

	return nil
}

func ipOAPIPermissionIDHash(sgID, ruleType string, ips []oapi.SecurityGroupRule) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", sgID))

	for _, ip := range ips {
		if ip.FromPortRange > 0 {
			buf.WriteString(fmt.Sprintf("%d-", ip.FromPortRange))
		}
		if ip.ToPortRange > 0 {
			buf.WriteString(fmt.Sprintf("%d-", ip.ToPortRange))
		}
		buf.WriteString(fmt.Sprintf("%s-", ip.IpProtocol))
		buf.WriteString(fmt.Sprintf("%s-", ruleType))

		// We need to make sure to sort the strings below so that we always
		// generate the same hash code no matter what is in the set.
		if len(ip.IpRanges) > 0 {
			s := make([]string, len(ip.IpRanges))
			copy(s, ip.IpRanges)
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.PrefixListIds) > 0 {
			s := make([]string, len(ip.PrefixListIds))
			copy(s, ip.PrefixListIds)
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.SecurityGroupsMembers) > 0 {
			sort.Sort(ByGroupsMember(ip.SecurityGroupsMembers))
			for _, pair := range ip.SecurityGroupsMembers {
				if pair.SecurityGroupId != "" {
					buf.WriteString(fmt.Sprintf("%s-", pair.SecurityGroupId))
				} else {
					buf.WriteString("-")
				}
				if pair.SecurityGroupName != "" {
					buf.WriteString(fmt.Sprintf("%s-", pair.SecurityGroupName))
				} else {
					buf.WriteString("-")
				}
			}
		}
	}

	return fmt.Sprintf("sgrule-%d", hashcode.String(buf.String()))
}

func findOAPIRuleMatch(p []oapi.SecurityGroupRule, rules []oapi.SecurityGroupRule, isVPC bool) *oapi.SecurityGroupRule {
	var rule *oapi.SecurityGroupRule
	for _, i := range p {
		for _, r := range rules {
			if i.ToPortRange != r.ToPortRange {
				continue
			}

			if i.FromPortRange != r.FromPortRange {
				continue
			}

			if i.IpProtocol != r.IpProtocol {
				continue
			}

			remaining := len(i.IpRanges)
			for _, ip := range i.IpRanges {
				for _, rip := range r.IpRanges {
					if ip == rip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(i.PrefixListIds)
			for _, pl := range i.PrefixListIds {
				for _, rpl := range r.PrefixListIds {
					if pl == rpl {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(i.SecurityGroupsMembers)
			for _, ip := range i.SecurityGroupsMembers {
				for _, rip := range r.SecurityGroupsMembers {
					if isVPC {
						if ip.SecurityGroupId == rip.SecurityGroupId {
							remaining--
						}
					} else {
						if ip.SecurityGroupName == rip.SecurityGroupName {
							remaining--
						}
					}
				}
			}

			if remaining > 0 {
				continue
			}

			rule = &r
		}
	}
	return rule
}

func setOAPIFromIPPerm(d *schema.ResourceData, sg *oapi.SecurityGroup, rules []oapi.SecurityGroupRule) ([]map[string]interface{}, error) {
	isVPC := sg.NetId != ""

	ips := make([]map[string]interface{}, len(rules))

	for k, rule := range rules {
		ip := make(map[string]interface{})

		ip["from_port_range"] = rule.FromPortRange
		ip["to_port_range"] = rule.ToPortRange
		ip["ip_protocol"] = rule.IpProtocol
		ip["ip_ranges"] = rule.IpRanges
		ip["prefix_list_ids"] = rule.PrefixListIds

		if len(rule.SecurityGroupsMembers) > 0 {
			s := rule.SecurityGroupsMembers[0]

			if isVPC {
				d.Set("destination_firewall_rules_set_account_id", s.SecurityGroupId)
			} else {
				d.Set("destination_firewall_rules_set_account_id", s.SecurityGroupName)
			}
		}

		ips[k] = ip
	}

	return ips, nil
}

type oapiSecurityGroupNotFound struct {
	id             string
	securityGroups []oapi.SecurityGroup
}

func (err oapiSecurityGroupNotFound) Error() string {
	if err.securityGroups == nil {
		return fmt.Sprintf("No security group with ID %q", err.id)
	}
	return fmt.Sprintf("Expected to find one security group with ID %q, got: %#v",
		err.id, err.securityGroups)
}

type ByGroupsMember []oapi.SecurityGroupsMember

func (b ByGroupsMember) Len() int      { return len(b) }
func (b ByGroupsMember) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByGroupsMember) Less(i, j int) bool {
	if b[i].SecurityGroupId != "" && b[j].SecurityGroupId != "" {
		return b[i].SecurityGroupId < b[j].SecurityGroupId
	}
	if b[i].SecurityGroupName != "" && b[j].SecurityGroupName != "" {
		return b[i].SecurityGroupName < b[j].SecurityGroupName
	}

	panic("mismatched security group rules, may be a terraform bug")
}
