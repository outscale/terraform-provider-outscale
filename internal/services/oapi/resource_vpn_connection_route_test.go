package oapi_test

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_VPNConnectionRoute_Basic(t *testing.T) {
	resourceName := "outscale_vpn_connection_route.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "destination_ip_range"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_connection_id"),
				),
			},
			testacc.ImportStepWithStateIdFunc(resourceName, vpnConnectionRouteStateIDFunc(resourceName), testacc.DefaultIgnores()...),
		},
	})
}

func vpnConnectionRouteStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["vpn_connection_id"], rs.Primary.Attributes["destination_ip_range"]), nil
	}
}

func TestAccOthers_VPNConnectionRoute_Migration(t *testing.T) {
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
		),
	})
}

func testAccOutscaleVPNConnectionRouteConfig(bgpAsn int, publicIP, destinationIPRange string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			connection_type = "ipsec.1"
		}

		resource "outscale_vpn_connection" "vpn_connection" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			connection_type    = "ipsec.1"
			static_routes_only  = true
		}

		resource "outscale_vpn_connection_route" "foo" {
			destination_ip_range = "%s"
			vpn_connection_id    = outscale_vpn_connection.vpn_connection.id
		}
	`, bgpAsn, publicIP, destinationIPRange)
}
