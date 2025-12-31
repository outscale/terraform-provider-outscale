package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_WithNicsDataSource(t *testing.T) {
	resourceName := "data.outscale_nics.data_nics"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleNicsDataSourceConfig(utils.GetRegion(), sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleNicsDataSourceConfig(subregion, sgName string) string {
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
			security_group_name = "%[2]s"
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
	`, subregion, sgName)
}
