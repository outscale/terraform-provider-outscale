package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/antihax/optional"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func TestAccOutscaleOAPIDhcpOptional_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_dhcp_option.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists("outscale_dhcp_option.foo"),
				),
			},
		},
	})
}

func TestAccOutscaleDHCPOption_importBasic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalConfig(),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleInternetServiceLinkImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleDHCPOptionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckOutscaleDHCPOptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Option id is set")
		}

		resp, _, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background(), &oscgo.ReadDhcpOptionsOpts{
			ReadDhcpOptionsRequest: optional.NewInterface(oscgo.ReadDhcpOptionsRequest{
				Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{rs.Primary.ID}},
			}),
		})
		if err != nil || len(resp.GetDhcpOptionsSets()) < 1 {
			return fmt.Errorf("DHCP Option is not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckOAPIDHCPOptionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_customer_endpoint" {
			continue
		}

		resp, _, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background(), &oscgo.ReadDhcpOptionsOpts{
			ReadDhcpOptionsRequest: optional.NewInterface(oscgo.ReadDhcpOptionsRequest{
				Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{rs.Primary.ID}},
			}),
		})
		if strings.Contains(fmt.Sprint(err), "InvalidDhcpID.NotFound") {
			continue
		}

		if err != nil {
			return err
		}

		if len(resp.GetDhcpOptionsSets()) > 0 {
			return fmt.Errorf("DHCP still exists: %v", resp.GetDhcpOptionsSets())
		}
	}

	return nil
}

func testAccOAPIDHCPOptionalConfig() string {
	return `
		resource "outscale_dhcp_option" "foo" {
			domain_name         = "test.fr"
			domain_name_servers = ["192.168.12.1"]
			ntp_servers         = ["192.0.0.2"]
			tags {
				key   = "name"
				value = "MyDHCPoptionsSet"
			}
		}
	`
}
