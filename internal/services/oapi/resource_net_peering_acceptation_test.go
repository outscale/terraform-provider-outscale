package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_PeeringConnectionAccepter_sameAccount(t *testing.T) {
	resourceName := "outscale_net_peering_acceptation.peer"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(oapihelpers.GetAccepterOwnerId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "accepter_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "state.0.name", "active"),
				),
			},
		},
	})
}

func TestAccNet_PeeringConnectionAccepter_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering_acceptation.peer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(oapihelpers.GetAccepterOwnerId()),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_PeeringConnectionAccepter_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.1.1", testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(oapihelpers.GetAccepterOwnerId())),
	})
}

func testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(accountId string) string {
	return fmt.Sprintf(`
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
		source_net_id   = outscale_net.foo.id
		accepter_net_id = outscale_net.bar.id
		accepter_owner_id = "%s"

		tags {
			key   = "Side"
			value = "Accepter"
		}
	}

	// Accepter's side of the connection.
	resource "outscale_net_peering_acceptation" "peer" {
		net_peering_id = outscale_net_peering.foo.id
	}
`, accountId)
}
