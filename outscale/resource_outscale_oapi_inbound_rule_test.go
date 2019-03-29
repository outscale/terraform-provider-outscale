package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPIInboundRule(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapiFlag, err := strconv.ParseBool(o)
	if err != nil {
		oapiFlag = false
	}

	if !oapiFlag {
		t.Skip()
	}
	var group oapi.SecurityGroup
	rInt := acctest.RandInt()

	testRuleCount := func(*terraform.State) error {
		if len(group.InboundRules) != 1 {
			return fmt.Errorf("Wrong Security Group rule count, expected %d, got %d",
				1, len(group.InboundRules))
		}

		rule := group.InboundRules[0]
		if rule.FromPortRange != int64(80) {
			return fmt.Errorf("Wrong Security Group port setting, expected %d, got %d",
				80, int(rule.FromPortRange))
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISecurityGroupRuleIngressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRuleExists("outscale_firewall_rules_set.web", &group),
					testAccCheckOutscaleOAPIRuleAttributes("outscale_inbound_rule.ingress_1", &group, nil, "ingress"),
					testRuleCount,
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISecurityGroupRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_firewall_rules_set" {
			continue
		}

		// Retrieve our group
		req := oapi.ReadSecurityGroupsRequest{
			Filters: oapi.FiltersSecurityGroup{
				SecurityGroupIds: []string{rs.Primary.ID},
			},
		}
		var resp *oapi.POST_ReadSecurityGroupsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadSecurityGroups(req)

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
			if len(resp.OK.SecurityGroups) > 0 && resp.OK.SecurityGroups[0].SecurityGroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists", rs.Primary.ID)
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

func testAccCheckOutscaleOAPIRuleExists(n string, group *oapi.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		req := oapi.ReadSecurityGroupsRequest{
			Filters: oapi.FiltersSecurityGroup{
				SecurityGroupIds: []string{rs.Primary.ID},
			},
		}

		var resp *oapi.POST_ReadSecurityGroupsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadSecurityGroups(req)

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

		if len(resp.OK.SecurityGroups) > 0 && resp.OK.SecurityGroups[0].SecurityGroupId == rs.Primary.ID {
			*group = resp.OK.SecurityGroups[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccCheckOutscaleOAPIRuleAttributes(n string, group *oapi.SecurityGroup, p *oapi.SecurityGroupRule, ruleType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Security Group Rule Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group Rule is set")
		}

		if p == nil {
			p = &oapi.SecurityGroupRule{
				FromPortRange: 80,
				ToPortRange:   8000,
				IpProtocol:    "tcp",
				IpRanges:      []string{"10.0.0.0/8"},
			}
		}

		var matchingRule *oapi.SecurityGroupRule
		var rules []oapi.SecurityGroupRule
		if ruleType == "ingress" {
			rules = group.InboundRules
		} else {
			rules = group.OutboundRules
		}

		if len(rules) == 0 {
			return fmt.Errorf("No IPPerms")
		}

		for _, r := range rules {
			if p.ToPortRange != r.ToPortRange {
				continue
			}

			if p.FromPortRange != r.FromPortRange {
				continue
			}

			if p.IpProtocol != r.IpProtocol {
				continue
			}

			remaining := len(p.IpRanges)
			for _, ip := range p.IpRanges {
				for _, rip := range r.IpRanges {
					if ip == rip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.SecurityGroupsMembers)
			for _, ip := range p.SecurityGroupsMembers {
				for _, rip := range r.SecurityGroupsMembers {
					if ip.SecurityGroupId == rip.SecurityGroupId {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.PrefixListIds)
			for _, pip := range p.PrefixListIds {
				for _, rpip := range r.PrefixListIds {
					if pip == rpip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			matchingRule = &r
		}

		if matchingRule != nil {
			return nil
		}

		return fmt.Errorf("Error here\n\tlooking for %+v, wasn't found in %+v", p, rules)
	}
}

func testAccOutscaleOAPISecurityGroupRuleIngressConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "web" {
		firewall_rules_set_name = "terraform_test_%d"
		description = "Used in the terraform acceptance tests"
					tag = {
									Name = "tf-acc-test"
					}
	}
	resource "outscale_inbound_rule" "ingress_1" {
		inbound_rule = {
				ip_protocol = "tcp"
				from_port_range = 80
				to_port_range = 8000
				ip_ranges = ["10.0.0.0/8"]
		}
		firewall_rules_set_id = "${outscale_firewall_rules_set.web.id}"
	}`, rInt)
}
