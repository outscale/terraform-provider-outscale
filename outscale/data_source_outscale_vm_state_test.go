package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleVmState(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleVmStateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleVMStateCheck("data.outscale_vm_state.state"),
					// testAccDataSourceOutscaleVMStateCheck("data.outscale_public_ip.by_public_ip"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVMStateCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		fmt.Printf("\n[DEBUG] TEST RS %s \n", s.RootModule().Resources)

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		vm, ok := s.RootModule().Resources["outscale_vm.basic"]
		if !ok {
			return fmt.Errorf("can't find outscale_public_ip.test in state")
		}

		state := rs.Primary.Attributes

		if state["instance_id"] != vm.Primary.Attributes["instance_id"] {
			return fmt.Errorf(
				"instance_id is %s; want %s",
				state["instance_id"],
				vm.Primary.Attributes["instance_id"],
			)
		}
		//
		// if attr["public_ip"] != eipRs.Primary.Attributes["public_ip"] {
		// 	return fmt.Errorf(
		// 		"public_ip is %s; want %s",
		// 		attr["public_ip"],
		// 		eipRs.Primary.Attributes["public_ip"],
		// 	)
		// }

		return nil
	}
}

const testAccDataSourceOutscaleVmStateConfig = `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-update"
}

data "outscale_vm_state" "state" {
  instance_id = ["${outscale_vm.basic.id}"]
}
`
