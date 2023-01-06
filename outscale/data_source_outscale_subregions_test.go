package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOutscaleOAPISubregionsDataSource_Basic(t *testing.T) {
	t.Parallel()
	subregionName := "eu-west-2b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISubregionsDataSourceConfig(subregionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISubregionsDataSource_All(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISubregionsDataSourceConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISubregionsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find subregions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("subregions data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleOAPISubregionsDataSourceConfig(subregionName string) string {
	return fmt.Sprintf(`
		data "outscale_subregions" "test" {
			filter {
				name = "subregion_names"
				values = ["%s"]
			}
		}
	`, subregionName)
}

var testAccCheckOutscaleOAPISubregionsDataSourceConfigAll = `
		data "outscale_subregions" "test" {
		}
	`
