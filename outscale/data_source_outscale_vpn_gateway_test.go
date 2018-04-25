package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleVpnGateway_unattached(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleVpnGatewayUnattachedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.outscale_vpn_gateway.test_by_id", "id",
						"outscale_vpn_gateway.unattached", "id"),
					resource.TestCheckResourceAttrSet("data.outscale_vpn_gateway.test_by_id", "state"),
					resource.TestCheckNoResourceAttr("data.outscale_vpn_gateway.test_by_id", "attached_vpc_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVpnGatewayUnattachedConfig(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vpn_gateway" "unattached" {
    tag {
		Name = "terraform-testacc-vpn-gateway-data-source-unattached-%d"
      	ABC  = "testacc-%d"
		XYZ  = "testacc-%d"
    }
}

data "outscale_vpn_gateway" "test_by_id" {
	vpn_gateway_id = "${outscale_vpn_gateway.unattached.id}"
}

`, rInt, rInt+1, rInt-1)
}
