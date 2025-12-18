package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_WithSubNet_Basic(t *testing.T) {
	resourceName := "outscale_subnet.subnet"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSubnetConfig(utils.GetRegion(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "map_public_ip_on_launch", "false"),
				),
			},
			{
				Config: testAccOutscaleSubnetConfig(utils.GetRegion(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "map_public_ip_on_launch", "true"),
				),
			},
		},
	})
}

func TestAccNet_WithSubNet_Basic_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.0", testAccOutscaleSubnetConfig(utils.GetRegion(), false)),
	})
}

func testAccOutscaleSubnetConfig(region string, mapPublicIpOnLaunch bool) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-subnet-rs"
			}
		}

		resource "outscale_subnet" "subnet" {
			ip_range       = "10.0.0.0/24"
			subregion_name = "%sb"
			net_id         = outscale_net.net.id
			map_public_ip_on_launch = %v
			tags {
				key   = "name"
				value = "terraform-subnet"
			}
		}
	`, region, mapPublicIpOnLaunch)
}
