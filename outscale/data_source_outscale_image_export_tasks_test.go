package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVM_withImageExportTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	imageName := acctest.RandomWithPrefix("test-image-name")

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleImageExportTasksDataSourceConfig(omi, utils.TestAccVmType, utils.GetRegion(), imageName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOutscaleImageExportTaskDataSourceID("outscale_image_export_task.outscale_image_export_task"),
					),
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func testAccOutscaleImageExportTasksDataSourceConfig(omi, vmType, region, imageName string) string {
	return fmt.Sprintf(`
	resource "outscale_vm" "basic" {
		image_id			      = "%s"
		vm_type             = "%s"
		keypair_name		    = "terraform-basic"
		placement_subregion_name = "%sa"
	}

	resource "outscale_image" "foo" {
		image_name  = "%s"
		vm_id       = outscale_vm.basic.id
		no_reboot   = "true"
		description = "terraform testing"
	}
	resource "outscale_image_export_task" "outscale_image_export_task" {
		image_id                     = outscale_image.foo.id
		osu_export {
			osu_bucket        = "terraform-export-%s"
			disk_image_format = "qcow2"
         }
	}

	data "outscale_image_export_tasks" "export_tasks" {
		filter {
			name = "tags_ids"
			values = [outscale_image_export_task.outscale_image_export_task.id]
		}
	}
	`, omi, vmType, region, imageName, imageName)
}
