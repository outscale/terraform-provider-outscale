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

func TestAccOutscaleOAPILBUUpdate_Listener(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var conf lbu.LoadBalancerDescription
	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer_listeners.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILBUListenersConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer_listeners.bar", &conf),
					testAccCheckOutscaleOAPILBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_listeners.bar", "listener.0.backend_port", "9000"),
				),
			},
		},
	})
}

func testAccOutscaleOAPILBUListenersConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
  sub_region_name = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-lbu-%d"
	 listener {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
	tag {
		bar = "baz"
	}

}

resource "outscale_load_balancer_listeners" "bar" {
	load_balancer_name               = "${outscale_load_balancer.lb.id}"
  listener {
    backend_port = 9000
    backend_protocol = "HTTP"
    load_balancer_port = 9000
    load_balancer_protocol = "HTTP"
  }
}
`, r)
}
