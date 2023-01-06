package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutscaleOAPIApiAccessPolicy_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIApiAccessPolicyConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIApiAccessPolicyConfig() string {
	return fmt.Sprintf(`
              data "outscale_api_access_policy" "api_access_policy" {}
	`)
}
