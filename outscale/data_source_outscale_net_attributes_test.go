package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_AttributesDataSource_basic(t *testing.T) {
	resourceName := "data.outscale_net_attributes.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleDSLinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

const testAccOutscaleDSLinAttrConfig = `
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
