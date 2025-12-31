package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_InternetServicesDatasource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServicesDatasourceConfig,
			},
		},
	})
}

const testAccOutscaleInternetServicesDatasourceConfig = `
	resource "outscale_internet_service" "gateway" {}

	data "outscale_internet_services" "outscale_internet_services" {
		filter {
			name = "internet_service_ids"
			values = [outscale_internet_service.gateway.id]
		}
	}
`
