package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleENIDataSource_basic(t *testing.T) {
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
			resource.TestStep{
				Config: testAccOutscaleENIDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIDataSourceExists("outscale_nic.outscale_nic"),
					resource.TestCheckResourceAttrSet("data.outscale_nic.outscale_nic", "availability_zone"),
				),
			},
		},
	})
}

func testAccCheckOutscaleENIDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		return nil
	}
}

const testAccOutscaleENIDataSourceConfig = `
resource "outscale_lin" "outscale_lin" {
    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.id}"
}

data "outscale_nic" "outscale_nic" {
		network_interface_id = "${outscale_nic.outscale_nic.id}"
}
`
