package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIPublicIPS(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIPublicIPSConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.outscale_public_ips.by_public_ips", "public_ips.0.public_ip"),
					resource.TestCheckResourceAttrSet(
						"data.outscale_public_ips.by_public_ips", "public_ips.1.public_ip"),
					resource.TestCheckResourceAttrSet(
						"data.outscale_public_ips.by_public_ips", "public_ips.2.public_ip"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIPublicIPSConfig = `
	resource "outscale_public_ip" "test" {}
	resource "outscale_public_ip" "test1" {}
	resource "outscale_public_ip" "test2" {}

	data "outscale_public_ips" "by_public_ips" {
		filter {
			name  = "public_ip"
			values = ["${outscale_public_ip.test.public_ip}", "${outscale_public_ip.test1.public_ip}", "${outscale_public_ip.test2.public_ip}"]
		}  
	}
`
