package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIINternetServiceDatasource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIINternetServiceDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_internet_service.outscale_internet_serviced", "internet_service_id"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIINternetServiceDatasourceConfig = `
	resource "outscale_internet_service" "outscale_internet_service" {}

	data "outscale_internet_service" "outscale_internet_serviced" {
		filter {
			name = "internet_service_ids"
			values = ["${outscale_internet_service.outscale_internet_service.internet_service_id}"]
		}
	}
`
