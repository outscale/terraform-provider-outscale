package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_AccessRule_Basic(t *testing.T) {
	resourceName := "outscale_api_access_rule.rule_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleApiAccessRuleConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(t.Context(), resourceName),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_AccessRule_Migration(t *testing.T) {
	ca_path := testAccCertPath

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.3.1",
			testAccOutscaleApiAccessRuleConfig(ca_path),
		),
	})
}

func testAccCheckOutscaleApiAccessRuleExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		client := testacc.ConfiguredClient.OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}
		req := osc.ReadApiAccessRulesRequest{
			Filters: &osc.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}
		exists := false

		resp, err := client.ReadApiAccessRules(ctx, req, options.WithRetryTimeout(DefaultTimeout))
		if err != nil || resp.ApiAccessRules == nil || len(*resp.ApiAccessRules) == 0 {
			return fmt.Errorf("api access rule not found (%s)", rs.Primary.ID)
		}

		for _, rule := range *resp.ApiAccessRules {
			if *rule.ApiAccessRuleId == rs.Primary.ID {
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
	client := testacc.ConfiguredClient.OSC

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_api_access_rule" {
			continue
		}
		req := osc.ReadApiAccessRulesRequest{
			Filters: &osc.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}

		exists := false
		resp, err := client.ReadApiAccessRules(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))
		if err != nil || resp.ApiAccessRules == nil || len(*resp.ApiAccessRules) == 0 {
			return fmt.Errorf("api access rule reading (%s)", rs.Primary.ID)
		}

		for _, r := range *resp.ApiAccessRules {
			if *r.ApiAccessRuleId == rs.Primary.ID {
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
