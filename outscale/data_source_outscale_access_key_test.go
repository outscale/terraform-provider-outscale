package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DatasourceAccessKey_basic(t *testing.T) {
	dataSourceName := "data.outscale_access_key.dataKeyBasic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeyDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
				),
			},
		},
	})
}

func TestAccOthers_AccessKey_withFilters(t *testing.T) {
	dataSourceName := "data.outscale_access_key.access_key_filters"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeyDataSourceWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
				),
			},
		},
	})
}

func testAccClientAccessKeyDataSourceBasic() string {
	return `
		resource "outscale_access_key" "accessKeyBasic" {}

		data "outscale_access_key" "dataKeyBasic" {
			access_key_id = outscale_access_key.accessKeyBasic.id
		}
	`
}

func testAccClientAccessKeyDataSourceWithFilters() string {
	return `
		resource "outscale_access_key" "keyFilters" {}

		data "outscale_access_key" "access_key_filters" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.keyFilters.id]
			}
		}
	`
}
