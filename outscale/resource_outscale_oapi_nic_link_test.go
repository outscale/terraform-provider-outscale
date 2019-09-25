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

func TestAccOutscaleOAPINetworkInterfaceAttachmentBasic(t *testing.T) {
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
		CheckDestroy:  testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetworkInterfaceAttachmentConfigBasic(rInt),
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

func testAccOutscaleOAPINetworkInterfaceAttachmentConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_instance" {                 
    image_id  = "ami-3e158364"
    vm_type   = "c4.large"
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name = "eu-west-2a"
    ip_range       = "10.0.0.0/16"
    net_id         = "${outscale_net.outscale_net.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_nic_link" "outscale_nic_link" {
		device_number   = 1
		vm_id             = "${outscale_vm.outscale_instance.id}"
		#vm_id = "i-fcefc32f"   
		nic_id    = "${outscale_nic.outscale_nic.id}"
}

`)
}
