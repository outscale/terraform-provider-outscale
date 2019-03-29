package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleDSLBSU_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.outscale_load_balancer",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleLBsUConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_load_balancers.outscale_load_balancers", "load_balancer_descriptions.#", "2"),
				)},
		},
	})
}

const testAccDSOutscaleLBsUConfig = `
resource "outscale_load_balancer" "outscale_load_balancer" {
  count = 1

  load_balancer_name = "foobar-terraform-elb"

  availability_zones = ["eu-west-2a"]

  listeners {
    instance_port = 1024

    instance_protocol = "HTTP"

    load_balancer_port = 25

    protocol = "HTTP"
  }
}

resource "outscale_load_balancer" "outscale_load_balancer2" {
  count = 1

  load_balancer_name = "foobar-terraform-elb2"

  availability_zones = ["eu-west-2a"]

  listeners {
    instance_port = 1024

    instance_protocol = "HTTP"

    load_balancer_port = 25

    protocol = "HTTP"
  }
}

data "outscale_load_balancers" "outscale_load_balancers" {
  load_balancer_name = ["${outscale_load_balancer.outscale_load_balancer.load_balancer_name}", "${outscale_load_balancer.outscale_load_balancer2.load_balancer_name}"]
}
`
