package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleAppCookieStickinessPolicy_basic(t *testing.T) {
	t.Skip()

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	lbName := fmt.Sprintf("tf-test-lb-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppCookieStickinessPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAppCookieStickinessPolicyConfig(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppCookieStickinessPolicy(
						"outscale_load_balancer.lb",
						"outscale_load_balancer_cookiepolicy.foo",
					),
				),
			},
			resource.TestStep{
				Config: testAccAppCookieStickinessPolicyConfigUpdate(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppCookieStickinessPolicy(
						"outscale_load_balancer.lb",
						"outscale_load_balancer_cookiepolicy.foo",
					),
				),
			},
		},
	})
}

func TestAccOutscaleAppCookieStickinessPolicy_missingLB(t *testing.T) {
	t.Skip()
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	lbName := fmt.Sprintf("tf-test-lb-%s", acctest.RandString(5))

	// check that we can destroy the policy if the LB is missing
	removeLB := func() {
		conn := testAccProvider.Meta().(*OutscaleClient).LBU

		deleteElbOpts := lbu.DeleteLoadBalancerInput{
			LoadBalancerName: aws.String(lbName),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.DeleteLoadBalancer(&deleteElbOpts)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Throttling") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("Error deleting ELB: %s", err)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppCookieStickinessPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAppCookieStickinessPolicyConfig(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppCookieStickinessPolicy(
						"outscale_load_balancer.lb",
						"outscale_load_balancer_cookiepolicy.foo",
					),
				),
			},
			resource.TestStep{
				PreConfig: removeLB,
				Config:    testAccAppCookieStickinessPolicyConfigDestroy(lbName),
			},
		},
	})
}

func testAccCheckAppCookieStickinessPolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer_cookiepolicy" {
			continue
		}
	}
	return nil
}

func testAccCheckAppCookieStickinessPolicy(elbResource string, policyResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[elbResource]
		if !ok {
			return fmt.Errorf("Not found: %s", elbResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if !ok {
			return fmt.Errorf("Not found: %s", policyResource)
		}

		return nil
	}
}

func testAccAppCookieStickinessPolicyConfig(rName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	 availability_zones = ["eu-west-2a"]
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}

resource "outscale_load_balancer_cookiepolicy" "foo" {
	policy_name = "foo-policy"
	load_balancer_name = "${outscale_load_balancer.lb.id}"
	cookie_name = "MyAppCookie"
}`, rName)
}

// Change the cookie_name to "MyOtherAppCookie".
func testAccAppCookieStickinessPolicyConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	 availability_zones = ["eu-west-2a"]
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}

resource "outscale_load_balancer_cookiepolicy" "foo" {
	policy_name = "foo-policy"
	load_balancer_name = "${outscale_load_balancer.lb.id}"
	cookie_name = "MyOtherAppCookie"
}`, rName)
}

// attempt to destroy the policy, but we'll delete the LB in the PreConfig
func testAccAppCookieStickinessPolicyConfigDestroy(rName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	 availability_zones = ["eu-west-2a"]
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}`, rName)
}
