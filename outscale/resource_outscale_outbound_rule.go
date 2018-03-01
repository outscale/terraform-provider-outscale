package outscale

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOutboundRuleCreate,
		Read:   resourceOutscaleOutboundRuleRead,
		Delete: resourceOutscaleOutboundRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				Computed: true,
				ForceNew: true,
			},
			"source_security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"source_security_group_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_permissions": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"groups": {
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"to_port": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateCIDRNetworkAddress,
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
			},
		},
	}
}

var awsMutexKV = mutexkv.NewMutexKV()

func resourceOutscaleOutboundRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)

	awsMutexKV.Lock(sg_id)
	defer awsMutexKV.Unlock(sg_id)

	sg, err := findResourceSecurityGroup(conn, sg_id)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 1 =>", err)

	if err != nil {
		return err
	}

	perm, err := expandIPPerm(d, sg)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 2 =>", err)

	if err != nil {
		return err
	}

	if err := validateOutscaleSecurityGroupRule(d); err != nil {
		fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 3 =>", err)

		return err
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""
	ruleType := "egress"
	var autherr error
	fmt.Printf("\n\n[DEBUG] Authorizing security group %s %s rule: %#v", sg_id, "Egress", perm)

	req := &fcu.AuthorizeSecurityGroupEgressInput{
		GroupId:       sg.GroupId,
		IpPermissions: []*fcu.IpPermission{perm},
	}

	_, autherr = conn.VM.AuthorizeSecurityGroupEgress(req)

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 4 =>", err)

	if autherr != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidPermission.Duplicate") {
			return fmt.Errorf(`[WARN] A duplicate Security Group rule was found on (%s). This may be
a side effect of a now-fixed Terraform issue causing two security groups with
identical attributes but different source_security_group_ids to overwrite each
other in the state. See https://github.com/hashicorp/terraform/pull/2376 for more
information and instructions for recovery. Error message: %s`, sg_id, fmt.Sprint(err))

		}

		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			ruleType, autherr)
	}

	id := ipPermissionIDHash(sg_id, ruleType, perm)
	fmt.Printf("\n\n[DEBUG] Computed group rule ID %s", id)

	retErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		sg, err := findResourceSecurityGroup(conn, sg_id)

		if err != nil {
			fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 6 =>", err)

			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("\n\n[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}

			fmt.Printf("\n\n[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sg_id, id, err)
			return resource.NonRetryableError(err)
		}

		var rules []*fcu.IpPermission
		rules = sg.IpPermissionsEgress

		rule := findRuleMatch(perm, rules, isVPC)

		if rule == nil {
			fmt.Printf("\n\n[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
				ruleType, id, sg_id)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		fmt.Printf("\n\n[DEBUG] Found rule for Security Group Rule (%s): %s", id, rule)
		return nil
	})

	fmt.Println("\n[DEDUG] ERROR resourceOutscaleOutboundRuleCreate 7 =>", retErr)

	if retErr != nil {
		return fmt.Errorf("Error finding matching %s Security Group Rule (%s) for Group %s",
			ruleType, id, sg_id)
	}

	d.SetId(id)
	return nil
}

func resourceOutscaleOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)
	sg, err := findResourceSecurityGroup(conn, sg_id)
	if _, notFound := err.(securityGroupNotFound); notFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sg_id, d.Id(), err)
	}

	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	var rule *fcu.IpPermission
	var rules []*fcu.IpPermission
	ruleType := "egress"
	rules = sg.IpPermissionsEgress

	p, err := expandIPPerm(d, sg)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		fmt.Printf("\n\n[WARN] No %s rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			ruleType, *sg.GroupName, d.Id())
		d.SetId("")
		return nil
	}

	rule = findRuleMatch(p, rules, isVPC)

	if rule == nil {
		fmt.Printf("\n\n[DEBUG] Unable to find matching %s Security Group Rule (%s) for Group %s",
			ruleType, d.Id(), sg_id)
		d.SetId("")
		return nil
	}

	fmt.Printf("\n\n[DEBUG] Found rule for Security Group Rule (%s): %s", d.Id(), rule)

	if err := setFromIPPerm(d, sg, p); err != nil {
		return errwrap.Wrapf("Error setting IP Permission for Security Group Rule: {{err}}", err)
	}

	return nil
}

func resourceOutscaleOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sg_id := d.Get("group_id").(string)

	awsMutexKV.Lock(sg_id)
	defer awsMutexKV.Unlock(sg_id)

	sg, err := findResourceSecurityGroup(conn, sg_id)
	if err != nil {
		return err
	}

	perm, err := expandIPPerm(d, sg)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n[DEBUG] Revoking security group %#v %s rule: %#v",
		sg_id, "egress", perm)
	req := &fcu.RevokeSecurityGroupEgressInput{
		GroupId:       sg.GroupId,
		IpPermissions: []*fcu.IpPermission{perm},
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.RevokeSecurityGroupEgress(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("\n\n[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf(
			"Error revoking security group %s rules: %s",
			sg_id, err)
	}

	d.SetId("")

	return nil
}

