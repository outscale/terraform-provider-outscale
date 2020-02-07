package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIInternetServicesDatasource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIInternetServicesDatasourceConfig,
			},
		},
	})
}

const testAccOutscaleOAPIInternetServicesDatasourceConfig = `
	resource "outscale_internet_service" "gateway" {}

	data "outscale_internet_services" "outscale_internet_services" {
		filter {
			name = "internet_service_id"
			values = ["${outscale_internet_service.gateway.id}"]
		}
	}
`
