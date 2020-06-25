package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSLBUListenerDesc_basic(t *testing.T) {
	t.Skip()

	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleOAPILBUListenerDescConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_listener_description.test", "listener.backend_port", "8000"),
				)},
		},
	})
}

const testAccDSOutscaleOAPILBUListenerDescConfig = `
	resource "outscale_load_balancer" "bar" {
		subregion_name = ["eu-west-2a"]
		load_balancer_name = "foobar-terraform-elb"

		listeners {
			backend_port           = 8000
			backend_protocol       = "HTTP"
			load_balancer_port     = 80
			load_balancer_protocol = "HTTP"
		}

		tag {
			bar = "baz"
		}
	}

	data "outscale_load_balancer_listener_description" "test" {
		load_balancer_name = "${outscale_load_balancer.bar.id}"
	}
`
