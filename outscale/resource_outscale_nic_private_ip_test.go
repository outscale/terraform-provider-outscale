package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNet_NICPrivateIPBasic(t *testing.T) {
	var conf oscgo.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttr("outscale_nic_private_ip.outscale_nic_private_ip", "private_ips.#", "1"),
					resource.TestCheckResourceAttr("outscale_nic_private_ip.outscale_nic_private_ip", "private_ips.0", "10.0.45.67"),
					resource.TestCheckResourceAttrSet("outscale_nic_private_ip.outscale_nic_private_ip", "primary_private_ip")),
			},
		},
	})
}

func TestAccNet_Import_NIC_PrivateIP_Basic(t *testing.T) {
	resourceName := "outscale_nic_private_ip.outscale_nic_private_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(utils.GetRegion()),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name" 
				value = "testacc-nic-private-ip-rs"
			}	
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
