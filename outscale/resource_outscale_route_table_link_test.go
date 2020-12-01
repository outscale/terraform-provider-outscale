package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPILinkRouteTable_basic(t *testing.T) {
	var v oscgo.RouteTable
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPILinkRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPILinkRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPILinkRouteTableExists(
						"outscale_route_table_link.foo", &v),
					resource.TestCheckResourceAttrSet(
						"outscale_route_table_link.foo", "link_route_table_id"),
					resource.TestCheckResourceAttrSet(
						"outscale_route_table_link.foo", "main"),
					resource.TestCheckResourceAttrSet(
						"outscale_route_table_link.foo", "request_id"),
				),
			},
		},
	})
}

func TestAccResourceOAPILinkRouteTable_importBasic(t *testing.T) {
	resourceName := "outscale_route_table_link.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPILinkRouteTableDestroy,
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
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
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
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
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
			value = "testacc-route-table-link-rs"
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
