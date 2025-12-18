package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_user_group_basic(t *testing.T) {
	resourceName := "outscale_user_group.basic_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

func TestAccOthers_userGroup_with_user(t *testing.T) {
	resourceName := "outscale_user_group.userGroupAcc"
	groupName := "groupWithUsers"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupWithUsers(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName),
				),
			},
		},
	})
}

func TestAccOthers_userGroup_update(t *testing.T) {
	resourceName := "outscale_user_group.userGroupTAcc1"
	groupName := "Gp1UpUser"
	userName := "userGp1"
	newGpName := "Gp2UpUsers"
	newUsName := "userGp2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupUpadate(groupName, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName),
				),
			},
			{
				Config: testAccUserGroupUpadate(newGpName, newUsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "user.#"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", newGpName),
				),
			},
		},
	})
}

const testAccUserGroupBasicConfig = `
	resource "outscale_user_group" "basic_group" {
	  user_group_name = "TestACC_group1"
	  path = "/"
	}`

func testAccUserGroupWithUsers(name string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "userToAdd1" {
			user_name = "userForGp1"
			path = "/"
		}
		resource "outscale_user" "userToAdd2" {
			user_name = "userForGp2"
			path = "/TestPath/"
		}

		resource "outscale_user_group" "userGroupAcc" {
			user_group_name = "%s"
			path = "/"
			user {
				user_name = outscale_user.userToAdd1.user_name
			}
			user {
				user_name = outscale_user.userToAdd2.user_name
				path = "/TestPath/"
			}
		}
	`, name)
}

func testAccUserGroupUpadate(name, userName string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "userUpToAdd01" {
			user_name = "userGp1"
			path = "/"
		}
		resource "outscale_user" "userUpToAdd02" {
			user_name = "userGp2"
			path = "/TestPath/"
		}
		resource "outscale_user_group" "userGroupTAcc1" {
			user_group_name = "%s"
			path = "/"
			user {
				user_name = "%s"
			}
			depends_on = [outscale_user.userUpToAdd01]
		}
	`, name, userName)
}
