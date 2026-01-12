package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_WithRouteTable_Basic(t *testing.T) {
	resourceName := "outscale_route_table.rtbTest"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.destination_ip_range", "10.1.0.0/16"),
				),
			},
		},
	})
}

func TestAccNet_RouteTable_Instance(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	resourceName := "outscale_route_table.rtbTest"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigInstance(sgName, omi, testAccVmType, utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.state", "active"),
				),
			},
		},
	})
}

func TestAccNet_WithRouteTable_tags(t *testing.T) {
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
	resourceName := "outscale_route_table.rtbTest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfigTags(value1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "Terraform-nic"),
				),
			},
			{
				Config: testAccOAPIRouteTableConfigTags(value2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNet_RouteTable_importBasic(t *testing.T) {
	resourceName := "outscale_route_table.rtbTest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIRouteTableConfig,
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_WithRouteTable_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps: testacc.FrameworkMigrationTestSteps("1.1.3",
			testAccOAPIRouteTableConfig,
			testAccOAPIRouteTableConfigInstance(sgName, omi, testAccVmType, utils.GetRegion())),
	})
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
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testacc.PreCheck(t) },
// 		Providers:    testacc.SDKProviders,
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
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:  func() { testacc.PreCheck(t) },
// 		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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

resource "outscale_route_table" "rtbTest" {
	net_id = outscale_net.foo.id
}
`

func testAccOAPIRouteTableConfigInstance(sgName, omi, vmType, region string) string {
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

		resource "outscale_security_group" "sg_route" {
			description           = "testAcc Terraform security group"
			security_group_name   = "%s"
			net_id                = outscale_net.foo.net_id

		}

		resource "outscale_vm" "foo" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			subnet_id                = outscale_subnet.foo.id
			placement_subregion_name = "%sa"
			placement_tenancy        = "default"
			security_group_ids = [outscale_security_group.sg_route.security_group_id]
		}

		resource "outscale_route_table" "rtbTest" {
			net_id = outscale_net.foo.id
		}
	`, sgName, omi, vmType, region)
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

resource "outscale_route_table" "rtbTest" {
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
