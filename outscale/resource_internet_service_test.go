package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_InternetService_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_internet_service.internet_service"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "internet_service_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				),
			},
		},
	})
}

func testAccOutscaleInternetServiceConfig() string {
	return `
		resource "outscale_internet_service" "internet_service" {
			tags {
				key = "Name"
				value = "testacc-internet-service"
			}
		}
	`
}
