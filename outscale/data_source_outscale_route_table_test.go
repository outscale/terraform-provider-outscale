package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleRouteTable_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleRouteTableGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleRouteTableCheck("data.outscale_route_table.by_filter"),
					testAccDataSourceOutscaleRouteTableCheck("data.outscale_route_table.by_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceOutscaleRouteTable_main(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleRouteTableMainRoute,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleRouteTableCheckMain("data.outscale_route_table.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleRouteTableCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		rts, ok := s.RootModule().Resources["outscale_route_table.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_route_table.test in state")
		}
		vpcRs, ok := s.RootModule().Resources["outscale_lin.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_lin.test in state")
		}
		subnetRs, ok := s.RootModule().Resources["outscale_subnet.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_subnet.test in state")
		}
		attr := rs.Primary.Attributes

		if attr["id"] != rts.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				rts.Primary.Attributes["id"],
			)
		}

		if attr["route_table_id"] != rts.Primary.Attributes["id"] {
			return fmt.Errorf(
				"route_table_id is %s; want %s",
				attr["route_table_id"],
				rts.Primary.Attributes["id"],
			)
		}

		if attr["vpc_id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"vpc_id is %s; want %s",
				attr["vpc_id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["association_set.0.subnet_id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"subnet_id is %v; want %s",
				attr["association_set.0.subnet_id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		return nil
	}
}

func testAccDataSourceOutscaleRouteTableCheckMain(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		// Verify attributes are set
		if _, ok := attr["id"]; !ok {
			return fmt.Errorf("id not set for main route table")
		}
		if _, ok := attr["vpc_id"]; !ok {
			return fmt.Errorf("vpc_id not set for main route table")
		}
		// Verify it's actually the main route table that's returned
		if attr["association_set.0.main"] != "true" {
			return fmt.Errorf("main route table not found")
		}

		return nil
	}
}

const testAccDataSourceOutscaleRouteTableGroupConfig = `
resource "outscale_lin" "test" {
  cidr_block = "172.16.0.0/16"

  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_subnet" "test" {
  cidr_block = "172.16.0.0/24"
  vpc_id     = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_route_table" "test" {
  vpc_id = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-routetable-data-source"
  }
}

resource "outscale_route_table_link" "a" {
    subnet_id = "${outscale_subnet.test.id}"
    route_table_id = "${outscale_route_table.test.id}"
}

data "outscale_route_table" "by_filter" {
  filter {
    name = "association.route-table-association-id"
    values = ["${outscale_route_table_link.a.id}"]
  }
  depends_on = ["outscale_route_table_link.a"]
}

data "outscale_route_table" "by_id" {
  route_table_id = "${outscale_route_table.test.id}"
  depends_on = ["outscale_route_table_link.a"]
}
`

const testAccDataSourceOutscaleRouteTableMainRoute = `

resource "outscale_lin" "test" {
  cidr_block = "172.16.0.0/16"

  tag {
    Name = "terraform-testacc-data-source"
  }
}

data "outscale_route_table" "by_filter" {
  filter {
    name = "association.main"
    values = ["true"]
  }
  filter {
    name = "vpc-id"
    values = ["${outscale_lin.test.id}"]
  }
}
`
