package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithLinkRouteTable_basic(t *testing.T) {
	t.Parallel()
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
				),
			},
		},
	})
}

func TestAccNet_ImportLinkRouteTable_Basic(t *testing.T) {
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
			rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
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
			rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
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
