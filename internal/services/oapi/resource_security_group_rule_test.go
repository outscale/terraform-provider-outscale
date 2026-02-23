package oapi_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_SecurityGroupRule_Basic(t *testing.T) {
	resourceName := "outscale_security_group_rule.outscale_security_group_rule_https"

	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleEgressConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ip_range"),
					resource.TestCheckResourceAttr(resourceName, "from_port_range", "443"),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroupRule_Import(t *testing.T) {
	resourceName := "outscale_security_group_rule.outscale_security_group_rule_https"

	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleImport(rInt),
			},
			// Ignore attributes related to the SG Rule, that gets populated after a refresh
			testacc.ImportStepWithStateIdFunc(resourceName, testAccCheckOutscaleRuleImportStateIDFunc(resourceName), testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_SecurityGroupRule_WithSecurityGroupMember(t *testing.T) {
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group_rule.outscale_security_group_rule-3"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupRuleWithGroupMembers(rInt, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rules.0.from_port_range", "22"),
					resource.TestCheckResourceAttrSet(resourceName, "rules.0.security_groups_members.#"),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroupRule_Migration(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.2.1", testAccOutscaleSecurityGroupRuleEgressConfig(rInt)),
	})
}

func TestAccOthers_SecurityGroupRule_WithSecurityGroupMember_Migration(t *testing.T) {
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestStepsWithExpectNonEmptyPlan("1.2.1",
			testAccOutscaleSecurityGroupRuleWithGroupMembers(rInt, accountID),
		),
	})
}

func testAccCheckOutscaleRuleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
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

func testAccOutscaleSecurityGroupRuleImport(rInt int) string {
	return fmt.Sprintf(`
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
