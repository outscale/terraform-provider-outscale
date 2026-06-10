package oapi_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_InternetService_Basic(t *testing.T) {
	resourceName := "outscale_internet_service.internet_service"
	resource.ParallelTest(t, resource.TestCase{
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
		Steps: testacc.FrameworkMigrationTestSteps("1.1.2", testAccOutscaleInternetServiceConfig()),
	})
}

func TestAccOthers_InternetService_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_internet_service.internet_service"
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-internet-service"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAccOutscaleInternetServiceConfigWithTag(invalidTagKey, tagValue),
			testAccOutscaleInternetServiceConfigWithTag("Name", tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "internet_service_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func testAccOutscaleInternetServiceConfig() string {
	return testAccOutscaleInternetServiceConfigWithTag("Name", "testacc-internet-service")
}

func testAccOutscaleInternetServiceConfigWithTag(tagKey, tagValue string) string {
	return `
		resource "outscale_internet_service" "internet_service" {
			tags {
				key = "` + tagKey + `"
				value = "` + tagValue + `"
			}
		}
	`
}
