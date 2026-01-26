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

func TestAccOthers_DataOutscaleApiAccessRule_basic(t *testing.T) {
	resourceName := "outscale_api_access_rule.rule_data"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccDataCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleApiAccessRuleConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(resourceName),
				),
			},
		},
	})
}

func testAccDataCheckOutscaleApiAccessRuleDestroy(s *terraform.State) error {
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
