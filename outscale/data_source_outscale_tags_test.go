package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOutscaleTagsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPITagsDataSourceConfig(),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccOAPITagsDataSourceConfig() string {
	return `
		data "outscale_tags" "web" {
			filter {
				name   = "resource_type"
				values = ["instance"]
			}
		}
	`
}
