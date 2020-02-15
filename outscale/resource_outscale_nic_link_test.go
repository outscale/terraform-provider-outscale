package outscale

import (
	"fmt"
	"os"
	"testing"

	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPINetworkInterfaceAttachmentBasic(t *testing.T) {
	//t.Skip()
	var conf oscgo.Nic
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfaceAttachmentConfigBasic(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttr(
						"outscale_nic_link.outscale_nic_link", "device_number", "1"),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_link.outscale_nic_link", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_link.outscale_nic_link", "nic_id"),
				),
			},
		},
	})
}

func testAccOutscaleOAPINetworkInterfaceAttachmentConfigBasic(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg1-test-group_test"
			net_id              = "${outscale_net.outscale_net.net_id}"
		}
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["${outscale_security_group.outscale_security_group.id}"]
			placement_subregion_name = "%[3]sa"
			subnet_id                = "${outscale_subnet.outscale_subnet.id}"
		}
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[3]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = "${outscale_net.outscale_net.id}"
		}
		resource "outscale_nic" "outscale_nic" {
			subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
		}
		resource "outscale_nic_link" "outscale_nic_link" {
			device_number = 1
			vm_id         = "${outscale_vm.vm.id}"
			nic_id        = "${outscale_nic.outscale_nic.id}"
		}
	`, omi, vmType, region)
}
