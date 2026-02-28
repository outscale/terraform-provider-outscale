package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceVirtualGateway_unattached(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVirtualGatewayUnattachedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.outscale_virtual_gateway.test_by_id", "id",
						"outscale_virtual_gateway.unattached", "id"),
					resource.TestCheckResourceAttrSet("data.outscale_virtual_gateway.test_by_id", "state"),
					resource.TestCheckNoResourceAttr("data.outscale_virtual_gateway.test_by_id", "attached_vpc_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVirtualGatewayUnattachedConfig() string {
	return `
		resource "outscale_virtual_gateway" "unattached" {
			connection_type = "ipsec.1"
		}

		data "outscale_virtual_gateway" "test_by_id" {
			virtual_gateway_id = outscale_virtual_gateway.unattached.id
		}
	`
}
