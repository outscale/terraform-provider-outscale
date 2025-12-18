package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_InternetService_Basic(t *testing.T) {
	resourceName := "outscale_internet_service.internet_service"
	resource.ParallelTest(t, resource.TestCase{
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

func TestAccOthers_InternetService_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.2", testAccOutscaleInternetServiceConfig()),
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
