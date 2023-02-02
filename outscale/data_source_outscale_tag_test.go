package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Tag_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_tag.vg"
	dataSourcesName := "data.outscale_tags.vg"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Tag_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key", "name"),
					resource.TestCheckResourceAttr(dataSourceName, "value", "test acc"),
					resource.TestCheckResourceAttr(dataSourceName, "resource_type", "virtual-private-gateway"),

					resource.TestCheckResourceAttrSet(dataSourcesName, "tags.#"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAcc_Tag_DataSource_Config() string {
	return fmt.Sprintf(`
	resource "outscale_virtual_gateway" "vg" {
		connection_type = "ipsec.1"	

		tags {
			key   = "name"
			value = "test acc"
		}
	}

		data "outscale_tag" "vg" {
			filter {
				name = "resource_ids"
				values = ["${outscale_virtual_gateway.vg.id}"]
			}
		}
		data "outscale_tags" "vg" {
			filter {
				name   = "resource_type"
				values = ["${outscale_virtual_gateway.vg.id}"]
			}
		}
	`)
}
