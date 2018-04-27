package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscalePrefixLists(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscalePrefixListsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_prefix_lists.s3_by_id", "prefix_list_set.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscalePrefixListsConfig = `
data "outscale_prefix_lists" "s3_by_id" {
  prefix_list_id = ["pl-a14a8cdc"]
}
`
