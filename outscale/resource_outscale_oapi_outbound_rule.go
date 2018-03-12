package outscale

// import (
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"

// 	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/awserr"
// 	"github.com/hashicorp/errwrap"
// 	"github.com/hashicorp/terraform/helper/resource"
// 	"github.com/hashicorp/terraform/helper/schema"
// )

// func resourceOutscaleOAPIOutboundRule() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceOutscaleOAPIOutboundRuleCreate,
// 		Read:   resourceOutscaleOAPIOutboundRuleRead,
// 		Delete: resourceOutscaleOAPIOutboundRuleDelete,

// 		Schema: map[string]*schema.Schema{
// 			"ip_range": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"from_port_range": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"firewall_rules_set_id": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"ip_protocol": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"destination_firewall_rules_set_name": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"destination_firewall_rules_set_account_id": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"to_port_range": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"inbound_rule": getOAPIIpPermissionsSchema(),
// 		},
// 	}
// }

// func getOAPIIpPermissionsSchema() *schema.Schema {
// 	return &schema.Schema{
// 		Type:     schema.TypeSet,
// 		Optional: true,
// 		ForceNew: true,
// 		Elem: &schema.Resource{
// 			Schema: map[string]*schema.Schema{
// 				"from_port_range": {
// 					Type:     schema.TypeInt,
// 					Optional: true,
// 					ForceNew: true,
// 				},
// 				"groups": {
// 					Type:     schema.TypeSet,
// 					Optional: true,
// 					Elem:     &schema.Schema{Type: schema.TypeString},
// 					Set:      schema.HashString,
// 				},
// 				"to_port_range": {
// 					Type:     schema.TypeInt,
// 					Optional: true,
// 					ForceNew: true,
// 				},
// 				"ip_protocol": {
// 					Type:     schema.TypeString,
// 					Optional: true,
// 					ForceNew: true,
// 				},
// 				"ip_ranges": {
// 					Type:     schema.TypeList,
// 					Optional: true,
// 					ForceNew: true,
// 					Elem: &schema.Schema{
// 						Type:         schema.TypeString,
// 						ValidateFunc: validateCIDRNetworkAddress,
// 					},
// 				},
// 				"prefix_list_ids": {
// 					Type:     schema.TypeList,
// 					Optional: true,
// 					ForceNew: true,
// 					Elem:     &schema.Schema{Type: schema.TypeString},
// 				},
// 			},
// 		},
// 	}
// }

// func resourceOutscaleOAPIOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
// 	conn := meta.(*OutscaleClient).FCU

// 	sg_id := d.Get("firewall_rules_set_id").(string)

// 	awsMutexKV.Lock(sg_id)
// 	defer awsMutexKV.Unlock(sg_id)

// 	sg, err := findResourceOAPISecurityGroup(conn, sg_id)
// 	if err != nil {
// 		return err
// 	}

// 	perm, err := expandOAPIIPPerm(d, sg)
// 	if err != nil {
// 		return err
// 	}

// 	if err := validateOutscaleOAPISecurityGroupRule(d); err != nil {
// 		return err
// 	}

// 	ruleType := "egress"
// 	isVPC := sg.VpcId != nil && *sg.VpcId != ""

// 	var autherr error
// 	fmt.Printf("[DEBUG] Authorizing security group %s %s rule: %#v",
// 		sg_id, "Egress", perm)

// 	req := &fcu.AuthorizeSecurityGroupEgressInput{
// 		GroupId:       sg.GroupId,
// 		IpPermissions: perm,
// 	}

// 	resource.Retry(5*time.Minute, func() *resource.RetryError {
// 		_, autherr = conn.VM.AuthorizeSecurityGroupEgress(req)

// 		if err != nil {
// 			if strings.Contains(err.Error(), "RequestLimitExceeded") || strings.Contains(err.Error(), "DependencyViolation") {
// 				return resource.RetryableError(err)
// 			}
// 			return resource.NonRetryableError(err)
// 		}

// 		return nil
// 	})

// 	if autherr != nil {
// 		if awsErr, ok := autherr.(awserr.Error); ok {
// 			if awsErr.Code() == "InvalidPermission.Duplicate" {
// 				return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
// a side effect of a now-fixed Terraform issue causing two security groups with
// identical attributes but different source_security_group_ids to overwrite each
// other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
// information and instructions for recovery. Error message: %s`, sg_id, awsErr.Message())
// 			}
// 		}

