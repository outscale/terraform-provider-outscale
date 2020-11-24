package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSLBUAttr_basic(t *testing.T) {
	t.Skip()

	r := acctest.RandIntRange(0, 10)

	var conf oscgo.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleOAPILBUAttrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_attributes.test", "load_balancer_attributes.0.access_log.is_enabled", "false"),
				)},
		},
	})
}

func testAccDSOutscaleOAPILBUAttrConfig(r int) string {
	return fmt.Sprintf(`
		resource "outscale_load_balancer" "bar" {
			subregion_name        = ["eu-west-2a"]
			load_balancer_name = "foobar-terraform-elb-%d"
		
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
		
		resource "outscale_load_balancer_attributes" "bar2" {
			is_enabled         = "false"
			osu_bucket_name    = "donustestbucket"
			load_balancer_name = "${outscale_load_balancer.bar.id}"
		}
		
		data "outscale_load_balancer_attributes" "test" {
			load_balancer_name = "${outscale_load_balancer.bar.id}"
		}
	`, r)
}
