package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleLBUUpdate_Listener(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var conf lbu.LoadBalancerDescription
	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer_listeners.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUListenersConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer_listeners.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_listeners.bar", "listeners.0.instance_port", "9000"),
				),
			},
		},
	})
}

func testAccOutscaleLBUListenersConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-lbu-%d"
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

resource "outscale_load_balancer_listeners" "bar" {
	load_balancer_name               = "${outscale_load_balancer.lb.id}"
  listeners {
    instance_port = 9000
    instance_protocol = "HTTP"
    load_balancer_port = 9000
    protocol = "HTTP"
  }
}
`, r)
}
