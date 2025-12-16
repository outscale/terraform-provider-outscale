package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_StatesDataSource(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVMStatesConfig(omi, utils.TestAccVmType),
			},
		},
	})
}

func testAccDataSourceOutscaleVMStatesConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vStates" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sg_volumes_link"
		}

		resource "outscale_vm" "basic" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vStates.security_group_id]

		}

		data "outscale_vm_states" "state" {
			vm_ids = [outscale_vm.basic.id]
		}
	`, omi, vmType)
}
