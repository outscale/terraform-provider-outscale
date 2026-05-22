package oapi_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_UserGroup_Basic(t *testing.T) {
	resourceName := "outscale_user_group.basic_group"
	groupName := acctest.RandomWithPrefix("testacc-ug")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupBasicConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_UserGroup_WithUser(t *testing.T) {
	resourceName := "outscale_user_group.userGroupAcc"
	groupName := acctest.RandomWithPrefix("testacc-ug")
	userName1 := acctest.RandomWithPrefix("testacc-user")
	userName2 := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupWithUsers(groupName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_UserGroup_CleanupAfterInvalidUser(t *testing.T) {
	groupName := acctest.RandomWithPrefix("testacc-ug")
	invalidUserName := acctest.RandomWithPrefix("testacc-missing-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccUserGroupWithInvalidUser(groupName, invalidUserName),
				ExpectError: regexp.MustCompile("Unable to add user to User Group"),
			},
			{
				RefreshState: true,
				Check:        testAccCheckOutscaleUserGroupDoesNotExist(groupName),
			},
		},
	})
}

func TestAccOthers_UserGroup_Update(t *testing.T) {
	resourceName := "outscale_user_group.userGroupTAcc1"
	groupName := acctest.RandomWithPrefix("testacc-ug")
	groupName2 := acctest.RandomWithPrefix("testacc-ug")
	userName1 := acctest.RandomWithPrefix("testacc-user")
	userName2 := acctest.RandomWithPrefix("testacc-user")
	userName3 := acctest.RandomWithPrefix("testacc-user")
	userName4 := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupUpdate(groupName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName),
				),
			},
			{
				Config: testAccUserGroupUpdate(groupName2, userName3, userName4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "user.#"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName2),
				),
			},
		},
	})
}

func TestAccOthers_UserGroup_WithPolicy(t *testing.T) {
	resourceName := "outscale_user_group.userGroupAccPolicy"
	groupName := acctest.RandomWithPrefix("testacc-ug")
	userName1 := acctest.RandomWithPrefix("testacc-user")
	userName2 := acctest.RandomWithPrefix("testacc-user")
	policyName := acctest.RandomWithPrefix("testacc-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupWithPolicy(groupName, userName1, userName2, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_group_name", groupName),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_UserGroup_CleanupAfterInvalidPolicy(t *testing.T) {
	groupName := acctest.RandomWithPrefix("testacc-ug")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccUserGroupWithInvalidPolicy(groupName),
				ExpectError: regexp.MustCompile("Unable to link policy to User Group"),
			},
			{
				RefreshState: true,
				Check:        testAccCheckOutscaleUserGroupDoesNotExist(groupName),
			},
		},
	})
}

func TestAccOthers_UserGroup_Migration(t *testing.T) {
	groupName := acctest.RandomWithPrefix("testacc-ug")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.3.1", testAccUserGroupBasicConfig(groupName)),
	})
}

func testAccCheckOutscaleUserGroupDoesNotExist(groupName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client := testacc.ConfiguredClient.OSC
		resp, err := client.ReadUserGroups(context.Background(), osc.ReadUserGroupsRequest{}, options.WithRetryTimeout(time.Minute))
		if err != nil {
			return err
		}
		if resp.UserGroups == nil {
			return nil
		}
		for _, group := range *resp.UserGroups {
			if group.Name != nil && *group.Name == groupName {
				return fmt.Errorf("user group %q still exists", groupName)
			}
		}
		return nil
	}
}

func testAccUserGroupBasicConfig(name string) string {
	return fmt.Sprintf(`
	resource "outscale_user_group" "basic_group" {
	  user_group_name = "%s"
	  path = "/"
	}`, name)
}

func testAccUserGroupWithInvalidUser(groupName, userName string) string {
	return fmt.Sprintf(`
	resource "outscale_user_group" "cleanup_group" {
	  user_group_name = "%s"
	  path = "/"
	  user {
	    user_name = "%s"
	  }
	}`, groupName, userName)
}

func testAccUserGroupWithInvalidPolicy(groupName string) string {
	return fmt.Sprintf(`
	resource "outscale_user_group" "cleanup_group" {
	  user_group_name = "%s"
	  path = "/"
	  policy {
	    policy_orn = "orn:ows:idauth::012345678910:policy/testacc-missing-policy"
	  }
	}`, groupName)
}

func testAccUserGroupWithUsers(groupName, userName1, userName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "userToAdd1" {
			user_name = "%s"
			path = "/"
		}
		resource "outscale_user" "userToAdd2" {
			user_name = "%s"
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
	`, userName1, userName2, groupName)
}

func testAccUserGroupWithPolicy(groupName, userName1, userName2, policyName string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "userToAdd1" {
			user_name = "%s"
			path = "/"
		}
		resource "outscale_user" "userToAdd2" {
			user_name = "%s"
			path = "/TestPath/"
		}

		resource "outscale_policy" "policy-2" {
			policy_name = "%s"
			description = "testacc-user-terraform"
   			document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
			path        = "/"
		}

		resource "outscale_policy_version" "policy-2-v2" {
			policy_orn = outscale_policy.policy-2.orn
			document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
		}

		resource "outscale_user_group" "userGroupAccPolicy" {
			user_group_name = "%s"
			path = "/"
			policy {
				policy_orn = outscale_policy.policy-2.orn
			}
			user {
				user_name = outscale_user.userToAdd1.user_name
			}
			user {
				user_name = outscale_user.userToAdd2.user_name
				path = "/TestPath/"
			}
		}
	`, userName1, userName2, policyName, groupName)
}

func testAccUserGroupUpdate(groupName, userName1, userName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "userUpToAdd01" {
			user_name = "%s"
			path = "/"
		}
		resource "outscale_user" "userUpToAdd02" {
			user_name = "%s"
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
	`, userName1, userName2, groupName, userName1)
}
