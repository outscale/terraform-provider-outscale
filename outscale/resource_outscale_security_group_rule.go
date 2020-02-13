package outscale

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/outscale/osc-go/oapi"
)

const OAPI_INBOUND_RULE = "Inbound"
const OAPI_OUTBOUND_RULE = "Outbound"

func sgProtocolIntegers() map[string]int {
	var protocolIntegers = make(map[string]int)
	protocolIntegers = map[string]int{
		"udp":  17,
		"tcp":  6,
		"icmp": 1,
		"all":  -1,
	}
	return protocolIntegers
}

func protocolForValue(v string) string {
	protocol := strings.ToLower(v)
	if protocol == "-1" || protocol == "all" {
		return "-1"
	}
	if _, ok := sgProtocolIntegers()[protocol]; ok {
		return protocol
	}
	p, err := strconv.Atoi(protocol)
	if err != nil {
		fmt.Printf("\n\n[WARN] Unable to determine valid protocol: %s", err)
		return protocol
	}

	for k, v := range sgProtocolIntegers() {
		if p == v {
			return strings.ToLower(k)
		}
	}

	return protocol
}

func resourceOutscaleOAPIOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIOutboundRuleCreate,
		Read:   resourceOutscaleOAPIOutboundRuleRead,
		Delete: resourceOutscaleOAPIOutboundRuleDelete,

		Schema: map[string]*schema.Schema{
			"flow": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					OAPI_INBOUND_RULE,
					OAPI_OUTBOUND_RULE,
				}, false),
			},
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
			"security_group_id": {
				Type: schema.TypeString,
				//Required: true,
				Optional: true,
				ForceNew: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_name_to_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_account_id_to_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"to_port_range": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"rules":          getIPOAPIPermissionsSchema(false),
			"inbound_rules":  getIPOAPIPermissionsSchema(true),
			"outbound_rules": getIPOAPIPermissionsSchema(true),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

var awsMutexKV = mutexkv.NewMutexKV()

func getIPOAPIPermissionsSchema(isForAttr bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: isForAttr,
		ForceNew: !isForAttr,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port_range": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: isForAttr,
					ForceNew: !isForAttr,
				},
				"ip_protocol": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: isForAttr,
					ForceNew: !isForAttr,
				},
				"ip_ranges": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: !isForAttr,
					Computed: isForAttr,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"service_ids": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: !isForAttr,
					Computed: isForAttr,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"to_port_range": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: !isForAttr,
					Computed: isForAttr,
				},
				"security_groups_members": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"security_group_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"security_group_name": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func moreThanOnePresent(ipProtocol, accountIdtoLink string, rules []interface{}) bool {
	argCount := 0
	if ipProtocol != "" {
		argCount++
	}
	if accountIdtoLink != "" {
		argCount++
	}
	if len(rules) > 0 {
		argCount++
	}
	return argCount > 1
}

func resourceOutscaleOAPIOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	sgID := d.Get("security_group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findOSCAPIResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	flow := d.Get("flow").(string)
	ipRange := d.Get("ip_range").(string)
	fromPortRange := d.Get("from_port_range").(int)
	ipProtocol := protocolForValue(d.Get("ip_protocol").(string))
	nameToLink := d.Get("security_group_name_to_link").(string)
	accountIdtoLink := d.Get("security_group_account_id_to_link").(string)
	toPortRange := d.Get("to_port_range").(int)
	rules := d.Get("rules").([]interface{})

	if moreThanOnePresent(ipProtocol, accountIdtoLink, rules) {
		return fmt.Errorf(
			"These parameters cannot be provided together: ip_protocol, rules, security_group_account_id_to_link. Expected at most: 1")
	}

	isOneRule := ipProtocol != ""
	expandedRules := []oscgo.SecurityGroupRule{}
	singleExpandedRule := []oscgo.SecurityGroupRule{}
	fPortRange := int64(fromPortRange)
	tPortRange := int64(toPortRange)
	if !isOneRule {
		expandedRules, err = expandOAPISecurityGroupRules(d, sg)
		if err != nil {
			return err
		}
		if err := validateOAPISecurityGroupRule(rules); err != nil {
			return err
		}
	} else {
		singleExpandedRule = []oscgo.SecurityGroupRule{
			{
				IpRanges:      &[]string{ipRange},
				FromPortRange: &fPortRange,
				IpProtocol:    &ipProtocol,
				ToPortRange:   &tPortRange,
			},
		}
	}

	req := oscgo.CreateSecurityGroupRuleRequest{
		SecurityGroupId:              sg.GetSecurityGroupId(),
		Rules:                        &expandedRules,
		Flow:                         flow,
		IpRange:                      &ipRange,
		FromPortRange:                &fPortRange,
		IpProtocol:                   &ipProtocol,
		SecurityGroupNameToLink:      &nameToLink,
		SecurityGroupAccountIdToLink: &accountIdtoLink,
		ToPortRange:                  &tPortRange,
	}

	//fmt.Printf("Req -> %+v\n", req)

	var autherr error
	log.Printf("[DEBUG] Authorizing security group %s %s rule: %#v", sgID, "Egress", expandedRules)

	var resp oscgo.CreateSecurityGroupRuleResponse
	autherr = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, _, err = conn.SecurityGroupRuleApi.CreateSecurityGroupRule(context.Background(), &oscgo.CreateSecurityGroupRuleOpts{CreateSecurityGroupRuleRequest: optional.NewInterface(req)})

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
			flow, autherr)
	}

	id := ipOSCAPIPermissionIDHash(flow, sgID, expandedRules)
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	var configRules []oscgo.SecurityGroupRule
	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, _, err = findOSCAPIResourceSecurityGroup(conn, sgID)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sgID, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []oscgo.SecurityGroupRule
		if OAPI_INBOUND_RULE == flow {
			rules = sg.GetInboundRules()
		} else {
			rules = sg.GetOutboundRules()
		}

		if isOneRule {
			configRules = singleExpandedRule
		} else {
			configRules = expandedRules
		}

		rule := findOSCAPIRuleMatch(configRules, rules)

		if len(rule) == 0 {
			log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				flow, id, sgID)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		return nil
	})

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			flow, id, sgID)
	}

	d.SetId(id)

	ips, err := setOSCAPIFromIPPerm(d, sg, findOSCAPIRuleMatch(configRules, sg.GetInboundRules()))

	if err != nil {
		return err
	}

	d.Set("inbound_rules", ips)

	ips, err = setOSCAPIFromIPPerm(d, sg, findOSCAPIRuleMatch(configRules, sg.GetOutboundRules()))

	if err != nil {
		return err
	}

	d.Set("outbound_rules", ips)
	d.Set("security_group_name", sg.SecurityGroupName)
	d.Set("net_id", sg.NetId)

	d.Set("request_id", resp.ResponseContext.GetRequestId())

	return nil
}

func resourceOutscaleOAPIOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	sgID := d.Get("security_group_id").(string)
	sg, requestID, err := findOSCAPIResourceSecurityGroup(conn, sgID)
	if _, notFound := err.(oapiSecurityGroupNotFound); notFound {
		// The security group containing this rule no longer exists.
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	if len(sg.GetInboundRules()) == 0 && len(sg.GetOutboundRules()) == 0 {
		log.Printf("[WARN] No rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			sg.GetSecurityGroupName(), d.Id())
		d.SetId("")
		return nil
	}

	flow := d.Get("flow").(string)
	ipRange := d.Get("ip_range").(string)
	fromPortRange := d.Get("from_port_range").(int)
	ipProtocol := protocolForValue(d.Get("ip_protocol").(string))
	//nameToLink := d.Get("security_group_name_to_link").(string)
	accountIdtoLink := d.Get("security_group_account_id_to_link").(string)
	toPortRange := d.Get("to_port_range").(int)
	rules := d.Get("rules").([]interface{})

	if moreThanOnePresent(ipProtocol, accountIdtoLink, rules) {
		return fmt.Errorf(
			"These parameters cannot be provided together: ip_protocol, rules, security_group_account_id_to_link. Expected at most: 1")
	}

	isOneRule := ipProtocol != ""
	expandedRules := []oscgo.SecurityGroupRule{}
	singleExpandedRule := []oscgo.SecurityGroupRule{}

	if !isOneRule {
		expandedRules, err = expandOAPISecurityGroupRules(d, sg)
		if err != nil {
			return err
		}
		if err := validateOAPISecurityGroupRule(rules); err != nil {
			return err
		}
	} else {
		fPortRange := int64(fromPortRange)
		tPortRange := int64(toPortRange)
		singleExpandedRule = []oscgo.SecurityGroupRule{
			{
				IpRanges:      &[]string{ipRange},
				FromPortRange: &fPortRange,
				IpProtocol:    &ipProtocol,
				ToPortRange:   &tPortRange,
			},
		}
	}

	var configRules []oscgo.SecurityGroupRule
	if isOneRule {
		configRules = singleExpandedRule
	} else {
		configRules = expandedRules
	}

	var existingRules []oscgo.SecurityGroupRule
	if OAPI_INBOUND_RULE == flow {
		existingRules = sg.GetInboundRules()
	} else {
		existingRules = sg.GetOutboundRules()
	}

	rule := findOSCAPIRuleMatch(configRules, existingRules)

	if len(rule) == 0 {
		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			flow, d.Id(), sgID)
		d.SetId("")
	}

	ips, err := setOSCAPIFromIPPerm(d, sg, findOSCAPIRuleMatch(configRules, sg.GetInboundRules()))

	if err != nil {
		return err
	}

	d.Set("inbound_rules", ips)

	ips, err = setOSCAPIFromIPPerm(d, sg, findOSCAPIRuleMatch(configRules, sg.GetOutboundRules()))

	if err != nil {
		return err
	}

	d.Set("outbound_rules", ips)

	d.Set("security_group_name", sg.GetSecurityGroupName())
	d.Set("net_id", sg.GetNetId())

	d.Set("request_id", requestID)

	return nil
}

func resourceOutscaleOAPIOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	sgID := d.Get("security_group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	flow := d.Get("flow").(string)
	ipRange := d.Get("ip_range").(string)
	fromPortRange := d.Get("from_port_range").(int)
	ipProtocol := protocolForValue(d.Get("ip_protocol").(string))
	nameToUnlink := d.Get("security_group_name_to_link").(string)
	accountIdtoUnlink := d.Get("security_group_account_id_to_link").(string)
	toPortRange := d.Get("to_port_range").(int)
	rules := d.Get("rules").([]interface{})

	sg, _, err := findOSCAPIResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	isOneRule := ipProtocol != ""
	expandedRules := []oscgo.SecurityGroupRule{}

	if !isOneRule {
		expandedRules, err = expandOAPISecurityGroupRules(d, sg)
		if err != nil {
			return err
		}
		if err := validateOAPISecurityGroupRule(rules); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
		sgID, flow, expandedRules)

	fPortRange := int64(fromPortRange)
	tPortRange := int64(toPortRange)
	req := oscgo.DeleteSecurityGroupRuleRequest{
		SecurityGroupId:                sg.GetSecurityGroupId(),
		Rules:                          &expandedRules,
		Flow:                           flow,
		IpRange:                        &ipRange,
		FromPortRange:                  &fPortRange,
		IpProtocol:                     &ipProtocol,
		SecurityGroupNameToUnlink:      &nameToUnlink,
		SecurityGroupAccountIdToUnlink: &accountIdtoUnlink,
		ToPortRange:                    &tPortRange,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.SecurityGroupRuleApi.DeleteSecurityGroupRule(context.Background(), &oscgo.DeleteSecurityGroupRuleOpts{DeleteSecurityGroupRuleRequest: optional.NewInterface(req)})

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

func findOSCAPIResourceSecurityGroup(conn *oscgo.APIClient, id string) (*oscgo.SecurityGroup, *string, error) {
	req := oscgo.ReadSecurityGroupsRequest{
		Filters: &oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &[]string{id},
		},
	}

	var err error
	var resp oscgo.ReadSecurityGroupsResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(context.Background(), &oscgo.ReadSecurityGroupsOpts{ReadSecurityGroupsRequest: optional.NewInterface(req)})

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
	if len(resp.GetSecurityGroups()) != 1 {
		return nil, nil, oapiSecurityGroupNotFound{id, resp.GetSecurityGroups()}
	}

	return &resp.GetSecurityGroups()[0], resp.ResponseContext.RequestId, nil
}

func expandOAPISecurityGroupRules(d *schema.ResourceData, sg *oscgo.SecurityGroup) ([]oscgo.SecurityGroupRule, error) {

	ippems := d.Get("rules").([]interface{})
	perms := make([]oscgo.SecurityGroupRule, len(ippems))

	return expandOSCAPIIPPerm(d, sg, perms, ippems)
}

func expandOSCAPIIPPerm(d *schema.ResourceData, sg *oscgo.SecurityGroup, perms []oscgo.SecurityGroupRule, ippems []interface{}) ([]oscgo.SecurityGroupRule, error) {

	for k, ip := range ippems {
		perm := oscgo.SecurityGroupRule{}
		v := ip.(map[string]interface{})

		perm.SetFromPortRange(int64(v["from_port_range"].(int)))
		perm.SetToPortRange(int64(v["to_port_range"].(int)))
		protocol := protocolForValue(v["ip_protocol"].(string))
		perm.SetIpProtocol(protocol)

		members := v["security_groups_members"].([]interface{})

		if len(members) > 0 {
			perm.SetSecurityGroupsMembers(make([]oscgo.SecurityGroupsMember, len(members)))
			for i, v := range members {
				member := v.(map[string]interface{})
				accountID := member["account_id"].(string)
				securityGroupID := member["security_group_id"].(string)
				securityGroupName := member["security_group_name"].(string)

				perm.GetSecurityGroupsMembers()[i] = oscgo.SecurityGroupsMember{
					AccountId:         &accountID,
					SecurityGroupId:   &securityGroupID,
					SecurityGroupName: &securityGroupName,
				}
			}
		}

		if raw, ok := v["ip_ranges"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.SetIpRanges(make([]string, len(list)))
				for i, v := range list {
					cidrIP, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in ip_ranges - consider using the compact function")
					}
					perm.GetIpRanges()[i] = cidrIP
				}
			}
		}

		if raw, ok := v["service_ids"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.SetServiceIds(make([]string, len(list)))
				for i, v := range list {
					prefixListID, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in service_ids - consider using the compact function")
					}
					perm.GetServiceIds()[i] = prefixListID
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

		members := v["security_groups_members"].([]interface{})

		if len(members) > 0 {
			for _, v := range members {
				member := v.(map[string]interface{})

				if member["security_group_id"].(string) == "" {
					return fmt.Errorf(
						"'security_group_id' must be set")
				}
			}
		}
	}

	return nil
}

