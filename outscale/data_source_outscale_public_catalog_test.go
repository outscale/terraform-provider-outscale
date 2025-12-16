package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DataSourcePublicCatalog_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscalePublicCatalogConfig(),
			},
		},
	})
}

func testAccDataSourceOutscalePublicCatalogConfig() string {
	return fmt.Sprintf(`
              data "outscale_public_catalog" "catalog" { }
	`)
}
