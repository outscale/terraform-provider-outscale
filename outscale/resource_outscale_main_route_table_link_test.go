package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_WithLinkMainRouteTable_basic(t *testing.T) {
	t.Parallel()
	var v oscgo.RouteTable
	resourceName := "outscale_main_route_table_link.main"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinkMainRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLinkMainRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinkMainRouteTableExists(
						resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "link_route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "main"),
					resource.TestCheckResourceAttr(resourceName, "main", "true"),
				),
			},
		},
	})
}

func testAccCheckLinkMainRouteTableDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_main_route_table_link" {
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

func testAccCheckLinkMainRouteTableExists(n string, v *oscgo.RouteTable) resource.TestCheckFunc {
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
