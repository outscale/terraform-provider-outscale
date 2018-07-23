package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSDirectLinkInterfaces_basic(t *testing.T) {
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
				Config: testAccOAPIDSDxPrivateVirtualInterfacesConfig(connectionID, vifName, bgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDSDirectLinkInterfacesExists("data.outscale_directlink_interfaces.outscale_directlink_interfaces"),
					resource.TestCheckResourceAttr("data.outscale_directlink_interfaces.outscale_directlink_interfaces", "virtual_interfaces.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIDSDirectLinkInterfacesExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccOAPIDSDxPrivateVirtualInterfacesConfig(cid, n string, bgpAsn int) string {
	return fmt.Sprintf(`
resource "outscale_vpn_gateway" "foo" {
  tag {
    Name = "%s"
  }
}

resource "outscale_directlink_interface" "foo" {
  connection_id    = "%s"

	new_private_virtual_interface {
		virtual_gateway_id = "${outscale_vpn_gateway.foo.id}"
		virtual_interface_name = "%s"
		vlan           = 4094
		asn        = %d
	}
}

data "outscale_directlink_interfaces" "outscale_directlink_interfaces" {
  connection_id = %s
}
`, n, cid, n, bgpAsn, cid)
}
