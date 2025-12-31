package oapi_test

import (
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_withpnGatewayAttachment_basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
				),
			},
		},
	})
}

func TestAccNet_VpnGatewayAttachment_importBasic(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
			},
			testacc.ImportStep(resourceName, "request_id"),
		},
	})
}

func TestAccNet_WithVpnGatewayAttachment_deleted(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.test"
	resourceName2 := "outscale_virtual_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName2, "net_to_virtual_gateway_links.#"),
					resource.TestCheckResourceAttr(resourceName, "net_to_virtual_gateway_links.0.state", "attached"),
				),
			},
			{
				Config: testAccNoOAPIVpnGatewayAttachmentConfig,
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName2, "net_to_virtual_gateway_links.0.state", "detached"),
				),
			},
			// Ignore attributes related to the Gateway Link, that gets populated after a refresh
			testacc.ImportStep(resourceName2, "net_to_virtual_gateway_links", "request_id"),
		},
	})
}

const testAccNoOAPIVpnGatewayAttachmentConfig = `
	resource "outscale_net" "test" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_virtual_gateway" "test" {
		connection_type = "ipsec.1"
}
`

const testAccOAPIVpnGatewayAttachmentConfig = `
resource "outscale_virtual_gateway" "test" {
 connection_type = "ipsec.1"
}
resource "outscale_net" "test" {
    ip_range = "10.0.0.0/18"
}
resource "outscale_virtual_gateway_link" "test" {
    virtual_gateway_id = outscale_virtual_gateway.test.virtual_gateway_id
    net_id              = outscale_net.test.net_id
}
`
