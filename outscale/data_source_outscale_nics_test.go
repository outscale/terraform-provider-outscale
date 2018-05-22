package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleNicsDataSource(t *testing.T) {
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
				Config: testAccCheckOutscaleNicsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNicsDataSourceID("data.outscale_nat_services.nat"),
					resource.TestCheckResourceAttr("data.outscale_nics.outscale_nics", "network_interface_set.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleNicsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find NICS data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NICS data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleNicsDataSourceConfig = `
resource "outscale_lin" "outscale_lin" {
    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

data "outscale_nics" "outscale_nics" {
	network_interface_id = ["${outscale_nic.outscale_nic.id}"]
}
`
