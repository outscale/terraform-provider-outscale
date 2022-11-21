package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSubregionsDataSource_Basic(t *testing.T) {
	t.Parallel()
	subregionName := "eu-west-2b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubregionsDataSourceConfig(subregionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func TestAccSubregionsDataSource_All(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubregionsDataSourceConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func testAccCheckSubregionsDataSourceID(n string) resource.TestCheckFunc {
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

func testAccCheckSubregionsDataSourceConfig(subregionName string) string {
	return fmt.Sprintf(`
		data "outscale_subregions" "test" {
			filter {
				name = "subregion_names"
				values = ["%s"]
			}
		}
	`, subregionName)
}

var testAccCheckSubregionsDataSourceConfigAll = `
		data "outscale_subregions" "test" {
		}
	`
