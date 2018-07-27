package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIDSLinAPIAccess_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIDSLinAPIAccessConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_lin_api_access.test", "service_name", "com.outscale.eu-west-2.osu"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIDSLinAPIAccessConfig = `
resource "outscale_lin" "foo" {
	ip_ranges = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
	lin_id = "${outscale_lin.foo.id}"
}

resource "outscale_lin_api_access" "link" {
	lin_id = "${outscale_lin.foo.id}"
	route_table_id = [
		"${outscale_route_table.foo.id}"
	]
	service_name = "com.outscale.eu-west-2.osu"
}

data "outscale_lin_api_access" "test" {
	lin_api_access_id = "${outscale_lin_api_access.link.id}"
}
`
