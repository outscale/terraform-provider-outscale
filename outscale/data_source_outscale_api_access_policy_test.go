package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ApiAccessPolicy_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_api_access_policy.api_access_policy"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_ApiAccessPolicy_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "max_access_key_expiration_seconds"),
				),
			},
		},
	})
}

func testAcc_ApiAccessPolicy_DataSource_Config() string {
	return fmt.Sprintf(`
              data "outscale_api_access_policy" "api_access_policy" {}
	`)
}
