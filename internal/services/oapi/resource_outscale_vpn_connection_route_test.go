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

func TestAccOthers_VPNclientectionRoute_basic(t *testing.T) {
	resourceName := "outscale_vpn_clientection_route.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testacc.PreCheck(t) },
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccOutscaleVPNclientectionRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionRouteExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "destination_ip_range"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_clientection_id"),
				),
			},
		},
	})
}

func TestAccOthers_ImportVPNclientectionRoute_basic(t *testing.T) {
	if os.Getenv("TEST_QUOTA") == "true" {
		resourceName := "outscale_vpn_clientection_route.foo"

		publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
		destinationIPRange := fmt.Sprintf("172.168.%d.0/24", utils.RandIntRange(1, 255))
		bgpAsn := oapihelpers.RandBgpAsn()

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:     func() { testacc.PreCheck(t) },
			Providers:    testacc.SDKProviders,
			CheckDestroy: testAccOutscaleVPNclientectionRouteDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccOutscaleVPNclientectionRouteConfig(bgpAsn, publicIP, destinationIPRange),
					Check: resource.ComposeTestCheckFunc(
						testAccOutscaleVPNclientectionRouteExists(t.Context(), resourceName),
						resource.TestCheckResourceAttrSet(resourceName, "destination_ip_range"),
						resource.TestCheckResourceAttrSet(resourceName, "vpn_clientection_id"),
					),
				},
				testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func testAccOutscaleVPNclientectionRouteExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpn clientection route id is set")
		}

		destinationIPRange, vpnclientectionID := oapihelpers.ParseVPNclientectionRouteID(rs.Primary.ID)

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VpnConnectionIds:         &[]string{vpnclientectionID},
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
			return fmt.Errorf("outscale vpn clientection route not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccOutscaleVPNclientectionRouteDestroy(s *terraform.State) error {
	client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_clientection_route" {
			continue
		}

		destinationIPRange, vpnclientectionID := oapihelpers.ParseVPNclientectionRouteID(rs.Primary.ID)

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{destinationIPRange},
				VirtualGatewayIds:        &[]string{vpnclientectionID},
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
			return fmt.Errorf("outscale vpn clientection route still exists (%s): %s", rs.Primary.ID, err)
		}

	}
	return nil
}

func testAccOutscaleVPNclientectionRouteConfig(bgpAsn int, publicIP, destinationIPRange string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "vpn_clientection" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			clientection_type    = "ipsec.1"
			static_routes_only  = true
		}

		resource "outscale_vpn_clientection_route" "foo" {
			destination_ip_range = "%s"
			vpn_clientection_id    = outscale_vpn_clientection.vpn_clientection.id
		}
	`, bgpAsn, publicIP, destinationIPRange)
}
