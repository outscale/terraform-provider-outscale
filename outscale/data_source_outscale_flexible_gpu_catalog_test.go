package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceFlexibleGpuCatalog_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlexibleGpuCatalogConfig(),
			},
		},
	})
}

func testAccDataSourceFlexibleGpuCatalogConfig() string {
	return fmt.Sprintf(`
              data "outscale_flexible_gpu_catalog" "catalog-fGPU" { }
	`)
}
