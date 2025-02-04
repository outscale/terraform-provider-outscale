package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_withImageExportTask_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	imageName := acctest.RandomWithPrefix("test-image-name")
	tags := `tags {
			key = "test"
			value = "test"
		}
		tags {
			key = "test-1"
			value = "test-1"
		}`
	if os.Getenv("TEST_QUOTA") == "true" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccOAPIImageExportTaskConfigBasic(omi, utils.TestAccVmType, region, imageName, ""),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOutscalemageExportTaskExists("outscale_image_export_task.outscale_image_export_task"),
					),
				},
				{
					Config: testAccOAPIImageExportTaskConfigBasic(omi, utils.TestAccVmType, region, imageName, tags),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOutscalemageExportTaskExists("outscale_image_export_task.outscale_image_export_task"),
					),
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func testAccCheckOutscalemageExportTaskExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image task id is set")
		}

		return nil
	}
}

func testAccOAPIImageExportTaskConfigBasic(omi, vmType, region, imageName, tags string) string {
	return fmt.Sprintf(`
	resource "outscale_security_group" "sg_task" {
		security_group_name = "sg_export_task"
		description         = "Used in the terraform acceptance tests"

		tags {
			key   = "Name"
			value = "tf-acc-test"
		}
	}

	resource "outscale_vm" "basic" {
		image_id	         = "%s"
		vm_type                  = "%s"
		keypair_name		 = "terraform-basic"
		placement_subregion_name = "%sa"
		security_group_ids = [outscale_security_group.sg_task.security_group_id]
	}

	resource "outscale_image" "foo" {
		image_name  = "%s"
		vm_id       = "outscale_vm.basic.id"
		no_reboot   = "true"
		description = "terraform testing"
	}
	resource "outscale_image_export_task" "outscale_image_export_task" {
		image_id                  = outscale_image.foo.id
		osu_export {
			osu_bucket        = "%s"
			disk_image_format = "qcow2"
		}
		%s
	}
	`, omi, vmType, region, imageName, imageName, tags)
}
