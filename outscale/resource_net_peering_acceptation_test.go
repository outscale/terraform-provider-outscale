package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"
)

func TestAccNet_PeeringConnectionAccepter_sameAccount(t *testing.T) {
	resourceName := "outscale_net_peering_acceptation.peer"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(utils.GetAccepterOwnerId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "accepter_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "state.0.name", "active"),
				),
			},
		},
	})
}

func TestAccNet_PeeringConnectionAccepter_sameAccount_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.1", testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(utils.GetAccepterOwnerId())),
	})
}

func TestAccNet_PeeringConnectionAccepter_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering_acceptation.peer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(utils.GetAccepterOwnerId()),
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

func TestAccNet_PeeringConnectionAccepter_importBasic_Migration(t *testing.T) {
	resourceName := "outscale_net_peering_acceptation.peer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinPeeringConnectionAccepterSameAccountConfig(utils.GetAccepterOwnerId()),
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
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
