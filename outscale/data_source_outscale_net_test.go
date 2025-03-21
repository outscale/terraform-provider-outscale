package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_DataSource_basic(t *testing.T) {
	t.Parallel()
	ipRange := utils.RandVpcCidr()
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%s", ipRange)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVpcConfig(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleVpcCheck("data.outscale_net.by_id", ipRange),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVpcCheck(name, ipRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		netRs, ok := s.RootModule().Resources["outscale_net.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_net.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != netRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				netRs.Primary.Attributes["id"],
			)
		}

		if attr["ip_range"] != ipRange {
			return fmt.Errorf("bad cidr_block %s, expected: %s", attr["ip_range"], ipRange)
		}

		return nil
	}
}

func testAccDataSourceOutscaleVpcConfig(ipRange, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "test" {
			ip_range = "%s"

			tags {
				key   = "Name"
				value = "%s"
			}
		}

		data "outscale_net" "by_id" {
			filter {
				name   = "net_ids"
				values = [outscale_net.test.id]
			}
		}
	`, ipRange, tag)
}
