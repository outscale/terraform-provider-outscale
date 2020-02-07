package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceOutscaleOAPIRegion(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIRegionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIRegionCheck("data.outscale_region.by_name_current", "eu-west-2", "true"),
					// testAccDataSourceOutscaleOAPIRegionCheck("data.outscale_region.by_name_other", "us-west-1", "false"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIRegionCheck(name, region, current string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["region_name"] != region {
			return fmt.Errorf("bad name %s", attr["region_name"])
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIRegionConfig = `
	data "outscale_region" "by_name_current" {
		filter {
			name = "region-name"
			values = ["eu-west-2"]
		}
	}
`
