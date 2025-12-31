package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_User_basic(t *testing.T) {
	resourceName := "outscale_user.basic_user"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
	resourceName := "outscale_user.update_user"
	name := "TestACC_user1"
	newName := "TestACC_user2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
