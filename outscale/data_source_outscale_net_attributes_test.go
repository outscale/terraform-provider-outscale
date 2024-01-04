package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_AttributesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
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
		tags {
			key = "Name"
			value = "testacc-net-attributes-ds-vpc"
		}
	}

	resource "outscale_net" "vpc2" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-net-attributes-ds-vpc2"
		}
	}

	resource "outscale_net_attributes" "outscale_net_attributes" {
		net_id = outscale_net.vpc.id
		dhcp_options_set_id = outscale_net.vpc2.dhcp_options_set_id
	}

	data "outscale_net_attributes" "test" {
		net_id = outscale_net.vpc.id
	}
`
