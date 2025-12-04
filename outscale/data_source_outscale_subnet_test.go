package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_WithSubnet_DataSource(t *testing.T) {
	resouceName := "data.outscale_subnet.by_filter"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSubnetConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resouceName, "state", "available"),
					resource.TestCheckResourceAttrSet(resouceName, "ip_range"),
				),
			},
		},
	})
}

func TestAccNet_SubnetDataSource_withAvailableIpsCountsFilter(t *testing.T) {
	resouceName := "data.outscale_subnet.by_filter"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSubnetWithAvailableIpsCountsFilter(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resouceName, "net_id"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSubnetConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name"
				value = "testacc-subet-ds"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id        = outscale_net.outscale_net.net_id
			ip_range      = "10.0.0.0/24"
			subregion_name = "%sa"
		}

		data "outscale_subnet" "by_id" {
			subnet_id = outscale_subnet.outscale_subnet.id
		}

		data "outscale_subnet" "by_filter" {
			filter {
				name   = "subnet_ids"
				values = [outscale_subnet.outscale_subnet.id]
			}
		}

        `, region)
}

func testAccDataSourceOutscaleSubnetWithAvailableIpsCountsFilter(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "Name"
				value = "Net1"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/24"
			net_id         = outscale_net.outscale_net.net_id
		}

		data "outscale_subnet" "by_filter" {
			filter {
				name   = "available_ips_counts"
				values = [outscale_subnet.outscale_subnet.available_ips_count]
			}
			filter {
				name   = "ip_ranges"
				values = ["10.0.0.0/24"]
			}
			filter {
				name   = "net_ids"
				values = [outscale_net.outscale_net.net_id]
			}
		}
	`, region)
}
