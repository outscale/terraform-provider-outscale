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
