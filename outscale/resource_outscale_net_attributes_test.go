package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNet_Attributes_basic(t *testing.T) {
	resourceName := "outscale_net_attributes.outscale_net_attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_set_id"),
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

func TestAccNet_Attributes_withoutDHCPID(t *testing.T) {
	resourceName := "outscale_net_attributes.outscale_net_attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILinAttrConfigwithoutDHCPID("outscale_net.vpc.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_set_id"),
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
			{
				Config: testAccOutscaleOAPILinAttrConfigwithoutDHCPID("outscale_net.vpc2.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_set_id"),
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinAttrConfig = `
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

func testAccOutscaleOAPILinAttrConfigwithoutDHCPID(vpc string) string {
	return fmt.Sprintf(`
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
			net_id = %s

			depends_on = ["outscale_net.vpc", "outscale_net.vpc2"]
		}
	`, vpc)
}
