package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleDataSourceAccessKeys(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.access_keys_with_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "filter.#"),
					resource.TestCheckResourceAttrSet("data.outscale_access_keys.access_keys", "access_keys.#"),
				),
			},
		},
	})
}

func testAccClientAccessKeysDataSource() string {
	return `
		resource "outscale_access_key" "outscale_access_key" {}

		data "outscale_access_keys" "access_keys_with_filter" {
			filter {
				name = "access_key_ids"
				values = ["${outscale_access_key.outscale_access_key.id}"]
			}
		}

		data "outscale_access_keys" "access_keys" {
			filter {
				name = "access_key_ids"
				values = ["${outscale_access_key.outscale_access_key.id}"]
			}
		}
	`
}
