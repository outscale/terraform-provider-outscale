package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/spf13/cast"
)

func TestAccOthers_ClientGateway_basic(t *testing.T) {
	resourceName := "outscale_client_gateway.foo"
	rBgpAsn := oapihelpers.RandBgpAsn()
	rBgpAsnUpdated := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccCheckClientGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),

					resource.TestCheckResourceAttr(resourceName, "bgp_asn", cast.ToString(rBgpAsn)),
				),
			},
			{
				Config: testAccClientGatewayConfig(rBgpAsnUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),

					resource.TestCheckResourceAttr(resourceName, "bgp_asn", cast.ToString(rBgpAsnUpdated)),
				),
			},
		},
	})
}

func TestAccOthers_ClientGateway_withTags(t *testing.T) {
	resourceName := "outscale_client_gateway.foo"
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))
	valueUpdated := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccCheckClientGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayConfigWithTags(value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
				),
			},
			{
				Config: testAccClientGatewayConfigWithTags(valueUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
				),
			},
		},
	})
}

func testAccCheckClientGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no client gateway id is set")
		}

		filter := osc.ReadClientGatewaysRequest{
			Filters: &osc.FiltersClientGateway{
				ClientGatewayIds: &[]string{rs.Primary.ID},
			},
		}
		resp, err := client.ReadClientGateways(context.Background(), filter, options.WithRetryTimeout(DefaultTimeout))

		if err != nil || resp.ClientGateways == nil || len(*resp.ClientGateways) < 1 {
			return fmt.Errorf("outscale client gateway not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckClientGatewayDestroy(s *terraform.State) error {
	client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_client_gateway" {
			continue
		}

		filter := osc.ReadClientGatewaysRequest{
			Filters: &osc.FiltersClientGateway{
				ClientGatewayIds: &[]string{rs.Primary.ID},
			},
		}
		resp, err := client.ReadClientGateways(context.Background(), filter, options.WithRetryTimeout(DefaultTimeout))

		if err != nil || resp.ClientGateways == nil || len(*resp.ClientGateways) > 0 && *(ptr.From(resp.ClientGateways)[0]).State != "deleted" {
			return fmt.Errorf("outscale client gateway still exists (%s): %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccClientGatewayConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = %d
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"
		}
	`, rBgpAsn)
}

func testAccClientGatewayConfigWithTags(value string) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = 3
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"

			tags {
				key = "Name"
				value = "%s"
			}
		}
	`, value)
}
