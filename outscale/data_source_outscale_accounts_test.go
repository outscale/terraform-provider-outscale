package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_DataSourceAccounts_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
