package outscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccOthers_DhcpOptional_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_dhcp_option.foo"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))
	updateValue := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(value, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
				),
			},
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(updateValue, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "log_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
					resource.TestCheckResourceAttr(resourceName, "log_servers.0", "192.0.0.12"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.2"),
				),
			},
		},
	})
}

func TestAccOthers_DhcpOptional_withEmptyAttrs(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"

	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))
	updateValue := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	ntpServers := []string{"192.0.0.1", "192.0.0.2"}
	ntpServersUpdated := []string{"192.0.0.1", "192.0.0.3"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfigWithEmptyAttrs(ntpServers, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.1", "192.0.0.2"),
				),
			},
			{
				Config: testAccOAPIDHCPOptionalBasicConfigWithEmptyAttrs(ntpServersUpdated, updateValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.1", "192.0.0.3"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccNet_withDhcpOptional(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_dhcp_option.outscale_dhcp_option"
	domainName := fmt.Sprintf("%s.compute%s.internal", utils.GetRegion(), acctest.RandString(3))
	domainServers := []string{"192.168.12.12", "192.168.12.132"}

	tags := &oscgo.Tag{}
	tags.SetKey(acctest.RandomWithPrefix("name"))
	tags.SetValue(acctest.RandomWithPrefix("test-MZI"))

	domainNameUpdated := fmt.Sprintf("%s.compute%s.internal", utils.GetRegion(), acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalWithNet(domainName, domainServers, tags),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", domainServers[0]),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.1", domainServers[1]),
				),
			},
			{
				Config: testAccOAPIDHCPOptionalWithNet(domainNameUpdated, []string{}, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainNameUpdated),
				),
			},
		},
	})
}

func TestAccOthers_DHCPOption_importBasic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalBasicConfig(value, true, true),
			},
			testutils.ImportStepSDKv2(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

func testAccOAPIDHCPOptionalBasicConfig(value string, ntpServers bool, logServers bool) string {
	var ntp string
	var log string

	if ntpServers {
		ntp = `ntp_servers = ["192.0.0.2"]`
	}

	if logServers {
		log = `log_servers = ["192.0.0.12"]`
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

		%s
	}
	`, value, ntp, log)

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
