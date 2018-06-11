package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
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

	//var conf lbu.LoadBalancerDescription

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
					//testAccCheckOutscaleOAPIHealthCheckExists("outscale_load_balancer_health_check.test", &conf),
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

func testAccCheckOutscaleOAPIHealthCheckExists(n string, res *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).LBU

		var err error
		var describe *lbu.DescribeLoadBalancersOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Throttling") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(describe.LoadBalancerDescriptions) != 1 ||
			*describe.LoadBalancerDescriptions[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("LBU not found")
		}

		*res = *describe.LoadBalancerDescriptions[0]

		if res.VPCId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("Expected to find source_security_group_id for LBU, but was empty")
			}
		}

		return nil
	}
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
