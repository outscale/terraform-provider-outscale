package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOthers_VPNConnectionRoute_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_vpn_connection_route.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionRouteConfig(publicIP, destinationIPRange),
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
	t.Parallel()
	if os.Getenv("TEST_QUOTA") == "true" {
		resourceName := "outscale_vpn_connection_route.foo"

		publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
		destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccOutscaleVPNConnectionRouteDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleVPNConnectionRouteConfig(publicIP, destinationIPRange),
					Check: resource.ComposeTestCheckFunc(
						testAccOutscaleVPNConnectionRouteExists(resourceName),
						resource.TestCheckResourceAttrSet(resourceName, "destination_ip_range"),
						resource.TestCheckResourceAttrSet(resourceName, "vpn_connection_id"),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateIdFunc:       testAccCheckOutscaleOAPIRouteImportStateIDFunc(resourceName),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"request_id"},
				},
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

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection Route ID is set")
		}

		destinationIPRange, vpnConnectionID := resourceOutscaleVPNConnectionRouteParseID(rs.Primary.ID)

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnConnectionID},
			},
		}
		var resp oscgo.ReadVpnConnectionsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection_route" {
			continue
		}

		destinationIPRange, vpnConnectionID := resourceOutscaleVPNConnectionRouteParseID(rs.Primary.ID)

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnConnectionID},
			},
		}
		var resp oscgo.ReadVpnConnectionsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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

func testAccOutscaleVPNConnectionRouteConfig(publicIP, destinationIPRange string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			connection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = 3
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
	`, publicIP, destinationIPRange)
}
