package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Config: testAccOutscaleOAPIFlexibleGpuLinkConfig(omi, "tinav3.c1r1p2", region),
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
                        vm_initiated_shutdown_behavior = "restart"
		}

                resource "outscale_flexible_gpu" "fGPU-1" { 
                        model_name             =  "nvidia-k2"
                        generation             =  "v3"
                        subregion_name         =  "%s"
                        delete_on_vm_deletion  =   true
                }

                resource "outscale_flexible_gpu_link" "link_fGPU" {
                         flexible_gpu_id = outscale_flexible_gpu.fGPU-1.flexible_gpu_id
                         vm_id           = outscale_vm.basic.vm_id
                }
		
	`, omi, vmType, region)
}
