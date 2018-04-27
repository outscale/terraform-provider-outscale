package outscale

import (
	"fmt"
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

func TestAccOutscaleNatService_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var natGateway fcu.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nat_service.outscale_nat_service",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("outscale_nat_service.outscale_nat_service", &natGateway),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nat_service" {
			continue
		}

		// Try to find the resource

		var resp *fcu.DescribeNatGatewaysOutput

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.DescribeNatGateways(&fcu.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err == nil {
			status := map[string]bool{
				"deleted":  true,
				"deleting": true,
				"failed":   true,
			}
			if _, ok := status[strings.ToLower(*resp.NatGateways[0].State)]; len(resp.NatGateways) > 0 && !ok {
				return fmt.Errorf("still exists")
			}

			return nil
		}
		if err != nil {
			if strings.Contains(err.Error(), "NatGatewayNotFound:") {
				return nil
			}

			return err
		}

	}

	return nil
}

func testAccCheckNatGatewayExists(n string, ng *fcu.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeNatGatewaysOutput
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error

			resp, err = conn.VM.DescribeNatGateways(&fcu.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{aws.String(rs.Primary.ID)},
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
		if len(resp.NatGateways) == 0 {
			return fmt.Errorf("NatGateway not found")
		}

		*ng = *resp.NatGateways[0]

		return nil
	}
}

const testAccNatGatewayConfig = `
resource "outscale_lin" "outscale_lin" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"

  cidr_block = "10.0.0.0/18"
}

resource "outscale_public_ip" "outscale_public_ip" {
  #domain               = "Standard"       # BUG doc API

  domain = ""
}

resource "outscale_nat_service" "outscale_nat_service" {
  depends_on = ["outscale_route.outscale_route"]

  subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"

  allocation_id = "${outscale_public_ip.outscale_public_ip.allocation_id}"
}

resource "outscale_route_table" "outscale_route_table" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_route" "outscale_route" {
  destination_cidr_block = "0.0.0.0/0"

  gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"

  route_table_id = "${outscale_route_table.outscale_route_table.id}"
}

resource "outscale_route_table_link" "outscale_route_table_link" {
  subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"

  route_table_id = "${outscale_route_table.outscale_route_table.id}"
}

resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}

resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"

  internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
}
`
