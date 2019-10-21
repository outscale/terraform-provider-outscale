package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVpnConnectionsDataSource_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVpnConnectionsDataSourceConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vpn_connections.test", "vpn_connection_set.#", "1"),
				),
			},
		},
	})
}

func testAccOutscaleOAPIVpnConnectionsDataSourceConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "vpn_gateway" {
			tag {
				Name = "vpn_gateway"
			}
		}
		
		resource "outscale_client_endpoint" "customer_gateway" {
			bgp_asn    = %d
			ip_address = "178.0.0.1"
			type       = "ipsec.1"
		
			tag {
				Name = "main-customer-gateway"
			}
		}
		
		resource "outscale_vpn_connection" "foo" {
			vpn_gateway_id     = "${outscale_vpn_gateway.vpn_gateway.id}"
			client_endpoint_id = "${outscale_client_endpoint.customer_gateway.id}"
			type               = "ipsec.1"
		
			vpn_connection_option {
				static_routes_only = true
			}
		}
		
		data "outscale_vpn_connections" "test" {
			vpn_connection_id = ["${outscale_vpn_connection.foo.id}"]
		}
	`, rBgpAsn)
}
