package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_Attributes_basic(t *testing.T) {
	resourceName := "outscale_net_attributes.outscale_net_attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_set_id"),
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

const testAccOutscaleLinAttrConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-attributes-rs-vpc"
		}
	}

	resource "outscale_net" "vpc2" {
		ip_range = "10.0.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-attributes-rs-vpc2"
		}
	}

	resource "outscale_net_attributes" "outscale_net_attributes" {
		net_id              = outscale_net.vpc.id
		dhcp_options_set_id = outscale_net.vpc2.dhcp_options_set_id
	}
`
