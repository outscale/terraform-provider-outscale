package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_user_groups_basic(t *testing.T) {
	resourceName := "data.outscale_user_groups.basicUGTest"
	groupName := acctest.RandomWithPrefix("testacc-usergroup")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupsBasicConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_groups.#"),
				),
			},
		},
	})
}

func testAccDataUserGroupsBasicConfig(groupName string) string {
	return fmt.Sprintf(`
		resource "outscale_user_group" "uGroupData" {
			user_group_name = "%s"
			path = "/"
		}
		data "outscale_user_groups" "basicUGTest" {
			filter {
				name   = "user_group_ids"
				values = [outscale_user_group.uGroupData.user_group_id]
			}
    }
`, groupName)
}
