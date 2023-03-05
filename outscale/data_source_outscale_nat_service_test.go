package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNet_WithNatServiceDataSource_Instance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPINatServiceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNatServiceDataSourceID("data.outscale_nat_service.nat"),
					resource.TestCheckResourceAttrSet("data.outscale_nat_service.nat", "subnet_id"),
				),
			},
		},
	})
}

func testAccCheckOutscaleNatServiceDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Nat Service data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Nat Service data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPINatServiceDataSourceConfig = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-nat-service-ds"
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
		net_id = outscale_net.outscale_net.net_id
	}

	resource "outscale_route" "outscale_route" {
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = outscale_subnet.outscale_subnet.subnet_id
		route_table_id = outscale_route_table.outscale_route_table.id
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = outscale_net.outscale_net.net_id
		internet_service_id = outscale_internet_service.outscale_internet_service.id
	}

	data "outscale_nat_service" "nat" {
		nat_service_id = outscale_nat_service.outscale_nat_service.id
	}
`
