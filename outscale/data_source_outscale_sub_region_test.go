package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIAvailabilityZone(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIAvailabilityZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIAvailabilityZoneCheck("data.outscale_sub_region.by_name"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIAvailabilityZoneCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["sub_region_name"] != "eu-west-2a" {
			return fmt.Errorf("bad name %s", attr["sub_region_name"])
		}
		if attr["region_name"] != "eu-west-2" {
			return fmt.Errorf("bad region %s", attr["region_name"])
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIAvailabilityZoneConfig = `
data "outscale_sub_region" "by_name" {
  sub_region_name = "eu-west-2a"
}
`
