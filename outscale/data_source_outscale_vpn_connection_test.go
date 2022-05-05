package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleVPNConnectionDataSource_basic(t *testing.T) {
	t.Parallel()
	publicIP := fmt.Sprintf("172.0.0.%d", acctest.RandIntRange(1, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionDataSourceConfigBasic(publicIP),
			},
		},
	})
}

func TestAccOutscaleVPNConnectionDataSource_withFilters(t *testing.T) {
	t.Parallel()
	publicIP := fmt.Sprintf("172.0.0.%d", acctest.RandIntRange(1, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionDataSourceConfigWithFilters(publicIP),
			},
		},
	})
}

func testAccOutscaleVPNConnectionDataSourceConfigBasic(publicIP string) string {
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
			client_gateway_id  = "${outscale_client_gateway.customer_gateway.id}"
			virtual_gateway_id = "${outscale_virtual_gateway.virtual_gateway.id}"
			connection_type    = "ipsec.1"
			static_routes_only  = true

			tags {
        key   = "Name"
        value = "test-VPN"
			}
		}

		data "outscale_vpn_connection" "test" {
			vpn_connection_id = "${outscale_vpn_connection.foo.id}"
		}
	`, publicIP)
}

func testAccOutscaleVPNConnectionDataSourceConfigWithFilters(publicIP string) string {
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
			client_gateway_id  = "${outscale_client_gateway.customer_gateway.id}"
			virtual_gateway_id = "${outscale_virtual_gateway.virtual_gateway.id}"
			connection_type    = "ipsec.1"
			static_routes_only  = true
		}

		data "outscale_vpn_connection" "test" {
			filter {
				name = "vpn_connection_ids"
				values = ["${outscale_vpn_connection.foo.id}"]
			}
		}
	`, publicIP)
}
