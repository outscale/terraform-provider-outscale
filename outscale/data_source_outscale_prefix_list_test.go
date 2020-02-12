package outscale

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceOutscaleOAPIPrefixList(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIPrefixListConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIPrefixListCheck("data.outscale_prefix_list.s3_by_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIPrefixListCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr["prefix_list_name"] != "com.outscale.eu-west-2.osu" {
			return fmt.Errorf("bad name %s", attr["prefix_list_name"])
		}
		if attr["prefix_list_id"] != "pl-a14a8cdc" {
			return fmt.Errorf("bad id %s", attr["prefix_list_id"])
		}

		var (
			cidrBlockSize int
			err           error
		)

		if cidrBlockSize, err = strconv.Atoi(attr["cidr_set.#"]); err != nil {
			return err
		}
		if cidrBlockSize < 1 {
			return fmt.Errorf("cidr_set seem suspiciously low: %d", cidrBlockSize)
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIPrefixListConfig = `
	data "outscale_prefix_list" "s3_by_id" {
		prefix_list_id = "pl-a14a8cdc"
	}
`
