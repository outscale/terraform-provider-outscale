package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIPublicIPS(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIPublicIPSConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_public_ips.by_public_ips", "addresses_set.0.domain", "standard"),
					resource.TestCheckResourceAttr(
						"data.outscale_public_ips.by_public_ips", "addresses_set.1.domain", "standard"),
					resource.TestCheckResourceAttr(
						"data.outscale_public_ips.by_public_ips", "addresses_set.2.domain", "standard"),
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
