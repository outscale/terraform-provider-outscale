package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOthers_DataSourceVirtualGateway_unattached(t *testing.T) {
	//t.Skip()
	t.Parallel()
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVirtualGatewayUnattachedConfig(rInt),
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

func testAccDataSourceOutscaleOAPIVirtualGatewayUnattachedConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "unattached" {
			connection_type = "ipsec.1"	
		}
		
		data "outscale_virtual_gateway" "test_by_id" {
			virtual_gateway_id = outscale_virtual_gateway.unattached.id
		}
	`)
}
