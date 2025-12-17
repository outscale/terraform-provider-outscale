package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"
)

func TestAccNet_WithVirtualGateway_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_virtual_gateway.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfigChangeVPC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "net_to_virtual_gateway_links.#"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_to_virtual_gateway_links.0.net_id"),
				),
			},
			// Ignore attributes related to the Gateway Link, that gets populated after a refresh
			testutils.ImportStep(resourceName, "net_to_virtual_gateway_links", "request_id"),
		},
	})
}

func TestAccOthers_VirtualGatewayChangeTags(t *testing.T) {
	resourceName := "outscale_virtual_gateway.outscale_virtual_gateway"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfigChangeTags("ipsec.1", "test-VGW"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "test-VGW"),
				),
			},
			{
				Config: testAccOAPIVirtualGatewayConfigChangeTags("ipsec.1", "test-VGW2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "test-VGW2"),
				),
			},
		},
	})
}

func TestAccOthres_ImportVirtualGateway_Basic(t *testing.T) {
	resourceName := "outscale_virtual_gateway.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVirtualGatewayConfig,
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

const testAccOAPIVirtualGatewayConfig = `
	resource "outscale_virtual_gateway" "foo" {
	    connection_type = "ipsec.1"
    }`

const testAccOAPIVirtualGatewayConfigChangeVPC = `
	resource "outscale_net" "bar" {
		ip_range = "10.2.0.0/16"
	}

	resource "outscale_virtual_gateway" "foo" {
	    connection_type = "ipsec.1"
    }
	resource "outscale_virtual_gateway_link" "test" {
        virtual_gateway_id = outscale_virtual_gateway.foo.virtual_gateway_id
        net_id             = outscale_net.bar.net_id
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
