package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPINetworkInterfacePrivateIPBasic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var conf fcu.NetworkInterface
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_private_ip.outscale_nic_private_ip", "nic_id"),
				),
			},
		},
	})
}

func testAccOutscaleOAPINetworkInterfacePrivateIPConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_instance" {                 
    image_id                    = "ami-880caa66"
    type               = "c4.large"
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_net" "outscale_net" {
    ip_range          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    ip_range          = "10.0.0.0/16"
    lin_id              = "${outscale_lin.outscale_lin.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
    	nic_id    = "${outscale_nic.outscale_nic.id}"
}
`)
}
