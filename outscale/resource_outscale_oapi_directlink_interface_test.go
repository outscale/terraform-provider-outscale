package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
)

func TestAccOutscaleOAPIDirectLinkInterface_basic(t *testing.T) {
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIDirectLinkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDxPrivateVirtualInterfaceConfig(connectionID, vifName, bgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDirectLinkInterfaceExists("outscale_directlink_interface.foo"),
					resource.TestCheckResourceAttr("outscale_directlink_interface.foo", "virtual_interface_name", vifName),
				),
			},
			{
				ResourceName:      "outscale_directlink_interface.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOutscaleOAPIDirectLinkInterfaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).DL

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_directlink_interface" {
			continue
		}

		input := &dl.DescribeVirtualInterfacesInput{
			VirtualInterfaceID: aws.String(rs.Primary.ID),
		}

		var err error
		var resp *dl.DescribeVirtualInterfacesOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.DescribeVirtualInterfaces(input)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}
		for _, v := range resp.VirtualInterfaces {
			if *v.VirtualInterfaceID == rs.Primary.ID && !(*v.VirtualInterfaceState == "deleted") {
				return fmt.Errorf("[DESTROY ERROR] Dx Private VIF (%s) not deleted", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIDirectLinkInterfaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccOAPIDxPrivateVirtualInterfaceConfig(cid, n string, bgpAsn int) string {
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
		direct_link_interface_name = "%s"
		vlan           = 4094
		bgp_asn        = %d
	}
}
`, n, cid, n, bgpAsn)
}
