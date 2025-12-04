package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_WithNicsDataSource(t *testing.T) {
	resourceName := "data.outscale_nics.data_nics"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleNicsDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleNicsDataSourceConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-nics-ds"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "sg_dataNic" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.sg_dataNic.security_group_id]
		}

		data "outscale_nics" "data_nics" {
			filter {
				name   = "nic_ids"
				values = [outscale_nic.outscale_nic.id]
			}
			depends_on = [outscale_nic.outscale_nic]
		}
	`, subregion)
}
