package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_WithRouteTableDataSource_basic(t *testing.T) {
	resourceName := "data.outscale_route_table.by_filter"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleRouteTableGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleRouteTableGroupConfig = `
	resource "outscale_net" "test" {
		ip_range = "172.16.0.0/16"

		tags {
			key   = "Name"
			value = "terraform-testacc-rt-data-source"
		}
	}

	resource "outscale_subnet" "test" {
		ip_range = "172.16.0.0/24"
		net_id   = outscale_net.test.id
	}

	resource "outscale_route_table" "test" {
		net_id = outscale_net.test.id

		tags {
			key   = "Name"
			value = "terraform-testacc-routetable-data-source"
		}
	}

	resource "outscale_route_table_link" "rtLink" {
		subnet_id      = outscale_subnet.test.id
		route_table_id = outscale_route_table.test.id
	}

	data "outscale_route_table" "by_filter" {
		filter {
			name   = "route_table_ids"
			values = [outscale_route_table.test.id]
		}

		depends_on = [outscale_route_table_link.rtLink]
	}

	data "outscale_route_table" "by_id" {
		route_table_id = outscale_route_table.test.id
		depends_on     = [outscale_route_table_link.rtLink]
	}
`
