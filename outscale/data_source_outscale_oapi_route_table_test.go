package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPIRouteTable_basic(t *testing.T) {
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
				Config: testAccDataSourceOutscaleOAPIRouteTableGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIRouteTableCheck("data.outscale_route_table.by_filter"),
					testAccDataSourceOutscaleOAPIRouteTableCheck("data.outscale_route_table.by_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceOutscaleOAPIRouteTable_main(t *testing.T) {
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
				Config: testAccDataSourceOutscaleOAPIRouteTableMainRoute,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPIRouteTableCheckMain("data.outscale_route_table.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIRouteTableCheck(name string) resource.TestCheckFunc {
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

		if attr["lin_id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"lin_id is %s; want %s",
				attr["lin_id"],
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

func testAccDataSourceOutscaleOAPIRouteTableCheckMain(name string) resource.TestCheckFunc {
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
		if _, ok := attr["lin_id"]; !ok {
			return fmt.Errorf("lin_id not set for main route table")
		}
		// Verify it's actually the main route table that's returned
		if attr["association_set.0.main"] != "true" {
			return fmt.Errorf("main route table not found")
		}

		return nil
	}
}

const testAccDataSourceOutscaleOAPIRouteTableGroupConfig = `
resource "outscale_lin" "test" {
  ip_range = "172.16.0.0/16"

  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_subnet" "test" {
  ip_range = "172.16.0.0/24"
  lin_id     = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_route_table" "test" {
  lin_id = "${outscale_lin.test.id}"
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

const testAccDataSourceOutscaleOAPIRouteTableMainRoute = `

resource "outscale_lin" "test" {
  ip_range = "172.16.0.0/16"

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
