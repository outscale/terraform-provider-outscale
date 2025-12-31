package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_InternetService_Basic(t *testing.T) {
	resourceName := "outscale_internet_service.internet_service"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.1.2", testAccOutscaleInternetServiceConfig()),
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
