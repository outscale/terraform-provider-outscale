package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIVpnGatewayAttachment_basic(t *testing.T) {
	t.Skip()

	var vpc fcu.Vpc
	var vgw fcu.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		IDRefreshName: "outscale_vpn_gateway_link.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckOutscaleOAPILinExists(
					// 	"outscale_net.test",
					// 	&vpc), TODO: fix once we develop this resource
					testAccCheckOAPIVpnGatewayExists(
						"outscale_vpn_gateway.test",
						&vgw),
					testAccCheckOAPIVpnGatewayAttachmentExists(
						"outscale_vpn_gateway_link.test",
						&vpc, &vgw),
				),
			},
		},
	})
}

func TestAccAWSOAPIVpnGatewayAttachment_deleted(t *testing.T) {
	t.Skip()

	var vpc fcu.Vpc
	var vgw fcu.VpnGateway

	testDeleted := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[n]
			if ok {
				return fmt.Errorf("expected vpn gateway attachment resource %q to be deleted", n)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		IDRefreshName: "outscale_vpn_gateway_link.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckOutscaleOAPILinExists(
					// 	"outscale_net.test",
					// 	&vpc),  TODO: Fix once we develop this resource
					testAccCheckOAPIVpnGatewayExists(
						"outscale_vpn_gateway.test",
						&vgw),
					testAccCheckOAPIVpnGatewayAttachmentExists(
						"outscale_vpn_gateway_link.test",
						&vpc, &vgw),
				),
			},
			resource.TestStep{
				Config: testAccNoOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("outscale_vpn_gateway_link.test"),
				),
			},
		},
	})
}

func testAccCheckOAPIVpnGatewayAttachmentExists(n string, vpc *fcu.Vpc, vgw *fcu.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		vpcID := rs.Primary.Attributes["lin_id"]
		vgwID := rs.Primary.Attributes["vpn_gateway_id"]

		if len(vgw.VpcAttachments) == 0 {
			return fmt.Errorf("vpn gateway %q has no attachments", vgwID)
		}

		if *vgw.VpcAttachments[0].State != "attached" {
			return fmt.Errorf("Expected VPN Gateway %q to be in attached state, but got: %q",
				vgwID, *vgw.VpcAttachments[0].State)
		}

		if *vgw.VpcAttachments[0].VpcId != *vpc.VpcId {
			return fmt.Errorf("Expected VPN Gateway %q to be attached to VPC %q, but got: %q",
				vgwID, vpcID, *vgw.VpcAttachments[0].VpcId)
		}

		return nil
	}
}

func testAccCheckOAPIVpnGatewayAttachmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_gateway_link" {
			continue
		}

		vgwID := rs.Primary.Attributes["vpn_gateway_id"]

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(vgwID)},
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

		vgw := resp.VpnGateways[0]
		if *vgw.VpcAttachments[0].State != "detached" {
			return fmt.Errorf("Expected VPN Gateway %q to be in detached state, but got: %q",
				vgwID, *vgw.VpcAttachments[0].State)
		}
	}

	return nil
}

const testAccNoOAPIVpnGatewayAttachmentConfig = `
	resource "outscale_net" "test" {
		cidr_block = "10.0.0.0/16"
	}

	resource "outscale_vpn_gateway" "test" {}
`

const testAccOAPIVpnGatewayAttachmentConfig = `
	resource "outscale_net" "test" {
		ip_range = "10.0.0.0/16"
	}

	resource "outscale_vpn_gateway" "test" { 
		type = "ipsec.1" 
	}

	resource "outscale_vpn_gateway_link" "test" {
		net_id         = "${outscale_net.test.id}"
		vpn_gateway_id = "${outscale_vpn_gateway.test.id}"
	}
`
