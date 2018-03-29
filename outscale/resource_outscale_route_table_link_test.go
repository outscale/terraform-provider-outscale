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

func TestAccOutscaleRouteTableAssociation_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var v, v2 fcu.RouteTable

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteTableAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(
						"outscale_route_table_link.foo", &v),
				),
			},

			resource.TestStep{
				Config: testAccRouteTableAssociationConfigChange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(
						"outscale_route_table_link.foo", &v2),
				),
			},
		},
	})
}

func testAccCheckRouteTableAssociationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_route_table_link" {
			continue
		}

		var resp *fcu.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
				RouteTableIds: []*string{aws.String(rs.Primary.Attributes["route_table_id"])},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
					log.Printf("[DEBUG] Trying to create route again: %q", err)
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
				return err
			}
			return nil
		}

		rt := resp.RouteTables[0]
		if len(rt.Associations) > 0 {
			return fmt.Errorf(
				"route table %s has associations", *rt.RouteTableId)

		}
	}

	return nil
}

func testAccCheckRouteTableAssociationExists(n string, v *fcu.RouteTable) resource.TestCheckFunc {
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
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
				RouteTableIds: []*string{aws.String(rs.Primary.Attributes["route_table_id"])},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
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

		if len(v.Associations) == 0 {
			return fmt.Errorf("no associations")
		}

		return nil
	}
}

const testAccRouteTableAssociationConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	cidr_block = "10.1.1.0/24"
}

resource "outscale_lin_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table_link" "foo" {
	route_table_id = "${outscale_route_table.foo.id}"
	subnet_id = "${outscale_subnet.foo.id}"
}
`

const testAccRouteTableAssociationConfigChange = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	cidr_block = "10.1.1.0/24"
}

resource "outscale_lin_internet_gateway" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "bar" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table_link" "foo" {
	route_table_id = "${outscale_route_table.bar.id}"
	subnet_id = "${outscale_subnet.foo.id}"
}
`
