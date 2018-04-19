package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVpnConnection_basic(t *testing.T) {
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	var vpn fcu.VpnConnection

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_connection.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVpnConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnConnectionConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVpnConnection(
						"outscale_lin.vpc",
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_client_endpoint.customer_gateway",
						"outscale_vpn_connection.foo",
						&vpn,
					),
				),
			},
		},
	})
}

func TestAccOutscaleVpnConnection_withoutStaticRoutes(t *testing.T) {
	rInt := acctest.RandInt()
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	var vpn fcu.VpnConnection
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_connection.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVpnConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnConnectionConfigUpdate(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVpnConnection(
						"outscale_lin.vpc",
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_client_endpoint.customer_gateway",
						"outscale_vpn_connection.foo",
						&vpn,
					),
					resource.TestCheckResourceAttr("outscale_vpn_connection.foo", "static_routes_only", "false"),
				),
			},
		},
	})
}

func TestAccOutscaleVpnConnection_disappears(t *testing.T) {
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	var vpn fcu.VpnConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOutscaleVpnConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnConnectionConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVpnConnection(
						"outscale_lin.vpc",
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_client_endpoint.customer_gateway",
						"outscale_vpn_connection.foo",
						&vpn,
					),
					testAccOutscaleVpnConnectionDisappears(&vpn),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccOutscaleVpnConnectionDisappears(connection *fcu.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.VM.DeleteVpnConnection(&fcu.DeleteVpnConnectionInput{
				VpnConnectionId: connection.VpnConnectionId,
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
			if fcuerr, ok := err.(awserr.Error); ok && fcuerr.Code() == "InvalidVpnConnectionID.NotFound" {
				return nil
			}
			if err != nil {
				return err
			}
		}

		return resource.Retry(40*time.Minute, func() *resource.RetryError {
			opts := &fcu.DescribeVpnConnectionsInput{
				VpnConnectionIds: []*string{connection.VpnConnectionId},
			}

			var resp *fcu.DescribeVpnConnectionsOutput
			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				resp, err = conn.VM.DescribeVpnConnections(opts)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return resource.NonRetryableError(err)
			})

			if err != nil {
				cgw, ok := err.(awserr.Error)
				if ok && cgw.Code() == "InvalidVpnConnectionID.NotFound" {
					return nil
				}
				if ok && cgw.Code() == "IncorrectState" {
					return resource.RetryableError(fmt.Errorf(
						"Waiting for VPN Connection to be in the correct state: %v", connection.VpnConnectionId))
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving VPN Connection: %s", err))
			}
			if *resp.VpnConnections[0].State == "deleted" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for VPN Connection: %v", connection.VpnConnectionId))
		})
	}
}

func testAccOutscaleVpnConnectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection" {
			continue
		}

		var resp *fcu.DescribeVpnConnectionsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
				VpnConnectionIds: []*string{aws.String(rs.Primary.ID)},
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
			if fcuerr, ok := err.(awserr.Error); ok && fcuerr.Code() == "InvalidVpnConnectionID.NotFound" {
				// not found
				return nil
			}
			return err
		}

		var vpn *fcu.VpnConnection
		for _, v := range resp.VpnConnections {
			if v.VpnConnectionId != nil && *v.VpnConnectionId == rs.Primary.ID {
				vpn = v
			}
		}

		if vpn == nil {
			// vpn connection not found
			return nil
		}

		if vpn.State != nil && *vpn.State == "deleted" {
			return nil
		}

	}

	return nil
}

func testAccOutscaleVpnConnection(
	vpcResource string,
	vpnGatewayResource string,
	customerGatewayResource string,
	vpnConnectionResource string,
	vpnConnection *fcu.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[vpnConnectionResource]
		if !ok {
			return fmt.Errorf("Not found: %s", vpnConnectionResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		connection, ok := s.RootModule().Resources[vpnConnectionResource]
		if !ok {
			return fmt.Errorf("Not found: %s", vpnConnectionResource)
		}

		fcuconn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeVpnConnectionsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = fcuconn.VM.DescribeVpnConnections(&fcu.DescribeVpnConnectionsInput{
				VpnConnectionIds: []*string{aws.String(connection.Primary.ID)},
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

		*vpnConnection = *resp.VpnConnections[0]

		return nil
	}
}

func testAccOutscaleVpnConnectionConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "vpn_gateway" {
		  tag {
		    Name = "vpn_gateway"
		  }
		}

		resource "outscale_client_endpoint" "customer_gateway" {
		  bgp_asn = %d
		  ip_address = "178.0.0.1"
		  type = "ipsec.1"
			tag {
				Name = "main-customer-gateway"
			}
		}

		resource "outscale_vpn_connection" "foo" {
		  vpn_gateway_id = "${outscale_vpn_gateway.vpn_gateway.id}"
		  customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
		  type = "ipsec.1"
		  options {
				static_routes_only = true
			}
		}
		`, rBgpAsn)
}

// Change static_routes_only to be false, forcing a refresh.
func testAccOutscaleVpnConnectionConfigUpdate(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
	resource "outscale_vpn_gateway" "vpn_gateway" {
	  tag {
	    Name = "vpn_gateway"
	  }
	}

	resource "outscale_client_endpoint" "customer_gateway" {
	  bgp_asn = %d
	  ip_address = "178.0.0.1"
	  type = "ipsec.1"
		tag {
	    Name = "main-customer-gateway-%d"
	  }
	}

	resource "outscale_vpn_connection" "foo" {
	  vpn_gateway_id = "${outscale_vpn_gateway.vpn_gateway.id}"
	  customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
	  type = "ipsec.1"
	  static_routes_only = false
	}
	`, rBgpAsn, rInt)
}

// Test our VPN tunnel config XML parsing
const testAccOutscaleVpnTunnelInfoXML = `
<vpn_connection id="vpn-abc123">
  <ipsec_tunnel>
    <vpn_gateway>
      <tunnel_outside_address>
        <ip_address>SECOND_ADDRESS</ip_address>
      </tunnel_outside_address>
    </vpn_gateway>
    <ike>
      <pre_shared_key>SECOND_KEY</pre_shared_key>
    </ike>
  </ipsec_tunnel>
  <ipsec_tunnel>
    <vpn_gateway>
      <tunnel_outside_address>
        <ip_address>FIRST_ADDRESS</ip_address>
      </tunnel_outside_address>
    </vpn_gateway>
    <ike>
      <pre_shared_key>FIRST_KEY</pre_shared_key>
    </ike>
  </ipsec_tunnel>
</vpn_connection>
`
