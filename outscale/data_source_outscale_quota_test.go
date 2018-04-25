package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleQuota(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleQuotaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleQuotaCheck("data.outscale_quota.s3_by_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleQuotaCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}

const testAccDataSourceOutscaleQuotaConfig = `
data "outscale_quota" "s3_by_id" {
  quota_name = "vm_limit"
}
`
