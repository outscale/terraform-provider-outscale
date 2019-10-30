package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPIOutboundRule(t *testing.T) {
	var group oapi.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISecurityGroupRuleEgressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRuleExists("outscale_security_group.outscale_security_group", &group),
					testAccCheckOutscaleOAPIRuleAttributes("outscale_security_group_rule.outscale_security_group_rule_https", &group, nil, "Inbound"),
				),
			},
		},
	})
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
				FromPortRange: 443,
				ToPortRange:   443,
				IpProtocol:    "tcp",
				IpRanges:      []string{"46.231.147.8/32"},
			}
		}

		var matchingRule *oapi.SecurityGroupRule
		var rules []oapi.SecurityGroupRule
		if ruleType == "Inbound" {
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

		return fmt.Errorf("Security Rules: looking for %+v, wasn't found in %+v", p, rules)
	}
}

func testAccOutscaleOAPISecurityGroupRuleEgressConfig(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_security_group_rule" "outscale_security_group_rule" {
	flow              = "Inbound"
	security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"
	from_port_range = "0"
	to_port_range = "0"
	ip_protocol = "tcp"
	ip_range = "0.0.0.0/0"
}

resource "outscale_security_group_rule" "outscale_security_group_rule_https" {
	flow = "Inbound"
	from_port_range = 443
	to_port_range = 443
	ip_protocol = "tcp"
	ip_range = "46.231.147.8/32"
	security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"
}

resource "outscale_security_group" "outscale_security_group" {
	description         = "test group"
	security_group_name = "sg1-test-group_test_%d"
}
`, rInt)
}
