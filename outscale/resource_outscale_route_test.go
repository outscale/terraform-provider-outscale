package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIRoute_noopdiff(t *testing.T) {
	var route oscgo.Route

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
				),
			},
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIRoute_importBasic(t *testing.T) {

	resourceName := "outscale_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleOAPIRouteImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func TestAccOutscaleOAPIRoute_importWithNatService(t *testing.T) {

	resourceName := "outscale_route.outscale_route_nat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIRouteWithNatService,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleOAPIRouteImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleOAPIRouteImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s_%s", rs.Primary.ID, rs.Primary.Attributes["destination_ip_range"]), nil
	}
}

func testAccCheckOutscaleOAPIRouteExists(n string, res *oscgo.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
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

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
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

var testAccOutscaleOAPIRouteWithNatService = fmt.Sprint(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key   = "name"
			value = "net"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id   = "${outscale_net.outscale_net.net_id}"
		ip_range = "10.0.0.0/18"
		tags {
			key   = "name"
			value = "subnet"
		}
	}

	resource "outscale_public_ip" "outscale_public_ip" {
		tags {
			key   = "name"
			value = "public_ip"
		}
	}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = "${outscale_net.outscale_net.net_id}"
		tags {
			key   = "name"
			value = "route_table"
		}
	}

	resource "outscale_route" "outscale_route" {
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = "${outscale_internet_service.outscale_internet_service.internet_service_id}"
		route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
		route_table_id = "${outscale_route_table.outscale_route_table.id}"
	}

	resource "outscale_internet_service" "outscale_internet_service" {
		tags {
			key   = "name"
			value = "internet_service"
		}
	}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = "${outscale_net.outscale_net.net_id}"
		internet_service_id = "${outscale_internet_service.outscale_internet_service.id}"
	}

	resource "outscale_nat_service" "outscale_nat_service" {
		depends_on   = ["outscale_route.outscale_route"]
		subnet_id    = "${outscale_subnet.outscale_subnet.subnet_id}"
		public_ip_id = "${outscale_public_ip.outscale_public_ip.public_ip_id}"
		tags {
			key   = "name"
			value = "nat"
		}
	}

	resource "outscale_route" "outscale_route_nat" {
		destination_ip_range = "40.0.0.0/16"
		nat_service_id       = "${outscale_nat_service.outscale_nat_service.nat_service_id}"
		route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
	}
`)
