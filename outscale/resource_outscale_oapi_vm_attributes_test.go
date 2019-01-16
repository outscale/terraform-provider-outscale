package outscale

import (
	"fmt"
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMATTRConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "true"),
				),
			},
			{
				Config: testAccCheckOutscaleOAPIVMATTRConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm_attributes.outscale_vm_attributes", "vm_initiated_shutdown_behavior", "terminate"),
				),
			},
		},
	})
}

var vmAmi = "ami-880caa66"
var vmType = "c4.large"

var testAccCheckOutscaleOAPIVMATTRConfigBasic = fmt.Sprintf(`
resource "outscale_vm" "outscale_vm" {
	image_id = "%s"
	vm_type = "%s"
	keypair_name = "integ_sut_keypair"
}`, vmAmi, vmType)

var testAccCheckOutscaleOAPIVMATTRConfigBasicUpdate = fmt.Sprintf(`
resource "outscale_vm" "outscale_vm" {
	image_id = "%s"
	vm_type = "%s"
	keypair_name = "integ_sut_keypair"
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
	vm_id = "${outscale_vm.outscale_vm.id}"
	vm_initiated_shutdown_behavior = "terminate"
}`, vmAmi, vmType)
