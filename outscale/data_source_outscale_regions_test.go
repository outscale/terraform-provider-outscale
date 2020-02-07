package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceOutscaleOAPIRegions(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_regions.by_name_current", "region_info.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIRegionsConfig = `
	data "outscale_regions" "by_name_current" {
		filter {
			name = "region-name"
			values = ["eu-west-2"]
		}
	}
`
