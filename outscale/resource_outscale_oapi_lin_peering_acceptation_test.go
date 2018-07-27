package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPILinPeeringConnectionAccepter_sameAccount(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOutscaleOAPILinPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinPeeringConnectionAccepterSameAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinPeeringConnectionAccepterExists("outscale_lin_peering_acceptation.peer"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPILinPeeringConnectionAccepterExists(n string) resource.TestCheckFunc {
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

func testAccOutscaleOAPILinPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccOutscaleOAPILinPeeringConnectionAccepterSameAccountConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleOAPILinPeeringConnection_basic"
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
