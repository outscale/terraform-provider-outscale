package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_FlexibleGpu_Basic(t *testing.T) {
	resourceName := "outscale_flexible_gpu.fGPU-1"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleGpuConfig(utils.GetRegion(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "model_name"),
					resource.TestCheckResourceAttrSet(resourceName, "generation"),
					resource.TestCheckResourceAttrSet(resourceName, "subregion_name"),
					resource.TestCheckResourceAttr(resourceName, "delete_on_vm_deletion", "false"),
				),
			},
			{
				Config: testAccFlexibleGpuConfig(utils.GetRegion(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "model_name"),
					resource.TestCheckResourceAttrSet(resourceName, "generation"),
					resource.TestCheckResourceAttrSet(resourceName, "subregion_name"),
					resource.TestCheckResourceAttr(resourceName, "delete_on_vm_deletion", "true"),
				),
			},
		},
	})
}

func TestAccOthers_FlexibleGpu_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.0.1", testAccFlexibleGpuConfig(utils.GetRegion(), false)),
	})
}

func testAccFlexibleGpuConfig(region string, deletion bool) string {
	return fmt.Sprintf(`
                resource "outscale_flexible_gpu" "fGPU-1" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%sa"
                        delete_on_vm_deletion  =  %v
                }

		data "outscale_flexible_gpu" "data_fGPU-1" {
			filter {
				name = "flexible_gpu_ids"
				values = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
			}
		}
	`, region, deletion)
}
