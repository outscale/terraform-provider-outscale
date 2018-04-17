package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleNatServiceDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleNatServiceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNatServiceDataSourceID("data.outscale_nat_service.nat"),
					resource.TestCheckResourceAttr("data.outscale_nat_service.nat", "subnet_id", "subnet-861fbecc"),
				),
			},
		},
	})
}

func testAccCheckOutscaleNatServiceDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Nat Service data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Nat Service data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleNatServiceDataSourceConfig = `
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

data "outscale_nat_service" "nat" {
	nat_gateway_id = "${outscale_nat_service.gateway.id}"
}
`
