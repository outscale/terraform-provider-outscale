package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_InternetServiceDatasource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
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
