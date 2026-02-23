package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_PeeringconnectionDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleLinPeeringconnectionConfig(oapihelpers.GetAccepterOwnerId()),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleLinPeeringconnectionCheck("outscale_net_peering.net_peering"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleLinPeeringconnectionCheck(name string) resource.TestCheckFunc {
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

func testAccDataSourceOutscaleLinPeeringconnectionConfig(accountId string) string {
	return fmt.Sprintf(`
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
		accepter_owner_id = "%s"
	}

	data "outscale_net_peering" "net_peering_data" {
		filter {
			name   = "net_peering_ids"
			values = [outscale_net_peering.net_peering.net_peering_id]
		}
	}
`, accountId)
}
