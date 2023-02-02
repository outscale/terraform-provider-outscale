package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_PublicCatalog_DataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_PublicCatalog_DataSource_Config(),
			},
		},
	})
}

func testAcc_PublicCatalog_DataSource_Config() string {
	return fmt.Sprintf(`
              data "outscale_public_catalog" "catalog" { }
	`)
}
