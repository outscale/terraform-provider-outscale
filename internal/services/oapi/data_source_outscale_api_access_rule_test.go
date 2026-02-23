package oapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataOutscaleApiAccessRule_basic(t *testing.T) {
	resourceName := "outscale_api_access_rule.rule_data"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccDataCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleApiAccessRuleConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(t.Context(), resourceName),
				),
			},
		},
	})
}

func testAccDataCheckOutscaleApiAccessRuleDestroy(s *terraform.State) error {
	client := testacc.ConfiguredClient.OSC

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_api_access_rule" {
			continue
		}
		req := osc.ReadApiAccessRulesRequest{
			Filters: &osc.FiltersApiAccessRule{ApiAccessRuleIds: &[]string{rs.Primary.ID}},
		}

		exists := false
		resp, err := client.ReadApiAccessRules(context.Background(), req, options.WithRetryTimeout(120*time.Second))
		if err != nil {
			return fmt.Errorf("api access rule reading (%s)", rs.Primary.ID)
		}

		for _, r := range ptr.From(resp.ApiAccessRules) {
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

func testAccDataOutscaleApiAccessRuleConfig(path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_rule" {
   ca_pem       = file(%q)
   description  = "Ca data test create"
}

resource "outscale_api_access_rule" "rule_data" {
  ca_ids      = [outscale_ca.ca_rule.id]
  ip_ranges   = ["192.4.2.32/16"]
  description = "test api access rule"
}

data "outscale_api_access_rule" "api_access_rule" {
  filter {
    name   = "api_access_rule_ids"
    values = [outscale_api_access_rule.rule_data.id]
  }

  filter {
    name   = "ip_ranges"
    values = ["192.4.2.32/16"]
  }

  filter {
    name   = "descriptions"
    values = ["test api access rule"]
  }
}`, path)
}
