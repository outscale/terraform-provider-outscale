package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNets_DataSource_basic(t *testing.T) {
	ipRange := oapihelpers.RandVpcCidr()
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%s", ipRange)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
