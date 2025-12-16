package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_user_groups_per_user_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_user_groups_per_user.groupList"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupsPerUserBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_groups.#"),
				),
			},
		},
	})
}

const testAccDataUserGroupsPerUserBasicConfig = `
resource "outscale_user" "userForGroup01" {
  user_name = "user1_per_group"
  path      = "/groupsPerUser/"
}

resource "outscale_user_group" "uGroupFORuser" {
		user_group_name = "Group1TestACC"
		path = "/"
		user {
           user_name = outscale_user.userForGroup01.user_name
           path      = "/groupsPerUser/"
        }
}
resource "outscale_user_group" "uGroup2FORuser" {
	user_group_name = "Group02TestACC"
	path = "/TestPath/"
	user {
        user_name = outscale_user.userForGroup01.user_name
    }
}
data "outscale_user_groups_per_user" "groupList" {
		user_name   = outscale_user.userForGroup01.user_name
		depends_on =[outscale_user_group.uGroup2FORuser]
}`
