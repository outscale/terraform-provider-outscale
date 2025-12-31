package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccVM_withImageExportTask_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	imageName := acctest.RandomWithPrefix("test-image-name")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	tags := `tags {
			key = "test"
			value = "test"
		}
		tags {
			key = "test-1"
			value = "test-1"
		}`
	if os.Getenv("TEST_QUOTA") == "true" {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck: func() {
				testacc.PreCheck(t)
			},
			ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccOAPIImageExportTaskConfigBasic(omi, oapi.TestAccVmType, region, imageName, "", sgName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOutscalemageExportTaskExists("outscale_image_export_task.outscale_image_export_task"),
					),
				},
				{
					Config: testAccOAPIImageExportTaskConfigBasic(omi, oapi.TestAccVmType, region, imageName, tags, sgName),
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

func testAccOAPIImageExportTaskConfigBasic(omi, vmType, region, imageName, tags, sgName string) string {
	return fmt.Sprintf(`
	resource "outscale_security_group" "sg_task" {
		security_group_name = "%[7]s"
		description         = "Used in the terraform acceptance tests"

		tags {
			key   = "Name"
			value = "tf-acc-test"
		}
	}

	resource "outscale_vm" "basic" {
		image_id	         = "%[1]s"
		vm_type                  = "%[2]s"
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
	`, omi, vmType, region, imageName, imageName, tags, sgName)
}
