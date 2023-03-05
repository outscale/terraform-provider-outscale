package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccVM_withFlexibleGpuLink_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIFlexibleGpuLinkConfig(omi, "tinav5.c2r2p2", utils.GetRegion()),
			},
		},
	})
}

func testAccOutscaleOAPIFlexibleGpuLinkConfig(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
			placement_subregion_name = "%[3]sa"

		}

                resource "outscale_flexible_gpu" "fGPU-1" { 
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%[3]sa"
                        delete_on_vm_deletion  =   true
                }

                resource "outscale_flexible_gpu_link" "link_fGPU" {
                         flexible_gpu_id = outscale_flexible_gpu.fGPU-1.flexible_gpu_id
                         vm_id           = outscale_vm.basic.vm_id
                }
		
	`, omi, vmType, region)
}
