package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_VMTypes_DataSource(t *testing.T) {
	t.Parallel()
	dataSourcesName := "data.outscale_vm_types.vm_types"
	dataSourcesAllName := "data.outscale_vm_types.all-types"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_VMTypes_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "vm_types.#"),
					resource.TestCheckResourceAttrSet(dataSourcesAllName, "vm_types.#"),
				),
			},
		},
	})
}

func testAcc_VMTypes_DataSource_Config() string {
	return fmt.Sprintf(`
		data "outscale_vm_types" "vm_types" {
			filter {
				name = "bsu_optimized"
				values = ["true"]
			}
		}

		data "outscale_vm_types" "all-types" { }
	`)
}
