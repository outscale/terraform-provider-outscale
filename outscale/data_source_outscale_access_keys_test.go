package outscale

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_DataSourceAccessKeys_basic(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.read_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientAccessKeysDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "access_keys.0.access_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_keys.0.expiration_date"),
				),
			},
		},
	})
}

func TestAccOthers_DataSourceAccessKeys_withFilters(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.filters_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
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
	creationDate := time.Now().AddDate(1, 1, 0).Format("2006-01-02")
	return fmt.Sprintf(`
		resource "outscale_access_key" "data_access_key" {
		expiration_date = "%s"
		}

		data "outscale_access_keys" "read_access_key" {
			access_key_ids = [outscale_access_key.data_access_key.id]
		}
	`, creationDate)
}

func testAccClientAccessKeysDataSourceWithFilters() string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "datas_access_key" {}

		data "outscale_access_keys" "filters_access_key" {
			filter {
				name = "access_key_ids"
				values = [outscale_access_key.datas_access_key.id]
			}
		}
	`)
}
