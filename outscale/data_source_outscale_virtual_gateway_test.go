package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_VirtualGateway_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_virtual_gateway.vg"
	dataSourcesName := "data.outscale_virtual_gateways.vgs"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_VirtualGateway_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "virtual_gateways.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_gateway_id"),
				),
			},
		},
	})
}

func testAcc_VirtualGateway_DataSource_Config() string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "unattached" {
			connection_type = "ipsec.1"	
		}
		
		data "outscale_virtual_gateway" "vg" {
			filter {
				name = "virtual_gateway_ids"
				values = ["${outscale_virtual_gateway.unattached.id}"]
			}
		}
		data "outscale_virtual_gateways" "vgs" {
			filter {
				name = "virtual_gateway_ids"
				values = ["${outscale_virtual_gateway.unattached.id}"]
			}
		}
	`)
}
