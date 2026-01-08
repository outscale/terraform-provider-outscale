package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_SubregionsDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSubregionsDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func TestAccOthers_SubregionsDataSource_All(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSubregionsDataSourceConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSubregionsDataSourceID("data.outscale_subregions.test"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSubregionsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find subregions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("subregions data source id not set")
		}
		return nil
	}
}

func testAccCheckOutscaleSubregionsDataSourceConfig(subregionName string) string {
	return fmt.Sprintf(`
		data "outscale_subregions" "test" {
			filter {
				name = "subregion_names"
				values = ["%[1]sa"]
			}
			filter {
				name = "region_names"
				values = ["%[1]s"]
			}
		}
	`, subregionName)
}

var testAccCheckOutscaleSubregionsDataSourceConfigAll = `
		data "outscale_subregions" "test" {
		}
	`
