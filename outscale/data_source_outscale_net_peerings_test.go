package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_PeeringsConnectionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPILinPeeringsConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_net_peerings.outscale_net_peerings", "net_peerings.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPILinPeeringsConnectionConfig = `
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
		accepter_net_id = "${outscale_net.outscale_net.net_id}"
		source_net_id   = "${outscale_net.outscale_net2.net_id}"
		tags {
			key = "okht"
			value = "testacc-peerings-ds"
		}

	}

	resource "outscale_net_peering" "outscale_net_peering2" {
		accepter_net_id = "${outscale_net.outscale_net.net_id}"
		source_net_id   = "${outscale_net.outscale_net2.net_id}"
	}

	data "outscale_net_peerings" "outscale_net_peerings" {
		filter {
			name   = "net_peering_ids"
			values = ["${outscale_net_peering.outscale_net_peering.net_peering_id}"]
		}
	}
`
