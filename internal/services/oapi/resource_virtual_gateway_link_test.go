package oapi_test

import (
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_VirtualGatewayLink(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.vgw_link"
	resourceName2 := "outscale_virtual_gateway.vgw"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetWithVGWLink,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "net_to_virtual_gateway_links.0.state", "attached"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnoresWith("net_to_virtual_gateway_links")...),
			{
				Config: testAccNetWithVGW,
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName2, "net_to_virtual_gateway_links.0.state", "detached"),
				),
			},
		},
	})
}

func TestAccNet_VirtualGatewayLink_NetsSwap(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.vgw_link"
	resourceName2 := "outscale_virtual_gateway.vgw"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetWithVGWLink,
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttr(resourceName2, "net_to_virtual_gateway_links.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "net_to_virtual_gateway_links.0.state", "attached"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnoresWith("net_to_virtual_gateway_links")...),
			{
				Config: testAccNetWithVGWLinkNetsSwap,
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttr(resourceName2, "net_to_virtual_gateway_links.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "net_to_virtual_gateway_links.0.state", "attached"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnoresWith("net_to_virtual_gateway_links")...),
		},
	})
}

func TestAccNet_VirtualGatewayLink_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.6.0", testAccNetWithVGWLink),
	})
}

const testAccNetWithVGW = `
resource "outscale_net" "net" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_virtual_gateway" "vgw" {
	connection_type = "ipsec.1"
}
`

const testAccNetWithVGWLink = testAccNetWithVGW + `
resource "outscale_virtual_gateway_link" "vgw_link" {
    virtual_gateway_id = outscale_virtual_gateway.vgw.id
    net_id              = outscale_net.net.id
}
`

const testAccNetWithVGWLinkNetsSwap = testAccNetWithVGW + `
resource "outscale_net" "net2" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_virtual_gateway_link" "vgw_link" {
    virtual_gateway_id = outscale_virtual_gateway.vgw.id
    net_id              = outscale_net.net2.id
}
`
