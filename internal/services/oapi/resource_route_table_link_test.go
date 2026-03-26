package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_WithLinkRouteTable_Basic(t *testing.T) {
	resourceName := "outscale_route_table_link.foo"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPILinkRouteTableConfig,
			},
			testacc.ImportStepWithStateIdFunc(resourceName, testAccCheckOAPILinkRouteTableImportStateIDFunc(resourceName), testacc.DefaultIgnores()...),
		},
	})
}

func testAccCheckOAPILinkRouteTableImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["route_table_id"], rs.Primary.ID), nil
	}
}

func TestAccNet_WithLinkRouteTable_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.1.3", testAccOAPILinkRouteTableConfig),
	})
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
