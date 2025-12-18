package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DataSourcesQuotas(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcesOutscaleQuotasConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourcesOutscaleQuotasConfig = `
data "outscale_quotas" "all-quotas" {}
`
