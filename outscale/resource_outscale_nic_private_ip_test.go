package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleNetworkInterfacePrivateIPBasic(t *testing.T) {
	var conf fcu.NetworkInterface
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNetworkInterfacePrivateIPConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_private_ip.outscale_nic_private_ip", "network_interface_id"),
				),
			},
		},
	})
}

func testAccOutscaleNetworkInterfacePrivateIPConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_instance" {                 
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_lin" "outscale_lin" {
    cidr_block          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
    	network_interface_id    = "${outscale_nic.outscale_nic.id}"
}
`)
}
