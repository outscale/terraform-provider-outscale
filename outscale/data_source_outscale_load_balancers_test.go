package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSLBSU_basic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")
	zone := fmt.Sprintf("%sa", region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleOAPILBsUConfig(zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_load_balancers.test", "load_balancer.#", "1"),
				)},
		},
	})
}

func testAccDSOutscaleOAPILBsUConfig(zone string) string {
	return fmt.Sprintf(`
	resource "outscale_load_balancer" "bar" {
		subregion_names = ["%s"]
		load_balancer_name        = "foobar-terraform-elb"

		listeners {
			backend_port      = 8000
			backend_protocol  = "HTTP"
			load_balancer_port = 80

			// Protocol should be case insensitive
			load_balancer_protocol = "HTTP"
		}

		tags {
			key = "name"
			value = "baz"
		}
	}

	data "outscale_load_balancers" "test" {
		load_balancer_name = ["${outscale_load_balancer.bar.id}"]
	}
`, zone)
}
