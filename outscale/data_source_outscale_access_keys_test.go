package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_DataSourceAccessKeys_basic(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_access_keys.read_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessKeyDestroy,
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
	dataSourceName := "data.outscale_access_keys.filters_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessKeyDestroy,
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

func testAccCheckAccessKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_access_key" {
			continue
		}
		req := oscgo.ReadAccessKeysRequest{}
		req.Filters = &oscgo.FiltersAccessKeys{
			AccessKeyIds: &[]string{rs.Primary.ID},
		}

		var resp oscgo.ReadAccessKeysResponse
		var err error
		exists := false
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.AccessKeyApi.ReadAccessKeys(context.Background()).ReadAccessKeysRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("AccessKeys reading (%s)", rs.Primary.ID)
		}

		for _, ca := range resp.GetAccessKeys() {
			if ca.GetAccessKeyId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("Access_Key still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
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
