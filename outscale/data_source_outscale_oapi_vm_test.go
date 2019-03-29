package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVMDataSource_basic(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMDataSourceConfig(omi, "c4.large"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "image_id", omi),
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "type", "c4.large"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccOAPIVMDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id               = "%s"
			vm_type                = "%s"
			keypair_name           = "terraform-basic"
			security_group_ids     = ["sg-9752b7a6"]
		}

		data "outscale_vm" "basic_web" {
			filter {
			name = "vm_ids"
			values = ["${outscale_vm.basic.id}"]
		  }
		}`, omi, vmType)
}
