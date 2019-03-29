package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSLinAttr_basic(t *testing.T) {
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
				Config: testAccOutscaleOAPIDSLinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVpcAttrCheck("data.outscale_net_attributes.test"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpcAttrCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["dns_support_enabled"] != "true" {
			return fmt.Errorf("bad dns_support_enabled %s, expected: %s", attr["dns_support_enabled"], "true")
		}

		return nil
	}
}

const testAccOutscaleOAPIDSLinAttrConfig = `
resource "outscale_net" "vpc" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_net_attributes" "outscale_net_attributes" {
	dns_support_enabled = true
	net_id = "${outscale_net.vpc.id}"
	attribute = "enableDnsSupport"
}

data "outscale_net_attributes" "test" {
	net_id = "${outscale_net.vpc.id}"
	attribute = "enableDnsSupport"
}
`
