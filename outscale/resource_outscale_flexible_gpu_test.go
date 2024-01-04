package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_FlexibleGpu_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIFlexibleGpuConfig(utils.GetRegion()),
			},
		},
	})
}

func testAccOutscaleOAPIFlexibleGpuConfig(region string) string {
	return fmt.Sprintf(`
                resource "outscale_flexible_gpu" "fGPU-1" { 
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%sa"
                        delete_on_vm_deletion  =   true
                }

		data "outscale_flexible_gpu" "data_fGPU-1" {
			filter {
				name = "flexible_gpu_ids"
				values = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id]
			}
		}
	`, region)
}
