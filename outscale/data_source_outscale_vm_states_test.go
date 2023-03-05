package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVM_StatesDataSource(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVMStatesConfig(omi, "tinav4.c2r2p2"),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMStatesConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
		}

		data "outscale_vm_states" "state" {
			vm_ids = [outscale_vm.basic.id]
		}
	`, omi, vmType)
}
