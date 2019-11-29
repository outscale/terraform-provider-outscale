package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSPublicCatalog_basic(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOAPIPublicCatalogConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicCatalogDataSourceID("data.outscale_public_catalog.test"),
					resource.TestCheckResourceAttrSet("data.outscale_public_catalog.test", "catalog.#"),
					resource.TestCheckResourceAttrSet("data.outscale_public_catalog.test", "catalog.0.catalog_entries.#"),
					resource.TestCheckResourceAttrSet("data.outscale_public_catalog.test", "catalog.0.catalog_attributes.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicCatalogDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Catalog data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Catalog data source ID not set")
		}
		return nil
	}
}

func testAccDSOAPIPublicCatalogConfig() string {
	return fmt.Sprintf(`
		data "outscale_public_catalog" "test" {}
	`)
}
