package outscale

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVpcs_basic(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(16)
	ipRange := fmt.Sprintf("172.%d.0.0/16", rInt)
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
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
			net_id = ["${outscale_net.test.id}"]
		}
	`, ipRange, tag)
}
