package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutscaleApiAccessPolicy_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleApiAccessPolicyConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleApiAccessPolicyConfig() string {
	return fmt.Sprintf(`
              data "outscale_api_access_policy" "api_access_policy" {}
	`)
}
