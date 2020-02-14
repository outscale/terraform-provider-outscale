package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIDSLinAttr_basic(t *testing.T) {
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

		if attr["dhcp_options_set_id"] == "" {
			return fmt.Errorf("bad dhcp_options_set_id is empty")
		}

		return nil
	}
}

const testAccOutscaleOAPIDSLinAttrConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_net_attributes" "outscale_net_attributes" {
		net_id = "${outscale_net.vpc.id}"
		dhcp_options_set_id = "dopt-ca98300d"
	}

	data "outscale_net_attributes" "test" {
		net_id = "${outscale_net.vpc.id}"
	}
`
