package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIRoute_noopdiff(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var route fcu.Route
	var routeTable fcu.RouteTable

	testCheck := func(s *terraform.State) error {
		return nil
	}

	testCheckChange := func(s *terraform.State) error {
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
					testCheck,
				),
			},
			{
				Config: testAccOutscaleOAPIRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIRouteExists("outscale_route.test", &route),
					testAccCheckRouteTableExists("outscale_route_table.test", &routeTable),
					testCheckChange,
				),
			},
		},
	})
}

// func TestAccOutscaleRoute_doesNotCrashWithVPCEndpoint(t *testing.T) {
// 	var route fcu.Route

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckOAPIOutscaleRouteDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccOutscaleRouteWithVPCEndpoint,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscaleOAPIRouteExists("outscale_route.bar", &route),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckOutscaleOAPIRouteExists(n string, res *fcu.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		r, err := findResourceRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_cidr_block"],
		)

		if err != nil {
			return err
		}

		if r == nil {
			return fmt.Errorf("Route not found")
		}

		*res = *r

		return nil
	}
}

func testAccCheckOAPIOutscaleRouteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		route, err := findResourceRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_cidr_block"],
		)

		if route == nil && err == nil {
			return nil
		}
	}

	return nil
}

var testAccOutscaleOAPIRouteNoopChange = fmt.Sprint(`
resource "outscale_lin" "test" {
  ip_range = "10.10.0.0/16"
}

resource "outscale_route_table" "test" {
  lin_id = "${outscale_lin.test.id}"
}

resource "outscale_subnet" "test" {
  lin_id = "${outscale_lin.test.id}"
  ip_range = "10.10.10.0/24"
}

resource "outscale_route" "test" {
  route_table_id = "${outscale_route_table.test.id}"
  destination_ip_range = "0.0.0.0/0"
  vm_id = "${outscale_vm.nat.id}"
}

resource "outscale_vm" "nat" {
	image_id = "ami-8a6a0120"
	type = "t2.micro"
  subnet_id = "${outscale_subnet.test.id}"
}
`)

// TODO: missing resource vpc_endpoint to make this test
// var testAccOutscaleRouteWithVPCEndpoint = fmt.Sprint(`
// resource "outscale_lin" "foo" {
//   ip_range = "10.1.0.0/16"
// }

// resource "outscale_lin_internet_gateway" "foo" {
//   lin_id = "${outscale_lin.foo.id}"
// }

// resource "outscale_route_table" "foo" {
//   lin_id = "${outscale_lin.foo.id}"
// }

// resource "outscale_route" "bar" {
//   route_table_id         = "${outscale_route_table.foo.id}"
//   destination_cidr_block = "10.3.0.0/16"
//   gateway_id             = "${outscale_lin_internet_gateway.foo.id}"

//   # Forcing endpoint to create before route - without this the crash is a race.
//   depends_on = ["aws_vpc_endpoint.baz"]
// }

// resource "aws_vpc_endpoint" "baz" {
//   lin_id          = "${outscale_lin.foo.id}"
//   service_name    = "com.amazonaws.us-west-2.s3"
//   route_table_ids = ["${outscale_route_table.foo.id}"]
// }
// `)
