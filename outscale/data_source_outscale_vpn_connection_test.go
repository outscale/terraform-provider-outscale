package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVpnConnectionDataSource_basic(t *testing.T) {
	t.Skip()

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
				Config: testAccOutscaleVpnConnectionDataSourceConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_vpn_connection.test", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.outscale_vpn_connection.test", "state"),
				),
			},
		},
	})
}

func testAccOutscaleVpnConnectionDataSourceConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_vpn_gateway" "vpn_gateway" {
		  tag {
		    Name = "vpn_gateway"
		  }
		}

		resource "outscale_client_endpoint" "customer_gateway" {
		  bgp_asn = %d
		  ip_address = "178.0.0.1"
		  type = "ipsec.1"
			tag {
				Name = "main-customer-gateway"
			}
		}

		resource "outscale_vpn_connection" "foo" {
		  vpn_gateway_id = "${outscale_vpn_gateway.vpn_gateway.id}"
		  customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
		  type = "ipsec.1"
		  options {
				static_routes_only = true
			}
		}

		data "outscale_vpn_connection" "test" {
    	vpn_connection_id = "${outscale_vpn_connection.foo.id}"
		}
`, rBgpAsn)
}
