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
resource "outscale_vm" "basic" {
	image_id               = "ami-5c450b62"
	vm_type                = "c4.large"
	keypair_name           = "testkp"
	security_group_ids     = ["sg-9752b7a6"]
}

data "outscale_vm_state" "state" {
  vm_id = ["${outscale_vm.basic.id}"]
}
`
