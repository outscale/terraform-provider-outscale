package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVM_tags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var v fcu.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.foo", &v),
					testAccCheckTags(&v.Tags, "foo", "bar"),
					testAccCheckTags(&v.Tags, "#", ""),
				),
			},
		},
	})
}

func testAccCheckTags(
	ts *[]*fcu.Tag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := tagsToMap(*ts)
		v, ok := m[0]["key"]
		if value != "" && !ok {
			return fmt.Errorf("Missing tag: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra tag: %s", key)
		}
		if value == "" {
			return nil
		}

		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}

		return nil
	}
}

const testAccCheckInstanceConfigTags = `
resource "outscale_tag" "outscale_tag" {
    resource_ids = ["i-9ea2e54f"]

    tag {                               # NOK should be tag not tags
        name7 = "testDataSource7"          # NOK delete doesn't delete tag
        #name8 = "testDataSource8"          # tfa doesn't display correctly
    }                                      # tfs displays nothing
}

`