func ipOAPIPermissionIDHash(ruleType, sgID string, ips []oapi.SecurityGroupRule) string {
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
			//sort.Sort(ByGroupsMember(ip.SecurityGroupsMembers))
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

func findOAPIRuleMatch(p []oapi.SecurityGroupRule, rules []oapi.SecurityGroupRule) []oapi.SecurityGroupRule {
	var rule = make([]oapi.SecurityGroupRule, 0)
	//fmt.Printf("Rules (from config) -> %+v\n", p)
	//fmt.Printf("Rules (from service) -> %+v\n", rules)
	for _, i := range p {
		for _, r := range rules {

			//fmt.Printf("Rule (from config) -> %+v\nRule (from service) -> %+v\n", i, r)
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
					if ip.SecurityGroupId == rip.SecurityGroupId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			rule = append(rule, r)
		}
	}
	return rule
}

func setOAPIFromIPPerm(d *schema.ResourceData, sg *oapi.SecurityGroup, rules []oapi.SecurityGroupRule) ([]map[string]interface{}, error) {
	ips := make([]map[string]interface{}, len(rules))

	for k, rule := range rules {
		ip := make(map[string]interface{})

		ip["from_port_range"] = rule.FromPortRange
		ip["to_port_range"] = rule.ToPortRange
		ip["ip_protocol"] = rule.IpProtocol
		ip["ip_ranges"] = rule.IpRanges
		ip["service_ids"] = rule.PrefixListIds

		if len(rule.SecurityGroupsMembers) > 0 {
			s := rule.SecurityGroupsMembers[0]

			d.Set("account_id", s.AccountId)
			d.Set("security_group_id", s.SecurityGroupId)
			d.Set("security_group_name", s.SecurityGroupName)
		}

		ips[k] = ip
	}

	return ips, nil
}

func ipOSCAPIPermissionIDHash(ruleType, sgID string, ips []oscgo.SecurityGroupRule) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", sgID))

	for _, ip := range ips {
		if ip.GetFromPortRange() > 0 {
			buf.WriteString(fmt.Sprintf("%d-", ip.GetFromPortRange()))
		}
		if ip.GetToPortRange() > 0 {
			buf.WriteString(fmt.Sprintf("%d-", ip.GetToPortRange()))
		}
		buf.WriteString(fmt.Sprintf("%s-", ip.GetIpProtocol()))
		buf.WriteString(fmt.Sprintf("%s-", ruleType))

		// We need to make sure to sort the strings below so that we always
		// generate the same hash code no matter what is in the set.
		if len(ip.GetIpRanges()) > 0 {
			s := make([]string, len(ip.GetIpRanges()))
			copy(s, ip.GetIpRanges())
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.GetServiceIds()) > 0 {
			s := make([]string, len(ip.GetServiceIds()))
			copy(s, ip.GetServiceIds())
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.GetSecurityGroupsMembers()) > 0 {
			sort.Sort(ByGroupsMember(ip.GetSecurityGroupsMembers()))
			for _, pair := range ip.GetSecurityGroupsMembers() {
				if pair.GetSecurityGroupId() != "" {
					buf.WriteString(fmt.Sprintf("%s-", pair.GetSecurityGroupId()))
				} else {
					buf.WriteString("-")
				}
				if pair.GetSecurityGroupName() != "" {
					buf.WriteString(fmt.Sprintf("%s-", pair.GetSecurityGroupName()))
				} else {
					buf.WriteString("-")
				}
			}
		}
	}

	return fmt.Sprintf("sgrule-%d", hashcode.String(buf.String()))
}

