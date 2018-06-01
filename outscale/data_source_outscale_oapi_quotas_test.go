package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIQuotas(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIQuotasConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_quotas.s3_by_id", "reference_quota_set.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIQuotasConfig = `
data "outscale_quotas" "s3_by_id" {
  quota_name = ["vm_limit"]
}
`
