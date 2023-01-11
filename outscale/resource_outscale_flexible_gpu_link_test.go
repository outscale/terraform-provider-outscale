package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOutscaleOAPIFlexibleGpuLink_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := fmt.Sprintf("%sa", os.Getenv("OUTSCALE_REGION"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIFlexibleGpuLinkConfig(omi, "tinav5.c2r2p1", region),
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
		}

                resource "outscale_flexible_gpu" "fGPU-1" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%[3]s"
                        delete_on_vm_deletion  =   true
                }
                resource "outscale_flexible_gpu" "fGPU-2" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%[3]s"
                        delete_on_vm_deletion  =   true
                }
                resource "outscale_flexible_gpu_link" "link_fGPU" {
                         flexible_gpu_ids = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id,outscale_flexible_gpu.fGPU-2.flexible_gpu_id]
                         vm_id           = outscale_vm.basic.vm_id
                }
	`, omi, vmType, region)
}
