package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleAvailabilityZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleAvailabilityZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleAvailabilityZoneCheck("data.outscale_sub_region.by_name"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleAvailabilityZoneCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["zone_name"] != "eu-west-2a" {
			return fmt.Errorf("bad name %s", attr["zone_name"])
		}
		if attr["region_name"] != "eu-west-2" {
			return fmt.Errorf("bad region %s", attr["region_name"])
		}

		return nil
	}
}

const testAccDataSourceOutscaleAvailabilityZoneConfig = `
data "outscale_sub_region" "by_name" {
  zone_name = "eu-west-2a"
}
`
