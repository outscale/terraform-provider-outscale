package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleDSLBU_basic(t *testing.T) {
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
				Config: testAccDSOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "availability_zones_member.#", "2"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "availability_zones_member.0", "eu-west-2a"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "availability_zones_member.1", "eu-west-2b"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "listeners_member.0.instance_port", "8000"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "listeners_member.0.instance_protocol", "http"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "listeners_member.0.load_balancer_port", "80"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.test", "listeners_member.0.protocol", "http"),
				)},
		},
	})
}

const testAccDSOutscaleLBUConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b"]
	load_balancer_name               = "foobar-terraform-elb"
  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    // Protocol should be case insensitive
    protocol = "http"
  }

	tag {
		bar = "baz"
	}
}

data "outscale_load_balancer" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`
