package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_PeeringsConnectionDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleLinPeeringsConnectionConfig(oapihelpers.GetAccepterOwnerId()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_net_peerings.outscale_net_peerings", "net_peerings.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleLinPeeringsConnectionConfig(accountId string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.10.0.0/24"
		tags {
			key = "Name"
			value = "testacc-net-peerings-ds-net"
		}
	}

	resource "outscale_net" "outscale_net2" {
		ip_range = "10.31.0.0/16"
		tags {
			key = "Name"
			value = "testacc-net-peerings-ds-net2"
		}
	}

	resource "outscale_net_peering" "outscale_net_peering" {
		accepter_net_id = outscale_net.outscale_net.net_id
		source_net_id   = outscale_net.outscale_net2.net_id
		accepter_owner_id = "%s"
		tags {
			key = "name"
			value = "testacc-peerings-ds"
		}

	}

	data "outscale_net_peerings" "outscale_net_peerings" {
		filter {
			name   = "net_peering_ids"
			values = [outscale_net_peering.outscale_net_peering.net_peering_id]
		}
	}
`, accountId)
}
