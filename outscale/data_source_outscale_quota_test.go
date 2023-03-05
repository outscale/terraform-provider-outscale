package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOthers_DataSourceQuota(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIQuotaConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIQuotaConfig = `
	data "outscale_quota" "lbu-quota" { 
  filter {
        name     = "quota_names"
        values   = ["lb_listeners_limit"]
    }
}
`
