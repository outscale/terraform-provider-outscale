package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNet_WithSecurityGroup(t *testing.T) {
	var group oscgo.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISecurityGroupRuleExists("outscale_security_group.web", &group),
					resource.TestCheckResourceAttr(
						"outscale_security_group.web", "security_group_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISGRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_security_group" {
			continue
		}

		sg, _, err := readSecurityGroups(conn, rs.Primary.ID)
		if sg != nil && err == nil {
			return fmt.Errorf("Outscale Security Group(%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOutscaleOAPISecurityGroupRuleExists(n string, group *oscgo.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		_, resp, err := readSecurityGroups(conn, rs.Primary.ID)
		if err != nil || len(resp.GetSecurityGroups()) < 1 {
			return fmt.Errorf("Outscale Security Group(%s) does not exists: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccOutscaleOAPISecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "web" {
			security_group_name = "terraform_test_%d"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = outscale_net.net.id
		}
	`, rInt)
}
