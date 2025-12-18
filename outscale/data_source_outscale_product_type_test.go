package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DataSourceProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleProductTypeConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourceOutscaleProductTypeConfig = `
 data "outscale_product_type" "test" {
   filter {
        name     = "product_type_ids"
        values   = ["0001"]
    }
}`
