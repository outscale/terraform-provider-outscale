package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_DataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	datasourcceName := "data.outscale_vm.basic_web"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMDataSourceConfig(omi, utils.TestAccVmType, utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourcceName, "image_id", omi),
					resource.TestCheckResourceAttr(datasourcceName, "vm_type", utils.TestAccVmType),
					resource.TestCheckResourceAttr(datasourcceName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccOAPIVMDataSourceConfig(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_data" {
			security_group_name = "sg_vm_basic"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "outscale_vm_data" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vm_data.security_group_id]
			tags {
				key   = "name"
				value = "Terraform-VM"
			}

		}

		data "outscale_vm" "basic_web" {
			filter {
				name   = "vm_ids"
				values = [outscale_vm.outscale_vm_data.vm_id]
			}
		}
	`, omi, vmType, region)
}
