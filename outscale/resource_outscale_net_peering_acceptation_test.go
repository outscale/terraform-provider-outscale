package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccLinPeeringConnectionAccepter_sameAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLinPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLinPeeringConnectionAccepterSameAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinPeeringConnectionAccepterExists("outscale_net_peering_acceptation.peer"),
				),
			},
		},
	})
}

func testAccCheckLinPeeringConnectionAccepterExists(n string) resource.TestCheckFunc {
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

func testAccLinPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccLinPeeringConnectionAccepterSameAccountConfig = `
	resource "outscale_net" "foo" {
		ip_range = "10.0.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-acceptation-rs-foo"
		}
	}

	resource "outscale_net" "bar" {
		ip_range = "10.1.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-acceptation-rs-bar"
		}
	}

	resource "outscale_net_peering" "foo" {
		source_net_id   = "${outscale_net.foo.id}"
		accepter_net_id = "${outscale_net.bar.id}"

		tags {
			key   = "Side"
			value = "Accepter"
		}
	}

	// Accepter's side of the connection.
	resource "outscale_net_peering_acceptation" "peer" {
		net_peering_id = "${outscale_net_peering.foo.id}"
	}
`
