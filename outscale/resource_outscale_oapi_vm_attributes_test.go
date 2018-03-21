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
				Config: testAccCheckOutscaleOAPIVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vm", "deletion_protection", "true"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVMATTRConfig_basic() string {
	return `
resource "outscale_vm" "outscale_vm" {
  count = 1

  image_id                = "ami-880caa66"
  type           = "c4.large"
  deletion_protection = false
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
  vm_id             = "${outscale_vm.outscale_vm.0.id}"
  attribute               = "disableApiTermination"
  deletion_protection = true
}`
}
