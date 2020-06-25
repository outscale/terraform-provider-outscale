package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIDSLBUH_basic(t *testing.T) {
	t.Skip()

	var conf oscgo.LoadBalancer
	rs := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: getTestAccDSOutscaleOAPILBUHConfig(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttrSet(
						"data.outscale_load_balancer_health_check.test", "healthy_threshold"),
					resource.TestCheckResourceAttrSet(
						"data.outscale_load_balancer_health_check.test", "check_interval"),
				)},
		},
	})
}
func getTestAccDSOutscaleOAPILBUHConfig(r string) string {
	return fmt.Sprintf(`
		resource "outscale_load_balancer" "bar" {
			subregion_name        = ["eu-west-2a"]
			load_balancer_name = "foobar-terraform-elb-%s"
		
			listener {
				backend_port           = 8000
				backend_protocol       = "HTTP"
				load_balancer_port     = 80
				load_balancer_protocol = "HTTP"
			}
		
			tag {
				bar = "baz"
			}
		}
		
		data "outscale_load_balancer_health_check" "test" {
			load_balancer_name = "${outscale_load_balancer.bar.id}"
		}
	`, r)
}
