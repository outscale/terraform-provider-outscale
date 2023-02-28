package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNet_WithRouteTableDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIRouteTableGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIRouteTableCheck("data.outscale_route_table.by_filter"),
					testAccDataSourceOutscaleOAPIRouteTableCheck("data.outscale_route_table.by_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIRouteTableCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		rts, ok := s.RootModule().Resources["outscale_route_table.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_route_table.test in state")
		}
		vpcRs, ok := s.RootModule().Resources["outscale_net.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_net.test in state")
		}
		// FIXME: Missing route_table_links on ReadRouteTables request
		// subnetRs, ok := s.RootModule().Resources["outscale_subnet.test"]
		// if !ok {
		// 	return fmt.Errorf("can't find outscale_subnet.test in state")
		// }
		attr := rs.Primary.Attributes

		if attr["id"] != rts.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				rts.Primary.Attributes["id"],
			)
		}

		if attr["route_table_id"] != rts.Primary.Attributes["id"] {
			return fmt.Errorf(
				"route_table_id is %s; want %s",
				attr["route_table_id"],
				rts.Primary.Attributes["id"],
			)
		}

		if attr["net_id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"net_id is %s; want %s",
				attr["net_id"],
				vpcRs.Primary.Attributes["id"],
			)
		}
		// FIXME: Missing route_table_links on ReadRouteTables request
		// if attr["route_table_links.0.subnet_id"] != subnetRs.Primary.Attributes["id"] {
		// 	return fmt.Errorf(
		// 		"subnet_id is %v; want %s",
		// 		attr["route_table_link.0.subnet_id"],
		// 		subnetRs.Primary.Attributes["id"],
		// 	)
		// }

		return nil
	}
}

const testAccDataSourceOutscaleOAPIRouteTableGroupConfig = `
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

	resource "outscale_route_table_link" "a" {
		subnet_id      = outscale_subnet.test.id
		route_table_id = outscale_route_table.test.id
	}

	data "outscale_route_table" "by_filter" {
		filter {
			name   = "route_table_ids"
			values = [outscale_route_table.test.id]
		}

		depends_on = ["outscale_route_table_link.a"]
	}

	data "outscale_route_table" "by_id" {
		route_table_id = outscale_route_table.test.id
		depends_on     = ["outscale_route_table_link.a"]
	}
`
