package outscale

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPILBUBasic(t *testing.T) {
	t.Parallel()
	var conf oscgo.LoadBalancer

	resourceName := "outscale_load_balancer.bar"
	r := utils.RandIntRange(0, 10)
	zone := fmt.Sprintf("%sa", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILBUConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subregion_names.0", zone),
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secured_cookies", "true"),
				)},
		},
	})
}

func TestAccOutscaleOAPILBUPublicIp(t *testing.T) {
	t.Skip("will be done soon")
	t.Parallel()
	var conf oscgo.LoadBalancer

	resourceName := "outscale_load_balancer.bar"

	r := utils.RandIntRange(10, 20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILBUPublicIpConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
				)},
		},
	})
}

func testAccCheckOutscaleOAPILBUDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		var err error
		var resp oscgo.ReadLoadBalancersResponse
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			filter := &oscgo.FiltersLoadBalancer{
				LoadBalancerNames: &[]string{rs.Primary.ID},
			}

			req := &oscgo.ReadLoadBalancersRequest{
				Filters: filter,
			}

			rp, httpResp, err := conn.LoadBalancerApi.ReadLoadBalancers(
				context.Background()).ReadLoadBalancersRequest(*req).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err == nil {
			if len(*resp.LoadBalancers) != 0 &&
				*(*resp.LoadBalancers)[0].LoadBalancerName ==
					rs.Primary.ID {
				return fmt.Errorf("LBU still exists")
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

func testAccCheckOutscaleOAPILBUExists(n string, res *oscgo.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		var err error
		var resp oscgo.ReadLoadBalancersResponse
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			filter := &oscgo.FiltersLoadBalancer{
				LoadBalancerNames: &[]string{rs.Primary.ID},
			}

			req := &oscgo.ReadLoadBalancersRequest{
				Filters: filter,
			}

			rp, httpResp, err := conn.LoadBalancerApi.ReadLoadBalancers(
				context.Background()).ReadLoadBalancersRequest(*req).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return err
		}

		if len(*resp.LoadBalancers) != 1 ||
			*(*resp.LoadBalancers)[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("LBU not found")
		}

		res = &(*resp.LoadBalancers)[0]

		if res.NetId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("Expected to find source_security_group_id for LBU, but was empty")
			}
		}

		return nil
	}
}

func testAccOutscaleOAPILBUConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
	subregion_names = ["%sa"]
	load_balancer_name               = "foobar-terraform-elb-%d"
	
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

func testAccOutscaleOAPILBUPublicIpConfig(r int) string {
	return fmt.Sprintf(`

	resource "outscale_public_ip" "my_public_ip" {
	}

	resource "outscale_load_balancer" "bar" {
		subregion_names = ["%[1]sa"]
		load_balancer_name = "foobar-terraform-elb-%[2]d"
	  
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
`, os.Getenv("OUTSCALE_REGION"), r)
}
