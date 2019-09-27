package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVMAttr_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
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

var region = os.Getenv("OUTSCALE_REGION")
var vmAmi = getOMIByRegion(region, "ubuntu").OMI
var vmType = "c4.large"

var testAccCheckOutscaleOAPIVMATTRConfigBasic = fmt.Sprintf(`
	resource "outscale_vm" "outscale_vm" {
		image_id           = "%s"
		vm_type            = "%s"
		keypair_name       = "terraform-basic"
	}
`, vmAmi, vmType)

var testAccCheckOutscaleOAPIVMATTRConfigBasicUpdate = fmt.Sprintf(`
	resource "outscale_vm" "outscale_vm" {
		image_id           = "%s"
		vm_type            = "%s"
		keypair_name       = "terraform-basic"
	}

	resource "outscale_vm_attributes" "outscale_vm_attributes" {
		vm_id = "${outscale_vm.outscale_vm.id}"
		vm_initiated_shutdown_behavior = "terminate"
	}
`, vmAmi, vmType)
