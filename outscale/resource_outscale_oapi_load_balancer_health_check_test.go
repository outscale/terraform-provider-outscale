package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIHealthCheck_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer_health_check.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleHealthCheckConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_health_check.test", "health_check.healthy_threshold", "2"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_health_check.test", "health_check.unhealthy_threshold", "4"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_health_check.test", "health_check.checked_vm", "HTTP:8000/index.html"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_health_check.test", "health_check.interval", "5"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_health_check.test", "health_check.timeout", "5"),
				)},
		},
	})
}

func testAccOutscaleOAPIHealthCheckConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  sub_region_name = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-%d"
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

resource "outscale_load_balancer_health_check" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
	health_check {
		healthy_threshold = 2
		unhealthy_threshold = 4
		check_interval = 5
		timeout = 5
		checked_vm = "HTTP:8000/index.html"
	}
}
`, r)
}
