package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOthers_DataSourceProductType(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIProductTypeConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIProductTypeConfig = `
 data "outscale_product_type" "test" { 
   filter {
        name     = "product_type_ids"
        values   = ["0001"]
    }
}`
