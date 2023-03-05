package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOthers_AccessRule_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_api_access_rule.rule_test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIApiAccessRuleConfig(utils.TestCaPem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckOutscaleApiAccessRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set")
		}
		req := oscgo.ReadApiAccessRulesRequest{
			Filters: &oscgo.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadApiAccessRulesResponse
		var err error
		exists := false
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetApiAccessRules()) == 0 {
			return fmt.Errorf("Api Access Rule not found (%s)", rs.Primary.ID)
		}

		for _, rule := range resp.GetApiAccessRules() {
			if rule.GetApiAccessRuleId() == rs.Primary.ID {
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("Api Access Rule not found (%s)", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckOutscaleApiAccessRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_api_access_rule" {
			continue
		}
		req := oscgo.ReadApiAccessRulesRequest{
			Filters: &oscgo.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadApiAccessRulesResponse
		var err error
		exists := false
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("Api Access Rule reading (%s)", rs.Primary.ID)
		}

		for _, r := range resp.GetApiAccessRules() {
			if r.GetApiAccessRuleId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("Api Access Rule still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleOAPIApiAccessRuleConfig(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_rule" { 
   ca_pem       = %[1]q
   description  = "Ca testacc create"
}

resource "outscale_api_access_rule" "rule_test" {
  ca_ids      = [outscale_ca.ca_rule.id]
  ip_ranges   = ["192.0.2.0/16"]
  description = "testing api access rule"
}
	`, ca_pem)
}
