package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVM_withImageExportTaskDataSource_basic(t *testing.T) {
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
					Config: testAccOutscaleOAPIImageExportTaskDataSourceConfig(omi, "tinav4.c2r2p2", utils.GetRegion(), imageName),
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

func testAccOutscaleOAPIImageExportTaskDataSourceConfig(omi, vmType, region, imageName string) string {
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
