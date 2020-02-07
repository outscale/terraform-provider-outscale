package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIVpnConnectionRoute_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccOutscaleOAPIVpnConnectionRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIVpnConnectionRouteConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleOAPIVpnConnectionRoute(
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_client_endpoint.customer_gateway",
						"outscale_vpn_connection.vpn_connection",
						"outscale_vpn_connection_route.foo",
					),
				),
			},
			resource.TestStep{
				Config: testAccOutscaleOAPIVpnConnectionRouteConfigUpdate(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleOAPIVpnConnectionRoute(
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_client_endpoint.customer_gateway",
						"outscale_vpn_connection.vpn_connection",
						"outscale_vpn_connection_route.foo",
					),
				),
			},
		},
	})
}

func testAccOutscaleOAPIVpnConnectionRouteDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection_route" {
			continue
		}

		cidrBlock, vpnConnectionID := resourceOutscaleOAPIVpnConnectionRouteParseID(rs.Primary.ID)

		routeFilters := []*fcu.Filter{
			&fcu.Filter{
				Name:   aws.String("route.destination-cidr-block"),
				Values: []*string{aws.String(cidrBlock)},
			},
			&fcu.Filter{
				Name:   aws.String("vpn-connection-id"),
				Values: []*string{aws.String(vpnConnectionID)},
			},
		}

		var resp *fcu.DescribeVpnConnectionsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
				Filters: routeFilters,
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVpnConnectionID.NotFound" {
				// not found, all good
				return nil
			}
			return err
		}

		var vpnc *fcu.VpnConnection
		if resp != nil {
			// range over the connections and isolate the one we created
			for _, v := range resp.VpnConnections {
				if *v.VpnConnectionId == vpnConnectionID {
					vpnc = v
				}
			}

			if vpnc == nil {
				// vpn connection not found, so that's good...
				return nil
			}

			if vpnc.State != nil && *vpnc.State == "deleted" {
				return nil
			}
		}

	}
	return fmt.Errorf("Fall through error, Check Destroy criteria not met")
}

func testAccOutscaleOAPIVpnConnectionRoute(
	vpnGatewayResource string,
	customerGatewayResource string,
	vpnConnectionResource string,
	vpnConnectionRouteResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[vpnConnectionRouteResource]
		if !ok {
			return fmt.Errorf("Not found: %s", vpnConnectionRouteResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		route, ok := s.RootModule().Resources[vpnConnectionRouteResource]
		if !ok {
			return fmt.Errorf("Not found: %s", vpnConnectionRouteResource)
		}

		cidrBlock, vpnConnectionID := resourceOutscaleOAPIVpnConnectionRouteParseID(route.Primary.ID)

		routeFilters := []*fcu.Filter{
			&fcu.Filter{
				Name:   aws.String("route.destination-cidr-block"),
				Values: []*string{aws.String(cidrBlock)},
			},
			&fcu.Filter{
				Name:   aws.String("vpn-connection-id"),
				Values: []*string{aws.String(vpnConnectionID)},
			},
		}

		FCU := testAccProvider.Meta().(*OutscaleClient).FCU

		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = FCU.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
				Filters: routeFilters,
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccOutscaleOAPIVpnConnectionRouteConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "vpn_gateway" {
			tag {
				Name = "vpn_gateway"
			}
		}
		
		resource "outscale_client_endpoint" "customer_gateway" {
			bgp_asn    = %d
			ip_address = "182.0.0.1"
			type       = "ipsec.1"
		}
		
		resource "outscale_vpn_connection" "vpn_connection" {
			vpn_gateway_id      = "${outscale_vpn_gateway.vpn_gateway.id}"
			customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
			type                = "ipsec.1"
		
			options {
				static_routes_only = true
			}
		}
		
		resource "outscale_vpn_connection_route" "foo" {
			destination_ip_range = "172.168.10.0/24"
			vpn_connection_id    = "${outscale_vpn_connection.vpn_connection.id}"
		}
	`, rBgpAsn)
}

// Change destination_ip_range
func testAccOutscaleOAPIVpnConnectionRouteConfigUpdate(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "vpn_gateway" {
			tag {
				Name = "vpn_gateway"
			}
		}
		
		resource "outscale_client_endpoint" "customer_gateway" {
			bgp_asn    = %d
			ip_address = "182.0.0.1"
			type       = "ipsec.1"
		}
		
		resource "outscale_vpn_connection" "vpn_connection" {
			vpn_gateway_id      = "${outscale_vpn_gateway.vpn_gateway.id}"
			customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
			type                = "ipsec.1"
		
			options {
				static_routes_only = true
			}
		}
		
		resource "outscale_vpn_connection_route" "foo" {
			destination_ip_range = "172.168.20.0/24"
			vpn_connection_id    = "${outscale_vpn_connection.vpn_connection.id}"
		}
	`, rBgpAsn)
}
