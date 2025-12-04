package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNets_DataSource_basic(t *testing.T) {

	ipRange := utils.RandVpcCidr()
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%s", ipRange)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVpcsConfig(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_nets.by_id", "nets.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVpcsConfig(ipRange, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "test" {
			ip_range = "%s"

			tags {
			key = "Name"
			value = "%s"
			}
		}

		data "outscale_nets" "by_id" {
                  filter {
                   name = "net_ids"
                   values = [outscale_net.test.id]
                 }
             }
	`, ipRange, tag)

}
