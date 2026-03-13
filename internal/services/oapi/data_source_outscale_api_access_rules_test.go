package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccDataOutscaleApiAccessRules_basic(t *testing.T) {
	resourceName := "outscale_api_access_rule.rule_data"
	ca_path := testAccCertPath
	resource.ParallelTest(t, resource.TestCase{
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccDataCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleApiAccessRulesConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(t.Context(), resourceName),
				),
			},
		},
	})
}

func testAccDataOutscaleApiAccessRulesConfig(cert_path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_rule" {
   ca_pem       = file("%q")
   description  = "Ca testacc create"
}

resource "outscale_api_access_rule" "rule_data" {
  ca_ids      = [outscale_ca.ca_rule.id]
  ip_ranges   = ["192.4.2.32/16"]
  description = "test api access rule"
}

data "outscale_api_access_rules" "filters_rules" {
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
}

data "outscale_api_access_rules" "all_rules" {}
	`, cert_path)
}
