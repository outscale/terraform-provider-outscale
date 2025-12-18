package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_VPNConnectionsDataSource_basic(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionsDataSourceConfigBasic(bgpAsn, publicIP),
			},
		},
	})
}

func TestAccOthers_VPNConnectionsDataSource_withFilters(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionsDataSourceConfigWithFilters(bgpAsn, publicIP),
			},
		},
	})
}

func testAccOutscaleVPNConnectionsDataSourceConfigBasic(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "foo" {
			client_gateway_id   = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id  = outscale_virtual_gateway.virtual_gateway.id
			connection_type     = "ipsec.1"
			static_routes_only  = true

			tags {
				key   = "Name"
				value = "test-VPN"
			}
		}

		data "outscale_vpn_connections" "test" {
			vpn_connection_ids = [outscale_vpn_connection.foo.id]
		}
	`, bgpAsn, publicIP)
}

func testAccOutscaleVPNConnectionsDataSourceConfigWithFilters(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "foo" {
			client_gateway_id   = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id  = outscale_virtual_gateway.virtual_gateway.id
			connection_type     = "ipsec.1"
		}

		data "outscale_vpn_connections" "test" {
			filter {
				name = "vpn_connection_ids"
				values = [outscale_vpn_connection.foo.id]
			}
		}
	`, bgpAsn, publicIP)
}
