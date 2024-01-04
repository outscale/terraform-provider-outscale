package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_PeeringConnectionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPILinPeeringConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPILinPeeringConnectionCheck("outscale_net_peering.net_peering"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPILinPeeringConnectionCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		pcxRs, ok := s.RootModule().Resources["outscale_net_peering.net_peering"]
		if !ok {
			return fmt.Errorf("can't find outscale_net_peering.net_peering in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != pcxRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				pcxRs.Primary.Attributes["id"],
			)
		}

		return nil
	}
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

	data "outscale_net_peering" "net_peering_data" {
		filter {
			name   = "net_peering_ids"
			values = [outscale_net_peering.net_peering.net_peering_id]
		}
	}
`
