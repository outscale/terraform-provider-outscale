package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Account_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_account.account"
	dataSourcesName := "data.outscale_accounts.accounts"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Account_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "accounts.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "account_id"),
				),
			},
		},
	})
}

func testAcc_Account_DataSource_Config() string {
	return fmt.Sprintf(`
        data "outscale_account" "account" { }

		data "outscale_accounts" "accounts" { }
	`)
}
