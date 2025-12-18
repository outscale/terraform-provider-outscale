package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_ClientGatewaysDatasource_basic(t *testing.T) {
	bgpAsn1 := utils.RandBgpAsn()
	bgpAsn2 := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleClientGatewaysDatasourceConfigBasic(bgpAsn1, bgpAsn2),
			},
		},
	})
}

func TestAccOthers_ClientGatewaysDatasource_withFilters(t *testing.T) {
	bgpAsn1 := utils.RandBgpAsn()
	bgpAsn2 := utils.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleClientGatewaysDatasourceConfigWithFilters(bgpAsn1, bgpAsn2),
			},
		},
	})
}

func testAccOutscaleClientGatewaysDatasourceConfigBasic(asn1, asn2 int) string {
	return fmt.Sprintf(`
	resource "outscale_client_gateway" "foo1" {
		bgp_asn         = %[1]d
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	resource "outscale_client_gateway" "foo2" {
		bgp_asn         = %[2]d
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	data "outscale_client_gateways" "test" {
		client_gateway_ids = [outscale_client_gateway.foo1.id, outscale_client_gateway.foo2.id]
	}
`, asn1, asn2)
}

func testAccOutscaleClientGatewaysDatasourceConfigWithFilters(asn1, asn2 int) string {
	return fmt.Sprintf(`
	resource "outscale_client_gateway" "foo1" {
		bgp_asn         = %[1]d
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	resource "outscale_client_gateway" "foo2" {
		bgp_asn         = %[2]d
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	data "outscale_client_gateways" "test" {
		filter {
			name = "client_gateway_ids"
			values = [outscale_client_gateway.foo1.id, outscale_client_gateway.foo2.id]
		}
	}
	`, asn1, asn2)
}
