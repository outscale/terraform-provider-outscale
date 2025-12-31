package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_WithInternetServiceLink_Basic(t *testing.T) {
	resourceName := "outscale_internet_service_link.outscale_internet_service_link"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_service_id"),
				),
			},
		},
	})
}

func TestAccNet_WithImportInternetServiceLink_Basic(t *testing.T) {
	resourceName := "outscale_internet_service_link.outscale_internet_service_link"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_WithInternetServiceLink_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.1.3", testAccOutscaleInternetServiceLinkConfig()),
	})
}

func testAccOutscaleInternetServiceLinkConfig() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-internet-service-link-rs"
			}
		}

		resource "outscale_internet_service" "outscale_internet_service" {}

		resource "outscale_internet_service_link" "outscale_internet_service_link" {
			net_id              = outscale_net.outscale_net.net_id
			internet_service_id = outscale_internet_service.outscale_internet_service.id
		}
	`
}
