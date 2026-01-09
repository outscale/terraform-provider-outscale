package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVM_withImageExportTaskDataSource_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	imageName := acctest.RandomWithPrefix("test-image-name")

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck: func() {
				testacc.PreCheck(t)
			},
			Providers: testacc.SDKProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleImageExportTaskDataSourceConfig(omi, testAccVmType, utils.GetRegion(), imageName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOutscaleImageExportTaskDataSourceID("data.outscale_image_export_task.test"),
					),
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func testAccCheckOutscaleImageExportTaskDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Image Export Task data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("image Export Task data source ID not set")
		}
		return nil
	}
}

func testAccOutscaleImageExportTaskDataSourceConfig(omi, vmType, region, imageName string) string {
	return fmt.Sprintf(`
	resource "outscale_vm" "basicExport" {
		image_id	         = "%s"
		vm_type                  = "%s"
		keypair_name	         = "terraform-basic"
		placement_subregion_name = "%sa"
	}

	resource "outscale_image" "foo" {
		image_name  = "%s"
		vm_id       = outscale_vm.basicExport.id
		no_reboot   = "true"
		description = "terraform testing"
	}
	resource "outscale_image_export_task" "outscale_image_export_task" {
		image_id                  = outscale_image.foo.id
		osu_export {
			osu_bucket        = "terraform-export-%s"
			disk_image_format = "qcow2"
         }
	}


	data "outscale_image_export_task" "test" {
		filter {
			name   = "task_ids"
			values = [outscale_image_export_task.outscale_image_export_task.id]
		}
	}
	`, omi, vmType, region, imageName, imageName)
}
