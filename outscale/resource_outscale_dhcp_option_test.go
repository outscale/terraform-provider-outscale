package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/openlyinc/pointy"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccOutscaleOAPIDhcpOptional_basic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))
	updateValue := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(value, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
					/*resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", value),*/
				),
			},
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(updateValue, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.2"),
					/*resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", updateValue),*/
				),
			},
		},
	})
}

func TestAccOutscaleOAPIDhcpOptional_withEmptyAttrs(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"

	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))
	updateValue := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	ntpServres := []string{"192.0.0.1", "192.0.0.2"}
	ntpServresUpdated := []string{"192.0.0.1", "192.0.0.3"}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfigWithEmptyAttrs(ntpServres, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "ntp_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.1", "192.0.0.2"),
					/*resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", value),*/
				),
			},
			{
				Config: testAccOAPIDHCPOptionalBasicConfigWithEmptyAttrs(ntpServresUpdated, updateValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "ntp_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.1", "192.0.0.3"),
					/*resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", updateValue),*/
				),
			},
		},
	})
}

func TestAccOutscaleOAPIDhcpOptional_withNet(t *testing.T) {
	resourceName := "outscale_dhcp_option.outscale_dhcp_option"

	domainName := fmt.Sprintf("eu-west-2.compute%s.internal", acctest.RandString(3))
	ntpServres := []string{"192.168.12.12", "192.168.12.132"}
	tags := &oscgo.Tag{
		Key:   pointy.String(acctest.RandomWithPrefix("name")),
		Value: pointy.String(acctest.RandomWithPrefix("test-MZI")),
	}

	domainNameUpdated := fmt.Sprintf("eu-west-2.compute%s.internal", acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalWithNet(domainName, ntpServres, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", ntpServres[0]),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.1", ntpServres[1]),
					/*resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", tags.GetKey()),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", tags.GetValue()),*/
				),
			},
			{
				Config: testAccOAPIDHCPOptionalWithNet(domainNameUpdated, []string{}, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDHCPOptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainNameUpdated),
				),
			},
		},
	})
}

func TestAccOutscaleDHCPOption_importBasic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(value, true),
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

		resp, _, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background()).ReadDhcpOptionsRequest(oscgo.ReadDhcpOptionsRequest{
			Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{rs.Primary.ID}},
		}).Execute()
		if err != nil || len(resp.GetDhcpOptionsSets()) < 1 {
			return fmt.Errorf("DHCP Option is not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckOAPIDHCPOptionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_dhcp_option" {
			continue
		}

		resp, _, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background()).ReadDhcpOptionsRequest(oscgo.ReadDhcpOptionsRequest{
			Filters: &oscgo.FiltersDhcpOptions{DhcpOptionsSetIds: &[]string{rs.Primary.ID}},
		}).Execute()
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

func testAccOAPIDHCPOptionalBasicConfig(value string, ntpServers bool) string {
	var ntp string

	if ntpServers {
		ntp = `ntp_servers = ["192.0.0.2"]`
	}

	tf := fmt.Sprintf(`
	resource "outscale_dhcp_option" "foo" {
		domain_name         = "test.fr"
		domain_name_servers = ["192.168.12.1"]


		tags {
			key   = "name"
			value = "%s"
		}

		%s

	}
	`, value, ntp)

	return tf
}

func testAccOAPIDHCPOptionalBasicConfigWithEmptyAttrs(ntpServers []string, value string) string {
	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "foo" {
			ntp_servers = %s

			tags {
				key   = "name"
				value = "%s"
			}
		}
	`, strings.ReplaceAll(fmt.Sprintf("%+q", ntpServers), " ", ","), value)
}

func testAccOAPIDHCPOptionalWithNet(domainName string, domainServers []string, tags *oscgo.Tag) string {
	var servers, dhcpTags string

	if len(domainServers) > 0 {
		servers = fmt.Sprintf(
			`domain_name_servers = %s`,
			strings.ReplaceAll(fmt.Sprintf("%+q", domainServers), " ", ","),
		)
	}

	if tags != nil {
		dhcpTags = fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, tags.GetKey(), tags.GetValue())
	}

	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "outscale_dhcp_option" {
			domain_name = "%s"

			%s

			%s
		}

		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "name"
				value = "net"
			}
		}

		resource "outscale_net" "vpc" {
			ip_range = "10.0.0.0/16"
			tags {
				key   = "name"
				value = "vpc"
			}
		}

		resource "outscale_net_attributes" "net_attr_with_net" {
			net_id              = outscale_net.net.id
			dhcp_options_set_id = outscale_dhcp_option.outscale_dhcp_option.id
		}

		resource "outscale_net_attributes" "net_attr_with_vpc" {
			net_id              = outscale_net.vpc.id
			dhcp_options_set_id = outscale_dhcp_option.outscale_dhcp_option.id
		}
	`, domainName, servers, dhcpTags)
}
