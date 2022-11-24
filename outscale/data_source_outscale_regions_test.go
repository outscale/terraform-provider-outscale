package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRegionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionsDataSourceID("data.outscale_regions.regions"),
				),
			},
		},
	})
}

func testAccCheckRegionsDataSourceID(n string) resource.TestCheckFunc {
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

var testAccCheckRegionsDataSourceConfig = `
		data "outscale_regions" "regions" {}
	`
