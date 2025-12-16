package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_InternetServicesDatasource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
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
