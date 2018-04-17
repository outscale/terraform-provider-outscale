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

func TestAccDataSourceOutscaleOAPISubnet(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 256)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISubnetCheck("data.outscale_subnet.by_id", rInt),
					testAccDataSourceOutscaleOAPISubnetCheck("data.outscale_subnet.by_filter", rInt),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetCheck(name string, rInt int) resource.TestCheckFunc {
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

		if attr["ip_range"] != fmt.Sprintf("172.%d.123.0/24", rInt) {
			return fmt.Errorf("bad ip_range %s", attr["ip_range"])
		}
		if attr["sub_region_name"] != "eu-west-2a" {
			return fmt.Errorf("bad sub_region_name %s", attr["sub_region_name"])
		}

		return nil
	}
}

func testAccDataSourceOutscaleOAPISubnetConfig(rInt int) string {
	return fmt.Sprintf(`
		
		resource "outscale_subnet" "test" {
		  lin_id            = "vpc-e9d09d63"
		  ip_range        = "172.%d.123.0/24"
		  sub_region_name = "eu-west-2a"
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
