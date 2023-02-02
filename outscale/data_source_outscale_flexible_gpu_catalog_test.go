package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_FlexibleGpuCatalog_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_flexible_gpu_catalog.catalog-fGPU"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_FlexibleGpuCatalog_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					// data source validations
					resource.TestCheckResourceAttrSet(dataSourceName, "flexible_gpu_catalog.#"),
				),
			},
		},
	})
}

func testAcc_FlexibleGpuCatalog_DataSource_Config() string {
	return fmt.Sprintf(`
              data "outscale_flexible_gpu_catalog" "catalog-fGPU" { }
	`)
}
