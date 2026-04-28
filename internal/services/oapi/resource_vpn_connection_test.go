package oapi_test

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_VPNConnection_Basic(t *testing.T) {
	resourceName := "outscale_vpn_connection.vpn_basic"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: vpnConnectionConfig(bgpAsn, publicIP, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttrSet(resourceName, "vgw_telemetries.#"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "true"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_VPNConnection_WithoutStaticRoutes(t *testing.T) {
	resourceName := "outscale_vpn_connection.foo"
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(0, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: vpnConnectionConfigWithoutStaticRoutes(bgpAsn, publicIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "false"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_VPNConnection_Migration(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			vpnConnectionConfig(bgpAsn, publicIP, true),
			vpnConnectionConfigWithoutStaticRoutes(bgpAsn, publicIP),
		),
	})
}

func vpnConnectionConfig(bgpAsn int, publicIP string, staticRoutesOnly bool) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "vpn_basic" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			connection_type    = "ipsec.1"
			static_routes_only = "%t"
		}
	`, bgpAsn, publicIP, staticRoutesOnly)
}

func vpnConnectionConfigWithoutStaticRoutes(bgpAsn int, publicIP string) string {
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
		}
	`, bgpAsn, publicIP)
}
