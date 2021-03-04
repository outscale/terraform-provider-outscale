package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPISubregionDataSource_basic(t *testing.T) {
	subregionName := "eu-west-2a"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISubregionDataSourceBasicConfig(subregionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubregionDataSourceID("data.outscale_subregion.test"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISubregionDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find subregion data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("subregion data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleOAPISubregionDataSourceBasicConfig(subregionName string) string {
	return fmt.Sprintf(`	
		data "outscale_subregion" "test" {
			filter {
				name   = "subregion_names"
				values = ["%s"]
			}
		}
	`, subregionName)
}
