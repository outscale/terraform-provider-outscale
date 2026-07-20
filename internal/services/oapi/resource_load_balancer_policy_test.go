package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_LBUCookieStickinessPolicy_Basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-lbu")
	policyName1 := acctest.RandomWithPrefix("test-policy")
	policyName2 := acctest.RandomWithPrefix("test-policy")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccCookieStickinessPolicyConfig(name, policyName1, policyName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_load_balancer_policy.app-policy", "cookie_name", "MyAppCookie"),
					resource.TestCheckResourceAttr("outscale_load_balancer_policy.lb-policy", "cookie_expiration_period", "180"),
				),
			},
			{
				Config: testAccCookieStickinessPolicyConfigUpdate(name, policyName1, policyName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_load_balancer_policy.app-policy", "cookie_name", "MyOtherAppCookie"),
					resource.TestCheckResourceAttr("outscale_load_balancer_policy.lb-policy", "cookie_expiration_period", "100"),
				),
			},
		},
	})
}

func TestAccOthers_LBUCookieStickinessPolicy_Migration(t *testing.T) {
	name := acctest.RandomWithPrefix("test-lbu")
	policyName1 := acctest.RandomWithPrefix("test-policy")
	policyName2 := acctest.RandomWithPrefix("test-policy")

	testacc.MigrationTest(t, "1.7.0",
		testAccCookieStickinessPolicyConfig(name, policyName1, policyName2),
		testAccCookieStickinessPolicyConfigUpdate(name, policyName1, policyName2),
	)
}

func testAccCookieStickinessPolicyConfig(name, policyName1, policyName2 string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	subregion_names = [var.subregion]
	listeners {
		backend_port = 8000
		backend_protocol = "HTTP"
		load_balancer_port = 80
		load_balancer_protocol = "HTTP"
	}
}

resource "outscale_load_balancer_policy" "app-policy" {
	policy_type = "app"
	policy_name = "%s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "%s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 180
}
`, name, policyName1, policyName2)
}

// Change the cookie_name to "MyOtherAppCookie".
func testAccCookieStickinessPolicyConfigUpdate(rName, policyName1, policyName2 string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lb" {
	load_balancer_name = "%s"
	subregion_names = [var.subregion]
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }

  timeouts {
    delete = "10m"
  }
}

resource "outscale_load_balancer_policy" "app-policy" {
	policy_type = "app"
	policy_name = "%s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyOtherAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "%s"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 100
}`, rName, policyName1, policyName2)
}
