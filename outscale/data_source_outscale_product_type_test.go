package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ProductType_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_product_type.test"
	dataSourcesName := "data.outscale_product_types.all"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_ProductType_DataSource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "product_type_id", "0001"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "product_types.#"),
				),
			},
		},
	})
}

const testAcc_ProductType_DataSource_Config = `
	data "outscale_product_type" "test" { 
   		filter {
        	name     = "product_type_ids"
        	values   = ["0001"]
    	}
	}

	data "outscale_product_types" "all" {}
`
