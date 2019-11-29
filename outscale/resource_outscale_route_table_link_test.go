package outscale

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPILinkRouteTable_basic(t *testing.T) {
	var v oapi.RouteTable
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPILinkRouteTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPILinkRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPILinkRouteTableExists(
						"outscale_route_table_link.foo", &v),
				),
			},
		},
	})
}

func testAccCheckOAPILinkRouteTableDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table_link" {
			continue
		}
		params := &oapi.ReadRouteTablesRequest{
			Filters: oapi.FiltersRouteTable{
				RouteTableIds: []string{rs.Primary.Attributes["route_table_id"]},
			},
		}
		var resp *oapi.POST_ReadRouteTablesResponses
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadRouteTables(*params)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					log.Printf("[DEBUG] Trying to create route again: %q", err)
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
				return nil
			}
			return err
		}

		if len(resp.OK.RouteTables) > 0 {
			return fmt.Errorf(
				"RouteTable: %s has LinkRouteTables", resp.OK.RouteTables[0].RouteTableId)
		}
	}
	return nil
}

func testAccCheckOAPILinkRouteTableExists(n string, v *oapi.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		params := &oapi.ReadRouteTablesRequest{
			Filters: oapi.FiltersRouteTable{
				RouteTableIds: []string{rs.Primary.Attributes["route_table_id"]},
			},
		}
		var resp *oapi.POST_ReadRouteTablesResponses
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadRouteTables(*params)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
					log.Printf("[DEBUG] Trying to create route again: %q", err)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}
		if len(resp.OK.RouteTables) == 0 {
			return fmt.Errorf("RouteTable not found")
		}

		*v = resp.OK.RouteTables[0]
		if len(v.LinkRouteTables) == 0 {
			return fmt.Errorf("RouteTable: %s has no LinkRouteTables", v.RouteTableId)
		}

		return nil
	}
}

const testAccOAPILinkRouteTableConfig = `
	resource "outscale_net" "foo" {
		ip_range = "10.1.0.0/16"

		tags {
			key = "Name"
			value = "outscale_net"
		}
	}

	resource "outscale_subnet" "foo" {
		net_id = "${outscale_net.foo.id}"
		ip_range = "10.1.1.0/24"
	}

	resource "outscale_route_table" "foo" {
		net_id = "${outscale_net.foo.id}"
	}

	resource "outscale_route_table_link" "foo" {
		route_table_id = "${outscale_route_table.foo.id}"
		subnet_id = "${outscale_subnet.foo.id}"
	}
`
