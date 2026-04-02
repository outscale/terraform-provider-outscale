package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_VPNConnectionRoute_basic(t *testing.T) {
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
		},
	})
}

func TestAccOthers_ImportVPNConnectionRoute_basic(t *testing.T) {
	if os.Getenv("TEST_QUOTA") == "true" {
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
				testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
			},
		})
	} else {
		t.Skip("will be done soon")
	}
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
