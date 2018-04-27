package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleAvailabilityZones(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_sub_regions.by_name", "availability_zone_info.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleAvailabilityZonesConfig = `
data "outscale_sub_regions" "by_name" {
	zone_name = ["eu-west-2a"]
}
`
