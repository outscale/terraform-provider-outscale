package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIFirewallRulesSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIFirewallRulesSetCreate,
		Read:   resourceOutscaleOAPIFirewallRulesSetRead,
		Delete: resourceOutscaleOAPIFirewallRulesSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tag": tagsSchema(),
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_rules_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lin_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_rules_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inbound_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateCIDRNetworkAddress,
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"outbound_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateCIDRNetworkAddress,
							},
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIFirewallRulesSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	securityGroupOpts := &fcu.CreateSecurityGroupInput{}

	gn, gnok := d.GetOk("firewall_rules_set_name")
	gd, gdok := d.GetOk("description")

	if gnok == false && gdok == false {
		return fmt.Errorf("group name and group description, are required attributes, and must be set")
	}

	if v, ok := d.GetOk("lin_id"); ok {
		securityGroupOpts.VpcId = aws.String(v.(string))
	}

	securityGroupOpts.GroupName = aws.String(gn.(string))
	securityGroupOpts.Description = aws.String(gd.(string))

	fmt.Printf(
		"[DEBUG] Security Group create configuration: %#v", securityGroupOpts)

	var err error
	var createResp *fcu.CreateSecurityGroupOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		createResp, err = conn.VM.CreateSecurityGroup(securityGroupOpts)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Security Group: %s", err)
	}

	d.SetId(*createResp.GroupId)
	d.Set("firewall_rules_set_id", *createResp.GroupId)

	fmt.Printf("[INFO] Security Group ID: %s", d.Id())

	fmt.Printf(
		"[DEBUG] Waiting for Security Group (%s) to exist",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists"},
		Refresh: SGOAPIStateRefreshFunc(conn, d.Id()),
		Timeout: 3 * time.Minute,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Security Group (%s) to become available: %s",
			d.Id(), err)
	}

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	group := resp.(*fcu.SecurityGroup)
	if group.VpcId != nil && *group.VpcId != "" {
		fmt.Printf("[DEBUG] Revoking default egress rule for Security Group for %s", d.Id())

		req := &fcu.RevokeSecurityGroupEgressInput{
			GroupId: createResp.GroupId,
			IpPermissions: []*fcu.IpPermission{
				{
					FromPort: aws.Int64(int64(0)),
					ToPort:   aws.Int64(int64(0)),
					IpRanges: []*fcu.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
					IpProtocol: aws.String("-1"),
				},
			},
		}

		if _, err = conn.VM.RevokeSecurityGroupEgress(req); err != nil {
			return fmt.Errorf(
				"Error revoking default egress rule for Security Group (%s): %s",
				d.Id(), err)
		}

	}

	return resourceOutscaleOAPISecurityGroupUpdate(d, meta)
}

func resourceOutscaleOAPIFirewallRulesSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sgRaw, _, err := SGOAPIStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	sg := sgRaw.(*fcu.SecurityGroup)

	remoteIngressRules := resourceOutscaleOAPISecurityGroupIPPermGather(d.Id(), sg.IpPermissions, sg.OwnerId)
	remoteEgressRules := resourceOutscaleOAPISecurityGroupIPPermGather(d.Id(), sg.IpPermissionsEgress, sg.OwnerId)

	localIngressRules := d.Get("inbound_rules").(*schema.Set).List()
	localEgressRules := d.Get("outbound_rules").(*schema.Set).List()

	ingressRules := matchOAPIRules("ingress", localIngressRules, remoteIngressRules)
	egressRules := matchOAPIRules("egress", localEgressRules, remoteEgressRules)

	d.Set("description", sg.Description)
	d.Set("firewall_rules_set_name", sg.GroupName)
	d.Set("lin_id", sg.VpcId)
	d.Set("account_id", sg.OwnerId)

	if err := d.Set("inbound_rules", ingressRules); err != nil {
		fmt.Printf("[WARN] Error setting Ingress rule set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("outbound_rules", egressRules); err != nil {
		fmt.Printf("[WARN] Error setting Egress rule set for (%s): %s", d.Id(), err)
	}

	if sg.Tags != nil {
		if err := d.Set("tags", tagsToMap(sg.Tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set("tags", []map[string]string{
			map[string]string{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceOutscaleOAPISecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sgRaw, _, err := SGOAPIStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	group := sgRaw.(*fcu.SecurityGroup)

	err = resourceOutscaleOAPISecurityGroupUpdateRules(d, "ingress", meta, group)
	if err != nil {
		return err
	}

	if d.Get("lin_id") != nil {
		err = resourceOutscaleOAPISecurityGroupUpdateRules(d, "egress", meta, group)
		if err != nil {
			return err
		}
	}

	if !d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	return resourceOutscaleOAPIFirewallRulesSetRead(d, meta)
}

func resourceOutscaleOAPISecurityGroupUpdateRules(
	d *schema.ResourceData, ruleset string,
	meta interface{}, group *fcu.SecurityGroup) error {

	if d.HasChange(ruleset) {
		o, n := d.GetChange(ruleset)
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove, err := expandOAPIIPPerms(group, os.Difference(ns).List())
		if err != nil {
			return err
		}
		add, err := expandOAPIIPPerms(group, ns.Difference(os).List())
		if err != nil {
			return err
		}

		if len(remove) > 0 || len(add) > 0 {
			conn := meta.(*OutscaleClient).FCU

			var err error
			if len(remove) > 0 {
				fmt.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
					group, ruleset, remove)

				if ruleset == "egress" {
					req := &fcu.RevokeSecurityGroupEgressInput{
						GroupId:       group.GroupId,
						IpPermissions: remove,
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

				} else {
					req := &fcu.RevokeSecurityGroupIngressInput{
						GroupId:       group.GroupId,
						IpPermissions: remove,
					}
					if group.VpcId == nil || *group.VpcId == "" {
						req.GroupId = nil
						req.GroupName = group.GroupName
					}

					err = resource.Retry(5*time.Minute, func() *resource.RetryError {
						_, err = conn.VM.RevokeSecurityGroupIngress(req)

						if err != nil {
							if strings.Contains(err.Error(), "RequestLimitExceeded") {
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}

						return nil
					})
				}

				if err != nil {
					return fmt.Errorf(
						"Error revoking security group %s rules: %s",
						ruleset, err)
				}
			}

			if len(add) > 0 {
				fmt.Printf("[DEBUG] Authorizing security group %#v %s rule: %#v",
					group, ruleset, add)
				// Authorize the new rules
				if ruleset == "egress" {
					req := &fcu.AuthorizeSecurityGroupEgressInput{
						GroupId:       group.GroupId,
						IpPermissions: add,
					}

					err = resource.Retry(5*time.Minute, func() *resource.RetryError {
						_, err = conn.VM.AuthorizeSecurityGroupEgress(req)

						if err != nil {
							if strings.Contains(err.Error(), "RequestLimitExceeded") {
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}

						return nil
					})
				} else {
					req := &fcu.AuthorizeSecurityGroupIngressInput{
						GroupId:       group.GroupId,
						IpPermissions: add,
					}
					if group.VpcId == nil || *group.VpcId == "" {
						req.GroupId = nil
						req.GroupName = group.GroupName
					}

					err = resource.Retry(5*time.Minute, func() *resource.RetryError {
						_, err = conn.VM.AuthorizeSecurityGroupIngress(req)

						if err != nil {
							if strings.Contains(err.Error(), "RequestLimitExceeded") {
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}

						return nil
					})
				}

				if err != nil {
					return fmt.Errorf(
						"Error authorizing security group %s rules: %s",
						ruleset, err)
				}
			}
		}
	}
	return nil
}

func resourceOutscaleOAPIFirewallRulesSetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	fmt.Printf("\n[DEBUG] Security Group destroy: %v", d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.DeleteSecurityGroup(&fcu.DeleteSecurityGroupInput{
			GroupId: aws.String(d.Id()),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") || strings.Contains(err.Error(), "DependencyViolation") {
				return resource.RetryableError(err)
			} else if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func SGOAPIStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(id)},
		}

		var resp *fcu.DescribeSecurityGroupsOutput
		var err error
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

		if err != nil {
			if strings.Contains(err.Error(), "InvalidSecurityGroupID.NotFound") || strings.Contains(err.Error(), "InvalidGroup.NotFound") {
				resp = nil
				err = nil
			}

			if err != nil {
				fmt.Printf("\nError on SGStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		group := resp.SecurityGroups[0]
		return group, "exists", nil
	}
}

func expandOAPIIPPerms(
	group *fcu.SecurityGroup, configured []interface{}) ([]*fcu.IpPermission, error) {
	vpc := group.VpcId != nil && *group.VpcId != ""

	perms := make([]*fcu.IpPermission, len(configured))
	for i, mRaw := range configured {
		var perm fcu.IpPermission
		m := mRaw.(map[string]interface{})

		perm.FromPort = aws.Int64(int64(m["from_port_range"].(int)))
		perm.ToPort = aws.Int64(int64(m["to_port_range"].(int)))
		perm.IpProtocol = aws.String(m["ip_protocol"].(string))

		if *perm.IpProtocol == "-1" && (*perm.FromPort != 0 || *perm.ToPort != 0) {
			return nil, fmt.Errorf(
				"from_port_range (%d) and to_port_range (%d) must both be 0 to use the 'ALL' \"-1\" protocol!",
				*perm.FromPort, *perm.ToPort)
		}

		var groups []string
		if raw, ok := m["groups"]; ok {
			list := raw.(*schema.Set).List()
			for _, v := range list {
				groups = append(groups, v.(string))
			}
		}

		if len(groups) > 0 {
			perm.UserIdGroupPairs = make([]*fcu.UserIdGroupPair, len(groups))
			for i, name := range groups {
				ownerId, id := "", name
				if items := strings.Split(id, "/"); len(items) > 1 {
					ownerId, id = items[0], items[1]
				}

				perm.UserIdGroupPairs[i] = &fcu.UserIdGroupPair{
					GroupId: aws.String(id),
				}

				if ownerId != "" {
					perm.UserIdGroupPairs[i].UserId = aws.String(ownerId)
				}

				if !vpc {
					perm.UserIdGroupPairs[i].GroupId = nil
					perm.UserIdGroupPairs[i].GroupName = aws.String(id)
				}
			}
		}

		if raw, ok := m["ip_ranges"]; ok {
			list := raw.([]interface{})
			for _, v := range list {
				perm.IpRanges = append(perm.IpRanges, &fcu.IpRange{CidrIp: aws.String(v.(string))})
			}
		}
		if raw, ok := m["prefix_list_ids"]; ok {
			list := raw.([]interface{})
			for _, v := range list {
				perm.PrefixListIds = append(perm.PrefixListIds, &fcu.PrefixListId{PrefixListId: aws.String(v.(string))})
			}
		}

		perms[i] = &perm
	}

	return perms, nil

}

func resourceOutscaleOAPISecurityGroupIPPermGather(groupId string, permissions []*fcu.IpPermission, ownerId *string) []map[string]interface{} {
	ruleMap := make(map[string]map[string]interface{})
	for _, perm := range permissions {
		var fromPort, toPort int64
		if v := perm.FromPort; v != nil {
			fromPort = *v
		}
		if v := perm.ToPort; v != nil {
			toPort = *v
		}

		k := fmt.Sprintf("%s-%d-%d", *perm.IpProtocol, fromPort, toPort)
		m, ok := ruleMap[k]
		if !ok {
			m = make(map[string]interface{})
			ruleMap[k] = m
		}

		m["from_port_range"] = fromPort
		m["to_port_range"] = toPort
		m["ip_protocol"] = *perm.IpProtocol

		if len(perm.IpRanges) > 0 {
			raw, ok := m["ip_ranges"]
			if !ok {
				raw = make([]string, 0, len(perm.IpRanges))
			}
			list := raw.([]string)

			for _, ip := range perm.IpRanges {
				list = append(list, *ip.CidrIp)
			}

			m["ip_ranges"] = list
		}

		if len(perm.PrefixListIds) > 0 {
			raw, ok := m["prefix_list_ids"]
			if !ok {
				raw = make([]string, 0, len(perm.PrefixListIds))
			}
			list := raw.([]string)

			for _, pl := range perm.PrefixListIds {
				list = append(list, *pl.PrefixListId)
			}

			m["prefix_list_ids"] = list
		}

		groups := flattenOAPISecurityGroups(perm.UserIdGroupPairs, ownerId)
		for i, g := range groups {
			if *g.GroupId == groupId {
				groups[i], groups = groups[len(groups)-1], groups[:len(groups)-1]
			}
		}

		if len(groups) > 0 {
			raw, ok := m["groups"]
			if !ok {
				raw = schema.NewSet(schema.HashString, nil)
			}
			list := raw.(*schema.Set)

			for _, g := range groups {
				if g.GroupName != nil {
					list.Add(*g.GroupName)
				} else {
					list.Add(*g.GroupId)
				}
			}

			m["groups"] = list
		}
	}
	rules := make([]map[string]interface{}, 0, len(ruleMap))
	for _, m := range ruleMap {
		rules = append(rules, m)
	}

	return rules
}

func flattenOAPISecurityGroups(list []*fcu.UserIdGroupPair, ownerId *string) []*fcu.GroupIdentifier {
	result := make([]*fcu.GroupIdentifier, 0, len(list))
	for _, g := range list {
		var userId *string
		if g.UserId != nil && *g.UserId != "" && (ownerId == nil || *ownerId != *g.UserId) {
			userId = g.UserId
		}

		vpc := g.GroupName == nil || *g.GroupName == ""
		var id *string
		if vpc {
			id = g.GroupId
		} else {
			id = g.GroupName
		}

		gi := &fcu.GroupIdentifier{}

		if userId != nil {
			id = aws.String(*userId + "/" + *id)
		}

		if vpc {
			gi.GroupId = id
			result = append(result, gi)
		} else {
			gi.GroupId = g.GroupId
			gi.GroupName = id
			result = append(result, gi)
		}
	}
	return result
}

func matchOAPIRules(rType string, local []interface{}, remote []map[string]interface{}) []map[string]interface{} {
	var saves []map[string]interface{}
	for _, raw := range local {
		l := raw.(map[string]interface{})

		localHash := idHash(rType, l["ip_protocol"].(string), int64(l["to_port_range"].(int)), int64(l["from_port_range"].(int)), l["self"].(bool))

		for _, r := range remote {

			rHash := idHash(rType, r["ip_protocol"].(string), r["to_port_range"].(int64), r["from_port_range"].(int64), r["self"].(bool))
			if rHash == localHash {
				var numExpectedCidrs, numExpectedPrefixLists, numExpectedSGs, numRemoteCidrs, numRemotePrefixLists, numRemoteSGs int
				var matchingCidrs []string
				var matchingSGs []string
				var matchingPrefixLists []string

				lcRaw, ok := l["ip_ranges"]
				if ok {
					numExpectedCidrs = len(l["ip_ranges"].([]interface{}))
				}
				lpRaw, ok := l["prefix_list_ids"]
				if ok {
					numExpectedPrefixLists = len(l["prefix_list_ids"].([]interface{}))
				}
				lsRaw, ok := l["groups"]
				if ok {
					numExpectedSGs = len(l["groups"].(*schema.Set).List())
				}
				rcRaw, ok := r["ip_ranges"]
				if ok {
					numRemoteCidrs = len(r["ip_ranges"].([]string))
				}
				rpRaw, ok := r["prefix_list_ids"]
				if ok {
					numRemotePrefixLists = len(r["prefix_list_ids"].([]string))
				}
				rsRaw, ok := r["groups"]
				if ok {
					numRemoteSGs = len(r["groups"].(*schema.Set).List())
				}

				if numExpectedCidrs > numRemoteCidrs {
					fmt.Printf("[DEBUG] Local rule has more CIDR blocks, continuing (%d/%d)", numExpectedCidrs, numRemoteCidrs)
					continue
				}
				if numExpectedPrefixLists > numRemotePrefixLists {
					fmt.Printf("[DEBUG] Local rule has more prefix lists, continuing (%d/%d)", numExpectedPrefixLists, numRemotePrefixLists)
					continue
				}
				if numExpectedSGs > numRemoteSGs {
					fmt.Printf("[DEBUG] Local rule has more Security Groups, continuing (%d/%d)", numExpectedSGs, numRemoteSGs)
					continue
				}

				var localCidrs []interface{}
				if lcRaw != nil {
					localCidrs = lcRaw.([]interface{})
				}
				localCidrSet := schema.NewSet(schema.HashString, localCidrs)

				var remoteCidrs []string
				if rcRaw != nil {
					remoteCidrs = rcRaw.([]string)
				}
				var list []interface{}
				for _, s := range remoteCidrs {
					list = append(list, s)
				}
				remoteCidrSet := schema.NewSet(schema.HashString, list)

				for _, s := range localCidrSet.List() {
					if remoteCidrSet.Contains(s) {
						matchingCidrs = append(matchingCidrs, s.(string))
					}
				}

				var localPrefixLists []interface{}
				if lpRaw != nil {
					localPrefixLists = lpRaw.([]interface{})
				}
				localPrefixListsSet := schema.NewSet(schema.HashString, localPrefixLists)

				var remotePrefixLists []string
				if rpRaw != nil {
					remotePrefixLists = rpRaw.([]string)
				}
				list = nil
				for _, s := range remotePrefixLists {
					list = append(list, s)
				}
				remotePrefixListsSet := schema.NewSet(schema.HashString, list)

				for _, s := range localPrefixListsSet.List() {
					if remotePrefixListsSet.Contains(s) {
						matchingPrefixLists = append(matchingPrefixLists, s.(string))
					}
				}

				var localSGSet *schema.Set
				if lsRaw == nil {
					localSGSet = schema.NewSet(schema.HashString, nil)
				} else {
					localSGSet = lsRaw.(*schema.Set)
				}

				var remoteSGSet *schema.Set
				if rsRaw == nil {
					remoteSGSet = schema.NewSet(schema.HashString, nil)
				} else {
					remoteSGSet = rsRaw.(*schema.Set)
				}

				for _, s := range localSGSet.List() {
					if remoteSGSet.Contains(s) {
						matchingSGs = append(matchingSGs, s.(string))
					}
				}

				if numExpectedCidrs == len(matchingCidrs) {
					if numExpectedPrefixLists == len(matchingPrefixLists) {
						if numExpectedSGs == len(matchingSGs) {

							diffCidr := remoteCidrSet.Difference(localCidrSet)
							var newCidr []string
							for _, cRaw := range diffCidr.List() {
								newCidr = append(newCidr, cRaw.(string))
							}

							if len(newCidr) > 0 {
								r["ip_ranges"] = newCidr
							} else {
								delete(r, "ip_ranges")
							}

							diffPrefixLists := remotePrefixListsSet.Difference(localPrefixListsSet)
							var newPrefixLists []string
							for _, pRaw := range diffPrefixLists.List() {
								newPrefixLists = append(newPrefixLists, pRaw.(string))
							}

							if len(newPrefixLists) > 0 {
								r["prefix_list_ids"] = newPrefixLists
							} else {
								delete(r, "prefix_list_ids")
							}

							diffSGs := remoteSGSet.Difference(localSGSet)
							if len(diffSGs.List()) > 0 {
								r["groups"] = diffSGs
							} else {
								delete(r, "groups")
							}

							saves = append(saves, l)
						}
					}

				}
			}
		}
	}

	for _, r := range remote {
		var lenCidr, lenPrefixLists, lenSGs int
		if rCidrs, ok := r["ip_ranges"]; ok {
			lenCidr = len(rCidrs.([]string))
		}
		if rPrefixLists, ok := r["prefix_list_ids"]; ok {
			lenPrefixLists = len(rPrefixLists.([]string))
		}
		if rawSGs, ok := r["groups"]; ok {
			lenSGs = len(rawSGs.(*schema.Set).List())
		}

		if lenSGs+lenCidr+lenPrefixLists > 0 {
			fmt.Printf("[DEBUG] Found a remote Rule that wasn't empty: (%#v)", r)
			saves = append(saves, r)
		}
	}

	return saves
}
