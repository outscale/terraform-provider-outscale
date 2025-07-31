package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_FlexibleGpu_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_flexible_gpu.fGPU-1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
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
