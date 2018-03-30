package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleRouteTable_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var v fcu.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.Routes) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]*fcu.Route)
		for _, r := range v.Routes {
			routes[*r.DestinationCidrBlock] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}
		// if _, ok := routes["10.2.0.0/16"]; !ok {
		// 	return fmt.Errorf("bad routes: %#v", v.Routes)
		// }

		return nil
	}

	testCheckChange := func(*terraform.State) error {
		if len(v.Routes) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]*fcu.Route)
		for _, r := range v.Routes {
			routes[*r.DestinationCidrBlock] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}
		// if _, ok := routes["10.3.0.0/16"]; !ok {
		// 	return fmt.Errorf("bad routes: %#v", v.Routes)
		// }
		// if _, ok := routes["10.4.0.0/16"]; !ok {
		// 	return fmt.Errorf("bad routes: %#v", v.Routes)
		// }

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(
						"outscale_route_table.foo", &v),
					testCheck,
				),
			},

			{
				Config: testAccRouteTableConfigChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(
						"outscale_route_table.foo", &v),
					testCheckChange,
				),
			},
		},
	})
}

func TestAccOutscaleRouteTable_instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var v fcu.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.Routes) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]*fcu.Route)
		for _, r := range v.Routes {
			routes[*r.DestinationCidrBlock] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}
		// if _, ok := routes["10.2.0.0/16"]; !ok {
		// 	return fmt.Errorf("bad routes: %#v", v.Routes)
		// }

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(
						"outscale_route_table.foo", &v),
					testCheck,
				),
			},
		},
	})
}

func TestAccOutscaleRouteTable_tags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var route_table fcu.RouteTable

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists("outscale_route_table.foo", &route_table),
					testAccCheckTags(&route_table.Tags, "foo", "bar"),
				),
			},
		},
	})
}

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table" {
			continue
		}

		var resp *fcu.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
				RouteTableIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") || strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
					log.Printf("[DEBUG] Trying to create route again: %q", err)
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			if len(resp.RouteTables) > 0 {
				return fmt.Errorf("still exist.")
			}

			return nil
		}

		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			return nil
		}
	}

	return nil
}

func testAccCheckRouteTableExists(n string, v *fcu.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
				RouteTableIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					log.Printf("[DEBUG] Trying to create route again: %q", err)
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}
		if len(resp.RouteTables) == 0 {
			return fmt.Errorf("RouteTable not found")
		}

		*v = *resp.RouteTables[0]

		return nil
	}
}

// VPC Peering connections are prefixed with pcx
// Right now there is no VPC Peering resource
// func TestAccOutscaleRouteTable_vpcPeering(t *testing.T) {
// 	var v fcu.RouteTable

// 	testCheck := func(*terraform.State) error {
// 		if len(v.Routes) != 2 {
// 			return fmt.Errorf("bad routes: %#v", v.Routes)
// 		}

// 		routes := make(map[string]*fcu.Route)
// 		for _, r := range v.Routes {
// 			routes[*r.DestinationCidrBlock] = r
// 		}

// 		if _, ok := routes["10.1.0.0/16"]; !ok {
// 			return fmt.Errorf("bad routes: %#v", v.Routes)
// 		}
// 		if _, ok := routes["10.2.0.0/16"]; !ok {
// 			return fmt.Errorf("bad routes: %#v", v.Routes)
// 		}

// 		return nil
// 	}
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckRouteTableDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccRouteTableVpcPeeringConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckRouteTableExists(
// 						"outscale_route_table.foo", &v),
// 					testCheck,
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccOutscaleRouteTable_vgwRoutePropagation(t *testing.T) {
// 	var v fcu.RouteTable
// 	var vgw fcu.VpnGateway

// 	testCheck := func(*terraform.State) error {
// 		if len(v.PropagatingVgws) != 1 {
// 			return fmt.Errorf("bad propagating vgws: %#v", v.PropagatingVgws)
// 		}

// 		propagatingVGWs := make(map[string]*fcu.PropagatingVgw)
// 		for _, gw := range v.PropagatingVgws {
// 			propagatingVGWs[*gw.GatewayId] = gw
// 		}

// 		if _, ok := propagatingVGWs[*vgw.VpnGatewayId]; !ok {
// 			return fmt.Errorf("bad propagating vgws: %#v", v.PropagatingVgws)
// 		}

// 		return nil

// 	}
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		CheckDestroy: resource.ComposeTestCheckFunc(
// 			testAccCheckVpnGatewayDestroy,
// 			testAccCheckRouteTableDestroy,
// 		),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccRouteTableVgwRoutePropagationConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckRouteTableExists(
// 						"outscale_route_table.foo", &v),
// 					testAccCheckVpnGatewayExists(
// 						"aws_vpn_gateway.foo", &vgw),
// 					testCheck,
// 				),
// 			},
// 		},
// 	})
// }

const testAccRouteTableConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
	#vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}
`

const testAccRouteTableConfigChange = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
	#vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}
`

const testAccRouteTableConfigInstance = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	cidr_block = "10.1.1.0/24"
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_vm" "foo" {
	# us-west-2
	ami = "ami-4fccb37f"
	instance_type = "m1.small"
	subnet_id = "${outscale_subnet.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}
`

const testAccRouteTableConfigTags = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"

	tag {
		foo = "bar"
	}
}
`

// TODO: missing resource vpc peering to make this test
// VPC Peering connections are prefixed with pcx
// const testAccRouteTableVpcPeeringConfig = `
// resource "outscale_lin" "foo" {
// 	cidr_block = "10.1.0.0/16"
// }

// resource "outscale_lin_internet_gateway" "foo" {
// 	vpc_id = "${outscale_lin.foo.id}"
// }

// resource "outscale_lin" "bar" {
// 	cidr_block = "10.3.0.0/16"
// }

// resource "outscale_lin_internet_gateway" "bar" {
// 	vpc_id = "${outscale_lin.bar.id}"
// }

// resource "aws_vpc_peering_connection" "foo" {
// 		vpc_id = "${outscale_lin.foo.id}"
// 		peer_vpc_id = "${outscale_lin.bar.id}"
// 		tags {
// 			foo = "bar"
// 		}
// }

// resource "outscale_route_table" "foo" {
// 	vpc_id = "${outscale_lin.foo.id}"

// 	route {
// 		cidr_block = "10.2.0.0/16"
// 		vpc_peering_connection_id = "${aws_vpc_peering_connection.foo.id}"
// 	}
// }
// `

// TODO: missing vpn_gateway to make this test
// const testAccRouteTableVgwRoutePropagationConfig = `
// resource "outscale_lin" "foo" {
// 	cidr_block = "10.1.0.0/16"
// }

// resource "aws_vpn_gateway" "foo" {
// 	vpc_id = "${outscale_lin.foo.id}"
// }

// resource "outscale_route_table" "foo" {
// 	vpc_id = "${outscale_lin.foo.id}"

// 	propagating_vgws = ["${aws_vpn_gateway.foo.id}"]
// }
// `
