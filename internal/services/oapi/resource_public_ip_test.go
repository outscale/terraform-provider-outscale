package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_PublicIP_Basic(t *testing.T) {
	resourceName := "outscale_public_ip.pip"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
				),
			},
			// Ignore attributes related to the Public IP Link, that gets populated after a refresh
			testacc.ImportStep(resourceName, "link_public_ip_id", "nic_account_id", "nic_id", "private_ip", "vm_id", "request_id"),
		},
	})
}

func TestAccVM_PublicIP_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.4.0", testAccPublicIPConfig),
	})
}

func TestAccOthers_PublicIP_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_public_ip.pip"
	invalidTagKey := strings.Repeat("a", 256)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAccPublicIPConfigWithTag(invalidTagKey, "public_ip_test"),
			testAccPublicIPConfigWithTag("Name", "public_ip_test_recovery"),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", "public_ip_test_recovery"),
			),
		),
	})
}

const testAccPublicIPConfig = `
resource "outscale_public_ip" "pip" {
	tags {
		key = "Name"
		value = "public_ip_test"
	}
}
`

func testAccPublicIPConfigWithTag(tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "outscale_public_ip" "pip" {
	tags {
		key = %q
		value = %q
	}
}
`, tagKey, tagValue)
}
