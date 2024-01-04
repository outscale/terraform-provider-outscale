package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_ClientGatewaysDatasource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleClientGatewaysDatasourceConfigBasic,
			},
		},
	})
}

func TestAccOthers_ClientGatewaysDatasource_withFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleClientGatewaysDatasourceConfigWithFilters,
			},
		},
	})
}

const testAccOutscaleClientGatewaysDatasourceConfigBasic = `
	resource "outscale_client_gateway" "foo1" {
		bgp_asn         = 3
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"

		tags {
			key   = "Name"
			value = "%s"
		}
	}

	resource "outscale_client_gateway" "foo2" {
		bgp_asn         = 4
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	data "outscale_client_gateways" "test" {
		client_gateway_ids = [outscale_client_gateway.foo1.id, outscale_client_gateway.foo2.id]
	}
`

const testAccOutscaleClientGatewaysDatasourceConfigWithFilters = `
	resource "outscale_client_gateway" "foo1" {
		bgp_asn         = 3
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"

		tags {
			key   = "Name"
			value = "%s"
		}
	}

	resource "outscale_client_gateway" "foo2" {
		bgp_asn         = 4
		public_ip       = "172.0.0.1"
		connection_type = "ipsec.1"
	}

	data "outscale_client_gateways" "test" {
		filter {
			name = "client_gateway_ids"
			values = [outscale_client_gateway.foo1.id, outscale_client_gateway.foo2.id]
		}
	}
`
