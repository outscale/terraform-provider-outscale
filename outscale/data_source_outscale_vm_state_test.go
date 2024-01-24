package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVM_StateDataSource(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVmStateConfig(omi, "tinav4.c2r2p2"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVMStateCheck("data.outscale_vm_state.state"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMStateCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		vm, ok := s.RootModule().Resources["outscale_vm.basic_state"]
		if !ok {
			return fmt.Errorf("can't find outscale_public_ip.test in state")
		}

		state := rs.Primary.Attributes

		if state["vm_id"] != vm.Primary.Attributes["vm_id"] {
			return fmt.Errorf(
				"vm_id is %s; want %s",
				state["vm_id"],
				vm.Primary.Attributes["vm_id"],
			)
		}
		return nil
	}
}

func testAccDataSourceOutscaleOAPIVmStateConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_state" {
			security_group_name = "sg_vmState"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic_state" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vm_state.security_group_id]
		}

		data "outscale_vm_state" "state" {
			all_vms = false
			filter {
				name   = "vm_ids"
				values = [outscale_vm.basic_state.id]
			}
			filter {
				name   = "vm_states"
				values = ["running"]
			}
		}
	`, omi, vmType)
}
