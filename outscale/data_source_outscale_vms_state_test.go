package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleVMSState(t *testing.T) {
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
				Config: testAccDataSourceOutscaleVMSStateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_vms_state.state", "instance_status_set.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleVMSStateConfig = `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-update"
}
resource "outscale_vm" "basic2" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-update"
}

data "outscale_vms_state" "state" {
  instance_id = ["${outscale_vm.basic.id}", "${outscale_vm.basic2.id}"]
}
`
