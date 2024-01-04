package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_DataSourceAccessKeys_basic(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.outscale_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_key_ids.#"),
				),
			},
		},
	})
}

func TestAccOthers_DataSourceAccessKeys_withFilters(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.outscale_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSourceWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "filter.#"),
				),
			},
		},
	})
}

func testAccClientAccessKeysDataSourceBasic() string {
	return `
		resource "outscale_access_key" "outscale_access_key" {}

		data "outscale_access_keys" "outscale_access_key" {
			access_key_ids = [outscale_access_key.outscale_access_key.id]
		}
	`
}

func testAccClientAccessKeysDataSourceWithFilters() string {
	return `
		resource "outscale_access_key" "outscale_access_key" {}

		data "outscale_access_keys" "outscale_access_key" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.outscale_access_key.id]
			}
		}
	`
}
