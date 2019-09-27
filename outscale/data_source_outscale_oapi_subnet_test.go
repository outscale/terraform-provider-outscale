package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPISubnet(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISubnetCheck("data.outscale_subnet.by_id"),
					testAccDataSourceOutscaleOAPISubnetCheck("data.outscale_subnet.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		subnetRs, ok := s.RootModule().Resources["outscale_subnet.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_subnet.test in state")
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

const testAccDataSourceOutscaleOAPISubnetConfig = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_subnet" "test" {
		net_id        = "${outscale_net.outscale_net.net_id}"
		ip_range      = "10.0.0.0/16"
		subregion_name = "eu-west-2a"
	}
	
	data "outscale_subnet" "by_id" {
		subnet_id = "${outscale_subnet.test.id}"
	}

	data "outscale_subnet" "by_filter" {
		filter {
			name   = "subnet_ids"
			values = ["${outscale_subnet.test.id}"]
		}
	}
`
