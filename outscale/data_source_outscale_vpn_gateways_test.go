package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVpnGateways_unattached(t *testing.T) {
	t.Skip()

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIVpnGatewaysUnattachedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vpn_gateways.test_by_id", "vpn_gateway_set.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpnGatewaysUnattachedConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "unattached" {
			tag {
				Name = "terraform-testacc-vpn-gateway-data-source-unattached-%d"
				ABC  = "testacc-%d"
				XYZ  = "testacc-%d"
			}
		}
		
		data "outscale_vpn_gateways" "test_by_id" {
			vpn_gateway_id = ["${outscale_vpn_gateway.unattached.id}"]
		}
	`, rInt, rInt+1, rInt-1)
}
