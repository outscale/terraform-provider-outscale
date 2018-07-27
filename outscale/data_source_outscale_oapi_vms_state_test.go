package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIVMSState(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err == nil {
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
				Config: testAccDataSourceOutscaleOAPIVMSStateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVMStateCheck("data.outscale_vm_state.state"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMSStateCheck(name string) resource.TestCheckFunc {
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

		if state["instance_id"] != vm.Primary.Attributes["instance_id"] {
			return fmt.Errorf(
				"instance_id is %s; want %s",
				state["instance_id"],
				vm.Primary.Attributes["instance_id"],
			)
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIVMSStateConfig = `
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}

data "outscale_vm_state" "state" {
  vm = ["${outscale_vm.basic.id}"]
}
`
