package outscale

import (
	"fmt"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPINetworkInterfacePrivateIPBasic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")
	var conf oscgo.Nic

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttr("outscale_nic_private_ip.outscale_nic_private_ip", "private_ips.#", "1"),
					resource.TestCheckResourceAttr("outscale_nic_private_ip.outscale_nic_private_ip", "private_ips.0", "10.0.45.67"),
					resource.TestCheckResourceAttrSet("outscale_nic_private_ip.outscale_nic_private_ip", "primary_private_ip")),
			},
		},
	})
}

func testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}
		
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
		}
		
		resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
			nic_id      = "${outscale_nic.outscale_nic.nic_id}"
			private_ips = ["10.0.45.67"]
		}
	`, region)
}
