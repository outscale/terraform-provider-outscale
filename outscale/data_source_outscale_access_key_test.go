package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_AccessKey_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_key.outscale_access_key_d"
	dataSourcesName := "data.outscale_access_keys.outscale_access_keys_d"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_AccessKey_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),

					resource.TestCheckResourceAttrSet(dataSourcesName, "access_keys.#"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "filter.#"),
				),
			},
		},
	})
}

func testAcc_AccessKey_DataSource_Config() string {
	return `
		resource "outscale_access_key" "outscale_access_key" {}

		data "outscale_access_key" "outscale_access_key_d" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.outscale_access_key.id]
			}
			filter {
				name = "states"
				values = [outscale_access_key.outscale_access_key.state]
			}
		}

		data "outscale_access_keys" "outscale_access_keys_d" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.outscale_access_key.id]
			}
			filter {
				name = "states"
				values = [outscale_access_key.outscale_access_key.state]
			}
		}
	`
}
