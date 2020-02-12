package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIVpnRoutePropagation_basic(t *testing.T) {
	t.Skip()
	rBgpAsn := acctest.RandIntRange(64512, 65534)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVpnRoutePropagationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVpnRoutePropagationConfig(rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleOAPIVpnRoutePropagation(
						"outscale_vpn_gateway_route_propagation.foo",
					),
				),
			},
		},
	})
}

func testAccCheckOAPIVpnRoutePropagationDestroy(s *terraform.State) error {
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
			return nil
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

func testAccOutscaleOAPIVpnRoutePropagation(routeProp string) resource.TestCheckFunc {
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

func testAccOutscaleOAPIVpnRoutePropagationConfig(rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			count    = 1
			ip_range = "10.0.0.0/16"
		}
		
		resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
			type = "ipsec.1"
    type = "ipsec.1" 
			type = "ipsec.1"
		}
		
		resource "outscale_vpn_gateway_link" "test" {
			lin_id         = "${outscale_net.outscale_net.id}"
			vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
		}
		
		resource "outscale_route_table" "outscale_route_table" {
			net_id = "${outscale_net.outscale_net.id}"
		}
		
		resource "outscale_vpn_gateway_route_propagation" "foo" {
			vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
			route_table_id = "${outscale_route_table.outscale_route_table.route_table_id}"
		}	
}
		}	
	`)
}
