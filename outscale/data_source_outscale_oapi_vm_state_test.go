package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIVmState(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIVmStateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVMStateCheck("data.outscale_vm_state.state"),
					// testAccDataSourceOutscaleVMStateCheck("data.outscale_public_ip.by_public_ip"),
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

		vm, ok := s.RootModule().Resources["outscale_vm.basic"]
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

const testAccDataSourceOutscaleOAPIVmStateConfig = `
resource "outscale_vm" "basic" {
	image_id               = "ami-5c450b62"
	vm_type                = "c4.large"
	keypair_name           = "testkp"
	security_group_ids     = ["sg-9752b7a6"]
}

data "outscale_vm_state" "state" {
  vm_id = "${outscale_vm.basic.id}"
}
`
