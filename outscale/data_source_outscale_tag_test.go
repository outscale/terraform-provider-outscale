package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(omi, "tinav4.c2r2p2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "key", "Name"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "value", "test-vm"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "resource_type", "vm"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccTagDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id            = "%s"
			vm_type             = "%s"
			keypair_name        = "terraform-basic"

			tags {
				key = "Name"
				value = "test-vm"
			}
		}

		data "outscale_tag" "web" {
			filter {
				name = "resource_ids"
				values = [outscale_vm.basic.id]
			}
		}
	`, omi, vmType)
}
