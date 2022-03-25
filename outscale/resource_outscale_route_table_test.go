package outscale

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIRouteTable_basic(t *testing.T) {
	var v oscgo.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.GetRoutes()) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]oscgo.Route)
		for _, r := range v.GetRoutes() {
			routes[r.GetDestinationIpRange()] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}
		return nil
	}

	testCheckChange := func(*terraform.State) error {
		if len(v.GetRoutes()) != 1 {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}

		routes := make(map[string]oscgo.Route)
		for _, r := range v.GetRoutes() {
			routes[r.GetDestinationIpRange()] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.Routes)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &v, nil),
					testCheck,
				),
			},

			{
				Config: testAccOAPIRouteTableConfigChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &v, nil),
					testCheckChange,
				),
			},
		},
	})
}

func TestAccOutscaleOAPIRouteTable_instance(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.RouteTable

	testCheck := func(*terraform.State) error {
		if len(v.GetRoutes()) != 1 {
			return fmt.Errorf("bad routes: %#v", v.GetRoutes())
		}

		routes := make(map[string]oscgo.Route)
		for _, r := range v.GetRoutes() {
			routes[r.GetDestinationIpRange()] = r
		}

		if _, ok := routes["10.1.0.0/16"]; !ok {
			return fmt.Errorf("bad routes: %#v", v.GetRoutes())
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigInstance(omi, "tinav4.c2r2p2", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists(
						"outscale_route_table.foo", &v, nil),
					testCheck,
				),
			},
		},
	})
}

func TestAccOutscaleOAPIRouteTable_tags(t *testing.T) {
	value1 := `
	tags {
		key = "name"
		value = "Terraform-nic"
	}`

	value2 := `
	tags{
		key = "name"
		value = "Terraform-RT"
	}
	tags{
		key = "name2"
		value = "Terraform-RT2"
	}`

	var rt oscgo.RouteTable
	rtTags := make([]oscgo.ResourceTag, 0)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigTags(value1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &rt, &rtTags),

					testAccCheckOAPITags(&rtTags, "name", "Terraform-nic"),
				),
			},
			{
				Config: testAccOAPIRouteTableConfigTags(value2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIRouteTableExists("outscale_route_table.foo", &rt, &rtTags),
					testAccCheckOAPITags(&rtTags, "name", "Terraform-RT"),
					testAccCheckOAPITags(&rtTags, "name2", "Terraform-RT2"),
				),
			},
		},
	})
}

func TestAccOutscaleRouteTable_importBasic(t *testing.T) {
	resourceName := "outscale_route_table.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfig,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleRouteTableImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleRouteTableImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccCheckOAPIRouteTableDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table" {
			continue
		}

		var resp oscgo.ReadRouteTablesResponse
		var err error
		params := oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{
				RouteTableIds: &[]string{rs.Primary.ID},
			},
		}

		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
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
			if len(resp.GetRouteTables()) > 0 {
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

func testAccCheckOAPIRouteTableExists(n string, v *oscgo.RouteTable, t *[]oscgo.ResourceTag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		var resp oscgo.ReadRouteTablesResponse
		var err error
		params := oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{
				RouteTableIds: &[]string{rs.Primary.ID},
			},
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(params).Execute()
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
		if len(resp.GetRouteTables()) == 0 {
			return fmt.Errorf("RouteTable not found")
		}

		*v = resp.GetRouteTables()[0]

		if t != nil {
			*t = resp.GetRouteTables()[0].GetTags()
			log.Printf("[DEBUG] Route Table Tags= %+v", t)
		}

		log.Printf("[DEBUG] RouteTable in Exist %+v", resp.GetRouteTables())

		return nil
	}
}

// VPC Peering connections are prefixed with pcx
// Right now there is no VPC Peering resource
// func TestAccOutscaleRouteTable_vpcPeering(t *testing.T) {
// 	var v oscgo.RouteTable

// 	testCheck := func(*terraform.State) error {
// 		if len(v.Routes) != 2 {
// 			return fmt.Errorf("bad routes: %#v", v.Routes)
// 		}

// 		routes := make(map[string]oscgo.Route)
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
// 	var v oscgo.RouteTable
// 	var vgw oscgo.VpnGateway

// 	testCheck := func(*terraform.State) error {
// 		if len(v.PropagatingVgws) != 1 {
// 			return fmt.Errorf("bad propagating vgws: %#v", v.PropagatingVgws)
// 		}

// 		propagatingVGWs := make(map[string]*oscgo.PropagatingVgw)
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

	tags {
		key = "Name"
		value = "testacc-route-table-rs"
	}
}

resource "outscale_internet_service" "foo" {}

resource "outscale_route_table" "foo" {
	net_id = outscale_net.foo.id
}
`

const testAccOAPIRouteTableConfigChange = `
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"

	tags {
		key = "Name"
		value = "testacc-route-table-rs"
	}
}

resource "outscale_internet_service" "foo" {}

resource "outscale_route_table" "foo" {
	net_id = outscale_net.foo.id
}
`

func testAccOAPIRouteTableConfigInstance(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "foo" {
			ip_range = "10.1.0.0/16"

			tags {
				key = "Name"
				value = "testacc-route-table-rs"
			}
		}

		resource "outscale_subnet" "foo" {
			ip_range = "10.1.1.0/24"
			net_id   = outscale_net.foo.id
		}

		resource "outscale_vm" "foo" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			subnet_id                = outscale_subnet.foo.id
			placement_subregion_name = "%sa"
			placement_tenancy        = "default"
		}

		resource "outscale_route_table" "foo" {
			net_id = outscale_net.foo.id
		}
	`, omi, vmType, region)
}

func testAccOAPIRouteTableConfigTags(value string) string {
	return fmt.Sprintf(`
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"

	tags {
		key = "Name"
		value = "testacc-route-table-rs"
	}
}

resource "outscale_route_table" "foo" {
	net_id = outscale_net.foo.id

	%s

}
`, value)
}

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
