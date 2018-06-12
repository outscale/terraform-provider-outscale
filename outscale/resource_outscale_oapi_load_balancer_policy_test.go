package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestOutscaleOAPILoadBalancerPolicy_basic(t *testing.T) {
	// var out eim.GetLoadBalancerOutput

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckOutscaleLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILoadBalancerPolicyPrefixNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_policy.outscale_load_balancer_policy", "load_balancer_port", "25"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILoadBalancerPolicyPrefixNameConfig = `
resource "outscale_load_balancer" "outscale_load_balancer" {
    load_balancer_name     = "foobar-terraform-elb"
    sub_region_name 	   = ["eu-west-2a"]

    listener {
    	backend_port = 8000
    	backend_protocol = "HTTP"
    	load_balancer_port = 80
    	load_balancer_protocol = "HTTP"
    }
}

resource "outscale_load_balancer_policy" "outscale_load_balancer_policy" {
    load_balancer_name = "${outscale_load_balancer.outscale_load_balancer.load_balancer_name}"
    
    load_balancer_port = "${outscale_load_balancer.outscale_load_balancer.listeners.0.load_balancer_port}"
}
`
