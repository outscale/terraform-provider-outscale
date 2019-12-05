package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPINetworkInterfaceAttachmentBasic(t *testing.T) {
	t.Skip()
	var conf oapi.Nic
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
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
		resource "outscale_vm" "outscale_instance" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%[3]sa"
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
			vm_id         = "${outscale_vm.outscale_instance.id}"
			nic_id        = "${outscale_nic.outscale_nic.id}"
		}
	`, omi, vmType, region)
}
