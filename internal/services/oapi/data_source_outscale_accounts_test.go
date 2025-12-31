package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceAccounts_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig(),
			},
		},
	})
}

func testAccDataSourceAccountsConfig() string {
	return fmt.Sprintf(`
              data "outscale_accounts" "accounts" { }
	`)
}
