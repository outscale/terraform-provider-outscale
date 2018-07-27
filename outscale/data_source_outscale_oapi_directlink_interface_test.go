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

func TestAccOutscaleOAPIDSDirectLinkInterface_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	key := "OUTSCALE_CONNECTION_ID"
	connectionID := os.Getenv(key)
	if connectionID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}
	vifName := fmt.Sprintf("terraform-testacc-dxvif-%s", acctest.RandString(5))
	bgpAsn := acctest.RandIntRange(64512, 65534)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDSDxPrivateVirtualInterfaceConfig(connectionID, vifName, bgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDSDirectLinkInterfaceExists("data.outscale_directlink_interface.outscale_directlink_interface"),
					resource.TestCheckResourceAttr("data.outscale_directlink_interface.outscale_directlink_interface", "virtual_interface_name", vifName),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIDSDirectLinkInterfaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccOAPIDSDxPrivateVirtualInterfaceConfig(cid, n string, bgpAsn int) string {
	return fmt.Sprintf(`
resource "outscale_vpn_gateway" "foo" {
  tag {
    Name = "%s"
  }
}

resource "outscale_directlink_interface" "foo" {
  connection_id    = "%s"

	new_private_virtual_interface {
		vpn_gateway_id = "${outscale_vpn_gateway.foo.id}"
		direct_link_Interface_name = "%s"
		vlan           = 4094
		bgp_asn        = %d
	}
}

data "outscale_directlink_interface" "outscale_directlink_interface" {
  direct_link_interface_id = "${outscale_directlink_interface.outscale_directlink_interface.id}"
}
`, n, cid, n, bgpAsn)
}
