package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_WithSecurityGroup(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group.web"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
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

func TestAccNet_WithSecurityGroup_Migration(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"outscale": {
						VersionConstraint: "1.2.1",
						Source:            "outscale/outscale",
					},
				},
				Config: testAccOutscaleSecurityGroupConfig(rInt),
			},
			{
				ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
				Config:                   testAccOutscaleSecurityGroupConfig(rInt),
				PlanOnly:                 true,
			},
		},
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
