package outscale

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleLinPeeringConnection_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var connection fcu.VpcPeeringConnection

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_lin_peering.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcPeeringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinPeeringConnectionExists(
						"outscale_lin_peering.foo",
						&connection),
				),
			},
		},
	})
}

func TestAccOutscaleLinPeeringConnection_plan(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var connection fcu.VpcPeeringConnection

	// reach out and DELETE the VPC Peering connection outside of Terraform
	testDestroy := func(*terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		log.Printf("[DEBUG] Test deleting the VPC Peering Connection.")
		_, err := conn.VM.DeleteVpcPeeringConnection(
			&fcu.DeleteVpcPeeringConnectionInput{
				VpcPeeringConnectionId: connection.VpcPeeringConnectionId,
			})
		if err != nil {
			return err
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcPeeringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinPeeringConnectionExists(
						"outscale_lin_peering.foo",
						&connection),
					testDestroy,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckOutscaleLinPeeringConnectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_lin_peering" {
			continue
		}

		var describe *fcu.DescribeVpcPeeringConnectionsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.VM.DescribeVpcPeeringConnections(
				&fcu.DescribeVpcPeeringConnectionsInput{
					VpcPeeringConnectionIds: []*string{aws.String(rs.Primary.ID)},
				})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		var pc *fcu.VpcPeeringConnection
		for _, c := range describe.VpcPeeringConnections {
			if rs.Primary.ID == *c.VpcPeeringConnectionId {
				pc = c
			}
		}

		if pc == nil {
			// not found
			return nil
		}

		if pc.Status != nil {
			if *pc.Status.Code == "deleted" {
				return nil
			}
			return fmt.Errorf("Found the VPC Peering Connection in an unexpected state: %s", pc)
		}

		return fmt.Errorf("fall through error for testAccCheckOutscaleLinPeeringConnectionDestroy")
	}

	return nil
}

func testAccCheckOutscaleLinPeeringConnectionExists(n string, connection *fcu.VpcPeeringConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Peering Connection ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeVpcPeeringConnectionsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpcPeeringConnections(
				&fcu.DescribeVpcPeeringConnectionsInput{
					VpcPeeringConnectionIds: []*string{aws.String(rs.Primary.ID)},
				})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if len(resp.VpcPeeringConnections) == 0 {
			return fmt.Errorf("VPC Peering Connection could not be found")
		}

		*connection = *resp.VpcPeeringConnections[0]

		return nil
	}
}

func testAccCheckOutscaleLinPeeringConnectionOptions(n, block string, options *fcu.VpcPeeringConnectionOptionsDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Peering Connection ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeVpcPeeringConnectionsOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpcPeeringConnections(
				&fcu.DescribeVpcPeeringConnectionsInput{
					VpcPeeringConnectionIds: []*string{aws.String(rs.Primary.ID)},
				})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		pc := resp.VpcPeeringConnections[0]

		o := pc.AccepterVpcInfo
		if block == "requester_vpc_info" {
			o = pc.RequesterVpcInfo
		}

		if !reflect.DeepEqual(o.PeeringOptions, options) {
			return fmt.Errorf("Expected the VPC Peering Connection Options to be %#v, got %#v",
				options, o.PeeringOptions)
		}

		return nil
	}
}

const testAccVpcPeeringConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_basic"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
}
`

const testAccVpcPeeringConfigTags = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_tags"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
	tag {
		foo = "bar"
	}
}
`

const testAccVpcPeeringConfigOptions = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_options"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.1.0.0/16"
	enable_dns_hostnames = true
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
}
`

const testAccVpcPeeringConfigFailedState = `
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_failedState"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
}
`
