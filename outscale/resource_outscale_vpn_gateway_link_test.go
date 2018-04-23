package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccAWSVpnGatewayAttachment_basic(t *testing.T) {
	var vpc fcu.Vpc
	var vgw fcu.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_gateway_link.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinExists(
						"outscale_lin.test",
						&vpc),
					testAccCheckVpnGatewayExists(
						"outscale_vpn_gateway.test",
						&vgw),
					testAccCheckVpnGatewayAttachmentExists(
						"outscale_vpn_gateway_link.test",
						&vpc, &vgw),
				),
			},
		},
	})
}

func TestAccAWSVpnGatewayAttachment_deleted(t *testing.T) {
	var vpc fcu.Vpc
	var vgw fcu.VpnGateway

	testDeleted := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[n]
			if ok {
				return fmt.Errorf("Expected VPN Gateway attachment resource %q to be deleted.", n)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_gateway_link.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinExists(
						"outscale_lin.test",
						&vpc),
					testAccCheckVpnGatewayExists(
						"outscale_vpn_gateway.test",
						&vgw),
					testAccCheckVpnGatewayAttachmentExists(
						"outscale_vpn_gateway_link.test",
						&vpc, &vgw),
				),
			},
			resource.TestStep{
				Config: testAccNoVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("outscale_vpn_gateway_link.test"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayAttachmentExists(n string, vpc *fcu.Vpc, vgw *fcu.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		vpcId := rs.Primary.Attributes["vpc_id"]
		vgwId := rs.Primary.Attributes["vpn_gateway_id"]

		if len(vgw.VpcAttachments) == 0 {
			return fmt.Errorf("VPN Gateway %q has no attachments.", vgwId)
		}

		if *vgw.VpcAttachments[0].State != "attached" {
			return fmt.Errorf("Expected VPN Gateway %q to be in attached state, but got: %q",
				vgwId, *vgw.VpcAttachments[0].State)
		}

		if *vgw.VpcAttachments[0].VpcId != *vpc.VpcId {
			return fmt.Errorf("Expected VPN Gateway %q to be attached to VPC %q, but got: %q",
				vgwId, vpcId, *vgw.VpcAttachments[0].VpcId)
		}

		return nil
	}
}

func testAccCheckVpnGatewayAttachmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_gateway_link" {
			continue
		}

		vgwId := rs.Primary.Attributes["vpn_gateway_id"]

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(vgwId)},
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
				vgwId, *vgw.VpcAttachments[0].State)
		}
	}

	return nil
}

const testAccNoVpnGatewayAttachmentConfig = `
resource "outscale_lin" "test" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_gateway" "test" { }
`

const testAccVpnGatewayAttachmentConfig = `
resource "outscale_lin" "test" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_gateway" "test" { }

resource "outscale_vpn_gateway_link" "test" {
	vpc_id = "${outscale_lin.test.id}"
	vpn_gateway_id = "${outscale_vpn_gateway.test.id}"
}
`
