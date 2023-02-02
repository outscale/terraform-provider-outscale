package outscale

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAcc_Net_DataSource(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UTC().UnixNano())
	ipRange := utils.RandVpcCidr()
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%s", ipRange)

	dataSourceName := "data.outscale_net.net"
	dataSourcesName := "data.outscale_nets.nets"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Net_DataSource_Config(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ip_range", ipRange),

					resource.TestCheckResourceAttr(dataSourcesName, "nets.#", "1"),
				),
			},
		},
	})
}

func testAcc_Net_DataSource_Config(ipRange, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "%s"
		
			tags {
				key   = "Name"
				value = "%s"
			}
		}
		
		data "outscale_net" "net" {
			filter {
				name   = "net_ids"
				values = ["${outscale_net.net.id}"]
			}
		}

		data "outscale_nets" "nets" {
            filter {
                name = "net_ids"
                values = [outscale_net.net.id]
            }
        }
	`, ipRange, tag)
}
