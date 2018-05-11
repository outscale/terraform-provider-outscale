package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPINICsDataSource_basic(t *testing.T) {
	var conf fcu.NetworkInterface

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPINICsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIDataSourceExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIDataSourceAttributes(&conf),
				),
			},
		},
	})
}

const testAccOutscaleOAPINICsDataSourceConfig = `
resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

data "outscale_nics" "nic" {
		network_interface_id = "NICID"
		subnet_id = "1"
}
`
