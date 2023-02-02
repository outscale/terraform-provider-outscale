package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ApiAccessRule_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_api_access_rule.rule"
	dataSourcesName := "data.outscale_api_access_rules.filters_rules"
	dataSourcesAllName := "data.outscale_api_access_rules.all_rules"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_ApiAccessRule_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "ip_ranges.#"),

					resource.TestCheckResourceAttrSet(dataSourcesName, "api_access_rules.#"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "filter.#"),

					resource.TestCheckResourceAttrSet(dataSourcesAllName, "api_access_rules.#"),
				),
			},
		},
	})
}

func testAcc_ApiAccessRule_DataSource_Config() string {
	return fmt.Sprintf(`
resource "outscale_api_access_rule" "rule_data" {
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

data "outscale_api_access_rule" "rule" {
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

`)
}
