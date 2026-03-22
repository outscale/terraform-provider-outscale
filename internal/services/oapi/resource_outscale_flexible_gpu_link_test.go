package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccVM_withFlexibleGpuLink_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleFlexibleGpuLinkConfig(omi, testAccVmType, utils.GetRegion(), sgName),
			},
		},
	})
}

func testAccOutscaleFlexibleGpuLinkConfig(omi, vmType, region, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_fgpu" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"
			placement_subregion_name = "%[3]sa"
			security_group_ids = [outscale_security_group.sg_fgpu.security_group_id]
		}

                resource "outscale_flexible_gpu" "fGPU-1" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%[3]sa"
                        delete_on_vm_deletion  =   true
                }
                resource "outscale_flexible_gpu" "fGPU-2" {
                        model_name             =  "nvidia-p6"
                        generation             =  "v5"
                        subregion_name         =  "%[3]sa"
                        delete_on_vm_deletion  =   true
                }
                resource "outscale_flexible_gpu_link" "link_fGPU" {
                         flexible_gpu_ids = [outscale_flexible_gpu.fGPU-1.flexible_gpu_id,outscale_flexible_gpu.fGPU-2.flexible_gpu_id]
                         vm_id           = outscale_vm.basic.vm_id
                }
	`, omi, vmType, region, sgName)
}
