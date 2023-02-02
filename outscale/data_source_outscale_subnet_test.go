package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Subnets_DataSource(t *testing.T) {
	t.Parallel()
	region := fmt.Sprintf("%sa", os.Getenv("OUTSCALE_REGION"))
	dataSourcesName := "data.outscale_subnets.subnets"
	dataSourceName := "data.outscale_subnet.subnet"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Subnet_DataSource_Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ip_range"),

					resource.TestCheckResourceAttr(dataSourcesName, "subnets.#", "2"),
				),
			},
		},
	})
}

func testAcc_Subnet_DataSource_Config(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "Name"
				value = "Net1"
			}
		}

		resource "outscale_subnet" "sub1" {
			subregion_name = "%[1]s"
			ip_range       = "10.0.1.0/24"
			net_id         = outscale_net.net.net_id
		}

		resource "outscale_subnet" "sub2" {
			subregion_name = "%[1]s"
			ip_range       = "10.0.2.0/24"
			net_id         = outscale_net.net.net_id
		}

		data "outscale_subnet" "subnet" {
			filter {
				name   = "subnet_ids"
				values = ["${outscale_subnet.sub1.id}"]
			}
		}

		data "outscale_subnets" "subnets" {
			filter {
				name   = "net_ids"
				values = ["${outscale_subnet.sub1.net_id}", "${outscale_subnet.sub2.net_id}"]
			}
		}
	`, region)
}
