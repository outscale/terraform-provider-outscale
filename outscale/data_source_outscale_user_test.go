package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleUserDS_basic(t *testing.T) {
	var conf eim.GetUserOutput

	name1 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	path1 := "/"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleUserDSConfig(name1, path1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserDSExists("data.outscale_user.outscale_user", &conf),
					resource.TestCheckResourceAttrSet("data.outscale_user.outscale_user", "path"),
				),
			},
		},
	})
}

func testAccCheckOutscaleUserDSExists(n string, res *eim.GetUserOutput) resource.TestCheckFunc {
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

func testAccOutscaleUserDSConfig(r, p string) string {
	return fmt.Sprintf(`
resource "outscale_user" "user" {
	user_name = "%s"
	path = "%s"
}

data "outscale_user" "outscale_user" {
	user_name = "${outscale_user.user.user_name}"
}`, r, p)
}
