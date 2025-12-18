package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOutscaleVPNConnection_basic(t *testing.T) {
	resourceName := "outscale_vpn_connection.vpn_basic"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionConfig(bgpAsn, publicIP, true),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttrSet(resourceName, "vgw_telemetries.#"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "true"),
				),
			},
			{
				Config: testAccOutscaleVPNConnectionConfig(bgpAsn, publicIP, false),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttrSet(resourceName, "vgw_telemetries.#"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "false"),
				),
			},
		},
	})
}

func TestAccOutscaleVPNConnection_withoutStaticRoutes(t *testing.T) {
	resourceName := "outscale_vpn_connection.foo"
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(0, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_connection.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionConfigWithoutStaticRoutes(bgpAsn, publicIP),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
				),
			},
		},
	})
}

func TestAccOutscaleVPNConnection_withTags(t *testing.T) {
	resourceName := "outscale_vpn_connection.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionConfigWithTags(bgpAsn, publicIP, value),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					// resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					// resource.TestCheckResourceAttr(resourceName, "tags.0.value", value),
				),
			},
		},
	})
}

func TestAccOutscaleVPNConnection_importBasic(t *testing.T) {
	resourceName := "outscale_vpn_connection.vpn_basic"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccOutscaleVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNConnectionConfig(bgpAsn, publicIP, true),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "true"),
				),
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

func testAccOutscaleVPNConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection ID is set")
		}

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				VpnConnectionIds: &[]string{rs.Primary.ID},
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

		if err != nil || len(resp.GetVpnConnections()) < 1 {
			return fmt.Errorf("Outscale VPN Connection not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccOutscaleVPNConnectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_connection" {
			continue
		}

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				VpnConnectionIds: &[]string{rs.Primary.ID},
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

		if err != nil ||
			len(resp.GetVpnConnections()) > 0 && resp.GetVpnConnections()[0].GetState() != "deleted" {
			return fmt.Errorf("Outscale VPN Connection still exists (%s): %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccOutscaleVPNConnectionConfig(bgpAsn int, publicIP string, staticRoutesOnly bool) string {
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

func testAccOutscaleVPNConnectionConfigWithoutStaticRoutes(bgpAsn int, publicIP string) string {
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

func testAccOutscaleVPNConnectionConfigWithTags(bgpAsn int, publicIP, value string) string {
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
			static_routes_only = true


			tags {
				key   = "Name"
				value = "%s"
			}
		}
	`, bgpAsn, publicIP, value)
}
