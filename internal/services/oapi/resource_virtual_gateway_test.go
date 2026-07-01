package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_VirtualGateway_Basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway.outscale_virtual_gateway"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualGatewayWithTags("name", "testacc-vgw"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "testacc-vgw"),
				),
			},
			{
				Config: testAccVirtualGatewayWithTags("name", "testacc-vgw-up"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "testacc-vgw-up"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_VirtualGateway_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_virtual_gateway.outscale_virtual_gateway"
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-resource-virtual-gateway"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAccVirtualGatewayWithTags(invalidTagKey, tagValue),
			testAccVirtualGatewayWithTags("name", tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "virtual_gateway_id"),
				resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func TestAccOthers_VirtualGateway_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.6.0", testAccVirtualGateway),
	})
}

const testAccVirtualGateway = `
resource "outscale_virtual_gateway" "foo" {
    connection_type = "ipsec.1"
}`

func testAccVirtualGatewayWithTags(tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
	connection_type = "ipsec.1"
	tags {
		key = %q
		value = %q
	}
}`, tagKey, tagValue)
}
