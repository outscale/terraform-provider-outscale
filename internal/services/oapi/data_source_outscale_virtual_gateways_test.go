package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataSourceVpnGateways_unattached(t *testing.T) {
	// t.Skip()

	resource.ParallelTest(t, resource.TestCase{
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVpnGatewaysUnattachedConfig(),
			},
		},
	})
}

func testAccDataSourceOutscaleVpnGatewaysUnattachedConfig() string {
	return `
		resource "outscale_virtual_gateway" "unattached" {
			connection_type = "ipsec.1"
		}

		data "outscale_virtual_gateways" "test_by_id" {
			virtual_gateway_id = [outscale_virtual_gateway.unattached.id]
		}
	`
}
