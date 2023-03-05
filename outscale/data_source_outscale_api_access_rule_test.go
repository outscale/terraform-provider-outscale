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

func TestAccOthers_DataOutscaleOAPIApiAccessRule_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_api_access_rule.rule_data"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDataCheckOutscaleApiAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleOAPIApiAccessRuleConfig(utils.TestCaPem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleApiAccessRuleExists(resourceName),
				),
			},
		},
	})
}

func testAccDataCheckOutscaleApiAccessRuleDestroy(s *terraform.State) error {
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

func testAccDataOutscaleOAPIApiAccessRuleConfig(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_rule" { 
   ca_pem       = %[1]q
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
}
	`, ca_pem)
}
