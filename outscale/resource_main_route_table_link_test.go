package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNet_WithLinkMainRouteTable_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_main_route_table_link.main"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkMainRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "link_route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "main"),
					resource.TestCheckResourceAttr(resourceName, "main", "true"),
				),
			},
		},
	})
}

const testAccLinkMainRouteTableConfig = `
	resource "outscale_net" "main_net" {
		ip_range = "10.1.0.0/16"
		tags {
			key = "Name"
			value = "testacc-mainRTable-link"
		}
	}

	resource "outscale_subnet" "mainSubnet" {
		net_id = outscale_net.main_net.net_id
		ip_range = "10.1.1.0/24"
	}

	resource "outscale_route_table" "mainRTable" {
		net_id = outscale_net.main_net.net_id
	}

	resource "outscale_main_route_table_link" "main" {
		net_id = outscale_net.main_net.net_id
		route_table_id = outscale_route_table.mainRTable.id
	}
`
