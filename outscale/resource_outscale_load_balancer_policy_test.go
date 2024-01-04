package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_CookieStickinessPolicy_basic(t *testing.T) {
	t.Parallel()
	lbName := fmt.Sprintf("tf-test-lb-%s", acctest.RandString(10))
	zone := fmt.Sprintf("%sa", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppCookieStickinessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCookieStickinessPolicyConfig(lbName, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppCookieStickinessPolicy(
						"outscale_load_balancer.lb",
						"outscale_load_balancer_cookiepolicy.foo",
					),
				),
			},
			{
				Config: testAccCookieStickinessPolicyConfigUpdate(lbName, zone),
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

func testAccCookieStickinessPolicyConfig(rName string, zone string) string {
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
	policy_name = "foo-policy"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "lb-policy"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 180
}
`, rName, zone)
}

// Change the cookie_name to "MyOtherAppCookie".
func testAccCookieStickinessPolicyConfigUpdate(rName string, zone string) string {
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
	policy_name = "foo-policy"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_name = "MyOtherAppCookie"
}

resource "outscale_load_balancer_policy" "lb-policy" {
	policy_type = "load_balancer"
	policy_name = "lb-policy"
	load_balancer_name = outscale_load_balancer.lb.id
	cookie_expiration_period = 100
}`, rName, zone)
}
