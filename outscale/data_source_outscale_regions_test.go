package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Regions_DataSource(t *testing.T) {
	t.Parallel()
	dataSourcesName := "data.outscale_regions.regions"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Regions_DataSource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "regions.#"),
				),
			},
		},
	})
}

var testAcc_Regions_DataSource_Config = `
		data "outscale_regions" "regions" {}
	`
