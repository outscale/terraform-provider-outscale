package outscale

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOutboundRuleCreate,
		Read:   resourceOutscaleOutboundRuleRead,
		Delete: resourceOutscaleOutboundRuleDelete,

		Schema: map[string]*schema.Schema{
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_security_group_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ip_permissions": getIPPermissionsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

var awsMutexKV = mutexkv.NewMutexKV()

func getIPPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_port": {
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
				"to_port": {
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

func resourceOutscaleOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandIPPermEgress(d, sg)
	if err != nil {
		return err
	}

	ippems := d.Get("ip_permissions").([]interface{})

	if err := validateAwsSecurityGroupRule(ippems); err != nil {
		return err
	}

	ruleType := "egress"
	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var autherr error
	log.Printf("[DEBUG] Authorizing security group %s %s rule: %#v", sgID, "Egress", perms)

	req := &fcu.AuthorizeSecurityGroupEgressInput{
		GroupId:       sg.GroupId,
		IpPermissions: perms,
	}

	autherr = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.AuthorizeSecurityGroupEgress(req)

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

	id := ipPermissionIDHash(sgID, ruleType, perms)
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, _, err := findResourceSecurityGroup(conn, sgID)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sgID, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []*fcu.IpPermission
		rules = sg.IpPermissionsEgress

		rule := findRuleMatch(perms, rules, isVPC)

		if rule == nil {
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

func resourceOutscaleOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)
	sg, reqID, err := findResourceSecurityGroup(conn, sgID)
	if _, notFound := err.(securityGroupNotFound); notFound {
		// The security group containing this rule no longer exists.
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var rule *fcu.IpPermission
	var rules []*fcu.IpPermission
	ruleType := "egress"
	rules = sg.IpPermissionsEgress

	p, err := expandIPPermEgress(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		log.Printf("[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, *sg.GroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findRuleMatch(p, rules, isVPC)

	if rule == nil {
		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sgID)
		d.SetId("")
		return nil
	}

	if ips, err := setFromIPPerm(d, sg, p); err != nil {
		return d.Set("ip_permissions", ips)
	}
	return d.Set("request_id", aws.StringValue(reqID))
}

func resourceOutscaleOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sgID := d.Get("group_id").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, _, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}

	perms, err := expandIPPermEgress(d, sg)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
		sgID, "egress", perms)
	req := &fcu.RevokeSecurityGroupEgressInput{
		GroupId:       sg.GroupId,
		IpPermissions: perms,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.RevokeSecurityGroupEgress(req)

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

func findResourceSecurityGroup(conn *fcu.Client, id string) (*fcu.SecurityGroup, *string, error) {
	req := &fcu.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	}

	var err error
	var resp *fcu.DescribeSecurityGroupsOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSecurityGroups(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err, ok := err.(awserr.Error); ok && err.Code() == "InvalidGroup.NotFound" {
		return nil, nil, securityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, securityGroupNotFound{id, nil}
	}
	if len(resp.SecurityGroups) != 1 || resp.SecurityGroups[0] == nil {
		return nil, nil, securityGroupNotFound{id, resp.SecurityGroups}
	}

	return resp.SecurityGroups[0], resp.RequestId, nil
}

func expandIPPermEgress(d *schema.ResourceData, sg *fcu.SecurityGroup) ([]*fcu.IpPermission, error) {

	ippems := d.Get("ip_permissions").([]interface{})
	perms := make([]*fcu.IpPermission, len(ippems))

	return expandIPPerm(d, sg, perms, ippems)
}

func expandIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup, perms []*fcu.IpPermission, ippems []interface{}) ([]*fcu.IpPermission, error) {

	for k, ip := range ippems {
		perm := fcu.IpPermission{}
		v := ip.(map[string]interface{})

		perm.FromPort = aws.Int64(int64(v["from_port"].(int)))
		perm.ToPort = aws.Int64(int64(v["to_port"].(int)))
		protocol := protocolForValue(v["ip_protocol"].(string))
		perm.IpProtocol = aws.String(protocol)

		groups := make(map[string]bool)
		if raw, ok := d.GetOk("source_security_group_owner_id"); ok {
			groups[raw.(string)] = true
		}

		if v, ok := d.GetOk("self"); ok && v.(bool) {
			if sg.VpcId != nil && *sg.VpcId != "" {
				groups[*sg.GroupId] = true
			} else {
				groups[*sg.GroupName] = true
			}
		}

		if len(groups) > 0 {
			perm.UserIdGroupPairs = make([]*fcu.UserIdGroupPair, len(groups))
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

				perm.UserIdGroupPairs[i] = &fcu.UserIdGroupPair{
					GroupId: aws.String(id),
					UserId:  aws.String(ownerID),
				}

				if sg.VpcId == nil || *sg.VpcId == "" {
					perm.UserIdGroupPairs[i].GroupId = nil
					perm.UserIdGroupPairs[i].GroupName = aws.String(id)
					perm.UserIdGroupPairs[i].UserId = nil
				}
			}
		}

		if raw, ok := v["ip_ranges"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.IpRanges = make([]*fcu.IpRange, len(list))
				for i, v := range list {
					cidrIP, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in cidr_blocks - consider using the compact function")
					}
					perm.IpRanges[i] = &fcu.IpRange{CidrIp: aws.String(cidrIP)}
				}
			}
		}

		if raw, ok := v["prefix_list_ids"]; ok {
			list := raw.([]interface{})
			if len(list) > 0 {
				perm.PrefixListIds = make([]*fcu.PrefixListId, len(list))
				for i, v := range list {
					prefixListID, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("empty element found in prefix_list_ids - consider using the compact function")
					}
					perm.PrefixListIds[i] = &fcu.PrefixListId{PrefixListId: aws.String(prefixListID)}
				}
			}
		}

		perms[k] = &perm
	}
	return perms, nil
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

