package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceAccounts_basic(t *testing.T) {
	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig(),
			},
		},
	})
}

func testAccDataSourceAccountsConfig() string {
	return `
        data "outscale_accounts" "accounts" { }
	`
}
