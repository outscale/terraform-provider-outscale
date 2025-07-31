package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithLinkRouteTable_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_route_table_link.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPILinkRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "link_route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "main"),
				),
			},
		},
	})
}

func TestAccNet_ImportLinkRouteTable_Basic(t *testing.T) {
	resourceName := "outscale_route_table_link.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPILinkRouteTableConfig,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOAPILinkRouteTableImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOAPILinkRouteTableImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["route_table_id"], rs.Primary.ID), nil
	}
}

const testAccOAPILinkRouteTableConfig = `
	resource "outscale_net" "foo" {
		ip_range = "10.1.0.0/16"

		tags {
			key = "Name"
			value = "testacc-route-table-link-rs"
		}
	}

	resource "outscale_subnet" "foo" {
		net_id = outscale_net.foo.id
		ip_range = "10.1.1.0/24"
	}

	resource "outscale_route_table" "foo" {
		net_id = outscale_net.foo.id
	}

	resource "outscale_route_table_link" "foo" {
		route_table_id = outscale_route_table.foo.id
		subnet_id = outscale_subnet.foo.id
	}
`
