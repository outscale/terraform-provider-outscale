package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

// Security group import fans out to multiple resources due to the
// security group rules. Instead of creating one resource with nested
// rules, we use the best practices approach of one resource per rule.
func resourceOutscaleInboundImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*OutscaleClient).FCU

	// First query the security group
	sgRaw, _, err := SGRStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return nil, err
	}
	if sgRaw == nil {
		return nil, fmt.Errorf("security group not found")
	}
	sg := sgRaw.(*fcu.SecurityGroup)

	// Start building our results
	results := make([]*schema.ResourceData, 1, 1+len(sg.IpPermissions))
	results[0] = d

	// Construct the rules
	permMap := map[string][]*fcu.IpPermission{
		"ip_permissions": sg.IpPermissions,
	}
	for ruleType, perms := range permMap {
		for _, perm := range perms {
			ds, err := resourceOutscaleInboundImportStatePerm(sg, ruleType, perm)
			if err != nil {
				return nil, err
			}
			results = append(results, ds...)
		}
	}

	return results, nil
}

func resourceOutscaleInboundImportStatePerm(sg *fcu.SecurityGroup, ruleType string, perm *fcu.IpPermission) ([]*schema.ResourceData, error) {
	var result []*schema.ResourceData

	if perm.IpRanges != nil {
		p := &fcu.IpPermission{
			FromPort:      perm.FromPort,
			IpProtocol:    perm.IpProtocol,
			PrefixListIds: perm.PrefixListIds,
			ToPort:        perm.ToPort,
			IpRanges:      perm.IpRanges,
		}

		r, err := resourceOutscaleInboundImportStatePermPair(sg, ruleType, p)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	if len(perm.UserIdGroupPairs) > 0 {
		for _, pair := range perm.UserIdGroupPairs {
			p := &fcu.IpPermission{
				FromPort:         perm.FromPort,
				IpProtocol:       perm.IpProtocol,
				PrefixListIds:    perm.PrefixListIds,
				ToPort:           perm.ToPort,
				UserIdGroupPairs: []*fcu.UserIdGroupPair{pair},
			}

			r, err := resourceOutscaleInboundImportStatePermPair(sg, ruleType, p)
			if err != nil {
				return nil, err
			}
			result = append(result, r)
		}
	}
	return result, nil
}

func resourceOutscaleInboundImportStatePermPair(sg *fcu.SecurityGroup, ruleType string, perm *fcu.IpPermission) (*schema.ResourceData, error) {
	// Construct the rule. We do this by populating the absolute
	// minimum necessary for Refresh on the rule to work. This
	// happens to be a lot of fields since they're almost all needed
	// for de-dupping.
	sgId := sg.GroupId
	id := ipPermissionIDHash(*sgId, ruleType, perm)
	ruleResource := resourceOutscaleInboundRule()
	d := ruleResource.Data(nil)
	d.SetId(id)
	d.SetType("outscale_inbound_rule")
	d.Set("group_id", sgId)

	// 'self' is false by default. Below, we range over the group ids and set true
	// if the parent sg id is found

	if len(perm.UserIdGroupPairs) > 0 {
		s := perm.UserIdGroupPairs[0]

		// Check for Pair that is the same as the Security Group, to denote self.
		// Otherwise, mark the group id in source_security_group_id
		isVPC := sg.VpcId != nil && *sg.VpcId != ""
		if isVPC {
			if *s.GroupId == *sg.GroupId {
				// prune the self reference from the UserIdGroupPairs, so we don't
				// have duplicate sg ids (both self and in source_security_group_id)
				perm.UserIdGroupPairs = append(perm.UserIdGroupPairs[:0], perm.UserIdGroupPairs[0+1:]...)
			}
		} else {
			if *s.GroupName == *sg.GroupName {
				// prune the self reference from the UserIdGroupPairs, so we don't
				// have duplicate sg ids (both self and in source_security_group_id)
				perm.UserIdGroupPairs = append(perm.UserIdGroupPairs[:0], perm.UserIdGroupPairs[0+1:]...)
			}
		}
	}

	if err := setFromIPPerm(d, sg, perm); err != nil {
		return nil, errwrap.Wrapf("Error importing Outscale Security Group: {{err}}", err)
	}

	return d, nil
}

func SGRStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := &fcu.DescribeSecurityGroupsInput{
			Filters: []*fcu.Filter{&fcu.Filter{
				Name:   aws.String("ip-permission.group-id"),
				Values: []*string{aws.String(id)},
			}},
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

		group := resp.SecurityGroups[0].IpPermissions
		return group, "exists", nil
	}
}
