package outscale

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOthers_SecurityGroupRule_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_security_group_rule.outscale_security_group_rule_https"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleEgressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ip_range"),
					resource.TestCheckResourceAttr(resourceName, "from_port_range", "443"),
				),
			},
			{
				ImportStateIdFunc:       testAccCheckOutscaleRuleImportStateIDFunc(resourceName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func TestAccOthers_SecurityGroupRule_withSecurityGroupMember(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleWithGroupMembers(rInt, accountID),
			},
		},
	})
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
