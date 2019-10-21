package outscale

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleOAPIGroupsDS_basic(t *testing.T) {
	t.Skip()

	var conf eim.GetGroupOutput
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIGroupsDSConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIGroupsDSExists("data.outscale_groups.outscale_groups", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIGroupsDSExists(n string, res *eim.GetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Group name is set")
		}

		return nil
	}
}

func testAccOutscaleOAPIGroupsDSConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_group" "group" {
			group_name = "test-group-%d"
			path       = "/"
		}
		
		
		data "outscale_groups" "outscale_groups" {
			path = "${outscale_group.group.path}"
		}
	`, rInt)
}
