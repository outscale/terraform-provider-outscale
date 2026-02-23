package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_data_user_group_basic(t *testing.T) {
	resourceName := "data.outscale_user_group.basicUTest"
	groupName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupBasicConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
				),
			},
		},
	})
}

func testAccDataUserGroupBasicConfig(name string) string {
	return fmt.Sprintf(`
		resource "outscale_user_group" "uGData" {
			user_group_name = "%s"
			path = "/"
		}
		data "outscale_user_group" "basicUTest" {
		    user_group_name = outscale_user_group.uGData.user_group_name
    }
`, name)
}
