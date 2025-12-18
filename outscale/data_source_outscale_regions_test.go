package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_RegionsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleRegionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRegionsDataSourceID("data.outscale_regions.regions"),
				),
			},
		},
	})
}

func testAccCheckOutscaleRegionsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find regions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("regions data source ID not set")
		}
		return nil
	}
}

var testAccCheckOutscaleRegionsDataSourceConfig = `
		data "outscale_regions" "regions" {}
	`
