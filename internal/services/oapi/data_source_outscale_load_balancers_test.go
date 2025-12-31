package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_LBUs_basic(t *testing.T) {
	region := fmt.Sprintf("%sa", utils.GetRegion())
	numLbu := utils.RandIntRange(0, 50)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleLBsUConfig(region, numLbu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_load_balancers.test", "load_balancer.#", "1"),
				),
			},
		},
	})
}

func testAccDSOutscaleLBsUConfig(region string, numLbu int) string {
	return fmt.Sprintf(`
	resource "outscale_load_balancer" "bars" {
		subregion_names = ["%s"]
		load_balancer_name        = "foobar-terraform-elb%d"

		listeners {
			backend_port      = 8000
			backend_protocol  = "HTTP"
			load_balancer_port = 80

			// Protocol should be case insensitive
			load_balancer_protocol = "HTTP"
		}

		tags {
			key = "name"
			value = "baz"
		}
	}

	data "outscale_load_balancers" "test" {
		load_balancer_name = [outscale_load_balancer.bars.id]
	}
`, region, numLbu)
}
