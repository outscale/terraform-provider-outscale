package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_VPNclientectionsDataSource_basic(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionsDataSourceConfigBasic(bgpAsn, publicIP),
			},
		},
	})
}

func TestAccOthers_VPNclientectionsDataSource_withFilters(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionsDataSourceConfigWithFilters(bgpAsn, publicIP),
			},
		},
	})
}

func testAccOutscaleVPNclientectionsDataSourceConfigBasic(bgpAsn int, publicIP string) string {
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
			client_gateway_id   = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id  = outscale_virtual_gateway.virtual_gateway.id
			clientection_type     = "ipsec.1"
			static_routes_only  = true

			tags {
				key   = "Name"
				value = "test-VPN"
			}
		}

		data "outscale_vpn_clientections" "test" {
			vpn_clientection_ids = [outscale_vpn_clientection.foo.id]
		}
	`, bgpAsn, publicIP)
}

func testAccOutscaleVPNclientectionsDataSourceConfigWithFilters(bgpAsn int, publicIP string) string {
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
			client_gateway_id   = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id  = outscale_virtual_gateway.virtual_gateway.id
			clientection_type     = "ipsec.1"
		}

		data "outscale_vpn_clientections" "test" {
			filter {
				name = "vpn_clientection_ids"
				values = [outscale_vpn_clientection.foo.id]
			}
		}
	`, bgpAsn, publicIP)
}
