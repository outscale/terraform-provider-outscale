package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVMAttr_Basic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")
	omi := getOMIByRegion(region, "ubuntu").OMI
	vmType := "c4.large"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMATTRConfigBasic(omi, vmType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "false"),
				),
			},
			{
				Config: testAccCheckOutscaleOAPIVMATTRConfigBasicUpdate(omi, vmType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm_attributes.outscale_vm_attributes", "vm_initiated_shutdown_behavior", "terminate"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVMATTRConfigBasic(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "outscale_vm" {
			image_id           = "%s"
			vm_type            = "%s"
			keypair_name       = "terraform-basic"
		}
	`, omi, vmType)
}

func testAccCheckOutscaleOAPIVMATTRConfigBasicUpdate(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "outscale_vm" {
			image_id           = "%s"
			vm_type            = "%s"
			keypair_name       = "terraform-basic"
		}

		resource "outscale_vm_attributes" "outscale_vm_attributes" {
			vm_id = "${outscale_vm.outscale_vm.id}"
			vm_initiated_shutdown_behavior = "terminate"
		}
	`, omi, vmType)
}
