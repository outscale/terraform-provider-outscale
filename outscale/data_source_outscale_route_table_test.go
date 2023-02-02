package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_RouteTable_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_route_table.route_table"
	dataSourcesName := "data.outscale_route_tables.route_tables"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_RouteTable_DataSource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "route_tables.#", "1"),

					resource.TestCheckResourceAttrSet(dataSourceName, "net_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "routes.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAcc_RouteTable_DataSource_Config = `
	resource "outscale_net" "test" {
		ip_range = "172.16.0.0/16"

		tags {
			key   = "Name"
			value = "terraform-testacc-rt-data-source"
		}
	}

	resource "outscale_subnet" "test" {
		ip_range = "172.16.0.0/24"
		net_id   = "${outscale_net.test.id}"
	}

	resource "outscale_route_table" "test" {
		net_id = "${outscale_net.test.id}"

		tags {
			key   = "Name"
			value = "terraform-testacc-routetable-data-source"
		}
	}

	resource "outscale_route_table_link" "a" {
		subnet_id      = "${outscale_subnet.test.id}"
		route_table_id = "${outscale_route_table.test.id}"
	}

	data "outscale_route_table" "route_table" {
		filter {
			name   = "route_table_ids"
			values = ["${outscale_route_table.test.id}"]
		}

		depends_on = ["outscale_route_table_link.a"]
	}

	data "outscale_route_tables" "route_tables" {
		filter {
			name   = "route_table_ids"
			values = ["${outscale_route_table.test.id}"]
		}
	}
`
