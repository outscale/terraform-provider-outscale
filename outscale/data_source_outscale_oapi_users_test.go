package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleOAPIUsersDS_basic(t *testing.T) {
	var conf eim.GetUserOutput

	name1 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	path1 := "/"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIUsersDSConfig(name1, path1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUsersDSExists("data.outscale_users.outscale_users", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIUsersDSExists(n string, res *eim.GetUserOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User name is set")
		}

		return nil
	}
}

func testAccOutscaleOAPIUsersDSConfig(r, p string) string {
	return fmt.Sprintf(`
resource "outscale_user" "user" {
	user_name = "%s"
	path = "%s"
}

data "outscale_users" "outscale_users" {
	path = "${outscale_user.user.path}"
}`, r, p)
}
