package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPINicsDataSource(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPINicsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNicsDataSourceID("data.outscale_nat_services.nat"),
					resource.TestCheckResourceAttr("data.outscale_nics.outscale_nics", "network_interface_set.#", "1"),
				),
			},
		},
	})
}

const testAccCheckOutscaleOAPINicsDataSourceConfig = `
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    ip_range          = "10.0.0.0/16"
    net_id              = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

data "outscale_nics" "outscale_nics" {
	nic_id = ["${outscale_nic.outscale_nic.id}"]
}
`
