package oapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_VPNConnectionRoute_basic(t *testing.T) {
	resourceName := "outscale_vpn_connection_route.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testacc.PreCheck(t) },
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionRouteExists(resourceName),
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
			PreCheck:     func() { testacc.PreCheck(t) },
			Providers:    testacc.SDKProviders,
			CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleVPNConnectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
					Check: resource.ComposeTestCheckFunc(
						testAccOutscaleVPNConnectionRouteExists(resourceName),
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

func testAccOutscaleVPNConnectionRouteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection Route ID is set")
		}

		destinationIPRange, vpnConnectionID := oapi.ParseVPNConnectionRouteID(rs.Primary.ID)

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnConnectionID},
			},
		}
		var resp oscgo.ReadVpnConnectionsResponse
		var err error
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			rp, httpResp, err := conn.VpnConnectionApi.ReadVpnConnections(context.Background()).ReadVpnConnectionsRequest(filter).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		vpnConnection := resp.GetVpnConnections()[0]

		var state string
		for _, route := range vpnConnection.GetRoutes() {
			if route.GetDestinationIpRange() == destinationIPRange {
				state = route.GetState()
			}
		}

		if err != nil || state == "deleted" {
			return fmt.Errorf("Outscale VPN Connection Route not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccOutscaleVPNConnectionRouteDestroy(s *terraform.State) error {
	conn := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSCAPI
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection_route" {
			continue
		}

		destinationIPRange, vpnConnectionID := oapi.ParseVPNConnectionRouteID(rs.Primary.ID)

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnConnectionID},
			},
		}
		var resp oscgo.ReadVpnConnectionsResponse
		var err error
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			rp, httpResp, err := conn.VpnConnectionApi.ReadVpnConnections(context.Background()).ReadVpnConnectionsRequest(filter).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		vpnConnection := resp.GetVpnConnections()[0]

		var state string
		for _, route := range vpnConnection.GetRoutes() {
			if route.GetDestinationIpRange() == destinationIPRange {
				state = route.GetState()
			}
		}

		if err != nil || state == "available" {
			return fmt.Errorf("Outscale VPN Connection Route still exists (%s): %s", rs.Primary.ID, err)
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
