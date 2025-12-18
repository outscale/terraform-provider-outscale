package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_LBU_basic(t *testing.T) {
	const (
		MIN_LB_NAME_SUFFIX int = 1000
		MAX_LB_NAME_SUFFIX int = 5000
	)

	var conf oscgo.LoadBalancer

	zone := fmt.Sprintf("%sa", utils.GetRegion())
	suffix := utils.RandIntRange(MIN_LB_NAME_SUFFIX, MAX_LB_NAME_SUFFIX)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleLBUConfig(zone, suffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.dataLb", &conf),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.dataTest", "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer.dataTest", "subregion_names.0", zone),
				),
			},
		},
	})
}

func testAccDSOutscaleLBUConfig(zone string, suffix int) string {
	return fmt.Sprintf(`
	resource "outscale_load_balancer" "dataLb" {
		subregion_names    = ["%s"]
		load_balancer_name = "data-terraform-elb-%d"

		listeners {
			backend_port           = 8000
			backend_protocol       = "HTTP"
			load_balancer_port     = 80
			load_balancer_protocol = "HTTP"
		}

		tags {
			key   = "name"
			value = "baz"
		}
	}

	data "outscale_load_balancer" "dataTest" {
		load_balancer_name = outscale_load_balancer.dataLb.id
	}
`, zone, suffix)
}
