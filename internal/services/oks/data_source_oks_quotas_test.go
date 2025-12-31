package oks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOKSQuotasDataSource_basic(t *testing.T) {
	resourceName := "data.outscale_oks_quotas.quotas"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

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
