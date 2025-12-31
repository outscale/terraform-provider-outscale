package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_WithSecurityGroup(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group.web"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroupWithoutName(t *testing.T) {
	resourceName := "outscale_security_group.noname"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupWithoutNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
		},
	})
}

func TestAccNet_WithSecurityGroup_Migration(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.2.1", testAccOutscaleSecurityGroupConfig(rInt), testAccOutscaleSecurityGroupWithoutNameConfig()),
	})
}

func testAccOutscaleSecurityGroupConfig(rInt int) string {
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

func testAccOutscaleSecurityGroupWithoutNameConfig() string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "noname" {
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test-no-name"
			}
		}
	`)
}
