package oapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_AccessRule_basic(t *testing.T) {
	resourceName := "outscale_api_access_rule.rule_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleApiAccessRuleConfig(ca_path),
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
			return fmt.Errorf("not found: %s", n)
		}

		conn := testacc.ConfiguredClient.OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}
		req := oscgo.ReadApiAccessRulesRequest{
			Filters: &oscgo.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadApiAccessRulesResponse
		var err error
		exists := false
		err = retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetApiAccessRules()) == 0 {
			return fmt.Errorf("api access rule not found (%s)", rs.Primary.ID)
		}

		for _, rule := range resp.GetApiAccessRules() {
			if rule.GetApiAccessRuleId() == rs.Primary.ID {
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("api access rule not found (%s)", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckOutscaleApiAccessRuleDestroy(s *terraform.State) error {
	conn := testacc.ConfiguredClient.OSCAPI

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
		err = retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.ApiAccessRuleApi.ReadApiAccessRules(context.Background()).ReadApiAccessRulesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("api access rule reading (%s)", rs.Primary.ID)
		}

		for _, r := range resp.GetApiAccessRules() {
			if r.GetApiAccessRuleId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("api access rule still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleApiAccessRuleConfig(path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_rule" {
   ca_pem       = file(%q)
   description  = "Ca testacc create"
}

resource "outscale_api_access_rule" "rule_test" {
  ca_ids      = [outscale_ca.ca_rule.id]
  ip_ranges   = ["192.0.2.0/16"]
  description = "testing api access rule"
}`, path)
}
