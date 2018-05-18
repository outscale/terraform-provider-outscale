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
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleLBsUConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_load_balancers.test", "load_balancer_descriptions_member.#", "1"),
				)},
		},
	})
}

const testAccDSOutscaleLBsUConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb"
  listeners_member {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    // Protocol should be case insensitive
    protocol = "HTTP"
  }

	tag {
		bar = "baz"
	}
}

data "outscale_load_balancers" "test" {
	load_balancer_name = ["${outscale_load_balancer.bar.id}"]
}
`
