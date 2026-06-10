package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_WithSubNet_Basic(t *testing.T) {
	resourceName := "outscale_subnet.subnet"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
		Steps: testacc.FrameworkMigrationTestSteps("1.1.0", testAccOutscaleSubnetConfig(utils.GetRegion(), false)),
	})
}

func TestAccNet_Subnet_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_subnet.subnet"
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-create-failure"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAccOutscaleSubnetConfigWithTag(utils.GetRegion(), false, invalidTagKey, tagValue),
			testAccOutscaleSubnetConfigWithTag(utils.GetRegion(), false, "name", tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func testAccOutscaleSubnetConfig(region string, mapPublicIpOnLaunch bool) string {
	return testAccOutscaleSubnetConfigWithTag(region, mapPublicIpOnLaunch, "name", "terraform-subnet")
}

func testAccOutscaleSubnetConfigWithTag(region string, mapPublicIpOnLaunch bool, tagKey, tagValue string) string {
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
				key   = %q
				value = %q
			}
		}
	`, region, mapPublicIpOnLaunch, tagKey, tagValue)
}
