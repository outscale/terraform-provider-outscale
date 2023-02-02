package outscale

import (
	"fmt"
	"os"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Nic_DataSource(t *testing.T) {
	t.Parallel()
	subregion := os.Getenv("OUTSCALE_REGION")
	dataSourceName := "data.outscale_nic.nic"
	dataSourcesName := "data.outscale_nics.nics"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Nic_DataSource_Config(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "nics.#", "1"),

					resource.TestCheckResourceAttrSet(dataSourceName, "nic_id"),
				),
			},
		},
	})
}

func testAcc_Nic_DataSource_Config(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-nic-ds"
			}
		}
			
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.id
			tags {
				value = "tf-value"
				key   = "tf-key"
			}
		}

		data "outscale_nic" "nic" {
			filter {
				name = "nic_ids"
				values = [outscale_nic.outscale_nic.nic_id]
			} 
		}

		data "outscale_nics" "nics" {
			filter {
				name = "nic_ids"
				values = [outscale_nic.outscale_nic.nic_id]
			} 
		}  
	`, subregion)
}
