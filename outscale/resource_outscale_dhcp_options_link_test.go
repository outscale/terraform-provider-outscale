package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleDHCPOptionsAssociation_basic(t *testing.T) {
	t.Skip()
	var v fcu.Vpc
	var d fcu.DhcpOptions

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDHCPOptionsAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDHCPOptionsAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDHCPOptionsExists("outscale_dhcp_option.foo", &d),
					//testAccCheckOutscaleOAPILinExists("outscale_lin.foo", &v), //TODO: fix once we refactor this resourceTestGetOMIByRegion
					testAccCheckDHCPOptionsAssociationExist("outscale_dhcp_option_link.foo", &v),
				),
			},
		},
	})
}

func testAccCheckDHCPOptionsAssociationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_dhcp_option_link" {
			continue
		}

		// Try to find the VPC associated to the DHCP Options set
		vpcs, err := findVPCsByDHCPOptionsID(conn, rs.Primary.Attributes["dhcp_options_id"])
		if err != nil {
			return err
		}

		if len(vpcs) > 0 {
			return fmt.Errorf("DHCP Options association is still associated to %d VPCs", len(vpcs))
		}
	}

	return nil
}

func testAccCheckDHCPOptionsAssociationExist(n string, vpc *fcu.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Options Set association ID is set")
		}

		if *vpc.DhcpOptionsId != rs.Primary.Attributes["dhcp_options_id"] {
			return fmt.Errorf("VPC %s does not have DHCP Options Set %s associated", *vpc.VpcId, rs.Primary.Attributes["dhcp_options_id"])
		}

		if *vpc.VpcId != rs.Primary.Attributes["vpc_id"] {
			return fmt.Errorf("DHCP Options Set %s is not associated with VPC %s", rs.Primary.Attributes["dhcp_options_id"], *vpc.VpcId)
		}

		return nil
	}
}

const testAccDHCPOptionsAssociationConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_dhcp_option" "foo" {}

resource "outscale_dhcp_option_link" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	dhcp_options_id = "${outscale_dhcp_option.foo.id}"
}
`
