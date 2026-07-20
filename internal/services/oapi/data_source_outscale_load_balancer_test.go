package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_LBU_basic(t *testing.T) {
	const (
		MIN_LB_NAME_SUFFIX int = 1000
		MAX_LB_NAME_SUFFIX int = 5000
	)

	var conf osc.LoadBalancer

	zone := fmt.Sprintf("%sa", utils.GetRegion())
	suffix := utils.RandIntRange(MIN_LB_NAME_SUFFIX, MAX_LB_NAME_SUFFIX)
	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testacc.SDKProviders,
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

// Legacy helper, will be removed once the datasource is migrated to the Plugin Framework
func testAccCheckOutscaleLBUExists(n string, res *osc.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no lbu id is set")
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

		var err error

		req := osc.ReadLoadBalancersRequest{
			Filters: &osc.FiltersLoadBalancer{
				LoadBalancerNames: &[]string{rs.Primary.ID},
			},
		}

		resp, err := client.ReadLoadBalancers(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			return err
		}

		if len(*resp.LoadBalancers) != 1 ||
			(*resp.LoadBalancers)[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("lbu not found")
		}

		res = &(*resp.LoadBalancers)[0]

		if res.NetId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("expected to find source_security_group_id for lbu, but was empty")
			}
		}

		return nil
	}
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
