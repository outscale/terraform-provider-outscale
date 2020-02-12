package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSLinAPIAccess_basic(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIDSLinAPIAccessConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_net_api_access.test", "service_name", "com.outscale.eu-west-2.osu"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIDSLinAPIAccessConfig = `
	resource "outscale_net "foo" {
		ip_ranges = "10.1.0.0/16"
	}

	resource "outscale_route_table" "foo" {
		net_id = "${outscale_net.foo.id}"
	}

	resource "outscale_net_api_access" "link" {
		net_id = "${outscale_net.foo.id}"
		route_table_id = [
			"${outscale_route_table.foo.id}"
		]
		service_name = "com.outscale.eu-west-2.osu"
	}

	data "outscale_net_api_access" "test" {
		net_api_access_id = "${outscale_net_api_access.link.id}"
	}
`
