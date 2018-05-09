package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleLinPeeringConnectionAccepter_sameAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOutscaleLinPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinPeeringConnectionAccepterExists("outscale_lin_peering_acceptation.peer"),
				),
			},
		},
	})
}

func testAccCheckOutscaleLinPeeringConnectionAccepterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccOutscaleLinPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_basic"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
}

// Accepter's side of the connection.
resource "outscale_lin_peering_acceptation" "peer" {
    vpc_peering_connection_id = "${outscale_lin_peering.foo.id}"

    tag {
       Side = "Accepter"
    }
}
`
