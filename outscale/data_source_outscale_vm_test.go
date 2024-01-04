package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccVM_DataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	datasourcceName := "data.outscale_vm.basic_web"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMDataSourceConfig(omi, "tinav4.c2r2p2", utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourcceName, "image_id", omi),
					resource.TestCheckResourceAttr(datasourcceName, "vm_type", "tinav4.c2r2p2"),
					resource.TestCheckResourceAttr(datasourcceName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccOAPIVMDataSourceConfig(omi, vmType, region string) string {
	return fmt.Sprintf(`

		resource "outscale_vm" "outscale_vm" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"

			tags {
				key   = "name"
				value = "Terraform-VM"
			}
		}

		data "outscale_vm" "basic_web" {
			filter {
				name   = "vm_ids"
				values = [outscale_vm.outscale_vm.vm_id]
			}
		}
	`, omi, vmType, region)
}
