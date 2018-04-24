package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleRegions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_regions.by_name", "region_info.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleRegionsConfig = `
data "outscale_regions" "by_name_current" {
  filter {
		name = "region-name"
		values = ["eu-west-2"]
	}
}


`
