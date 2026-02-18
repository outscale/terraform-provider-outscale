package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOutscaleVPNclientection_basic(t *testing.T) {
	resourceName := "outscale_vpn_clientection.vpn_basic"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testacc.PreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccOutscaleVPNclientectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionConfig(bgpAsn, publicIP, true),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "clientection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttrSet(resourceName, "vgw_telemetries.#"),

					resource.TestCheckResourceAttr(resourceName, "clientection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "true"),
				),
			},
			{
				Config: testAccOutscaleVPNclientectionConfig(bgpAsn, publicIP, false),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "clientection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttrSet(resourceName, "vgw_telemetries.#"),

					resource.TestCheckResourceAttr(resourceName, "clientection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "false"),
				),
			},
		},
	})
}

func TestAccOutscaleVPNclientection_withoutStaticRoutes(t *testing.T) {
	resourceName := "outscale_vpn_clientection.foo"
	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(0, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testacc.PreCheck(t) },
		IDRefreshName: "outscale_vpn_clientection.foo",
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccOutscaleVPNclientectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionConfigWithoutStaticRoutes(bgpAsn, publicIP),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "clientection_type"),

					resource.TestCheckResourceAttr(resourceName, "clientection_type", "ipsec.1"),
				),
			},
		},
	})
}

func TestAccOutscaleVPNclientection_withTags(t *testing.T) {
	resourceName := "outscale_vpn_clientection.foo"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testacc.PreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccOutscaleVPNclientectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionConfigWithTags(bgpAsn, publicIP, value),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "clientection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "clientection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					// resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					// resource.TestCheckResourceAttr(resourceName, "tags.0.value", value),
				),
			},
		},
	})
}

func TestAccOutscaleVPNclientection_importBasic(t *testing.T) {
	resourceName := "outscale_vpn_clientection.vpn_basic"

	publicIP := fmt.Sprintf("172.0.0.%d", utils.RandIntRange(1, 255))
	bgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testacc.PreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccOutscaleVPNclientectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVPNclientectionConfig(bgpAsn, publicIP, true),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVPNclientectionExists(t.Context(), resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "clientection_type"),
					resource.TestCheckResourceAttrSet(resourceName, "static_routes_only"),
					resource.TestCheckResourceAttr(resourceName, "clientection_type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "true"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func testAccOutscaleVPNclientectionExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpn clientection id is set")
		}

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				VpnConnectionIds: &[]string{rs.Primary.ID},
			},
		}

		resp, err := client.ReadVpnConnections(ctx, filter, options.WithRetryTimeout(DefaultTimeout))

		if err != nil || resp.VpnConnections == nil || len(*resp.VpnConnections) < 1 {
			return fmt.Errorf("outscale vpn clientection not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccOutscaleVPNclientectionDestroy(s *terraform.State) error {
	client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_clientection" {
			continue
		}

		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				VpnConnectionIds: &[]string{rs.Primary.ID},
			},
		}
		resp, err := client.ReadVpnConnections(context.Background(), filter, options.WithRetryTimeout(DefaultTimeout))

		if err != nil || resp.VpnConnections == nil || len(*resp.VpnConnections) > 0 && *(*resp.VpnConnections)[0].State != "deleted" {
			return fmt.Errorf("outscale vpn connection still exists (%s): %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccOutscaleVPNclientectionConfig(bgpAsn int, publicIP string, staticRoutesOnly bool) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "vpn_basic" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			clientection_type    = "ipsec.1"
			static_routes_only = "%t"
		}
	`, bgpAsn, publicIP, staticRoutesOnly)
}

func testAccOutscaleVPNclientectionConfigWithoutStaticRoutes(bgpAsn int, publicIP string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "foo" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			clientection_type    = "ipsec.1"
		}
	`, bgpAsn, publicIP)
}

func testAccOutscaleVPNclientectionConfigWithTags(bgpAsn int, publicIP, value string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "virtual_gateway" {
			clientection_type = "ipsec.1"
		}

		resource "outscale_client_gateway" "customer_gateway" {
			bgp_asn         = %d
			public_ip       = "%s"
			clientection_type = "ipsec.1"
		}

		resource "outscale_vpn_clientection" "foo" {
			client_gateway_id  = outscale_client_gateway.customer_gateway.id
			virtual_gateway_id = outscale_virtual_gateway.virtual_gateway.id
			clientection_type    = "ipsec.1"
			static_routes_only = true


			tags {
				key   = "Name"
				value = "%s"
			}
		}
	`, bgpAsn, publicIP, value)
}
