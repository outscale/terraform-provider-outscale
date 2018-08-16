package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPILinPeeringConnection_basic(t *testing.T) {
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
				Config: testAccDataSourceOutscaleOAPILinPeeringConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPILinPeeringConnectionCheck("data.outscale_net_peering.test_by_id"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceOutscaleOAPILinPeeringConnectionCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		pcxRs, ok := s.RootModule().Resources["outscale_net_peering.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_net_peering.test in state")
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

const testAccDataSourceOutscaleOAPILinPeeringConnectionConfig = `
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
	lin_id = "${outscale_lin.foo.id}"
	peer_net_id = "${outscale_lin.bar.id}"

    tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    }
}

data "outscale_net_peering" "test_by_id" {
	lin_peering_id = "${outscale_net_peering.test.id}"
}
`
