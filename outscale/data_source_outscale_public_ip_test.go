package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PublicIP_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_public_ip.by_public_ip"
	dataSourcesName := "data.outscale_public_ips.by_public_ips"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_PublicIP_DataSource_Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "public_ips.#", "2"),

					resource.TestCheckResourceAttrSet(dataSourceName, "public_ip"),
				),
			},
		},
	})
}

const testAcc_PublicIP_DataSource_Config = `
	resource "outscale_public_ip" "test1" {}
	resource "outscale_public_ip" "test2" {}

	data "outscale_public_ip" "by_public_ip" {
		filter {
			name = "public_ips"
			values = [outscale_public_ip.test1.public_ip]
		}
	}

	data "outscale_public_ips" "by_public_ips" {
		filter {
			name  = "public_ips"
			values = [outscale_public_ip.test1.public_ip, outscale_public_ip.test2.public_ip]
		}
	}
`
