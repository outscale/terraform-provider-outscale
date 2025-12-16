package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_user_groups_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_user_groups.basicUGTest"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupsBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_groups.#"),
				),
			},
		},
	})
}

const testAccDataUserGroupsBasicConfig = `
	resource "outscale_user_group" "uGroupData" {
		user_group_name = "TestACC_uGdata"
		path = "/"
	}
	data "outscale_user_groups" "basicUGTest" {
		filter {
			name   = "user_group_ids"
			values = [outscale_user_group.uGroupData.user_group_id]
		}
    }
`
