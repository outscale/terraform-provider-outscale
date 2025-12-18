package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOutscaleApiAccessPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
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
