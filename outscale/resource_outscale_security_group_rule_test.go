package outscale

import (
	"fmt"
	"os"
	"strings"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOthers_SecurityGroupRule_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_security_group_rule.outscale_security_group_rule_https"

	rInt := acctest.RandInt()

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleSecurityGroupRuleEgressConfig(rInt),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceName, "ip_range"),
						resource.TestCheckResourceAttr(resourceName, "from_port_range", "443"),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateIdFunc:       testAccCheckOutscaleRuleImportStateIDFunc(resourceName),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"request_id"},
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func TestAccNet_AddSecurityGroupRuleMembersWithSgName(t *testing.T) {
	rInt := acctest.RandInt()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	resourceName := "outscale_security_group_rule.rule_group"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddSecurityGroupRuleMembersWithSgName(rInt, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttr(resourceName, "flow", "Inbound"),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroupRule_withSecurityGroupMember(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleWithGroupMembers(rInt, accountID),
			},
		},
	})
}

func testAccCheckOutscaleSecurityGroupRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_security_group_rule" {
			continue
		}

		sg, _, err := readSecurityGroups(conn, rs.Primary.ID)
		if sg != nil && err == nil {
			return fmt.Errorf("Outscale Security Group Rule(%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOutscaleRuleAttributes(n string, group *oscgo.SecurityGroup, p *oscgo.SecurityGroupRule, ruleType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Security Group Rule Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group Rule is set")
		}

		if p == nil {
			fromPortRange := int32(443)
			toPortRange := int32(443)
			ipProtocol := "tcp"
			p = &oscgo.SecurityGroupRule{
				FromPortRange: &fromPortRange,
				ToPortRange:   &toPortRange,
				IpProtocol:    &ipProtocol,
				IpRanges:      &[]string{"46.231.147.8/32"},
			}
		}

		var matchingRule *oscgo.SecurityGroupRule
		var rules []oscgo.SecurityGroupRule
		if ruleType == "Inbound" {
			rules = group.GetInboundRules()
		} else {
			rules = group.GetOutboundRules()
		}

		if len(rules) == 0 {
			return fmt.Errorf("No IPPerms")
		}

		for _, r := range rules {
			if p.GetToPortRange() != r.GetToPortRange() {
				continue
			}

			if p.GetFromPortRange() != r.GetFromPortRange() {
				continue
			}

			if p.GetIpProtocol() != r.GetIpProtocol() {
				continue
			}

			remaining := len(p.GetIpRanges())
			for _, ip := range p.GetIpRanges() {
				for _, rip := range r.GetIpRanges() {
					if ip == rip {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.GetSecurityGroupsMembers())
			for _, ip := range p.GetSecurityGroupsMembers() {
				for _, rip := range r.GetSecurityGroupsMembers() {
					if ip.GetSecurityGroupId() == rip.GetSecurityGroupId() {
						remaining--
					}
				}
			}

			if remaining > 0 {
				continue
			}

			remaining = len(p.GetServiceIds())
			for _, pip := range p.GetServiceIds() {
				for _, rpip := range r.GetServiceIds() {
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

func testAccCheckOutscaleRuleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s_%s_%s_%s_%s_%s", rs.Primary.ID, strings.ToLower(rs.Primary.Attributes["flow"]), rs.Primary.Attributes["ip_protocol"], rs.Primary.Attributes["from_port_range"], rs.Primary.Attributes["to_port_range"], rs.Primary.Attributes["ip_range"]), nil
	}
}

func testAccOutscaleSecurityGroupRuleEgressConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_security_group_rule" "outscale_security_group_rule" {
			flow              = "Inbound"
			security_group_id = outscale_security_group.outscale_security_group.security_group_id
                        from_port_range = "0"
			to_port_range   = "0"
			ip_protocol     = "tcp"
			ip_range        = "0.0.0.0/0"
		}

		resource "outscale_security_group_rule" "outscale_security_group_rule_https" {
			flow              = "Inbound"
			from_port_range   = 443
			to_port_range     = 443
			ip_protocol       = "tcp"
			ip_range          = "46.231.147.8/32"
			security_group_id = outscale_security_group.outscale_security_group.security_group_id
		}

		resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg1-test-group_test_%d"
		}
	`, rInt)
}

func testAccOutscaleSecurityGroupRuleWithGroupMembers(rInt int, accountID string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg3-terraform-test_%[2]d"
			tags {
				key   = "Name"
				value = "outscale_sg"
			}
		}

		resource "outscale_security_group" "outscale_security_group2" {
			description         = "test group"
			security_group_name = "sg4-terraform-test_%[2]d"
			tags {
				key   = "Name"
				value = "outscale_sg2"
			}
		}

		resource "outscale_security_group_rule" "outscale_security_group_rule-3" {
			flow              = "Inbound"
			security_group_id = outscale_security_group.outscale_security_group.id
			rules {
				from_port_range = "22"
				to_port_range   = "22"
				ip_protocol     = "tcp"
				security_groups_members {
					account_id          = "%[1]s"
					security_group_name = outscale_security_group.outscale_security_group2.security_group_name
				}
			}
                     depends_on = [outscale_security_group.outscale_security_group2]
		}
	`, accountID, rInt)
}

func testAccAddSecurityGroupRuleMembersWithSgName(rInt int, accountID string) string {
	return fmt.Sprintf(`

resource "outscale_net" "netSgtest" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "security_group" {
    description         = "testing security group"
    security_group_name = "terraform-test_%[2]d"
    net_id              = outscale_net.netSgtest.net_id
}
resource "outscale_security_group_rule" "rule_group" {
    security_group_id = outscale_security_group.security_group.security_group_id
    flow              = "Inbound"
    rules {
        from_port_range   = 22
        to_port_range     = 22
        ip_protocol       = "tcp"
        security_groups_members {
            account_id          = "%[1]s"
            security_group_name = outscale_security_group.security_group.security_group_name
        }
    }
}
	`, accountID, rInt)
}
