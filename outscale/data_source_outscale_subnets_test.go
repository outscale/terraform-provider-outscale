package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccNet_WithSubnetsDataSource(t *testing.T) {
	t.Parallel()
	rInt := utils.RandIntRange(16, 31)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsConfig(rInt, utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_subnets.by_filter", "subnets.#", "1"),
				),
			},
		},
	})
}

func TestAccNet_Subnets_withAvailableIpsCountsFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsWithAvailableIpsCountsFilter(utils.GetRegion()),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetsConfig(rInt int, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.%[1]d.123.0/24"

			tags {
				key = "Name"
				value = "testacc-subets-ds"
			}
		}

		resource "outscale_subnet" "subnet" {
			ip_range       = "10.%[1]d.123.0/24"
			subregion_name = "%[2]sa"
			net_id         = "${outscale_net.net.id}"

			tags {
				key   = "name"
				value = "terraform-subnet"
			}
		}

		data "outscale_subnets" "by_filter" {
			filter {
				name   = "subnet_ids"
				values = ["${outscale_subnet.subnet.id}"]
			}
		}
	`, rInt, region)
}

func testAccDataSourceOutscaleOAPISubnetsWithAvailableIpsCountsFilter(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net1" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "Name"
				value = "Net1"
			}
		}

		resource "outscale_net" "outscale_net2" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "Name"
				value = "Net1"
			}
		}

		resource "outscale_subnet" "sub1" {
			subregion_name = "%[1]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net1.net_id
		}

		resource "outscale_subnet" "sub2" {
			subregion_name = "%[1]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net2.net_id
		}


		data "outscale_subnets" "by_filter" {
			filter {
				name   = "available_ips_counts"
				values = [outscale_subnet.sub1.available_ips_count, outscale_subnet.sub2.available_ips_count]
			}
		}
	`, region)
}
