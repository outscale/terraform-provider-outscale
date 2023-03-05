package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNet_WithNatService_basic(t *testing.T) {
	var natService oscgo.NatService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPINatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPINatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPINatGatewayExists("outscale_nat_service.outscale_nat_service", &natService),
				),
			},
		},
	})
}

func TestAccNet_NatService_basicWithDataSource(t *testing.T) {
	var natService oscgo.NatService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPINatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPINatGatewayConfigWithDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPINatGatewayExists("outscale_nat_service.outscale_nat_service", &natService),
				),
			},
		},
	})
}

func testAccCheckOAPINatGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nat_service" {
			continue
		}

		filterReq := oscgo.ReadNatServicesRequest{
			Filters: &oscgo.FiltersNatService{NatServiceIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadNatServicesResponse
		err := resource.Retry(60*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.NatServiceApi.ReadNatServices(context.Background()).ReadNatServicesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetNatServices()) > 0 {
			return fmt.Errorf("Nat Services still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOAPINatGatewayExists(n string, ns *oscgo.NatService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		filterReq := oscgo.ReadNatServicesRequest{
			Filters: &oscgo.FiltersNatService{NatServiceIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadNatServicesResponse
		err := resource.Retry(60*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.NatServiceApi.ReadNatServices(context.Background()).ReadNatServicesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetNatServices()) < 1 {
			return fmt.Errorf("Nat Services not found (%s)", rs.Primary.ID)
		}

		ns = &resp.GetNatServices()[0]

		return nil
	}
}

const testAccOAPINatGatewayConfig = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags {
			key = "Name"
			value = "testacc-nat-service-rs"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id   = outscale_net.outscale_net.net_id
		ip_range = "10.0.0.0/18"
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_nat_service" "outscale_nat_service" {
		depends_on   = ["outscale_route.outscale_route"]
		subnet_id    = outscale_subnet.outscale_subnet.subnet_id
		public_ip_id = outscale_public_ip.outscale_public_ip.public_ip_id
	}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = "${outscale_net.outscale_net.net_id}"
	}

	resource "outscale_route" "outscale_route" {
		depends_on   = [outscale_route_table_link.outscale_route_table_link]
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = "${outscale_internet_service_link.outscale_internet_service_link.internet_service_id}"
		route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
		route_table_id = "${outscale_route_table.outscale_route_table.id}"
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = "${outscale_net.outscale_net.net_id}"
		internet_service_id = "${outscale_internet_service.outscale_internet_service.id}"
	}
`

const testAccOAPINatGatewayConfigWithDataSource = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags {
			key = "Name"
			value = "testacc-nat-service-rs"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id   = "${outscale_net.outscale_net.net_id}"
		ip_range = "10.0.0.0/18"
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_nat_service" "outscale_nat_service" {
		depends_on   = ["outscale_route.outscale_route"]
		subnet_id    = "${outscale_subnet.outscale_subnet.subnet_id}"
		public_ip_id = "${outscale_public_ip.outscale_public_ip.public_ip_id}"
	}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = "${outscale_net.outscale_net.net_id}"
	}

	resource "outscale_route" "outscale_route" {
		depends_on   = [outscale_route_table_link.outscale_route_table_link]
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = "${outscale_internet_service_link.outscale_internet_service_link.internet_service_id}"
		route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
		route_table_id = "${outscale_route_table.outscale_route_table.id}"
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = "${outscale_net.outscale_net.net_id}"
		internet_service_id = "${outscale_internet_service.outscale_internet_service.id}"
	}

	data "outscale_nat_service" "outscale_nat_service" {
		filter {
			name   = "nat_service_ids"
			values = ["${outscale_nat_service.outscale_nat_service.nat_service_id}"]
		}
	}
`
