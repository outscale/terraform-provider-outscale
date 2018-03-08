package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleInboundRule(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var group fcu.SecurityGroup
	rInt := acctest.RandInt()

	testRuleCount := func(*terraform.State) error {
		if len(group.IpPermissions) != 2 {
			return fmt.Errorf("Wrong Security Group rule count, expected %d, got %d",
				1, len(group.IpPermissions))
		}

		rule := group.IpPermissions[0]
		if *rule.FromPort != int64(22) {
			return fmt.Errorf("Wrong Security Group port setting, expected %d, got %d",
				22, int(*rule.FromPort))
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleIngressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRuleExists("outscale_firewall_rules_set.outscale_firewall_rules_set", &group),
					testAccCheckOutscaleRuleAttributes("outscale_inbound_rule.outscale_inbound_rule1", &group, nil, "ingress"),
					testRuleCount,
				),
			},
		},
	})
}

func testAccCheckOutscaleSecurityGroupRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_firewall_rules_set" {
			continue
		}

		// Retrieve our group
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
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
		if err == nil {
			if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		if strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
			return nil
		}

		return err
	}

	return nil
}

func testAccCheckOutscaleRuleExists(n string, group *fcu.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		req := &fcu.DescribeSecurityGroupsInput{
			GroupIds: []*string{aws.String(rs.Primary.ID)},
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

		if err != nil {
			return err
		}

		if len(resp.SecurityGroups) > 0 && *resp.SecurityGroups[0].GroupId == rs.Primary.ID {
			*group = *resp.SecurityGroups[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccCheckOutscaleRuleAttributes(n string, group *fcu.SecurityGroup, p []*fcu.IpPermission, ruleType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Security Group Rule Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group Rule is set")
		}

		if p == nil {
			p = []*fcu.IpPermission{
				&fcu.IpPermission{
					FromPort:   aws.Int64(22),
					ToPort:     aws.Int64(22),
					IpProtocol: aws.String("tcp"),
					IpRanges:   []*fcu.IpRange{{CidrIp: aws.String("46.231.147.8/32")}},
				},
				&fcu.IpPermission{
					FromPort:   aws.Int64(443),
					ToPort:     aws.Int64(443),
					IpProtocol: aws.String("tcp"),
					IpRanges:   []*fcu.IpRange{{CidrIp: aws.String("46.231.147.8/32")}},
				},
			}
		}

		var matchingRule *fcu.IpPermission
		var rules []*fcu.IpPermission
		if ruleType == "ingress" {
			rules = group.IpPermissions
		} else {
			rules = group.IpPermissionsEgress
		}

		if len(rules) == 0 {
			return fmt.Errorf("No IPPerms")
		}

		for i, r := range rules {
			if r.ToPort != nil && *p[i].ToPort != *r.ToPort {
				continue
			}

			if r.FromPort != nil && *p[i].FromPort != *r.FromPort {
				continue
			}

			if r.IpProtocol != nil && *p[i].IpProtocol != *r.IpProtocol {
				continue
			}

			remaining := len(p[i].IpRanges)
			for _, ip := range p[i].IpRanges {
				for _, rip := range r.IpRanges {
					if *ip.CidrIp == *rip.CidrIp {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p[i].UserIdGroupPairs)
			for _, ip := range p[i].UserIdGroupPairs {
				for _, rip := range r.UserIdGroupPairs {
					if *ip.GroupId == *rip.GroupId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p[i].PrefixListIds)
			for _, pip := range p[i].PrefixListIds {
				for _, rpip := range r.PrefixListIds {
					if *pip.PrefixListId == *rpip.PrefixListId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			fmt.Printf("RULES 2 %s", r)

			matchingRule = r
		}

		if matchingRule != nil {
			log.Printf("[DEBUG] Matching rule found : %s", matchingRule)
			return nil
		}

		return fmt.Errorf("Error here\n\tlooking for %s, wasn't found in %s", p, rules)
	}
}

func testAccOutscaleSecurityGroupRuleIngressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_inbound_rule" "outscale_inbound_rule1" {
     count = 1 

    ip_permissions = {
        from_port = 22
        to_port = 22
        ip_protocol = "tcp"
        ip_ranges = ["46.231.147.8/32"]
    }

    group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule2" {
count = 1

     ip_permissions = {
        from_port = 443
        to_port = 443
        ip_protocol = "tcp"
        ip_ranges = ["46.231.147.8/32"]
    }

    group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_firewall_rules_set" "outscale_firewall_rules_set" {
    count = 1

    group_description = "test group tf"
    group_name = "sg1-test-group_test-%d"
}

data "outscale_firewall_rules_set" "by_filter" {}`, rInt)
}
