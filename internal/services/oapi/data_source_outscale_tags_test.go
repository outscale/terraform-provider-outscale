package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOutscaleTagsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
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
