package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscalePublicIP(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscalePublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscalePublicIPCheck("data.outscale_public_ip.by_allocation_id"),
					testAccDataSourceOutscalePublicIPCheck("data.outscale_public_ip.by_public_ip"),
				),
			},
		},
	})
}

func testAccDataSourceOutscalePublicIPCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		fmt.Printf("\n[DEBUG] TEST RS %s \n", s.RootModule().Resources)

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		eipRs, ok := s.RootModule().Resources["outscale_public_ip.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_public_ip.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["allocation_id"] != eipRs.Primary.Attributes["allocation_id"] {
			return fmt.Errorf(
				"allocation_id is %s; want %s",
				attr["allocation_id"],
				eipRs.Primary.Attributes["allocation_id"],
			)
		}

		if attr["public_ip"] != eipRs.Primary.Attributes["public_ip"] {
			return fmt.Errorf(
				"public_ip is %s; want %s",
				attr["public_ip"],
				eipRs.Primary.Attributes["public_ip"],
			)
		}

		return nil
	}
}

const testAccDataSourceOutscalePublicIPConfig = `
resource "outscale_public_ip" "test" {}

data "outscale_public_ip" "by_allocation_id" {
  allocation_id = "${outscale_public_ip.test.allocation_id}"
}
data "outscale_public_ip" "by_public_ip" {
  public_ip = "${outscale_public_ip.test.public_ip}"
}
`
