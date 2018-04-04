package outscale

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIVpc_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("172.%d.0.0/16", rInt)
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVpcConfig(cidr, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVpcCheck("data.outscale_lin.by_id", cidr, tag),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpcCheck(name, cidr, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		vpcRs, ok := s.RootModule().Resources["outscale_lin.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_lin.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr_block"] != cidr {
			return fmt.Errorf("bad cidr_block %s, expected: %s", attr["cidr_block"], cidr)
		}

		return nil
	}
}

func testAccDataSourceOutscaleOAPIVpcConfig(cidr, tag string) string {
	return fmt.Sprintf(`

resource "outscale_lin" "test" {
  ip_range = "%s"

  tag {
    Name = "%s"
  }
}

data "outscale_lin" "by_id" {
  lin_id = "${outscale_lin.test.id}"
}`, cidr, tag)
}
