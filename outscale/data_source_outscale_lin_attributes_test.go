package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleDSLinAttr_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleDSLinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleVpcAttrCheck("data.outscale_lin_attributes.test"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVpcAttrCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["enable_dns_support"] != "true" {
			return fmt.Errorf("bad enable_dns_support %s, expected: %s", attr["enable_dns_support"], "true")
		}

		return nil
	}
}

const testAccOutscaleDSLinAttrConfig = `
resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_lin_attributes" "outscale_lin_attributes" {
	enable_dns_support = true
	vpc_id = "${outscale_lin.vpc.id}"
	attribute = "enableDnsSupport"
}

data "outscale_lin_attributes" "test" {
	vpc_id = "${outscale_lin.vpc.id}"
	attribute = "enableDnsSupport"
}
`
