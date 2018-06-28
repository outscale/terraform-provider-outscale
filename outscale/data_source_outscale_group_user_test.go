package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleDSGroupUser_basic(t *testing.T) {
	var group eim.GetGroupOutput

	rInt := acctest.RandInt()
	configBase := fmt.Sprintf(testAccOutscaleDSGroupUserConfig, rInt, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDSGroupUserExists("data.outscale_groups_for_user.team", &group),
				),
			},
		},
	})
}

func testAccCheckOutscaleDSGroupUserExists(n string, g *eim.GetGroupOutput) resource.TestCheckFunc {
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

const testAccOutscaleDSGroupUserConfig = `
resource "outscale_group" "group" {
	group_name = "test-group-%d"
	path = "/"
}

resource "outscale_user" "user" {
	user_name = "test-user-%d"
	path = "/"
}

resource "outscale_group_user" "team" {
	user_name = "${outscale_user.user.user_name}"
	group_name = "${outscale_group.group.group_name}"
}

data "outscale_groups_for_user" "team" {
	user_name = "${outscale_user.user.user_name}"
}
`
