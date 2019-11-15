package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPIRoute_noopdiff(t *testing.T) {
	var route oapi.Route

	testCheck := func(s *terraform.State) error {
		return nil
	}

	testCheckChange := func(s *terraform.State) error {
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
					testCheck,
				),
			},
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
					testCheckChange,
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIRouteExists(n string, res *oapi.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		r, _, err := findResourceOAPIRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_ip_range"],
		)

		if err != nil {
			return err
		}

		if r == nil {
			return fmt.Errorf("Route not found")
		}

		*res = *r

		return nil
	}
}

func testAccCheckOAPIOutscaleRouteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		route, _, err := findResourceOAPIRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_ip_range"],
		)

		if route == nil && err == nil {
			return nil
		}
	}

	return nil
}

var testAccOutscaleOAPIRouteNoopChange = fmt.Sprint(`
	resource "outscale_net" "test" {
		ip_range = "10.0.0.0/24"
	}

	resource "outscale_route_table" "test" {
		net_id = "${outscale_net.test.net_id}"
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		internet_service_id = "${outscale_internet_service.outscale_internet_service.id}"
		net_id              = "${outscale_net.test.net_id}"
	}

	resource "outscale_route" "test" {
		gateway_id           = "${outscale_internet_service.outscale_internet_service.id}"
		destination_ip_range = "10.0.0.0/16"
		route_table_id       = "${outscale_route_table.test.route_table_id}"
	}
`)
