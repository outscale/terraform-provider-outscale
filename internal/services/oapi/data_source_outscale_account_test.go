package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceAccount_basic(t *testing.T) {
	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountConfig(),
			},
		},
	})
}

func testAccDataSourceAccountConfig() string {
	return `
    	data "outscale_account" "account" { }
	`
}
