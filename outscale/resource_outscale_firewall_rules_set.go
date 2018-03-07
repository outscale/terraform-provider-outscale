package outscale

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleFirewallRulesSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleFirewallRulesSetCreate,
		Read:   resourceOutscaleFirewallRulesSetRead,
		Delete: resourceOutscaleFirewallRulesSetDelete,
		Importer: &schema.ResourceImporter{
			// State: resourceOutscaleSecurityGroupImportState,
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tag": tagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"group_description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port": {
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
			"ip_permissions_egress": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"to_port": {
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
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": {
				Type: schema.TypeMap,
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

func resourceOutscaleFirewallRulesSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	securityGroupOpts := &fcu.CreateSecurityGroupInput{}

	gn, gnok := d.GetOk("group_name")
	gd, gdok := d.GetOk("group_description")

	if gnok == false && gdok == false {
		return fmt.Errorf("group name and group description, are required attributes, and must be set")
	}

	if v, ok := d.GetOk("vpc_id"); ok {
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

	fmt.Printf("[INFO] Security Group ID: %s", d.Id())

	// Wait for the security group to truly exist
	fmt.Printf(
		"[DEBUG] Waiting for Security Group (%s) to exist",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists"},
		Refresh: SGStateRefreshFunc(conn, d.Id()),
		Timeout: 3 * time.Minute,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Security Group (%s) to become available: %s",
			d.Id(), err)
	}

	if err := setTags(conn, d); err != nil {
		return err
	}

	// defaults all Security Groups to have an ALLOW ALL egress rule. Here we
	// revoke that rule, so users don't unknowingly have/use it.
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

	return resourceOutscaleSecurityGroupUpdate(d, meta)
}

func resourceOutscaleFirewallRulesSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sgRaw, _, err := SGStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	sg := sgRaw.(*fcu.SecurityGroup)

	remoteIngressRules := resourceOutscaleSecurityGroupIPPermGather(d.Id(), sg.IpPermissions, sg.OwnerId)
	remoteEgressRules := resourceOutscaleSecurityGroupIPPermGather(d.Id(), sg.IpPermissionsEgress, sg.OwnerId)

	localIngressRules := d.Get("ip_permissions").(*schema.Set).List()
	localEgressRules := d.Get("ip_permissions_egress").(*schema.Set).List()

	// Loop through the local state of rules, doing a match against the remote
	// ruleSet we built above.
	ingressRules := matchRules("ingress", localIngressRules, remoteIngressRules)
	egressRules := matchRules("egress", localEgressRules, remoteEgressRules)

	d.Set("group_description", sg.Description)
	d.Set("group_name", sg.GroupName)
	d.Set("vpc_id", sg.VpcId)
	d.Set("owner_id", sg.OwnerId)

	if err := d.Set("ip_permissions", ingressRules); err != nil {
		fmt.Printf("[WARN] Error setting Ingress rule set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("ip_permissions_egress", egressRules); err != nil {
		fmt.Printf("[WARN] Error setting Egress rule set for (%s): %s", d.Id(), err)
	}

	d.Set("tag_set", tagsToMap(sg.Tags))

	return nil
}

func resourceOutscaleSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sgRaw, _, err := SGStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if sgRaw == nil {
		d.SetId("")
		return nil
	}

	group := sgRaw.(*fcu.SecurityGroup)

	err = resourceOutscaleSecurityGroupUpdateRules(d, "ingress", meta, group)
	if err != nil {
		return err
	}

	if d.Get("vpc_id") != nil {
		err = resourceOutscaleSecurityGroupUpdateRules(d, "egress", meta, group)
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

	return resourceOutscaleFirewallRulesSetRead(d, meta)
}

func resourceOutscaleSecurityGroupUpdateRules(
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

		remove, err := expandIPPerms(group, os.Difference(ns).List())
		if err != nil {
			return err
		}
		add, err := expandIPPerms(group, ns.Difference(os).List())
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

func resourceOutscaleFirewallRulesSetDelete(d *schema.ResourceData, meta interface{}) error {
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

func SGStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
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

func expandIPPerms(
	group *fcu.SecurityGroup, configured []interface{}) ([]*fcu.IpPermission, error) {
	vpc := group.VpcId != nil && *group.VpcId != ""

	perms := make([]*fcu.IpPermission, len(configured))
	for i, mRaw := range configured {
		var perm fcu.IpPermission
		m := mRaw.(map[string]interface{})

		perm.FromPort = aws.Int64(int64(m["from_port"].(int)))
		perm.ToPort = aws.Int64(int64(m["to_port"].(int)))
		perm.IpProtocol = aws.String(m["ip_protocol"].(string))

		if *perm.IpProtocol == "-1" && (*perm.FromPort != 0 || *perm.ToPort != 0) {
			return nil, fmt.Errorf(
				"from_port (%d) and to_port (%d) must both be 0 to use the 'ALL' \"-1\" protocol!",
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

func resourceOutscaleSecurityGroupIPPermGather(groupId string, permissions []*fcu.IpPermission, ownerId *string) []map[string]interface{} {
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

		m["from_port"] = fromPort
		m["to_port"] = toPort
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

		groups := flattenSecurityGroups(perm.UserIdGroupPairs, ownerId)
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

func flattenSecurityGroups(list []*fcu.UserIdGroupPair, ownerId *string) []*fcu.GroupIdentifier {
	result := make([]*fcu.GroupIdentifier, 0, len(list))
	for _, g := range list {
		var userId *string
		if g.UserId != nil && *g.UserId != "" && (ownerId == nil || *ownerId != *g.UserId) {
			userId = g.UserId
		}
		// userid nil here for same vpc groups

		vpc := g.GroupName == nil || *g.GroupName == ""
		var id *string
		if vpc {
			id = g.GroupId
		} else {
			id = g.GroupName
		}

		// id is groupid for vpcs
		// id is groupname for non vpc

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

func matchRules(rType string, local []interface{}, remote []map[string]interface{}) []map[string]interface{} {
	var saves []map[string]interface{}
	for _, raw := range local {
		l := raw.(map[string]interface{})

		// matching against self is required to detect rules that only include self
		// as the rule. resourceAwsSecurityGroupIPPermGather parses the group out
		// and replaces it with self if it's ID is found
		localHash := idHash(rType, l["ip_protocol"].(string), int64(l["to_port"].(int)), int64(l["from_port"].(int)))

		// loop remote rules, looking for a matching hash
		for _, r := range remote {

			// hash this remote rule and compare it for a match consideration with the
			// local rule we're examining
			rHash := idHash(rType, r["ip_protocol"].(string), r["to_port"].(int64), r["from_port"].(int64))
			if rHash == localHash {
				var numExpectedCidrs, numExpectedPrefixLists, numExpectedSGs, numRemoteCidrs, numRemotePrefixLists, numRemoteSGs int
				var matchingCidrs []string
				var matchingSGs []string
				var matchingPrefixLists []string

				// grab the local/remote cidr and sg groups, capturing the expected and
				// actual counts
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

				// check some early failures
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

				// match CIDRs by converting both to sets, and using Set methods
				var localCidrs []interface{}
				if lcRaw != nil {
					localCidrs = lcRaw.([]interface{})
				}
				localCidrSet := schema.NewSet(schema.HashString, localCidrs)

				// remote cidrs are presented as a slice of strings, so we need to
				// reformat them into a slice of interfaces to be used in creating the
				// remote cidr set
				var remoteCidrs []string
				if rcRaw != nil {
					remoteCidrs = rcRaw.([]string)
				}
				// convert remote cidrs to a set, for easy comparisons
				var list []interface{}
				for _, s := range remoteCidrs {
					list = append(list, s)
				}
				remoteCidrSet := schema.NewSet(schema.HashString, list)

				// Build up a list of local cidrs that are found in the remote set
				for _, s := range localCidrSet.List() {
					if remoteCidrSet.Contains(s) {
						matchingCidrs = append(matchingCidrs, s.(string))
					}
				}

				// match prefix lists by converting both to sets, and using Set methods
				var localPrefixLists []interface{}
				if lpRaw != nil {
					localPrefixLists = lpRaw.([]interface{})
				}
				localPrefixListsSet := schema.NewSet(schema.HashString, localPrefixLists)

				// remote prefix lists are presented as a slice of strings, so we need to
				// reformat them into a slice of interfaces to be used in creating the
				// remote prefix list set
				var remotePrefixLists []string
				if rpRaw != nil {
					remotePrefixLists = rpRaw.([]string)
				}
				// convert remote prefix lists to a set, for easy comparison
				list = nil
				for _, s := range remotePrefixLists {
					list = append(list, s)
				}
				remotePrefixListsSet := schema.NewSet(schema.HashString, list)

				// Build up a list of local prefix lists that are found in the remote set
				for _, s := range localPrefixListsSet.List() {
					if remotePrefixListsSet.Contains(s) {
						matchingPrefixLists = append(matchingPrefixLists, s.(string))
					}
				}

				// match SGs. Both local and remote are already sets
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

				// Build up a list of local security groups that are found in the remote set
				for _, s := range localSGSet.List() {
					if remoteSGSet.Contains(s) {
						matchingSGs = append(matchingSGs, s.(string))
					}
				}

				// compare equalities for matches.
				// If we found the number of cidrs and number of sgs, we declare a
				// match, and then remove those elements from the remote rule, so that
				// this remote rule can still be considered by other local rules
				if numExpectedCidrs == len(matchingCidrs) {
					if numExpectedPrefixLists == len(matchingPrefixLists) {
						if numExpectedSGs == len(matchingSGs) {

							// pop local cidrs from remote
							diffCidr := remoteCidrSet.Difference(localCidrSet)
							var newCidr []string
							for _, cRaw := range diffCidr.List() {
								newCidr = append(newCidr, cRaw.(string))
							}

							// reassigning
							if len(newCidr) > 0 {
								r["ip_ranges"] = newCidr
							} else {
								delete(r, "ip_ranges")
							}

							// pop local prefix lists from remote
							diffPrefixLists := remotePrefixListsSet.Difference(localPrefixListsSet)
							var newPrefixLists []string
							for _, pRaw := range diffPrefixLists.List() {
								newPrefixLists = append(newPrefixLists, pRaw.(string))
							}

							// reassigning
							if len(newPrefixLists) > 0 {
								r["prefix_list_ids"] = newPrefixLists
							} else {
								delete(r, "prefix_list_ids")
							}

							// pop local sgs from remote
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
	// Here we catch any remote rules that have not been stripped of all self,
	// cidrs, and security groups. We'll add remote rules here that have not been
	// matched locally, and let the graph sort things out. This will happen when
	// rules are added externally to Terraform
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

func idHash(rType, protocol string, toPort, fromPort int64) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", rType))
	buf.WriteString(fmt.Sprintf("%d-", toPort))
	buf.WriteString(fmt.Sprintf("%d-", fromPort))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(protocol)))

	return fmt.Sprintf("rule-%d", hashcode.String(buf.String()))
}

type ByGroupPair []*fcu.UserIdGroupPair

func ipPermissionIDHash(sg_id, ruleType string, ip *fcu.IpPermission) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", sg_id))
	if ip.FromPort != nil && *ip.FromPort > 0 {
		buf.WriteString(fmt.Sprintf("%d-", *ip.FromPort))
	}
	if ip.ToPort != nil && *ip.ToPort > 0 {
		buf.WriteString(fmt.Sprintf("%d-", *ip.ToPort))
	}
	buf.WriteString(fmt.Sprintf("%s-", *ip.IpProtocol))
	buf.WriteString(fmt.Sprintf("%s-", ruleType))

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

	return fmt.Sprintf("sgrule-%d", hashcode.String(buf.String()))
}

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

func resourceOutscaleSecurityGroupRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["from_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["to_port"].(int)))
	p := protocolForValue(m["ip_protocol"].(string))
	buf.WriteString(fmt.Sprintf("%s-", p))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["ip_ranges"]; ok {
		vs := v.([]interface{})
		s := make([]map[string]interface{}, len(vs))
		for i, raw := range vs {
			s[i] = raw.(map[string]interface{})
		}
		// sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["prefix_list_ids"]; ok {
		vs := v.([]interface{})
		s := make([]map[string]interface{}, len(vs))
		for i, raw := range vs {
			s[i] = raw.(map[string]interface{})
		}
		// sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["groups"]; ok {
		vs := v.(map[string]interface{})
		s := make([]map[string]interface{}, len(vs))
		for _, raw := range vs {
			s = append(s, raw.(map[string]interface{}))
		}
		// sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	return hashcode.String(buf.String())
}
