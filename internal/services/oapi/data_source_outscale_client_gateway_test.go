package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_GatewayDatasource_basic(t *testing.T) {
	rBgpAsn := oapihelpers.RandBgpAsn()
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceBasic(rBgpAsn, value),
			},
		},
	})
}

func TestAccOthers_GatewayDatasource_withFilters(t *testing.T) {
	// datasourceName := "data.outscale_client_gateway.test"
	rBgpAsn := oapihelpers.RandBgpAsn()
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceWithFilters(rBgpAsn, value),
			},
		},
	})
}

func TestAccOthers_GatewayDatasource_withFiltersNoLocalhost(t *testing.T) {
	bgpAsn := oapihelpers.RandBgpAsn()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceWithFiltersNoLocalhost(bgpAsn),
			},
		},
	})
}

func testAccClientGatewayDatasourceBasic(rBgpAsn int, value string) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = %d
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"

			tags {
				key = "Name"
				value = "%s"
			}
		}

		data "outscale_client_gateway" "test" {
			client_gateway_id = outscale_client_gateway.foo.id
		}
	`, rBgpAsn, value)
}

func testAccClientGatewayDatasourceWithFilters(rBgpAsn int, value string) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = %d
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"

			tags {
				key = "Name"
				value = "%s"
			}
		}

		data "outscale_client_gateway" "test" {
			filter {
				name = "client_gateway_ids"
				values = [outscale_client_gateway.foo.id]
			}
		}
	`, rBgpAsn, value)
}

func testAccClientGatewayDatasourceWithFiltersNoLocalhost(asn int) string {
	return fmt.Sprintf(`
	resource "outscale_client_gateway" "outscale_client_gateway" {
		bgp_asn     = %d
		public_ip  = "171.33.75.123"
		connection_type        = "ipsec.1"
		tags {
		 key = "name-mzi"
		 value = "CGW_1_mzi"
		}
	}

	data "outscale_client_gateway" "outscale_client_gateway_2" {
		filter {
		   name   = "client_gateway_ids"
		   values = [outscale_client_gateway.outscale_client_gateway.client_gateway_id]
		}
	}
	`, asn)
}
