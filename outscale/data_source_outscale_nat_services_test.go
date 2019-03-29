package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleNatServicesDataSource_Instance(t *testing.T) {
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
				Config: testAccCheckOutscaleNatServicesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNatServiceDataSourceID("data.outscale_nat_services.nat"),
					resource.TestCheckResourceAttr("data.outscale_nat_services.nat", "nat_gateway.#", "1"),
					resource.TestCheckResourceAttr("data.outscale_nat_services.nat", "nat_gateway.0.subnet_id", "subnet-861fbecc"),
				),
			},
		},
	})
}

const testAccCheckOutscaleNatServicesDataSourceConfig = `
resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	cidr_block = "10.0.0.0/16"
	vpc_id = "${outscale_lin.vpc.id}"
}

resource "outscale_public_ip" "bar" {}

resource "outscale_nat_service" "gateway" {
    allocation_id = "${outscale_public_ip.bar.id}"
    subnet_id = "${outscale_subnet.subnet.id}"
}

data "outscale_nat_services" "nat" {
	nat_gateway_id = ["${outscale_nat_service.gateway.id}"]
}
`
