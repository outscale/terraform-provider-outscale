package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPILinPeeringsConnection_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPILinPeeringsConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_net_peerings.test_by_id", "net_peering.#", "1"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccDataSourceOutscaleOAPILinPeeringsConnectionConfig = `
resource "outscale_net" "foo" {
  ip_range = "10.1.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-foo"
  }
}

resource "outscale_net" "bar" {
  ip_range = "10.2.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-bar"
  }
}

resource "outscale_net_peering" "test" {
	net_id = "${outscale_net.foo.id}"
	peer_net_id = "${outscale_net.bar.id}"

    tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    }
}

data "outscale_net_peerings" "test_by_id" {
	net_peering_id = "[${outscale_net_peering.test.id}]"
}
`
