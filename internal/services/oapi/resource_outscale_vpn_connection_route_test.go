package oapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_VPNConnectionRoute_basic(t *testing.T) {
	resourceName := "outscale_vpn_connection_route.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionRouteExists(t.Context(), resourceName),
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
			Providers:    testacc.SDKProviders,
			CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
					Check: resource.ComposeTestCheckFunc(
						testAccOutscaleVPNConnectionRouteExists(t.Context(), resourceName),
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

func testAccOutscaleVPNConnectionRouteExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpn connection route id is set")
		}

		destinationIPRange, vpnconnectionID := oapihelpers.ParseVPNConnectionRouteID(rs.Primary.ID)

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnconnectionID},
			},
		}
		resp, err := client.ReadVpnConnections(ctx, filter, options.WithRetryTimeout(DefaultTimeout))

		vpnConnection := ptr.From(resp.VpnConnections)[0]

		var state string
		for _, route := range ptr.From(vpnConnection.Routes) {
			if route.DestinationIpRange == destinationIPRange {
				state = route.State
			}
		}

		if err != nil || state == "deleted" {
			return fmt.Errorf("outscale vpn connection route not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccOutscaleVPNConnectionRouteDestroy(s *terraform.State) error {
	client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection_route" {
			continue
		}

		destinationIPRange, vpnConnectionID := oapihelpers.ParseVPNConnectionRouteID(rs.Primary.ID)

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnConnectionID},
			},
		}
		resp, err := client.ReadVpnConnections(context.Background(), filter, options.WithRetryTimeout(DefaultTimeout))

		vpnConnection := ptr.From(resp.VpnConnections)[0]

		var state string
		for _, route := range ptr.From(vpnConnection.Routes) {
			if route.DestinationIpRange == destinationIPRange {
				state = route.State
			}
		}

		if err != nil || state == "available" {
			return fmt.Errorf("outscale vpn connection route still exists (%s): %s", rs.Primary.ID, err)
		}

	}
	return nil
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