func findResourceSecurityGroup(conn *fcu.Client, id string) (*fcu.SecurityGroup, error) {
	req := &fcu.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeSecurityGroupsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSecurityGroups(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("\n\n[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
		return nil, securityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, securityGroupNotFound{id, nil}
	}
	if len(resp.SecurityGroups) != 1 || resp.SecurityGroups[0] == nil {
		return nil, securityGroupNotFound{id, resp.SecurityGroups}
	}

	return resp.SecurityGroups[0], nil
}

func expandIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup) (*fcu.IpPermission, error) {
	var perm fcu.IpPermission

	if raw, ok := d.GetOk("ip_permissions"); ok {

		fmt.Printf("\n\nRAW #v", raw.(map[string]interface{}))

		ipp := raw.(map[string]interface{})

		fp, fpe := strconv.Atoi(ipp["from_port"].(string))
		if fpe != nil {
			return nil, fpe
		}
		tp, tpe := strconv.Atoi(ipp["to_port"].(string))
		if tpe != nil {
			return nil, tpe
		}

		perm.FromPort = aws.Int64(int64(fp))
		perm.ToPort = aws.Int64(int64(tp))
		protocol := protocolForValue(ipp["ip_protocol"].(string))
		perm.IpProtocol = aws.String(protocol)

		groups := make(map[string]bool)
		if raw, ok := ipp["groups"]; ok {
			groups[raw.(string)] = true
		}

		if len(groups) > 0 {
			perm.UserIdGroupPairs = make([]*fcu.UserIdGroupPair, len(groups))
			var gl []string
			for k, _ := range groups {
				gl = append(gl, k)
			}

			for i, name := range gl {
				ownerId, id := "", name
				if items := strings.Split(id, "/"); len(items) > 1 {
					ownerId, id = items[0], items[1]
				}

				perm.UserIdGroupPairs[i] = &fcu.UserIdGroupPair{
					GroupId: aws.String(id),
					UserId:  aws.String(ownerId),
				}

				if sg.VpcId == nil || *sg.VpcId == "" {
					perm.UserIdGroupPairs[i].GroupId = nil
					perm.UserIdGroupPairs[i].GroupName = aws.String(id)
					perm.UserIdGroupPairs[i].UserId = nil
				}
			}
		}

		fmt.Printf("[DEBUG] => IP_RANGES #v =>", ipp["ip_ranges"])

		if raw, ok := ipp["ip_ranges"]; ok {
			// list := raw.([]interface{})
			perm.IpRanges = make([]*fcu.IpRange, 1)
			// for i, v := range list {
			cidrIP, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("empty element found in cidr_blocks - consider using the compact function")
			}
			perm.IpRanges[0] = &fcu.IpRange{CidrIp: aws.String(cidrIP)}
			// }
		}

		if raw, ok := ipp["prefix_list_ids"]; ok {
			list := raw.([]interface{})
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

	return &perm, nil
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

	fmt.Printf("\n\n[WARN] Unable to determine valid protocol: no matching protocols found")
	return protocol
}

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

func validateOutscaleSecurityGroupRule(d *schema.ResourceData) error {
	if ipp, ippemOk := d.GetOk("ip_permissions"); ippemOk {
		ippem := ipp.(map[string]interface{})

		_, blocksOk := ippem["ip_ranges"].(string)
		_, sourceOk := d.GetOk("source_security_group_owner_id")
		_, prefixOk := ippem["prefix_list_ids"].(string)
		if !blocksOk && !sourceOk && !prefixOk {
			return fmt.Errorf(
				"One of ['ip_ranges', 'source_security_group_owner_id', 'prefix_list_ids'] must be set to create an Outscale Security Group Rule")
		}
	}

	return nil
}

func findRuleMatch(p *fcu.IpPermission, rules []*fcu.IpPermission, isVPC bool) *fcu.IpPermission {
	var rule *fcu.IpPermission
	for _, r := range rules {
		if r.ToPort != nil && *p.ToPort != *r.ToPort {
			continue
		}

		if r.FromPort != nil && *p.FromPort != *r.FromPort {
			continue
		}

		if r.IpProtocol != nil && *p.IpProtocol != *r.IpProtocol {
			continue
		}

		remaining := len(p.IpRanges)
		for _, ip := range p.IpRanges {
			for _, rip := range r.IpRanges {
				if *ip.CidrIp == *rip.CidrIp {
					remaining--
				}
			}
		}

		if remaining > 0 {
			continue
		}

		remaining = len(p.PrefixListIds)
		for _, pl := range p.PrefixListIds {
			for _, rpl := range r.PrefixListIds {
				if *pl.PrefixListId == *rpl.PrefixListId {
					remaining--
				}
			}
		}

		if remaining > 0 {
			continue
		}

		remaining = len(p.UserIdGroupPairs)
		for _, ip := range p.UserIdGroupPairs {
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
	return rule
}

func setFromIPPerm(d *schema.ResourceData, sg *fcu.SecurityGroup, rule *fcu.IpPermission) error {
	isVPC := sg.VpcId != nil && *sg.VpcId != ""

	ippem := make(map[string]interface{})
	ippem["from_port"] = rule.FromPort
	ippem["to_port"] = rule.ToPort
	ippem["ip_protocol"] = rule.IpProtocol

	var cb []string
	for _, c := range rule.IpRanges {
		cb = append(cb, *c.CidrIp)
	}

	if len(cb) > 0 {
		ippem["ip_ranges"] = cb
		d.Set("cidr_ip", cb[0])
	}

	var g []map[string]interface{}
	for _, v := range rule.UserIdGroupPairs {
		g = append(g, map[string]interface{}{
			"group_name": v.GroupName,
			"group_id":   v.GroupId,
			"user_id":    v.UserId,
		})
	}

	ippem["groups"] = g

	var pl []string
	for _, p := range rule.PrefixListIds {
		pl = append(pl, *p.PrefixListId)
	}

	ippem["prefix_list_ids"] = pl

	if len(rule.UserIdGroupPairs) > 0 {
		s := rule.UserIdGroupPairs[0]

		if isVPC {
			d.Set("source_security_group_owner_id", *s.GroupId)
		} else {
			d.Set("source_security_group_name", *s.GroupName)
		}
	}

	d.Set("from_port", rule.FromPort)
	d.Set("ip_protocol", rule.IpProtocol)
	d.Set("to_port", rule.ToPort)
	d.Set("ip_permissions", ippem)

	return nil
}

func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}
