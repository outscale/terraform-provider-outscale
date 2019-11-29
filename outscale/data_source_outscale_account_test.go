package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestOutscaleOAPIAccount(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOAPIAccouuntDSConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_account.account", "account_pid"),
					resource.TestCheckResourceAttrSet("data.outscale_account.account", "company_name"),
					resource.TestCheckResourceAttrSet("data.outscale_account.account", "email"),
					resource.TestCheckResourceAttrSet("data.outscale_account.account", "country"),
				),
			},
		},
	})
}

func testOAPIAccouuntDSConfig() string {
	return fmt.Sprintf(`
		data "outscale_account" "account" {}
	`)
}
