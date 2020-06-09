package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleGatewayDatasource_basic(t *testing.T) {
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceBasic(rBgpAsn, value),
			},
		},
	})
}

func TestAccOutscaleGatewayDatasource_withFilters(t *testing.T) {
	// datasourceName := "data.outscale_client_gateway.test"
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceWithFilters(rBgpAsn, value),
			},
		},
	})
}

func TestAccOutscaleGatewayDatasource_withFiltersNoLocalhost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientGatewayDatasourceWithFiltersNoLocalhost(),
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
			client_gateway_id = "${outscale_client_gateway.foo.id}"
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
				values = ["${outscale_client_gateway.foo.id}"]
			}
		}
	`, rBgpAsn, value)
}

func testAccClientGatewayDatasourceWithFiltersNoLocalhost() string {
	return fmt.Sprintf(`
	resource "outscale_client_gateway" "outscale_client_gateway" {
		bgp_asn     = 571
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
	`)
}
