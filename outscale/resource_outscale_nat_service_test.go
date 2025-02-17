package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_WithNatService_basic(t *testing.T) {
	resourceName := "outscale_nat_service.nat_service"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPINatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "available"),
				),
			},
		},
	})
}

const testAccOAPINatGatewayConfig = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags {
			key = "Name"
			value = "testacc-nat-service-rs"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id   = outscale_net.outscale_net.net_id
		ip_range = "10.0.0.0/18"
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_nat_service" "nat_service" {
		depends_on   = [outscale_route.outscale_route]
		subnet_id    = outscale_subnet.outscale_subnet.subnet_id
		public_ip_id = outscale_public_ip.outscale_public_ip.public_ip_id
	}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = outscale_net.outscale_net.net_id
	}

	resource "outscale_route" "outscale_route" {
		depends_on   = [outscale_route_table_link.outscale_route_table_link]
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = outscale_internet_service_link.outscale_internet_service_link.internet_service_id
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = outscale_subnet.outscale_subnet.subnet_id
		route_table_id = outscale_route_table.outscale_route_table.id
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = outscale_net.outscale_net.net_id
		internet_service_id = outscale_internet_service.outscale_internet_service.id
	}
`
