package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_users_basic(t *testing.T) {
	resourceName := "data.outscale_users.basicTestUsers"
	userName := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleDataUserBasicConfig(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
				),
			},
		},
	})
}

func testAccOutscaleDataUserBasicConfig(userName string) string {
	return fmt.Sprintf(`
	resource "outscale_user" "basic_data_users" {
	  user_name = "%s"
	  path = "/"
	}
        data "outscale_users" "basicTestUsers" {
          depends_on = [outscale_user.basic_data_users]
        }`, userName)
}
