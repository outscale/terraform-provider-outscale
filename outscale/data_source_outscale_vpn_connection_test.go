package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAcc_VPNConnection_DataSource(t *testing.T) {
	t.Parallel()
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	dataSourceName := "data.outscale_vpn_connection.test"
	dataSourcesName := "data.outscale_vpn_connections.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_VPNConnection_DataSource_Config(publicIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "vpn_connection_id"),
					resource.TestCheckResourceAttr(dataSourcesName, "vpn_connections.#", "1"),
				),
			},
		},
	})
}

func testAcc_VPNConnection_DataSource_Config(publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = 3
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "foo" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			connection_type    = "ipsec.1"
			static_routes_only  = true
		}

		data "outscale_vpn_connection" "test" {
			filter {
				name = "vpn_connection_ids"
				values = [outscale_vpn_connection.foo.id]
			}
		}

		data "outscale_vpn_connections" "test" {
			filter {
				name = "vpn_connection_ids"
				values = [outscale_vpn_connection.foo.id]
			}
		}
	`, publicIP)
}
