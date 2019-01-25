package outscale

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPIRouteTable_basic(t *testing.T) {
	var v oapi.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.Routes) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]oapi.Route)
		for _, r := range v.Routes {
			routes[r.DestinationIpRange] = r
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

		routes := make(map[string]oapi.Route)
		for _, r := range v.Routes {
			routes[r.DestinationIpRange] = r
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
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &v),
					testCheck,
				),
			},

			{
				Config: testAccOAPIRouteTableConfigChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &v),
					testCheckChange,
				),
			},
		},
	})
}

func TestAccOutscaleOAPIRouteTable_instance(t *testing.T) {
	var v oapi.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.Routes) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]oapi.Route)
		for _, r := range v.Routes {
			routes[r.DestinationIpRange] = r
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
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists(
						"outscale_route_table.foo", &v),
					testCheck,
				),
			},
		},
	})
}

func TestAccOutscaleOAPIRouteTable_tags(t *testing.T) {
	var rt oapi.RouteTable

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &rt),
					testAccCheckOAPITags(rt.Tags, "foo", "bar"),
				),
			},
		},
	})
}

func testAccCheckOAPIRouteTableDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table" {
			continue
		}

		var resp *oapi.POST_ReadRouteTablesResponses
		var err error
		params := &oapi.ReadRouteTablesRequest{
			Filters: oapi.FiltersRouteTable{
				RouteTableIds: []string{rs.Primary.ID},
			},
		}

		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadRouteTables(*params)
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
			if len(resp.OK.RouteTables) > 0 {
				return fmt.Errorf("still exist")
			}

			return nil
		}

		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			return nil
		}
	}

	return nil
}

func testAccCheckOAPIRouteTableExists(n string, v *oapi.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		var resp *oapi.POST_ReadRouteTablesResponses
		var err error
		params := &oapi.ReadRouteTablesRequest{
			Filters: oapi.FiltersRouteTable{
				RouteTableIds: []string{rs.Primary.ID},
			},
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadRouteTables(*params)
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
		if len(resp.OK.RouteTables) == 0 {
			return fmt.Errorf("RouteTable not found")
		}

		*v = resp.OK.RouteTables[0]

		return nil
	}
}

// VPC Peering connections are prefixed with pcx
// Right now there is no VPC Peering resource
// func TestAccOutscaleRouteTable_vpcPeering(t *testing.T) {
// 	var v oapi.RouteTable

// 	testCheck := func(*terraform.State) error {
// 		if len(v.Routes) != 2 {
// 			return fmt.Errorf("bad routes: %#v", v.Routes)
// 		}

// 		routes := make(map[string]oapi.Route)
// 		for _, r := range v.Routes {
// 			routes[r.DestinationIpRange] = r
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
// 		CheckDestroy: testAccCheckOAPIRouteTableDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccRouteTableVpcPeeringConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOAPIRouteTableExists(
// 						"outscale_route_table.foo", &v),
// 					testCheck,
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccOutscaleRouteTable_vgwRoutePropagation(t *testing.T) {
// 	var v oapi.RouteTable
// 	var vgw oapi.VpnGateway

// 	testCheck := func(*terraform.State) error {
// 		if len(v.PropagatingVgws) != 1 {
// 			return fmt.Errorf("bad propagating vgws: %#v", v.PropagatingVgws)
// 		}

// 		propagatingVGWs := make(map[string]*oapi.PropagatingVgw)
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
// 			testAccCheckOAPIRouteTableDestroy,
// 		),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccRouteTableVgwRoutePropagationConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOAPIRouteTableExists(
// 						"outscale_route_table.foo", &v),
// 					testAccCheckVpnGatewayExists(
// 						"aws_vpn_gateway.foo", &vgw),
// 					testCheck,
// 				),
// 			},
// 		},
// 	})
// }

const testAccOAPIRouteTableConfig = `
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"
}

resource "outscale_internet_service" "foo" {}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"
}
`

const testAccOAPIRouteTableConfigChange = `
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"
}

resource "outscale_internet_service" "foo" {}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"
}
`

const testAccOAPIRouteTableConfigInstance = `
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	ip_range = "10.1.1.0/24"
	net_id = "${outscale_net.foo.id}"
}

resource "outscale_vm" "foo" {
	# us-west-2
	image_id = "ami-b4bd8de2"
	vm_type = "m1.small"
	subnet_id = "${outscale_subnet.foo.id}"
}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"
}
`

const testAccOAPIRouteTableConfigTags = `
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"

	tag {
		foo = "bar"
	}
}
`

// TODO: missing resource vpc peering to make this test
// VPC Peering connections are prefixed with pcx
// const testAccRouteTableVpcPeeringConfig = `
// resource "outscale_net" "foo" {
// 	ip_range = "10.1.0.0/16"
// }

// resource "outscale_internet_service" "foo" {
// 	net_id = "${outscale_net.foo.id}"
// }

// resource "outscale_net" "bar" {
// 	ip_range = "10.3.0.0/16"
// }

// resource "outscale_internet_service" "bar" {
// 	net_id = "${outscale_net.bar.id}"
// }

// resource "aws_vpc_peering_connection" "foo" {
// 		net_id = "${outscale_net.foo.id}"
// 		peer_vpc_id = "${outscale_net.bar.id}"
// 		tags {
// 			foo = "bar"
// 		}
// }

// resource "outscale_route_table" "foo" {
// 	net_id = "${outscale_net.foo.id}"

// 	route {
// 		ip_range = "10.2.0.0/16"
// 		vpc_peering_connection_id = "${aws_vpc_peering_connection.foo.id}"
// 	}
// }
// `

// TODO: missing vpn_gateway to make this test
// const testAccRouteTableVgwRoutePropagationConfig = `
// resource "outscale_net" "foo" {
// 	ip_range = "10.1.0.0/16"
// }

// resource "aws_vpn_gateway" "foo" {
// 	net_id = "${outscale_net.foo.id}"
// }

// resource "outscale_route_table" "foo" {
// 	net_id = "${outscale_net.foo.id}"

// 	propagating_vgws = ["${aws_vpn_gateway.foo.id}"]
// }
// `
