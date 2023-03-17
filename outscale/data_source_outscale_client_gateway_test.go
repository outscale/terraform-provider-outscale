package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAcc_ClientGateway_Datasource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_client_gateway.test"
	dataSourcesName := "data.outscale_client_gateways.test"
	rBgpAsn := utils.RandIntRange(64512, 65534)
	value := fmt.Sprintf("testacc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_ClientGateway_Datasource(rBgpAsn, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "client_gateway_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bgp_asn"),

					resource.TestCheckResourceAttrSet(dataSourcesName, "client_gateways.#"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "filter.#"),
				),
			},
		},
	})
}

func testAcc_ClientGateway_Datasource(rBgpAsn int, value string) string {
	return fmt.Sprintf(`
		resource "outscale_client_gateway" "foo" {
			bgp_asn         = %d
			public_ip       = "171.33.75.123"
			connection_type = "ipsec.1"

			tags {
				key = "Name"
				value = "%s"
			}
		}
	
		resource "outscale_client_gateway" "foo2" {
			bgp_asn         = 4
			public_ip       = "172.0.0.1"
			connection_type = "ipsec.1"
		}
	
		data "outscale_client_gateways" "test" {
			filter {
				name = "client_gateway_ids"
				values = [outscale_client_gateway.foo.id, outscale_client_gateway.foo2.id]
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
