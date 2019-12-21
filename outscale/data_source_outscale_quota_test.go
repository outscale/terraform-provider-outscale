package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIQuota(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIQuotaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleQuotaOAPICheck("data.outscale_quota.s3_by_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleQuotaOAPICheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIQuotaConfig = `
	data "outscale_quota" "s3_by_id" {
		quota_name = "vm_limit"
	}
`
