package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccDataSourceOutscaleApiAccessPolicy_basic(t *testing.T) {
	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleApiAccessPolicyConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleApiAccessPolicyConfig() string {
	return `
              data "outscale_api_access_policy" "api_access_policy" {}
	`
}
