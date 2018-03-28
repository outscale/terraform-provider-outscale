package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleRoute_basic(t *testing.T) {
	var route fcu.Route

	//aws creates a default route
	testCheck := func(s *terraform.State) error {
		if *route.DestinationCidrBlock != "10.3.0.0/16" {
			return fmt.Errorf("Destination Cidr (Expected=%s, Actual=%s)\n", "10.3.0.0/16", *route.DestinationCidrBlock)
		}

		name := "outscale_lin_internet_gateway.foo"
		gwres, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s\n", name)
		}

		if *route.GatewayId != gwres.Primary.ID {
			return fmt.Errorf("Internet Gateway Id (Expected=%s, Actual=%s)\n", gwres.Primary.ID, *route.GatewayId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.bar", &route),
					testCheck,
				),
			},
		},
	})
}

func TestAccOutscaleRoute_ipv6Support(t *testing.T) {
	var route fcu.Route

	//aws creates a default route
	testCheck := func(s *terraform.State) error {

		name := "aws_egress_only_internet_gateway.foo"
		gwres, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s\n", name)
		}

		if *route.EgressOnlyInternetGatewayId != gwres.Primary.ID {
			return fmt.Errorf("Egress Only Internet Gateway Id (Expected=%s, Actual=%s)\n", gwres.Primary.ID, *route.EgressOnlyInternetGatewayId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.bar", &route),
					testCheck,
				),
			},
		},
	})
}

func TestAccOutscaleRoute_changeCidr(t *testing.T) {
	var route fcu.Route
	var routeTable fcu.RouteTable

	//aws creates a default route
	testCheck := func(s *terraform.State) error {
		if *route.DestinationCidrBlock != "10.3.0.0/16" {
			return fmt.Errorf("Destination Cidr (Expected=%s, Actual=%s)\n", "10.3.0.0/16", *route.DestinationCidrBlock)
		}

		name := "outscale_lin_internet_gateway.foo"
		gwres, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s\n", name)
		}

		if *route.GatewayId != gwres.Primary.ID {
			return fmt.Errorf("Internet Gateway Id (Expected=%s, Actual=%s)\n", gwres.Primary.ID, *route.GatewayId)
		}

		return nil
	}

	testCheckChange := func(s *terraform.State) error {
		if *route.DestinationCidrBlock != "10.2.0.0/16" {
			return fmt.Errorf("Destination Cidr (Expected=%s, Actual=%s)\n", "10.2.0.0/16", *route.DestinationCidrBlock)
		}

		name := "outscale_lin_internet_gateway.foo"
		gwres, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s\n", name)
		}

		if *route.GatewayId != gwres.Primary.ID {
			return fmt.Errorf("Internet Gateway Id (Expected=%s, Actual=%s)\n", gwres.Primary.ID, *route.GatewayId)
		}

		if rtlen := len(routeTable.Routes); rtlen != 2 {
			return fmt.Errorf("Route Table has too many routes (Expected=%d, Actual=%d)\n", rtlen, 2)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.bar", &route),
					testCheck,
				),
			},
			{
				Config: testAccOutscaleRouteBasicConfigChangeCidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.bar", &route),
					testAccCheckRouteTableExists("outscale_route_table.foo", &routeTable),
					testCheckChange,
				),
			},
		},
	})
}

func TestAccOutscaleRoute_noopdiff(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.test", &route),
					testCheck,
				),
			},
			{
				Config: testAccOutscaleRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleRouteExists("outscale_route.test", &route),
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
// 		CheckDestroy: testAccCheckOutscaleRouteDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccOutscaleRouteWithVPCEndpoint,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscaleRouteExists("outscale_route.bar", &route),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckOutscaleRouteExists(n string, res *fcu.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
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

func testAccCheckOutscaleRouteDestroy(s *terraform.State) error {
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

var testAccOutscaleRouteBasicConfig = fmt.Sprint(`
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route" "bar" {
	route_table_id = "${outscale_route_table.foo.id}"
	destination_cidr_block = "10.3.0.0/16"
	gateway_id = "${outscale_lin_internet_gateway.foo.id}"
}
`)

var testAccOutscaleRouteConfigIpv6 = fmt.Sprintf(`
resource "outscale_lin" "foo" {
  cidr_block = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true
}

resource "aws_egress_only_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route" "bar" {
	route_table_id = "${outscale_route_table.foo.id}"
	destination_ipv6_cidr_block = "::/0"
	egress_only_gateway_id = "${aws_egress_only_internet_gateway.foo.id}"
}


`)

var testAccOutscaleRouteBasicConfigChangeCidr = fmt.Sprint(`
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route" "bar" {
	route_table_id = "${outscale_route_table.foo.id}"
	destination_cidr_block = "10.2.0.0/16"
	gateway_id = "${outscale_lin_internet_gateway.foo.id}"
}
`)

// Acceptance test if mixed inline and external routes are implemented
var testAccOutscaleRouteMixConfig = fmt.Sprint(`
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"

	route {
		cidr_block = "10.2.0.0/16"
		gateway_id = "${outscale_lin_internet_gateway.foo.id}"
	}
}

resource "outscale_route" "bar" {
	route_table_id = "${outscale_route_table.foo.id}"
	destination_cidr_block = "0.0.0.0/0"
	gateway_id = "${outscale_lin_internet_gateway.foo.id}"
}
`)

var testAccOutscaleRouteNoopChange = fmt.Sprint(`
resource "outscale_lin" "test" {
  cidr_block = "10.10.0.0/16"
}

resource "outscale_route_table" "test" {
  vpc_id = "${outscale_lin.test.id}"
}

resource "outscale_subnet" "test" {
  vpc_id = "${outscale_lin.test.id}"
  cidr_block = "10.10.10.0/24"
}

resource "outscale_route" "test" {
  route_table_id = "${outscale_route_table.test.id}"
  destination_cidr_block = "0.0.0.0/0"
  instance_id = "${outscale_vm.nat.id}"
}

resource "outscale_vm" "nat" {
  ami = "ami-9abea4fb"
  instance_type = "t2.nano"
  subnet_id = "${outscale_subnet.test.id}"
}
`)

var testAccOutscaleRouteWithVPCEndpoint = fmt.Sprint(`
resource "outscale_lin" "foo" {
  cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
  vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
  vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route" "bar" {
  route_table_id         = "${outscale_route_table.foo.id}"
  destination_cidr_block = "10.3.0.0/16"
  gateway_id             = "${outscale_lin_internet_gateway.foo.id}"

  # Forcing endpoint to create before route - without this the crash is a race.
  depends_on = ["aws_vpc_endpoint.baz"]
}

resource "aws_vpc_endpoint" "baz" {
  vpc_id          = "${outscale_lin.foo.id}"
  service_name    = "com.amazonaws.us-west-2.s3"
  route_table_ids = ["${outscale_route_table.foo.id}"]
}
`)
