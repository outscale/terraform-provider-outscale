package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceOutscaleOAPIProductTypes(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIProductTypesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_product_types.test", "product_type_set.#", "3"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIProductTypesConfig = `
data "outscale_product_types" "test" {}
`
