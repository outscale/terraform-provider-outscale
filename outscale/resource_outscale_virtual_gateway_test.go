package outscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNet_WithVirtualGateway_basic(t *testing.T) {
	t.Parallel()
	var v, v2 oscgo.VirtualGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_virtual_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIVirtualGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVirtualGatewayExists(
						"outscale_virtual_gateway.foo", &v),
				),
			},

			{
				Config: testAccOAPIVirtualGatewayConfigChangeVPC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVirtualGatewayExists(
						"outscale_virtual_gateway.foo", &v2),
				),
			},
		},
	})
}

func TestAccNet_WithVirtualGatewayChangeTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVirtualGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfigChangeTags("ipsec.1", "test-VGW"),
			},
			{
				Config: testAccOAPIVirtualGatewayConfigChangeTags("ipsec.1", "test-VGW2"),
			},
		},
	})
}

func TestAccNet_withVirtualGateway_delete(t *testing.T) {
	var virtualGateway oscgo.VirtualGateway

	testDeleted := func(r string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[r]
			if ok {
				return fmt.Errorf("VPN Gateway %q should have been deleted", r)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_virtual_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIVirtualGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVirtualGatewayExists("outscale_virtual_gateway.foo", &virtualGateway)),
			},
			{
				Config: testAccOAPINoVirtualGatewayConfig,
				Check:  resource.ComposeTestCheckFunc(testDeleted("outscale_virtual_gateway.foo")),
			},
		},
	})
}

func TestAccNet_ImportVirtualGateway_Basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVirtualGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfig,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccOutscaleOAPIVirtualGatewayDisappears(gateway *oscgo.VirtualGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		var err error

		opts := oscgo.DeleteVirtualGatewayRequest{
			VirtualGatewayId: gateway.GetVirtualGatewayId(),
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.VirtualGatewayApi.DeleteVirtualGateway(context.Background()).DeleteVirtualGatewayRequest(opts).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		return resource.Retry(40*time.Minute, func() *resource.RetryError {
			opts := oscgo.ReadVirtualGatewaysRequest{
				Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{gateway.GetVirtualGatewayId()}},
			}

			resp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(opts).Execute()
			if err != nil {
				if httpResp.StatusCode == http.StatusNotFound {
					return nil
				}
				if strings.Contains(err.Error(), utils.InvalidState) {
					return resource.RetryableError(fmt.Errorf(
						"Waiting for VPN Gateway to be in the correct state: %v", gateway.VirtualGatewayId))
				}
				return utils.CheckThrottling(httpResp, err)
			}
			if resp.GetVirtualGateways()[0].GetState() == "deleted" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for VPN Gateway: %v", gateway.VirtualGatewayId))
		})
	}
}

func testAccCheckOAPIVirtualGatewayDestroy(s *terraform.State) error {
	OSCAPI := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_virtual_gateway" {
			continue
		}

		var resp oscgo.ReadVirtualGatewaysResponse
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := OSCAPI.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
				Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{rs.Primary.ID}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err == nil {
			var v *oscgo.VirtualGateway
			for _, g := range resp.GetVirtualGateways() {
				if g.GetVirtualGatewayId() == rs.Primary.ID {
					v = &g
				}
			}

			if v == nil {
				// wasn't found
				return nil
			}

			if v.GetState() != "deleted" {
				return fmt.Errorf("Expected VPN Gateway to be in deleted state, but was not: %v", v)
			}
			return nil
		}
		return err
	}
	return nil
}

func testAccCheckOAPIVirtualGatewayExists(n string, ig *oscgo.VirtualGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		OSCAPI := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		var resp oscgo.ReadVirtualGatewaysResponse
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := OSCAPI.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
				Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{rs.Primary.ID}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return err
		}
		if len(resp.GetVirtualGateways()) == 0 {
			return fmt.Errorf("VPN Gateway not found")
		}

		*ig = resp.GetVirtualGateways()[0]

		return nil
	}
}

const testAccOAPINoVirtualGatewayConfig = `
	resource "outscale_net" "foo" {
		ip_range = "10.1.0.0/16"
	}
`

const testAccOAPIVirtualGatewayConfig = `
	resource "outscale_net" "foo" {
		ip_range = "10.1.0.0/16"
	}

	resource "outscale_virtual_gateway" "foo" {
	connection_type = "ipsec.1"	
	}

`

const testAccOAPIVirtualGatewayConfigChangeVPC = `
	resource "outscale_net" "bar" {
		ip_range = "10.2.0.0/16"
	}

	resource "outscale_virtual_gateway" "foo" {
	connection_type = "ipsec.1"	
}
`

func testAccOAPIVirtualGatewayConfigChangeTags(connectionType, name string) string {
	return fmt.Sprintf(`
		resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
		 connection_type = "%s"
		 tags {
		  key = "name"
		  value = "%s"
		  }
		}

	`, connectionType, name)
}

func testAccCheckOutscaleVirtualGatewayImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}