// 		return fmt.Errorf(
// 			"Error authorizing security group rule type %s: %s",
// 			ruleType, autherr)
// 	}

// 	id := ipPermissionIDHash(sg_id, ruleType, perm)
// 	log.Printf("[DEBUG] Computed group rule ID %s", id)

// 	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
// 		sg, err := findResourceOAPISecurityGroup(conn, sg_id)

// 		if err != nil {
// 			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sg_id, id, err)
// 			return resource.NonRetryableError(err)
// 		}

// 		var rules []*fcu.IpPermission
// 		rules = sg.IpPermissionsEgress

// 		rule := findOAPIRuleMatch(perm, rules, isVPC)

// 		if rule == nil {
// 			log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
// 				ruleType, id, sg_id)
// 			return resource.RetryableError(fmt.Errorf("No match found"))
// 		}

// 		log.Printf("[DEBUG] Found rule for Security Group Rule (%s): %s", id, rule)
// 		return nil
// 	})

// 	if retErr != nil {
// 		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
// 			ruleType, id, sg_id)
// 	}

// 	d.SetId(id)
// 	return nil
// }

// func resourceOutscaleOAPIOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
// 	conn := meta.(*OutscaleClient).FCU
// 	sg_id := d.Get("firewall_rules_set_id").(string)
// 	sg, err := findResourceOAPISecurityGroup(conn, sg_id)
// 	if _, notFound := err.(securityGroupNotFound); notFound {
// 		// The security group containing this rule no longer exists.
// 		d.SetId("")
// 		return nil
// 	}
// 	if err != nil {
// 		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sg_id, d.Id(), err)
// 	}

// 	isVPC := sg.VpcId != nil && *sg.VpcId != ""

// 	var rule *fcu.IpPermission
// 	var rules []*fcu.IpPermission
// 	ruleType := "egress"
// 	rules = sg.IpPermissionsEgress

// 	p, err := expandOAPIIPPerm(d, sg)
// 	if err != nil {
// 		return err
// 	}

// 	if len(rules) == 0 {
// 		log.Printf("[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
// 			ruleType, *sg.GroupName, d.Id())
// 		d.SetId("")
// 		return nil
// 	}

// 	rule = findOAPIRuleMatch(p, rules, isVPC)

// 	if rule == nil {
// 		log.Printf("[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
// 			ruleType, d.Id(), sg_id)
// 		d.SetId("")
// 		return nil
// 	}

// 	if err := setOAPIFromIPPerm(d, sg, p); err != nil {
// 		return errwrap.Wrapf("Error setting IP Permission for Security Group Rule: {{err}}", err)
// 	}

// 	log.Printf("[DEBUG] Found rule for Security Group Rule (%s): %s", d.Id(), rule)

// 	return nil
// }

// func resourceOutscaleOAPIOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
// 	conn := meta.(*OutscaleClient).FCU
// 	sg_id := d.Get("firewall_rules_set_id").(string)

// 	awsMutexKV.Lock(sg_id)
// 	defer awsMutexKV.Unlock(sg_id)

// 	sg, err := findResourceOAPISecurityGroup(conn, sg_id)
// 	if err != nil {
// 		return err
// 	}

// 	perm, err := expandOAPIIPPerm(d, sg)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("\n\n[DEBUG] Revoking security group %#v %s rule: %#v",
// 		sg_id, "egress", perm)
// 	req := &fcu.RevokeSecurityGroupEgressInput{
// 		GroupId:       sg.GroupId,
// 		IpPermissions: perm,
// 	}

// 	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 		_, err = conn.VM.RevokeSecurityGroupEgress(req)

