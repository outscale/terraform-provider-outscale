package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVpnRoutePropagation_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rBgpAsn := acctest.RandIntRange(64512, 65534)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnRoutePropagationConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVpnRoutePropagation(
						"outscale_lin.vpc",
						"outscale_vpn_gateway.vpn_gateway",
						"outscale_lin_internet_gateway.foo",
						"outscale_route_table.foo",
						"outscale_vpn_gateway_route_propagation.foo",
					),
				),
			},
		},
	})
}

func testAccOutscaleVpnRoutePropagation(
	vpcResource string,
	vpnGatewayResource string,
	linInterGatewayResource string,
	routeTable string,
	routeProp string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[routeProp]
		if !ok {
			return fmt.Errorf("Not found: %s", routeProp)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		// fcuconn := testAccProvider.Meta().(*OutscaleClient).FCU

		return nil
	}
}

func testAccOutscaleVpnRoutePropagationConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
    count = 1

    type = "ipsec.1" 
}

resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_vpn_gateway_route_propagation" "outscale_vpn_gateway_route_propagation" {
    gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
		route_table_id  = "${outscale_route_table.outscale_route_table.route_table_id}"
		depends_on = ["outscale_vpn_gateway.outscale_vpn_gateway", "outscale_route_table.outscale_route_table"]
}
`)
}
