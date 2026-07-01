package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_VirtualRoutePropagation_Basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway_route_propagation.outscale_virtual_gateway_route_propagation"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpnRoutePropagation + testAccVpnRoutePropagationSwitch(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
				),
			},
			{
				Config: testAccVpnRoutePropagation + testAccVpnRoutePropagationSwitch(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
				),
			},
		},
	})
}

func TestAccNet_VirtualRoutePropagation_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.6.0", testAccVpnRoutePropagation+testAccVpnRoutePropagationSwitch(true)),
	})
}

const testAccVpnRoutePropagation = `
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
}
`

func testAccVpnRoutePropagationSwitch(enabled bool) string {
	return fmt.Sprintf(`
resource "outscale_virtual_gateway_route_propagation" "outscale_virtual_gateway_route_propagation" {
	virtual_gateway_id = outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id
    route_table_id  = outscale_route_table.outscale_route_table.route_table_id
    enable = %t
}
`, enabled)
}
