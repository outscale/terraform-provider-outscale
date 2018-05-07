package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleLinPeeringsConnection_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleLinPeeringsConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_lin_peerings.test_by_id", "vpc_peering_connection_set.#", "1"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccDataSourceOutscaleLinPeeringsConnectionConfig = `
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

data "outscale_lin_peerings" "test_by_id" {
	vpc_peering_connection_id = ["${outscale_lin_peering.test.id}"]
}
`
