package outscale

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_WithNic_basic(t *testing.T) {
	subregion := utils.GetRegion()
	resourceName := "outscale_nic.outscale_nic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "private_ips.#", "2"),
				),
			},
			{
				Config: testAccOutscaleENIConfigUpdate(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subregion_name", fmt.Sprintf("%sa", subregion)),
					resource.TestCheckResourceAttr(resourceName, "private_ips.#", "3"),
				),
			},
		},
	})
}

func testAccOutscaleENIConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}

			private_ips	{
				is_primary = false
				private_ip = "10.0.0.46"
			}
		}
	`, subregion)
}

func testAccOutscaleENIConfigUpdate(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-nic-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]

			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}

			private_ips {
				is_primary = false
				private_ip = "10.0.0.46"
			}

			private_ips {
				is_primary = false
				private_ip = "10.0.0.69"
			}
		}
	`, subregion)
}
