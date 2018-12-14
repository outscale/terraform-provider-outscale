package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