// 		if err != nil {
// 			if strings.Contains(err.Error(), "RequestLimitExceeded") {
// 				fmt.Printf("\n\n[INFO] Request limit exceeded")
// 				return resource.RetryableError(err)
// 			}
// 			return resource.NonRetryableError(err)
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return fmt.Errorf(
// 			"Error revoking security group %s rules: %s",
// 			sg_id, err)
// 	}

// 	d.SetId("")

// 	return nil
// }

// func findResourceOAPISecurityGroup(conn *fcu.Client, id string) (*fcu.SecurityGroup, error) {
// 	req := &fcu.DescribeSecurityGroupsInput{
// 		GroupIds: []*string{aws.String(id)},
// 	}

// 	var resp *fcu.DescribeSecurityGroupsOutput
// 	var err error
// 	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 		resp, err = conn.VM.DescribeSecurityGroups(req)

// 		if err != nil {
// 			if strings.Contains(err.Error(), "RequestLimitExceeded") {
// 				fmt.Printf("\n\n[INFO] Request limit exceeded")
// 				return resource.RetryableError(err)
// 			}
// 			return resource.NonRetryableError(err)
// 		}

// 		return nil
// 	})

// 	if strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
// 		return nil, securityGroupNotFound{id, nil}
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp == nil {
// 		return nil, securityGroupNotFound{id, nil}
// 	}
// 	if len(resp.SecurityGroups) != 1 || resp.SecurityGroups[0] == nil {
// 		return nil, securityGroupNotFound{id, resp.SecurityGroups}
// 	}

// 	return resp.SecurityGroups[0], nil
// }

// func expandOAPIIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup) ([]*fcu.IpPermission, error) {
// 	ipp := d.Get("inbound_rule")
// 	ippem := ipp.(*schema.Set).List()

// 	perms := make([]*fcu.IpPermission, len(ippem))

// 	for i, v := range ippem {
// 		perm := fcu.IpPermission{}
// 		values := v.(map[string]interface{})

// 		if raw, ok := values["from_port_range"]; ok {
// 			perm.FromPort = aws.Int64(int64(raw.(int)))
// 		}
// 		if raw, ok := values["to_port_range"]; ok {
// 			perm.ToPort = aws.Int64(int64(raw.(int)))
// 		}
// 		if raw, ok := values["ip_protocol"]; ok {
// 			protocol := protocolForValue(raw.(string))
// 			perm.IpProtocol = aws.String(protocol)
// 		}

// 		if raw, ok := values["ip_ranges"]; ok {
// 			list := raw.([]interface{})
// 			perm.IpRanges = make([]*fcu.IpRange, len(list))
// 			for i, v := range list {
// 				cidrIP, ok := v.(string)
// 				if !ok {
// 					return nil, fmt.Errorf("empty element found in ip_ranges - consider using the compact function")
// 				}
// 				perm.IpRanges[i] = &fcu.IpRange{CidrIp: aws.String(cidrIP)}
// 			}
// 		}

// 		if raw, ok := values["prefix_list_ids"]; ok {
// 			list := raw.([]interface{})
// 			if len(list) > 0 {
// 				perm.PrefixListIds = make([]*fcu.PrefixListId, len(list))
// 				for i, v := range list {
// 					prefixListID, ok := v.(string)
// 					if !ok {
// 						return nil, fmt.Errorf("empty element found in prefix_list_ids - consider using the compact function")
// 					}
// 					perm.PrefixListIds[i] = &fcu.PrefixListId{PrefixListId: aws.String(prefixListID)}
// 				}
// 			}
// 		}

// 		groups := make(map[string]bool)
// 		if raw, ok := d.GetOk("source_security_group_owner_id"); ok {
// 			groups[raw.(string)] = true
// 		}

// 		if len(groups) > 0 {
// 			perm.UserIdGroupPairs = make([]*fcu.UserIdGroupPair, len(groups))
// 			// build string list of group name/ids
// 			var gl []string
// 			for k, _ := range groups {
// 				gl = append(gl, k)
// 			}

// 			for i, name := range gl {
// 				ownerId, id := "", name
// 				if items := strings.Split(id, "/"); len(items) > 1 {
// 					ownerId, id = items[0], items[1]
// 				}

// 				perm.UserIdGroupPairs[i] = &fcu.UserIdGroupPair{
// 					GroupId: aws.String(id),
// 					UserId:  aws.String(ownerId),
// 				}

// 				if sg.VpcId == nil || *sg.VpcId == "" {
// 					perm.UserIdGroupPairs[i].GroupId = nil
// 					perm.UserIdGroupPairs[i].GroupName = aws.String(id)
// 					perm.UserIdGroupPairs[i].UserId = nil
// 				}
// 			}
// 		}

// 		perms[i] = &perm
// 	}

// 	return perms, nil
// }

// func validateOutscaleOAPISecurityGroupRule(d *schema.ResourceData) error {
// 	if ipp, ippemOk := d.GetOk("inbound_rule"); ippemOk {
// 		ippem := ipp.(*schema.Set).List()

// 		for _, v := range ippem {
// 			values := v.(map[string]interface{})

// 			_, blocksOk := values["ip_ranges"]
// 			_, sourceOk := values["destination_firewall_rules_set_account_id"]
// 			_, prefixOk := values["prefix_list_ids"]
// 			if !blocksOk && !sourceOk && !prefixOk {
// 				return fmt.Errorf(
// 					"One of ['ip_ranges', 'destination_firewall_rules_set_account_id', 'prefix_list_ids'] must be set to create an Outscale Security Group Rule")
// 			}
// 		}
// 	}

// 	return nil
// }

// func findOAPIRuleMatch(p []*fcu.IpPermission, rules []*fcu.IpPermission, isVPC bool) *fcu.IpPermission {
// 	var rule *fcu.IpPermission

// 	for _, r := range rules {

// 		if r.ToPort != nil && *p[0].ToPort != *r.ToPort {
// 			continue
// 		}

// 		if r.FromPort != nil && *p[0].FromPort != *r.FromPort {
// 			continue
// 		}

// 		if r.IpProtocol != nil && *p[0].IpProtocol != *r.IpProtocol {
// 			continue
// 		}

// 		remaining := len(p[0].IpRanges)
// 		for _, ip := range p[0].IpRanges {
// 			for _, rip := range r.IpRanges {
// 				if *ip.CidrIp == *rip.CidrIp {
// 					remaining--
// 				}
// 			}
// 		}

// 		if remaining > 0 {
// 			continue
// 		}

// 		remaining = len(p[0].PrefixListIds)
// 		for _, pl := range p[0].PrefixListIds {
// 			for _, rpl := range r.PrefixListIds {
// 				if *pl.PrefixListId == *rpl.PrefixListId {
// 					remaining--
// 				}
// 			}
// 		}

// 		if remaining > 0 {
// 			continue
// 		}

// 		remaining = len(p[0].UserIdGroupPairs)
// 		for _, ip := range p[0].UserIdGroupPairs {
// 			for _, rip := range r.UserIdGroupPairs {
// 				if isVPC {
// 					if *ip.GroupId == *rip.GroupId {
// 						remaining--
// 					}
// 				} else {
// 					if *ip.GroupName == *rip.GroupName {
// 						remaining--
// 					}
// 				}
// 			}
// 		}

// 		if remaining > 0 {
// 			continue
// 		}

// 		rule = r
// 	}
// 	return rule
// }

// func setOAPIFromIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup, rules []*fcu.IpPermission) error {
// 	isVPC := sg.VpcId != nil && *sg.VpcId != ""

// 	ippems := make([]map[string]interface{}, len(rules))

// 	for i, rule := range rules {
// 		ippem := make(map[string]interface{})
// 		ippem["from_port_range"] = *rule.FromPort
// 		ippem["to_port_range"] = *rule.ToPort
// 		ippem["ip_protocol"] = *rule.IpProtocol

// 		cb := make([]*fcu.IpRange, len(rule.IpRanges))
// 		for i, c := range rule.IpRanges {
// 			cb[i] = &fcu.IpRange{CidrIp: c.CidrIp}
// 		}

// 		if len(cb) > 0 {
// 			ippem["ip_ranges"] = cb
// 		}

// 		var g []map[string]interface{}
// 		for _, v := range rule.UserIdGroupPairs {
// 			g = append(g, map[string]interface{}{
// 				"firewall_rules_set_name": *v.GroupName,
// 				"firewall_rules_set_id":   *v.GroupId,
// 				"account_id":              *v.UserId,
// 			})
// 		}

// 		if len(g) > 0 {
// 			ippem["groups"] = g
// 		}

// 		pl := make([]*fcu.PrefixListId, len(rule.PrefixListIds))
// 		for i, c := range rule.PrefixListIds {
// 			pl[i] = &fcu.PrefixListId{PrefixListId: c.PrefixListId}
// 		}

// 		if len(pl) > 0 {
// 			ippem["prefix_list_ids"] = pl
// 		}

// 		if len(rule.UserIdGroupPairs) > 0 {
// 			s := rule.UserIdGroupPairs[0]

// 			if isVPC {
// 				d.Set("destination_firewall_rules_set_account_id", *s.GroupId)
// 			} else {
// 				d.Set("destination_firewall_rules_set_name", *s.GroupName)
// 			}
// 		}

// 		ippems[i] = ippem
// 	}

// 	d.Set("inbound_rule", ippems)

// 	return nil
// }
