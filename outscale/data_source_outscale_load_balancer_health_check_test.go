package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleDSLBUH_basic(t *testing.T) {
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
				Config: testAccDSOutscaleLBUHConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttrSet(
						"data.outscale_load_balancer_health_check.test", "healthy_threshold"),
					resource.TestCheckResourceAttrSet(
						"data.outscale_load_balancer_health_check.test", "interval"),
				)},
		},
	})
}

const testAccDSOutscaleLBUHConfig = `
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

data "outscale_load_balancer_health_check" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`
