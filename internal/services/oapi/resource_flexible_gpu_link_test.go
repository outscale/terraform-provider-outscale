package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccVM_withFlexibleGpuLink_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_flexible_gpu_link.link_fGPU"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// single GPU linked
			{
				Config: testAccOutscaleFlexibleGpuLinkUpdateConfig(omi, sgName, false),
			},
			// add the second GPU to the link
			{
				Config: testAccOutscaleFlexibleGpuLinkUpdateConfig(omi, sgName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "flexible_gpu_ids.#", "2"),
				),
			},
			// remove the second GPU
			{
				Config: testAccOutscaleFlexibleGpuLinkUpdateConfig(omi, sgName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "flexible_gpu_ids.#", "1"),
				),
			},
			testacc.ImportStepWithStateIdFunc(resourceName, fgpuLinkStateIDFunc(resourceName)),
		},
	})
}

func fgpuLinkStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["vm_id"], nil
	}
}

func TestAccVM_withFlexibleGpuLink_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.4.0",
			testAccOutscaleFlexibleGpuLinkUpdateConfig(omi, sgName, false),
		),
	})
}

func testAccOutscaleFlexibleGpuLinkUpdateConfig(omi, sgName string, bothGpus bool) string {
	gpuIds := "outscale_flexible_gpu.fGPU-1.id"
	if bothGpus {
		gpuIds += ", outscale_flexible_gpu.fGPU-2.id"
	}
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_fgpu" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"
		}

		resource "outscale_vm" "basic" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"
			placement_subregion_name = "%[3]sa"
			security_group_ids = [outscale_security_group.sg_fgpu.security_group_id]

			lifecycle { ignore_changes = [state] }
		}

		resource "outscale_flexible_gpu" "fGPU-1" {
			model_name            = "%[6]s"
			generation            = "%[7]s"
			subregion_name        = "%[3]sa"
			delete_on_vm_deletion = true
		}

		resource "outscale_flexible_gpu" "fGPU-2" {
			model_name            = "%[6]s"
			generation            = "%[7]s"
			subregion_name        = "%[3]sa"
			delete_on_vm_deletion = true
		}

		resource "outscale_flexible_gpu_link" "link_fGPU" {
			flexible_gpu_ids = [%[5]s]
			vm_id            = outscale_vm.basic.vm_id
		}
	`, omi, testAccVmTypefGPU, utils.GetRegion(), sgName, gpuIds, testAccfGPUModel, testAccfGPUGeneration)
}
