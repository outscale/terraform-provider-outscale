package outscale

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIVpc_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(16)
	ipRange := fmt.Sprintf("172.%d.0.0/16", rInt)
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVpcConfig(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIVpcCheck("data.outscale_net.by_id", ipRange, tag),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpcCheck(name, ipRange, tag string) resource.TestCheckFunc {
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

func testAccDataSourceOutscaleOAPIVpcConfig(ipRange, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "test" {
			ip_range = "%s"
		
			#not supported yet
			tags {
				key   = "Name"
				value = "%s"
			}
		}
		
		data "outscale_net" "by_id" {
			#  net_id = "${outscale_net.test.id}"
			filter {
				name   = "net_ids"
				values = ["${outscale_net.test.id}"]
			}
		}
	`, ipRange, tag)
}
