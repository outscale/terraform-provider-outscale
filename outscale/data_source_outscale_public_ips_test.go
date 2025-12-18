package outscale

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_DataSourcePublicIPS(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscalePublicIPSConfig,
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

func TestAccOthers_DataSourcePublicIPS_withTags(t *testing.T) {
	if os.Getenv("TEST_QUOTA") == "true" {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceOutscalePublicIPSConfigWithTags,
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

const testAccDataSourceOutscalePublicIPSConfig = `
	resource "outscale_public_ip" "test" {}
	resource "outscale_public_ip" "test1" {}
	resource "outscale_public_ip" "test2" {}

	data "outscale_public_ips" "by_public_ips" {
		filter {
			name  = "public_ips"
			values = [outscale_public_ip.test.public_ip, outscale_public_ip.test1.public_ip, outscale_public_ip.test2.public_ip]
		}
	}
`

const testAccDataSourceOutscalePublicIPSConfigWithTags = `
	resource "outscale_public_ip" "outscale_public_ip" {
		tags {
			key   = "name"
			value = "public_ip-data"
		}
	}

	resource "outscale_public_ip" "outscale_public_ip2" {
		tags {
			key   = "name"
			value = "public_ip-data"
		}
	}

	data "outscale_public_ips" "outscale_public_ips" {
		filter {
			name   = "tags"
			values = ["name=public_ip-data"]
		}
         depends_on = [outscale_public_ip.outscale_public_ip, outscale_public_ip.outscale_public_ip2]
	}
`
