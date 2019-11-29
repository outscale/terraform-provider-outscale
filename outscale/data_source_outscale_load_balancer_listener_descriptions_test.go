package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleOAPIDSLBUListenerDescs_basic(t *testing.T) {
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
				Config: testAccDSOutscaleOAPILBUListenerDescsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_listener_descriptions.test", "listener_descriptions.0.listener.0.backend_port", "8000"),
				)},
		},
	})
}

const testAccDSOutscaleOAPILBUListenerDescsConfig = `
	resource "outscale_load_balancer" "bar" {
		sub_region         = ["eu-west-2a"]
		load_balancer_name = "foobar-terraform-elb"

		listener {
			backend_port           = 8000
			backend_protocol       = "HTTP"
			load_balancer_port     = 80
			load_balancer_protocol = "HTTP"
		}

		tag {
			bar = "baz"
		}
	}

	data "outscale_load_balancer_listener_descriptions" "test" {
		load_balancer_names = ["${outscale_load_balancer.bar.id}"]
	}
`
