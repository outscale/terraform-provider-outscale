package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleDSLBUListenerDesc_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleLBUListenerDescConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_listener_description.test", "listener_descriptions.0.instance_port", "8000"),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_listener_description.test", "listener_descriptions.1.instance_port", "8080"),
				)},
		},
	})
}

const testAccDSOutscaleLBUListenerDescConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb"
	listeners {
		instance_port = 8000
		instance_protocol = "HTTP"
		load_balancer_port = 80
		protocol = "HTTP"
	}

	listeners {
		instance_port = 8080
		instance_protocol = "HTTP"
		load_balancer_port = 8080
		protocol = "HTTP"
	}

	tag {
		bar = "baz"
	}
}

data "outscale_load_balancer_listener_description" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`
