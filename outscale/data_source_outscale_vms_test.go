package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVMS_DataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMSDataSourceConfig(omi, "tinav4.c2r2p2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "vms.0.image_id", omi),
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "vms.0.vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func testAccOAPIVMSDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id			= "%s"
			vm_type				= "%s"
			keypair_name	= "terraform-basic"
		}

		data "outscale_vms" "basic_web" {
			filter {
				name   = "vm_ids"
				values = [outscale_vm.basic.id]
			}
		}`, omi, vmType)
}
