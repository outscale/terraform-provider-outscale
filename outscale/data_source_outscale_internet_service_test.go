package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccInternetServiceDatasource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternetServiceDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_internet_service.outscale_internet_serviced", "internet_service_id"),
				),
			},
		},
	})
}

const testAccInternetServiceDatasourceConfig = `
	resource "outscale_internet_service" "outscale_internet_service" {}

	data "outscale_internet_service" "outscale_internet_serviced" {
		filter {
			name = "internet_service_ids"
			values = ["${outscale_internet_service.outscale_internet_service.internet_service_id}"]
		}
	}
`
