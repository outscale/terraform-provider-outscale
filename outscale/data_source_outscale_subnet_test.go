package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleSubnet(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 256)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSubnetConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleSubnetCheck("data.outscale_subnet.by_id", rInt),
					testAccDataSourceOutscaleSubnetCheck("data.outscale_subnet.by_filter", rInt),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSubnetCheck(name string, rInt int) resource.TestCheckFunc {
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

		if attr["cidr_block"] != fmt.Sprintf("10.0.%d.0/24", rInt) {
			return fmt.Errorf("bad cidr_block %s", attr["cidr_block"])
		}
		if attr["availability_zone"] != "eu-west-2a" {
			return fmt.Errorf("bad availability_zone %s", attr["availability_zone"])
		}

		return nil
	}
}

func testAccDataSourceOutscaleSubnetConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_lin" "vpc" {
			cidr_block = "10.0.0.0/16"
		}
		resource "outscale_subnet" "test" {
		  vpc_id            = "${outscale_lin.vpc.id}"
		  cidr_block        = "10.0.%d.0/24"
		  availability_zone = "eu-west-2a"
		}
		data "outscale_subnet" "by_id" {
		  subnet_id = "${outscale_subnet.test.id}"
		}
	
		data "outscale_subnet" "by_filter" {
		  filter {
		    name = "subnet-id"
		    values = ["${outscale_subnet.test.id}"]
		  }
		}
		`, rInt)
}
