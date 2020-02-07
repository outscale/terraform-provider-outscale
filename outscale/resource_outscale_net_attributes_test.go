package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPILinAttr_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("outscale_net_attributes.outscale_net_attributes", "dhcp_options_set_id"),
					resource.TestCheckResourceAttrSet("outscale_net_attributes.outscale_net_attributes", "net_id"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinAttrConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_net" "vpc2" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_net_attributes" "outscale_net_attributes" {
		net_id = "${outscale_net.vpc.id}"
		dhcp_options_set_id = "${outscale_net.vpc2.dhcp_options_set_id}"
	}
`
