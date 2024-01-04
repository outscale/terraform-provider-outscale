package outscale

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccNets_DataSource_basic(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())
	ipRange := utils.RandVpcCidr()
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%s", ipRange)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVpcsConfig(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_nets.by_id", "nets.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpcsConfig(ipRange, tag string) string {
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
