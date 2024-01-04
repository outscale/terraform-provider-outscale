package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_DataSourceVpnGateways_unattached(t *testing.T) {
	//t.Skip()
	t.Parallel()
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVpnGatewaysUnattachedConfig(rInt),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpnGatewaysUnattachedConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "unattached" {
			connection_type = "ipsec.1"	
		}
		
		data "outscale_virtual_gateways" "test_by_id" {
			virtual_gateway_id = [outscale_virtual_gateway.unattached.id]
		}
	`)
}
