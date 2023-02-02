package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Subregions_DataSource(t *testing.T) {
	t.Parallel()
	subregionName := "eu-west-2b"
	dataSourcesName := "data.outscale_subregions.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Subregions_DataSource_Config(subregionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "subregions.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Subregions_DataSource_All(t *testing.T) {
	t.Parallel()
	dataSourcesName := "data.outscale_subregions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Subregions_DataSource_All_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "subregions.#"),
				),
			},
		},
	})
}

func testAcc_Subregions_DataSource_Config(subregionName string) string {
	return fmt.Sprintf(`
		data "outscale_subregions" "test" {
			filter {
				name = "subregion_names"
				values = ["%s"]
			}
		}
	`, subregionName)
}

var testAcc_Subregions_DataSource_All_Config = `
		data "outscale_subregions" "test" {}
	`
