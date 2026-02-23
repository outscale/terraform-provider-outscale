package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_InternetServiceDatasource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_internet_service.outscale_internet_serviced", "internet_service_id"),
				),
			},
		},
	})
}

const testAccOutscaleInternetServiceDatasourceConfig = `
	resource "outscale_internet_service" "outscale_internet_service" {}

	data "outscale_internet_service" "outscale_internet_serviced" {
		filter {
			name = "internet_service_ids"
			values = [outscale_internet_service.outscale_internet_service.internet_service_id]
		}
	}
`
