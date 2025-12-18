package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVMS_DataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMSDataSourceConfig(omi, utils.TestAccVmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "vms.0.image_id", omi),
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "vms.0.vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func testAccOAPIVMSDataSourceConfig(omi, vmType, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vms" {
			security_group_name = "%[3]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "databasic" {
			image_id			= "%[1]s"
			vm_type				= "%[2]s"
			keypair_name	= "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vms.security_group_id]
		}

		data "outscale_vms" "basic_web" {
			filter {
				name   = "vm_ids"
				values = [outscale_vm.databasic.id]
			}
			filter {
				name   = "keypair_names"
				values = ["terraform-basic"]
			}
			filter {
				name   = "image_ids"
				values = ["%[1]s"]
			}

		}`, omi, vmType, sgName)
}
