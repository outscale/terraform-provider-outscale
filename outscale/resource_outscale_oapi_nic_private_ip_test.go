package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPINetworkInterfacePrivateIPBasic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	var conf oapi.Nic
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
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
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
	vm_type   = "c4.large"
    image_id  = "ami-880caa66"
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_net" "outscale_net" {
    ip_range          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name = "eu-west-2a"
    ip_range       = "10.0.0.0/16"
    net_id         = "${outscale_net.outscale_net.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
	nic_id    = "${outscale_nic.outscale_nic.nic_id}"
	private_ips = ["10.0.45.67"]
}
`)
}
