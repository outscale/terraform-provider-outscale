package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPILinkRouteTable_basic(t *testing.T) {
	var v oscgo.RouteTable
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
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table_link" {
			continue
		}
		params := oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{
				RouteTableIds: &[]string{rs.Primary.Attributes["route_table_id"]},
			},
		}
		var resp oscgo.ReadRouteTablesResponse
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background(), &oscgo.ReadRouteTablesOpts{ReadRouteTablesRequest: optional.NewInterface(params)})
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

		if len(resp.GetRouteTables()) > 0 {
			return fmt.Errorf(
				"RouteTable: %s has LinkRouteTables", resp.GetRouteTables()[0].GetRouteTableId())
		}
	}
	return nil
}

func testAccCheckOAPILinkRouteTableExists(n string, v *oscgo.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		params := oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{
				RouteTableIds: &[]string{rs.Primary.Attributes["route_table_id"]},
			},
		}
		var resp oscgo.ReadRouteTablesResponse
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background(), &oscgo.ReadRouteTablesOpts{ReadRouteTablesRequest: optional.NewInterface(params)})
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
		if len(resp.GetRouteTables()) == 0 {
			return fmt.Errorf("RouteTable not found")
		}

		*v = resp.GetRouteTables()[0]
		if len(v.GetLinkRouteTables()) == 0 {
			return fmt.Errorf("RouteTable: %s has no LinkRouteTables", v.GetRouteTableId())
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
