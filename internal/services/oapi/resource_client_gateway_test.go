package oapi_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_ClientGateway_Basic(t *testing.T) {
	resourceName := "outscale_client_gateway.foo"
	rBgpAsn := oapihelpers.RandBgpAsn()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "bgp_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_type"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", strconv.Itoa(rBgpAsn)),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_ClientGateway_Migration(t *testing.T) {
	rBgpAsn := oapihelpers.RandBgpAsn()

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccClientGatewayConfig(rBgpAsn),
		),
	})
}

func TestAccOthers_ClientGateway_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_client_gateway.foo"
	rBgpAsn := oapihelpers.RandBgpAsn()
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-client-gateway"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAccClientGatewayConfigWithTag(rBgpAsn, invalidTagKey, tagValue),
			testAccClientGatewayConfigWithTag(rBgpAsn, "Name", tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "client_gateway_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func testAccClientGatewayConfig(rBgpAsn int) string {
	return testAccClientGatewayConfigWithTag(rBgpAsn, "Name", "testacc-client-gateway")
}

func testAccClientGatewayConfigWithTag(rBgpAsn int, tagKey, tagValue string) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = %d
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"

			tags {
				key = %q
				value = %q
			}
		}
	`, rBgpAsn, tagKey, tagValue)
}
