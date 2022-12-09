package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleDataSourceAccessKey_withFilters(t *testing.T) {
	t.Parallel()
	dataSourceName := "outscale_access_key.access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeyDataSourceWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
				),
			},
		},
	})
}

func testAccClientAccessKeyDataSourceWithFilters() string {
	return `
		resource "outscale_access_key" "outscale_access_key" {}

		data "outscale_access_key" "access_key" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.outscale_access_key.id]
			}
		}
	`
}
