package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_DHCPOption_DataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.outscale_dhcp_option.test"
	dataSourcesName := "data.outscale_dhcp_options.test"
	value := fmt.Sprintf("test-acc-value-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_DHCPOption_DataSource_Config(value),
				Check: resource.ComposeTestCheckFunc(

					// data source validations
					resource.TestCheckResourceAttrSet(dataSourceName, "dhcp_options_set_id"),
					resource.TestCheckResourceAttr(dataSourcesName, "dhcp_options.#", "2"),
				),
			},
		},
	})
}

func testAcc_DHCPOption_DataSource_Config(value string) string {
	return fmt.Sprintf(`
		resource "outscale_dhcp_option" "foo1" {
			domain_name         = "test.fr"
			domain_name_servers = ["192.168.12.1"]
			ntp_servers         = ["192.0.0.2"]
			log_servers         = ["192.0.0.12"]

			tags {
				key   = "name"
				value = "%s"
			}
		}

		resource "outscale_dhcp_option" "foo2" {
			domain_name         = "test.fr"
			domain_name_servers = ["192.168.12.2"]
			ntp_servers         = ["192.0.0.3"]
			log_servers         = ["192.0.0.13"]
			
			tags {
				key   = "name"
				value = "%[1]s"
			}
		}

		data "outscale_dhcp_option" "test" {
			filter {
				name = "dhcp_options_set_ids"
				values = ["${outscale_dhcp_option.foo1.id}"]
			}
		}

		data "outscale_dhcp_options" "test" {
			filter {
				name = "dhcp_options_set_ids"
				values = ["${outscale_dhcp_option.foo1.id}", "${outscale_dhcp_option.foo2.id}"]
			}
		}
	`, value)
}
