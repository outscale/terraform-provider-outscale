package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTagsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsDataSourceConfig(),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccTagsDataSourceConfig() string {
	return `
		data "outscale_tags" "web" {
			filter {
				name   = "resource_type"
				values = ["instance"]
			}
		}
	`
}
