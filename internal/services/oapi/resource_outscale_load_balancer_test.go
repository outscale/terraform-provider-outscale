package oapi_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_LBUBasic(t *testing.T) {
	var conf osc.LoadBalancer

	lbResourceName := "outscale_load_balancer.barRes"
	r := acctest.RandString(5)
	zone := fmt.Sprintf("%sa", utils.GetRegion())

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: lbResourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists(lbResourceName, &conf),
					resource.TestCheckResourceAttr(lbResourceName, "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(lbResourceName, "subregion_names.0", zone),
					resource.TestCheckResourceAttr(lbResourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttr(lbResourceName, "secured_cookies", "true"),
				),
			},
		},
	})
}

func TestAccOthers_LBUPublicIp(t *testing.T) {
	t.Skip("Conflict UnlinkPublicIp: will be done soon")

	var conf osc.LoadBalancer

	resourceName := "outscale_load_balancer.barIp"
	r := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: resourceName,
		Providers:     testacc.SDKProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUPublicIpConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
				),
			},
		},
	})
}

func testAccCheckOutscaleLBUDestroy(s *terraform.State) error {
	client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		var err error
		req := osc.ReadLoadBalancersRequest{
			Filters: &osc.FiltersLoadBalancer{
				LoadBalancerNames: &[]string{rs.Primary.ID},
			},
		}
		resp, err := client.ReadLoadBalancers(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))

		if err == nil {
			if len(*resp.LoadBalancers) != 0 &&
				(*resp.LoadBalancers)[0].LoadBalancerName ==
					rs.Primary.ID {
				return fmt.Errorf("lbu still exists")
			}
		}

		if strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound") {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

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

func testAccOutscaleLBUConfig(r string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "barRes" {
	subregion_names = ["%sa"]
	load_balancer_name               = "foobar-terraform-elb-%s"

	secured_cookies                  = true

	listeners {
		backend_port = 8000
		backend_protocol = "HTTP"
		load_balancer_port = 80
		load_balancer_protocol = "HTTP"
	}

	tags {
		key = "name"
		value = "baz"
	}

}
`, utils.GetRegion(), r)
}

func testAccOutscaleLBUPublicIpConfig(r string) string {
	return fmt.Sprintf(`

	resource "outscale_public_ip" "my_public_ip" {
	}

	resource "outscale_load_balancer" "barIp" {
		subregion_names = ["%[1]sa"]
		load_balancer_name = "foobar-terraform-elb-%s"

		listeners {
		  backend_port           = 80
		  backend_protocol       = "HTTP"
		  load_balancer_protocol = "HTTP"
		  load_balancer_port     = 80
		}

		public_ip = outscale_public_ip.my_public_ip.public_ip

		tags {
		  key = "name"
		  value = "terraform-internet-facing-lb-with-eip"
		}
	  }
`, utils.GetRegion(), r)
}
