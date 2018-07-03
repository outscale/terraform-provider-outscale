package outscale

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleGroupsDS_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var conf eim.GetGroupOutput
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleGroupsDSConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleGroupsDSExists("data.outscale_groups.outscale_groups", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleGroupsDSExists(n string, res *eim.GetGroupOutput) resource.TestCheckFunc {
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

func testAccOutscaleGroupsDSConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_group" "group" {
		group_name = "test-group-%d"
		path = "/"
	}
	
	
	data "outscale_groups" "outscale_groups" {
		path_prefix = "${outscale_group.group.path}"
	}
	
	`, rInt)
}
