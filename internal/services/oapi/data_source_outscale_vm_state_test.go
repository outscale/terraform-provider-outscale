package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccVM_StateDataSource(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVmStateConfig(omi, testAccVmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleVMStateCheck("data.outscale_vm_state.state"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVMStateCheck(name string) resource.TestCheckFunc {
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

func testAccDataSourceOutscaleVmStateConfig(omi, vmType, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_state" {
			security_group_name = "%[3]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic_state" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
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
	`, omi, vmType, sgName)
}
