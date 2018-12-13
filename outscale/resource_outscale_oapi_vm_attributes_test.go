package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVMAttr_Basic(t *testing.T) {
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
		//CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMATTRConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vm", "deletion_protection", "true"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVMATTRConfigBasic() string {
	return `
resource "outscale_vm" "outscale_vm" {
  count = 1

  image_id               = "ami-5c450b62"
	vm_type                = "c4.large"
	security_group_ids     = ["sg-9752b7a6"]
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
  vm_id             = "${outscale_vm.outscale_vm.id}"
  keypair_name           = "testkp"
}`
}
