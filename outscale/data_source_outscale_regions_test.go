package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOthers_RegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIRegionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRegionsDataSourceID("data.outscale_regions.regions"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIRegionsDataSourceID(n string) resource.TestCheckFunc {
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

var testAccCheckOutscaleOAPIRegionsDataSourceConfig = `
		data "outscale_regions" "regions" {}
	`
