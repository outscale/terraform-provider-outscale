package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_User_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_user.basic_user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

func TestAccOthers_User_update(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_user.update_user"
	name := "TestACC_user1"
	newName := "TestACC_user2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserUpdatedConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
				),
			},
			{
				Config: testAccOutscaleUserUpdatedConfig(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_name", newName),
				),
			},
		},
	})
}

const testAccOutscaleUserBasicConfig = `
    resource "outscale_access_key" "access_key01" {
		state       = "ACTIVE"
		user_name   = outscale_user.basic_user.user_name
		depends_on  = [outscale_user.basic_user]
	}
	resource "outscale_user" "basic_user" {
		user_name = "ACC_test1"
		path = "/"
	}`

func testAccOutscaleUserUpdatedConfig(name string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "update_user" {
			user_name = "%s"
			path = "/"
		}
	`, name)
}
