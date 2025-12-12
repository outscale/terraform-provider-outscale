package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DataSourceFlexibleGpuCatalog_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleFlexibleGpuCatalogConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleFlexibleGpuCatalogConfig() string {
	return fmt.Sprintf(`
              data "outscale_flexible_gpu_catalog" "catalog-fGPU" { }
	`)
}
