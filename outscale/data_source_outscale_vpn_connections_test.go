package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_VPNConnectionsDataSource_basic(t *testing.T) {
	t.Parallel()
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionsDataSourceConfigBasic(publicIP),
			},
		},
	})
}

func TestAccOthers_VPNConnectionsDataSource_withFilters(t *testing.T) {
	t.Parallel()
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionsDataSourceConfigWithFilters(publicIP),
			},
		},
	})
}

func testAccOutscaleVPNConnectionsDataSourceConfigBasic(publicIP string) string {
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
			client_gateway_id   = "${outscale_client_gateway.customer_gateway.id}"
			virtual_gateway_id  = "${outscale_virtual_gateway.virtual_gateway.id}"
			connection_type     = "ipsec.1"
			static_routes_only  = true

			tags {
        key   = "Name"
        value = "test-VPN"
			}
		}

		data "outscale_vpn_connections" "test" {
			vpn_connection_ids = ["${outscale_vpn_connection.foo.id}"]
		}
	`, publicIP)
}

func testAccOutscaleVPNConnectionsDataSourceConfigWithFilters(publicIP string) string {
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
			client_gateway_id   = "${outscale_client_gateway.customer_gateway.id}"
			virtual_gateway_id  = "${outscale_virtual_gateway.virtual_gateway.id}"
			connection_type     = "ipsec.1"
		}

		data "outscale_vpn_connections" "test" {
			filter {
				name = "vpn_connection_ids"
				values = ["${outscale_vpn_connection.foo.id}"]
			}
		}
	`, publicIP)
}
