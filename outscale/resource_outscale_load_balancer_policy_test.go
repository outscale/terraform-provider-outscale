package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestOutscaleLoadBalancerPolicy_basic(t *testing.T) {
	// var out eim.GetLoadBalancerOutput

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckOutscaleLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleLoadBalancerPolicyPrefixNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_policy.outscale_load_balancer_policy", "load_balancer_port", "25"),
				),
			},
		},
	})
}

const testAccOutscaleLoadBalancerPolicyPrefixNameConfig = `
resource "outscale_load_balancer" "outscale_load_balancer" {
    load_balancer_name     = "foobar-terraform-elb"
    availability_zones     = ["eu-west-2a"]

    listeners {
        instance_port      = 1024
        instance_protocol  = "HTTP"
        load_balancer_port = 25
        protocol           = "HTTP"
    }
}

resource "outscale_load_balancer_policy" "outscale_load_balancer_policy" {
    load_balancer_name = "${outscale_load_balancer.outscale_load_balancer.load_balancer_name}"
    
    load_balancer_port = "${outscale_load_balancer.outscale_load_balancer.listeners.0.load_balancer_port}"
}
`
