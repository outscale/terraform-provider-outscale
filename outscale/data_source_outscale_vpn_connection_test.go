package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_VPNConnectionDataSource_basic(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionDataSourceConfigBasic(bgpAsn, publicIP),
			},
		},
	})
}

func TestAccOthers_VPNConnectionDataSource_withFilters(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionDataSourceConfigWithFilters(bgpAsn, publicIP),
			},
		},
	})
}

func testAccOutscaleVPNConnectionDataSourceConfigBasic(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway1" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway1" {
			bgp_asn         = %d
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "foo1" {
			client_gateway_id  = outscale_client_gateway.customer_gateway1.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway1.id
			connection_type    = "ipsec.1"
			static_routes_only  = true

			tags {
			      key   = "Name"
			      value = "test-VPN"
			}
		}

		data "outscale_vpn_connection" "test" {
			vpn_connection_id = outscale_vpn_connection.foo1.id
		}
	`, bgpAsn, publicIP)
}

func testAccOutscaleVPNConnectionDataSourceConfigWithFilters(bgpAsn int, publicIP string) string {
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
	`, bgpAsn, publicIP)
}
