package outscale

import (
	"fmt"
	oscgo "github.com/outscale/osc-sdk-go/osc"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPILBUUpdate_Listener(t *testing.T) {
	var conf oscgo.LoadBalancer
	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
  subregion_names = ["eu-west-2a"]
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
