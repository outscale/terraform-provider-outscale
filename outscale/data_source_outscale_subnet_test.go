package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceSubnet(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSubnetCheck("data.outscale_subnet.by_id"),
					testAccDataSourceSubnetCheck("data.outscale_subnet.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceSubnet_withAvailableIpsCountsFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubnetWithAvailableIpsCountsFilter(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSubnetCheck("data.outscale_subnet.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceSubnetCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		subnetRs, ok := s.RootModule().Resources["outscale_subnet.outscale_subnet"]
		if !ok {
			return fmt.Errorf("can't find outscale_subnet.outscale_subnet in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		if attr["ip_range"] != "10.0.0.0/16" {
			return fmt.Errorf("bad ip_range %s", attr["ip_range"])
		}
		if attr["subregion_name"] != "eu-west-2a" {
			return fmt.Errorf("bad subregion_name %s", attr["subregion_name"])
		}

		return nil
	}
}

const testAccDataSourceSubnetConfig = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-subet-ds"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id        = "${outscale_net.outscale_net.net_id}"
		ip_range      = "10.0.0.0/16"
		subregion_name = "eu-west-2a"
	}

	data "outscale_subnet" "by_id" {
		subnet_id = "${outscale_subnet.outscale_subnet.id}"
	}

	data "outscale_subnet" "by_filter" {
		filter {
			name   = "subnet_ids"
			values = ["${outscale_subnet.outscale_subnet.id}"]
		}
	}
`

func testAccDataSourceSubnetWithAvailableIpsCountsFilter() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "Name"
				value = "Net1"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "eu-west-2a"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}

		data "outscale_subnet" "by_filter" {
			filter {
				name   = "available_ips_counts"
				values = ["${outscale_subnet.outscale_subnet.available_ips_count}"]
			}
		}
	`
}
