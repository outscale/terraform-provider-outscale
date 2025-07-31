package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_withpnGatewayAttachment_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_virtual_gateway_link.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckVpnGatewayAttachmentImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckVpnGatewayAttachmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func TestAccNet_WithVpnGatewayAttachment_deleted(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_virtual_gateway_link.test"
	resourceName2 := "outscale_virtual_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		IDRefreshName:            resourceName2,
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