func ipPermissionIDHash(sgID, ruleType string, ips []*fcu.IpPermission) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", sgID))

	for _, ip := range ips {
		if ip.FromPort != nil && *ip.FromPort > 0 {
			buf.WriteString(fmt.Sprintf("%d-", *ip.FromPort))
		}
		if ip.ToPort != nil && *ip.ToPort > 0 {
			buf.WriteString(fmt.Sprintf("%d-", *ip.ToPort))
		}
		buf.WriteString(fmt.Sprintf("%s-", *ip.IpProtocol))
		buf.WriteString(fmt.Sprintf("%s-", ruleType))

		// We need to make sure to sort the strings below so that we always
		// generate the same hash code no matter what is in the set.
		if len(ip.IpRanges) > 0 {
			s := make([]string, len(ip.IpRanges))
			for i, r := range ip.IpRanges {
				s[i] = *r.CidrIp
			}
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.PrefixListIds) > 0 {
			s := make([]string, len(ip.PrefixListIds))
			for i, pl := range ip.PrefixListIds {
				s[i] = *pl.PrefixListId
			}
			sort.Strings(s)

			for _, v := range s {
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		if len(ip.UserIdGroupPairs) > 0 {
			sort.Sort(ByGroupPair(ip.UserIdGroupPairs))
			for _, pair := range ip.UserIdGroupPairs {
				if pair.GroupId != nil {
					buf.WriteString(fmt.Sprintf("%s-", *pair.GroupId))
				} else {
					buf.WriteString("-")
				}
				if pair.GroupName != nil {
					buf.WriteString(fmt.Sprintf("%s-", *pair.GroupName))
				} else {
					buf.WriteString("-")
				}
			}
		}
	}

	return fmt.Sprintf("sgrule-%d", hashcode.String(buf.String()))
}

func findRuleMatch(p []*fcu.IpPermission, rules []*fcu.IpPermission, isVPC bool) *fcu.IpPermission {
	var rule *fcu.IpPermission
	for _, i := range p {
		for _, r := range rules {
			if r.ToPort != nil && *i.ToPort != *r.ToPort {
				continue
			}

			if r.FromPort != nil && *i.FromPort != *r.FromPort {
				continue
			}

			if r.IpProtocol != nil && *i.IpProtocol != *r.IpProtocol {
				continue
			}

			remaining := len(i.IpRanges)
			for _, ip := range i.IpRanges {
				for _, rip := range r.IpRanges {
					if *ip.CidrIp == *rip.CidrIp {
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
					if *pl.PrefixListId == *rpl.PrefixListId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(i.UserIdGroupPairs)
			for _, ip := range i.UserIdGroupPairs {
				for _, rip := range r.UserIdGroupPairs {
					if isVPC {
						if *ip.GroupId == *rip.GroupId {
							remaining--
						}
					} else {
						if *ip.GroupName == *rip.GroupName {
							remaining--
						}
					}
				}
			}

			if remaining > 0 {
				continue
			}

			rule = r
		}
	}
	return rule
}

type securityGroupNotFound struct {
	id             string
	securityGroups []*fcu.SecurityGroup
}

func (err securityGroupNotFound) Error() string {
	if err.securityGroups == nil {
		return fmt.Sprintf("No security group with ID %q", err.id)
	}
	return fmt.Sprintf("Expected to find one security group with ID %q, got: %#v",
		err.id, err.securityGroups)
}

// ByGroupPair ...
type ByGroupPair []*fcu.UserIdGroupPair

func (b ByGroupPair) Len() int      { return len(b) }
func (b ByGroupPair) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByGroupPair) Less(i, j int) bool {
	if b[i].GroupId != nil && b[j].GroupId != nil {
		return *b[i].GroupId < *b[j].GroupId
	}
	if b[i].GroupName != nil && b[j].GroupName != nil {
		return *b[i].GroupName < *b[j].GroupName
	}

	panic("mismatched security group rules, may be a terraform bug")
}

func setFromIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup, rules []*fcu.IpPermission) ([]map[string]interface{}, error) {
	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	ips := make([]map[string]interface{}, len(rules))

	for k, rule := range rules {
		ip := make(map[string]interface{})

		if rule.FromPort != nil {
			ip["from_port"] = *rule.FromPort
		}
		if rule.ToPort != nil {
			ip["to_port"] = *rule.ToPort
		}
		if rule.IpProtocol != nil {
			ip["ip_protocol"] = *rule.IpProtocol
		}
		if rule.IpRanges != nil {
			var cb []string
			for _, c := range rule.IpRanges {
				cb = append(cb, *c.CidrIp)
			}
			ip["ip_ranges"] = cb
		}
		if rule.PrefixListIds != nil {
			var pl []string
			for _, p := range rule.PrefixListIds {
				pl = append(pl, *p.PrefixListId)
			}
			ip["prefix_list_ids"] = pl
		}

		if len(rule.UserIdGroupPairs) > 0 {
			s := rule.UserIdGroupPairs[0]

			if isVPC {
				d.Set("source_security_group_owner_id", *s.GroupId)
			} else {
				d.Set("source_security_group_owner_id", *s.GroupName)
			}
		}

		ips[k] = ip
	}

	return ips, nil
}
