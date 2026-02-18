package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_VPNclientectionDataSource_basic(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionDataSourceConfigBasic(bgpAsn, publicIP),
			},
		},
	})
}

func TestAccOthers_VPNclientectionDataSource_withFilters(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionDataSourceConfigWithFilters(bgpAsn, publicIP),
			},
		},
	})
}

func testAccOutscaleVPNclientectionDataSourceConfigBasic(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway1" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway1" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "foo1" {
			client_gateway_id  = outscale_client_gateway.customer_gateway1.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway1.id
			clientection_type    = "ipsec.1"
			static_routes_only  = true

			tags {
			      key   = "Name"
			      value = "test-VPN"
			}
		}

		data "outscale_vpn_clientection" "test" {
			vpn_clientection_id = outscale_vpn_clientection.foo1.id
		}
	`, bgpAsn, publicIP)
}

func testAccOutscaleVPNclientectionDataSourceConfigWithFilters(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "foo" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			clientection_type    = "ipsec.1"
			static_routes_only  = true
		}

		data "outscale_vpn_clientection" "test" {
			filter {
				name = "vpn_clientection_ids"
				values = [outscale_vpn_clientection.foo.id]
			}
		}
	`, bgpAsn, publicIP)
}
