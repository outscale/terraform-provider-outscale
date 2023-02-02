package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_NetPeering_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_net_peering.net_peering"
	dataSourcesName := "data.outscale_net_peerings.net_peerings"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPILinPeeringConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "net_peering_id"),

					resource.TestCheckResourceAttr(dataSourcesName, "net_peerings.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPILinPeeringConnectionConfig = `
	resource "outscale_net" "net" {
		ip_range = "10.10.0.0/24"
		tags {
			key = "Name"
			value = "testacc-net-peering-ds-net"
		}
	}

	resource "outscale_net" "net2" {
		ip_range = "10.11.0.0/24"
		tags {
			key = "Name"
			value = "testacc-net-peering-ds-net2"
		}
	}

	resource "outscale_net_peering" "net_peering" {
		accepter_net_id = outscale_net.net.net_id
		source_net_id   = outscale_net.net2.net_id
	}

	data "outscale_net_peering" "net_peering" {
		filter {
			name   = "net_peering_ids"
			values = [outscale_net_peering.net_peering.net_peering_id]
		}
	}

	data "outscale_net_peerings" "net_peerings" {
		filter {
			name   = "net_peering_ids"
			values = [outscale_net_peering.net_peering.net_peering_id]
		}
	}
`
