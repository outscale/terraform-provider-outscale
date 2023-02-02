package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_InternetService_Datasource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_internet_service.internet_service"
	dataSourcesName := "data.outscale_internet_services.internet_services"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_InternetService_Datasource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "internet_services.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "internet_service_id"),
				),
			},
		},
	})
}

const testAcc_InternetService_Datasource_Config = `
	resource "outscale_internet_service" "gateway" {}

	data "outscale_internet_service" "internet_service" {
		filter {
			name = "internet_service_ids"
			values = ["${outscale_internet_service.gateway.internet_service_id}"]
		}
	}

	data "outscale_internet_services" "internet_services" {
		filter {
			name = "internet_service_id"
			values = ["${outscale_internet_service.gateway.internet_service_id}"]
		}
	}
`
