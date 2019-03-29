package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleLinPeeringConnection_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleLinPeeringConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleLinPeeringConnectionCheck("data.outscale_lin_peering.test_by_id"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceOutscaleLinPeeringConnectionCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		pcxRs, ok := s.RootModule().Resources["outscale_lin_peering.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_lin_peering.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != pcxRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				pcxRs.Primary.Attributes["id"],
			)
		}

		return nil
	}
}

const testAccDataSourceOutscaleLinPeeringConnectionConfig = `
resource "outscale_lin" "foo" {
  cidr_block = "10.1.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-foo"
  }
}

resource "outscale_lin" "bar" {
  cidr_block = "10.2.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-bar"
  }
}

resource "outscale_lin_peering" "test" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"

    tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    }
}

data "outscale_lin_peering" "test_by_id" {
	vpc_peering_connection_id = "${outscale_lin_peering.test.id}"
}
`
