package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_user_groups_per_user_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_user_groups_per_user.groupList"
	groupName1 := acctest.RandomWithPrefix("testacc-usergroup")
	groupName2 := acctest.RandomWithPrefix("testacc-usergroup")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupsPerUserBasicConfig(groupName1, groupName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_groups.#"),
				),
			},
		},
	})
}

func testAccDataUserGroupsPerUserBasicConfig(name1, name2 string) string {
	return fmt.Sprintf(`
resource "outscale_user" "userForGroup01" {
  user_name = "user1_per_group"
  path      = "/groupsPerUser/"
}

resource "outscale_user_group" "uGroupFORuser" {
			user_group_name = "%s"
			path = "/"
			user {
           user_name = outscale_user.userForGroup01.user_name
           path      = "/groupsPerUser/"
        }
}
resource "outscale_user_group" "uGroup2FORuser" {
		user_group_name = "%s"
		path = "/TestPath/"
		user {
        user_name = outscale_user.userForGroup01.user_name
    }
}
data "outscale_user_groups_per_user" "groupList" {
			user_name   = outscale_user.userForGroup01.user_name
			depends_on =[outscale_user_group.uGroup2FORuser]
}`, name1, name2)
}
