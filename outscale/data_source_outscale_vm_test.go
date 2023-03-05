package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIVMDataSource_basic(t *testing.T) {
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
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-vm-ds"
			}
		}

 		resource "outscale_subnet" "outscale_subnet" {
			net_id         = outscale_net.outscale_net.net_id
			ip_range       = "10.0.0.0/24"
			subregion_name = "%[3]sa"
		}

		resource "outscale_vm" "outscale_vm" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"
			subnet_id    = outscale_subnet.outscale_subnet.subnet_id

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
