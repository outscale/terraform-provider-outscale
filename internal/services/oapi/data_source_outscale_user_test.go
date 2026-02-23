package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_data_user_basic(t *testing.T) {
	resourceName := "data.outscale_user.basicTestUser"
	userName := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserBasicConfig(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
				),
			},
		},
	})
}

func testAccDataUserBasicConfig(userName string) string {
	return fmt.Sprintf(`
	resource "outscale_user" "basic_dataUser" {
		user_name = "%s"
		path = "/"
	}
        data "outscale_user" "basicTestUser" {
		filter {
			name = "user_ids"
			values = [outscale_user.basic_dataUser.user_id]
		}
        }
`, userName)
}
