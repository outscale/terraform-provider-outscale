package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceFlexibleGpuCatalog_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleFlexibleGpuCatalogConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleFlexibleGpuCatalogConfig() string {
	return `
              data "outscale_flexible_gpu_catalog" "catalog-fGPU" { }
	`
}
