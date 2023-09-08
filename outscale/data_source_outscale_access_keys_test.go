package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOthers_DataSourceAccessKeys_basic(t *testing.T) {
	t.Parallel()
	DataSourceName := "data.outscale_access_keys.outscale_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(DataSourceName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(DataSourceName, "access_key_ids.#"),
				),
			},
		},
	})
}

func TestAccOthers_DataSourceAccessKeys_withFilters(t *testing.T) {
	t.Parallel()
	DataSourceName := "data.outscale_access_keys.outscale_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSourceWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(DataSourceName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(DataSourceName, "filter.#"),
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
