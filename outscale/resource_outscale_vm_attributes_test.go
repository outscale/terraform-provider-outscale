package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVMAttr_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm_attributes.outscale_vm_attributes", "disable_api_termination", "false"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVMAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return nil
	}
}

func testAccCheckOutscaleVMConfig_basic() string {
	return `
resource "outscale_vm_attributes" "outscale_vm_attributes" {
  instance_id 	= "i-aebb385b"
	attribute   	= "disableApiTermination"
	disable_api_termination = false
}
`
}
