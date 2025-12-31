package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DHCPOption_basic(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"
	dataSourceName := "data.outscale_dhcp_option.test"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientDHCPOptionBasic(value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "log_servers.#"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "log_servers.0", "192.0.0.12"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dhcp_options_set_id"),
				),
			},
		},
	})
}

func TestAccOthers_DHCPOption_withFilters(t *testing.T) {
	resourceName := "outscale_dhcp_option.foo"
	dataSourceName := "data.outscale_dhcp_option.test"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClientDHCPOptionWithFilters(value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "log_servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),

					resource.TestCheckResourceAttr(resourceName, "domain_name", "test.fr"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.12.1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "log_servers.0", "192.0.0.12"),
					resource.TestCheckResourceAttrSet(dataSourceName, "filter.#"),
					resource.TestCheckResourceAttr(dataSourceName, "filter.#", "2"),
				),
			},
		},
	})
}

func testAccClientDHCPOptionBasic(value string) string {
	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "foo" {
			domain_name         = "test.fr"
			domain_name_servers = ["192.168.12.1"]
			ntp_servers         = ["192.0.0.2"]
			log_servers         = ["192.0.0.12"]

			tags {
				key   = "name"
				value = "%s"
			}
		}

		data "outscale_dhcp_option" "test" {
			filter {
				name = "dhcp_options_set_ids"
				values = [outscale_dhcp_option.foo.id]
			}
		}

	`, value)
}

func testAccClientDHCPOptionWithFilters(value string) string {
	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "foo" {
			domain_name         = "test.fr"
			domain_name_servers = ["192.168.12.1"]
			ntp_servers         = ["192.0.0.2"]
			log_servers         = ["192.0.0.12"]

			tags {
				key   = "name"
				value = "%s"
			}
		}

		data "outscale_dhcp_option" "test" {
			filter {
				name = "dhcp_options_set_ids"
				values = [outscale_dhcp_option.foo.id]
			}
			filter {
				name = "tag_keys"
				values = ["name"]
			}
		}
	`, value)
}
