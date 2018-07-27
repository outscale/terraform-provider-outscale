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
					resource.TestCheckResourceAttr("data.outscale_lin_peerings.test_by_id", "lin_peering.#", "1"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccDataSourceOutscaleOAPILinPeeringsConnectionConfig = `
resource "outscale_lin" "foo" {
  ip_range = "10.1.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-foo"
  }
}

resource "outscale_lin" "bar" {
  ip_range = "10.2.0.0/16"

  tag {
	  Name = "terraform-testacc-vpc-peering-connection-data-source-bar"
  }
}

resource "outscale_lin_peering" "test" {
	lin_id = "${outscale_lin.foo.id}"
	peer_lin_id = "${outscale_lin.bar.id}"

    tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    }
}

data "outscale_lin_peerings" "test_by_id" {
	lin_peering_id = "[${outscale_lin_peering.test.id}]"
}
`
