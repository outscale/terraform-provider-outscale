package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_DHCPOptional_Basic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDHCPOptionalBasicConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
				),
			},
			{
				Config: testAccDHCPOptionalBasicConfig(true, true),
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
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_withDHCPOptional(t *testing.T) {
	resourceName := "outscale_dhcp_option.outscale_dhcp_option"
	domainServers := []string{"192.168.12.12", "192.168.12.132"}
	domainName := fmt.Sprintf("%s.compute%s.internal", utils.GetRegion(), acctest.RandString(3))
	domainNameUpdated := fmt.Sprintf("%s.compute%s.internal", utils.GetRegion(), acctest.RandString(3))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDHCPOptionalWithNet(domainName, domainServers),
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
				Config: testAccOAPIDHCPOptionalWithNet(domainNameUpdated, []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainNameUpdated),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_DHCPOption_Migration(t *testing.T) {
	domainServers := []string{"192.168.12.12", "192.168.12.132"}
	domainName := fmt.Sprintf("%s.compute%s.internal", utils.GetRegion(), acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccDHCPOptionalBasicConfig(false, false),
			testAccOAPIDHCPOptionalWithNet(domainName, domainServers),
		),
	})
}

func testAccDHCPOptionalBasicConfig(ntpServers bool, logServers bool) string {
	var ntp string
	var log string

	if ntpServers {
		ntp = `ntp_servers = ["192.0.0.2"]`
	}

	if logServers {
		log = `log_servers = ["192.0.0.12"]`
	}

	return fmt.Sprintf(`
	resource "outscale_dhcp_option" "foo" {
		domain_name         = "test.fr"
		domain_name_servers = ["192.168.12.1"]

		%s

		%s
	}
	`, ntp, log)
}

func testAccOAPIDHCPOptionalWithNet(domainName string, domainServers []string) string {
	var servers string

	if len(domainServers) > 0 {
		servers = fmt.Sprintf(
			`domain_name_servers = %s`,
			strings.ReplaceAll(fmt.Sprintf("%+q", domainServers), " ", ","),
		)
	}

	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "outscale_dhcp_option" {
			domain_name = "%s"

			%s
		}

		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_net" "vpc" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_net_attributes" "net_attr_with_net" {
			net_id              = outscale_net.net.id
			dhcp_options_set_id = outscale_dhcp_option.outscale_dhcp_option.id
		}

		resource "outscale_net_attributes" "net_attr_with_vpc" {
			net_id              = outscale_net.vpc.id
			dhcp_options_set_id = outscale_dhcp_option.outscale_dhcp_option.id
		}
	`, domainName, servers)
}
