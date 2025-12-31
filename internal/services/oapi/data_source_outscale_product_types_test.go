package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceProductTypes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleProductTypesConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourceOutscaleProductTypesConfig = `
 data "outscale_product_types" "all" {}`