func findOSCAPIRuleMatch(p []oscgo.SecurityGroupRule, rules []oscgo.SecurityGroupRule) []oscgo.SecurityGroupRule {
	var rule = make([]oscgo.SecurityGroupRule, 0)
	//fmt.Printf("Rules (from config) -> %+v\n", p)
	//fmt.Printf("Rules (from service) -> %+v\n", rules)
	for _, i := range p {
		for _, r := range rules {

			//fmt.Printf("Rule (from config) -> %+v\nRule (from service) -> %+v\n", i, r)
			if i.GetToPortRange() != r.GetToPortRange() {
				continue
			}

			if i.GetFromPortRange() != r.GetFromPortRange() {
				continue
			}

			if i.GetIpProtocol() != r.GetIpProtocol() {
				continue
			}

			remaining := len(i.GetIpRanges())
			for _, ip := range i.GetIpRanges() {
				for _, rip := range r.GetIpRanges() {
					if ip == rip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(i.GetServiceIds())
			for _, pl := range i.GetServiceIds() {
				for _, rpl := range r.GetServiceIds() {
					if pl == rpl {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(i.GetSecurityGroupsMembers())
			for _, ip := range i.GetSecurityGroupsMembers() {
				for _, rip := range r.GetSecurityGroupsMembers() {
					if ip.GetSecurityGroupId() == rip.GetSecurityGroupId() {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			rule = append(rule, r)
		}
	}
	return rule
}

func setOSCAPIFromIPPerm(d *schema.ResourceData, sg *oscgo.SecurityGroup, rules []oscgo.SecurityGroupRule) ([]map[string]interface{}, error) {
	ips := make([]map[string]interface{}, len(rules))

	for k, rule := range rules {
		ip := make(map[string]interface{})

		ip["from_port_range"] = rule.FromPortRange
		ip["to_port_range"] = rule.ToPortRange
		ip["ip_protocol"] = rule.IpProtocol
		ip["ip_ranges"] = rule.IpRanges
		ip["service_ids"] = rule.ServiceIds

		if len(rule.GetSecurityGroupsMembers()) > 0 {
			s := rule.GetSecurityGroupsMembers()[0]

			d.Set("account_id", s.GetAccountId())
			d.Set("security_group_id", s.GetSecurityGroupId())
			d.Set("security_group_name", s.GetSecurityGroupName())
		}

		ips[k] = ip
	}

	return ips, nil
}

type oapiSecurityGroupNotFound struct {
	id             string
	securityGroups []oscgo.SecurityGroup
}

func (err oapiSecurityGroupNotFound) Error() string {
	if err.securityGroups == nil {
		return fmt.Sprintf("No security group with ID %q", err.id)
	}
	return fmt.Sprintf("Expected to find one security group with ID %q, got: %#v",
		err.id, err.securityGroups)
}

// ByGroupsMember ..
type ByGroupsMember []oscgo.SecurityGroupsMember

func (b ByGroupsMember) Len() int      { return len(b) }
func (b ByGroupsMember) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByGroupsMember) Less(i, j int) bool {
	if b[i].GetSecurityGroupId() != "" && b[j].GetSecurityGroupId() != "" {
		return b[i].GetSecurityGroupId() < b[j].GetSecurityGroupId()
	}
	if b[i].GetSecurityGroupName() != "" && b[j].GetSecurityGroupName() != "" {
		return b[i].GetSecurityGroupName() < b[j].GetSecurityGroupName()
	}

	panic("mismatched security group rules, may be a terraform bug")
}

func validateAwsSecurityGroupRule(ippems []interface{}) error {

	for _, value := range ippems {
		v := value.(map[string]interface{})

		_, blocksOk := v["ip_ranges"]
		_, sourceOk := v["source_security_group_owner_id"]
		_, selfOk := v["self"]
		_, prefixOk := v["prefix_list_ids"]
		if !blocksOk && !sourceOk && !selfOk && !prefixOk {
			return fmt.Errorf(
				"One of ['cidr_blocks', 'self', 'source_security_group_id', 'prefix_list_ids'] must be set to create an AWS Security Group Rule")
		}
	}

	return nil
}
