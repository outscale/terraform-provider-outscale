package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleDSLBUListenerDescs_basic(t *testing.T) {
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
				Config: testAccDSOutscaleLBUListenerDescsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckState("data.outscale_load_balancer_listener_descriptions.test"),
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_listener_descriptions.test", "load_balancers.0.listener_descriptions.0.instance_port", "8000"),
				)},
		},
	})
}

const testAccDSOutscaleLBUListenerDescsConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb"
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }

  tag {
	bar = "baz"
  }
}

resource "outscale_load_balancer" "foo" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-foo"
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
		foo = "foo"
	}
}

data "outscale_load_balancer_listener_descriptions" "test" {
	load_balancer_names = ["${outscale_load_balancer.bar.id}", "${outscale_load_balancer.foo.id}"]
}
`
