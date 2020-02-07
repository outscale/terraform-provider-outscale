package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIDSCustomerGateway_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_client_endpoint.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPICustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPICustomerGatewayDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_client_endpoint.test", "ip_address", "172.0.0.1"),
				),
			},
		},
	})
}

func testAccOAPICustomerGatewayDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_client_endpoint" "foo" {
			bgp_asn = %d
			public_ip = "172.0.0.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-%d"
			}
		}

		data "outscale_client_endpoint" "test" {
			client_endpoint_id = "${outscale_client_endpoint.foo.id}"
		}
	`, rBgpAsn, rInt)
}
