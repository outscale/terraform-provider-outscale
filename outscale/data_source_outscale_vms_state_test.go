package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVMSState(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIVMSStateConfig(omi, "c4.large"),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMSStateConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
		}

		data "outscale_vms_state" "state" {
			vm_ids = ["${outscale_vm.basic.id}"]
		}
	`, omi, vmType)
}
