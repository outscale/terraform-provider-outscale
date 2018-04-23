package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVpnRoutePropagation_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rBgpAsn := acctest.RandIntRange(64512, 65534)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnRoutePropagationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVpnRoutePropagationConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleVpnRoutePropagation(
						"outscale_vpn_gateway_route_propagation.foo",
					),
				),
			},
		},
	})
}

func testAccCheckVpnRoutePropagationDestroy(s *terraform.State) error {
	FCU := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_gateway_route_propagation" {
			continue
		}

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = FCU.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(rs.Primary.Attributes["gateway_id"])},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			return err
		}

		if len(resp.VpnGateways) > 0 {

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := FCU.VM.DeleteVpnGateway(&fcu.DeleteVpnGatewayInput{
					VpnGatewayId: resp.VpnGateways[0].VpnGatewayId,
				})
				if err == nil {
					return nil
				}

				ec2err, ok := err.(awserr.Error)
				if !ok {
					return resource.RetryableError(err)
				}

				switch ec2err.Code() {
				case "InvalidVpnGatewayID.NotFound":
					return nil
				case "IncorrectState":
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			})

			if err != nil {
				return fmt.Errorf("ERROR => %s", err)
			}

		} else {
			return nil
		}
	}

	return nil
}

func testAccOutscaleVpnRoutePropagation(routeProp string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[routeProp]
		if !ok {
			return fmt.Errorf("Not found: %s", routeProp)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccOutscaleVpnRoutePropagationConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
    type = "ipsec.1" 
}

resource "outscale_vpn_gateway_link" "test" {
	vpc_id = "${outscale_lin.outscale_lin.id}"
	vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_vpn_gateway_route_propagation" "foo" {
    gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
		route_table_id  = "${outscale_route_table.outscale_route_table.route_table_id}"
}
`)
}
