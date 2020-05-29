package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIVpnGatewayAttachment_basic(t *testing.T) {
	//var vpc oscgo.NetToVirtualGatewayLink
	//var vgw oscgo.VirtualGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		//IDRefreshName: "outscale_virtual_gateway_link.test",
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckOutscaleOAPILinExists(
				// 	"outscale_net.test",
				// 	&vpc), TODO: fix once we develop this resource
				//testAccCheckOAPIVirtualGatewayExists(
				//	"outscale_virtual_gateway.outscale_virtual_gateway",
				//	&vgw),
				//testAccCheckOAPIVpnGatewayAttachmentExists(
				//	"outscale_virtual_gateway_link.outscale_virtual_gateway",
				//	&vpc, &vgw),
				),
			},
		},
	})
}

func TestAccResourceVpnGatewayAttachment_importBasic(t *testing.T) {
	resourceName := "outscale_virtual_gateway_link.outscale_virtual_gateway_link"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckVpnGatewayAttachmentImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckVpnGatewayAttachmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func TestAccAWSOAPIVpnGatewayAttachment_deleted(t *testing.T) {

	var vpc oscgo.NetToVirtualGatewayLink
	var vgw oscgo.VirtualGateway

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
		},
		IDRefreshName: "outscale_virtual_gateway_link.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIVpnGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckOutscaleOAPILinExists(
					// 	"outscale_net.test",
					// 	&vpc),  TODO: Fix once we develop this resource
					testAccCheckOAPIVirtualGatewayExists(
						"outscale_virtual_gateway.test",
						&vgw),
					testAccCheckOAPIVpnGatewayAttachmentExists(
						"outscale_virtual_gateway_link.test",
						&vpc, &vgw),
				),
			},
			resource.TestStep{
				Config: testAccNoOAPIVpnGatewayAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("outscale_virtual_gateway_link.test"),
				),
			},
		},
	})
}

func testAccCheckOAPIVpnGatewayAttachmentExists(n string, vpc *oscgo.NetToVirtualGatewayLink, vgw *oscgo.VirtualGateway) resource.TestCheckFunc {
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

		if len(vgw.GetNetToVirtualGatewayLinks()) == 0 {
			return fmt.Errorf("vpn gateway %q has no attachments", vgwID)
		}

		if vgw.GetNetToVirtualGatewayLinks()[0].GetState() != "attached" {
			return fmt.Errorf("Expected VPN Gateway %q to be in attached state, but got: %q",
				vgwID, vgw.GetNetToVirtualGatewayLinks()[0].GetState())
		}

		if vgw.GetNetToVirtualGatewayLinks()[0].GetNetId() != vpc.GetNetId() {
			return fmt.Errorf("Expected VPN Gateway %q to be attached to VPC %q, but got: %q",
				vgwID, vpcID, vgw.GetNetToVirtualGatewayLinks()[0].GetNetId())
		}

		return nil
	}
}

func testAccCheckOAPIVpnGatewayAttachmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_virtual_gateway_link" {
			continue
		}

		vgwID := rs.Primary.Attributes["vpn_gateway_id"]

		var resp oscgo.ReadVirtualGatewaysResponse
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.VirtualGatewayApi.ReadVirtualGateways(context.Background(), &oscgo.ReadVirtualGatewaysOpts{ReadVirtualGatewaysRequest: optional.NewInterface(
				oscgo.ReadVirtualGatewaysRequest{
					Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{vgwID}},
				})})
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

		if len(resp.GetVirtualGateways()) > 1 {
			vgw := resp.GetVirtualGateways()[0]
			if vgw.GetNetToVirtualGatewayLinks()[0].GetState() != "detached" {
				return fmt.Errorf("Expected VPN Gateway %q to be in detached state, but got: %q",
					vgwID, vgw.GetNetToVirtualGatewayLinks()[0].GetState())
			}
		}
	}

	return nil
}

const testAccNoOAPIVpnGatewayAttachmentConfig = `
	resource "outscale_net" "test" {
		cidr_block = "10.0.0.0/16"
	}

	resource "outscale_virtual_gateway" "test" {
		connection_type = "ipsec.1"
}
`

const testAccOAPIVpnGatewayAttachmentConfig = `
resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
 connection_type = "ipsec.1"
}
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/18"
}
resource "outscale_virtual_gateway_link" "outscale_virtual_gateway_link" {
    virtual_gateway_id = outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id
    net_id              = outscale_net.outscale_net.net_id
}
`
