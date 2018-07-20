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

func TestAccOutscaleDSOAPIDLVPNGateways_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDLVPNGatewaysDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDLVPNGatewaysDataSourceID("data.outscale_directlink_vpn_gateways.test"),
					resource.TestCheckResourceAttrSet("data.outscale_directlink_vpn_gateways.test", "virtual_gateways.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIDLVPNGatewaysDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DL VPN Gateways data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DL VPN Gateways data source ID not set")
		}
		return nil
	}
}

func testAccOAPIDLVPNGatewaysDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_directlink_vpn_gateways" "test" {}
	`)
}
