package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccDataSourceOutscaleOAPISubnets(t *testing.T) {
	t.Parallel()
	rInt := utils.RandIntRange(16, 31)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_subnets.by_filter", "subnets.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceOutscaleOAPISubnets_withAvailableIpsCountsFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsWithAvailableIpsCountsFilter(),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetsConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "172.%[1]d.123.0/24"

			tags {
				key = "Name"
				value = "testacc-subets-ds"
			}
		}

		resource "outscale_subnet" "subnet" {
			ip_range       = "172.%[1]d.123.0/24"
			subregion_name = "eu-west-2b"
			net_id         = outscale_net.net.id

			tags {
				key   = "name"
				value = "terraform-subnet"
			}
		}

		data "outscale_subnets" "by_filter" {
			filter {
				name   = "subnet_ids"
				values = [outscale_subnet.subnet.id]
			}
		}
	`, rInt)
}

func testAccDataSourceOutscaleOAPISubnetsWithAvailableIpsCountsFilter() string {
	return `
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
			subregion_name = "eu-west-2a"
			ip_range       = "10.0.0.0/24"
			net_id         = outscale_net.outscale_net1.net_id
		}

		resource "outscale_subnet" "sub2" {
			subregion_name = "eu-west-2b"
			ip_range       = "10.0.0.0/24"
			net_id         = outscale_net.outscale_net2.net_id
		}


		data "outscale_subnets" "by_filter" {
			filter {
				name   = "available_ips_counts"
				values = [outscale_subnet.sub1.available_ips_count, outscale_subnet.sub2.available_ips_count]
			}
		}
	`
}
