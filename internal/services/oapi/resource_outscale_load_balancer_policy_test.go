package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_CookieStickinessPolicy_basic(t *testing.T) {
	lbName := fmt.Sprintf("tf-test-lb-%s", acctest.RandString(10))
	policyName1 := acctest.RandomWithPrefix("test-policy")
	policyName2 := acctest.RandomWithPrefix("test-policy")
	zone := fmt.Sprintf("%sa", utils.GetRegion())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testacc.PreCheck(t) },
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccCheckAppCookieStickinessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCookieStickinessPolicyConfig(lbName, zone, policyName1, policyName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppCookieStickinessPolicy(
						"outscale_load_balancer.lb",
						"outscale_load_balancer_cookiepolicy.foo",
					),
				),
			},
			{
				Config: testAccCookieStickinessPolicyConfigUpdate(lbName, zone, policyName1, policyName2),
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

func testAccCookieStickinessPolicyConfig(rName, zone, policyName1, policyName2 string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	subregion_names = ["%s"]
	listeners {
		backend_port = 8000
		backend_protocol = "HTTP"
		load_balancer_port = 80
		load_balancer_protocol = "HTTP"
	}
}

resource "outscale_load_balancer_policy" "app-policy" {
	policy_type = "app"
	policy_name = "%[3]s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "%[4]s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 180
}
`, rName, zone, policyName1, policyName2)
}

// Change the cookie_name to "MyOtherAppCookie".
func testAccCookieStickinessPolicyConfigUpdate(rName, zone, policyName1, policyName2 string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	subregion_names = ["%s"]
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_policy" "app-policy" {
	policy_type = "app"
	policy_name = "%[3]s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyOtherAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "%[4]s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 100
}`, rName, zone, policyName1, policyName2)
}
