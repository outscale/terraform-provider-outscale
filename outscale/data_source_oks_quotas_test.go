package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOKSQuotasDataSource_basic(t *testing.T) {
	resourceName := "data.outscale_oks_quotas.quotas"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: `data "outscale_oks_quotas" "quotas" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "clusters_per_project"),
					resource.TestCheckResourceAttrSet(resourceName, "projects"),
				),
			},
		},
	})
}
