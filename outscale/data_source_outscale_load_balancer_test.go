package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSLBU_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.LoadBalancer

	zone := fmt.Sprintf("%sa", utils.GetRegion())
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleOAPILBUConfig(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "subregion_names.0", zone),
				)},
		},
	})
}

func testAccDSOutscaleOAPILBUConfig(zone string) string {
	return fmt.Sprintf(`
	resource "outscale_load_balancer" "bar" {
		subregion_names    = ["%s"]
		load_balancer_name = "foobar-terraform-elb"

		listeners {
			backend_port           = 8000
			backend_protocol       = "HTTP"
			load_balancer_port = 80
			load_balancer_protocol = "HTTP"
		}

		tags {
			key = "name"
			value = "baz"
		}
	}

	data "outscale_load_balancer" "test" {
		load_balancer_name = outscale_load_balancer.bar.id
	}
`, zone)
}
