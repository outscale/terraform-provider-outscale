package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_WithVirtualRoutePropagation_basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway_route_propagation.outscale_virtual_gateway_route_propagation"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnRoutePropagationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
				),
			},
		},
	})
}

func testAccOutscaleVpnRoutePropagationConfig() string {
	return `
		resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
 connection_type = "ipsec.1"
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/18"
}
resource "outscale_virtual_gateway_link" "outscale_virtual_gateway_link" {
    virtual_gateway_id = outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "terraform-RT"
    }
}

resource "outscale_virtual_gateway_route_propagation" "outscale_virtual_gateway_route_propagation" {
virtual_gateway_id = outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id
    route_table_id  = outscale_route_table.outscale_route_table.route_table_id
    enable = true
}
	`
}
