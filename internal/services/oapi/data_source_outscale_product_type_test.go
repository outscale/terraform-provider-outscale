package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testacc.SDKProviders,
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
