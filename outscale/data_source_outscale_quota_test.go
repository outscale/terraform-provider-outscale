package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Quota_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_quota.lbu-quota"
	dataSourcesName := "data.outscale_quotas.all-quotas"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Quota_DataSource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "quotas.#"),

					resource.TestCheckResourceAttr(dataSourceName, "name", "lb_listeners_limit"),
				),
			},
		},
	})
}

const testAcc_Quota_DataSource_Config = `
	data "outscale_quota" "lbu-quota" { 		
		filter {
        	name     = "quota_names"
        	values   = ["lb_listeners_limit"]
    	}
	}

	data "outscale_quotas" "all-quotas" {}
`
