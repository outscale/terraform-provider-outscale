package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPITagDataSource(t *testing.T) {
	t.Skip()
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPITagDataSourceConfig(omi, "c4.large"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "key", "foo"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "value", "bar"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "resource_type", "instance"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccOAPITagDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id            = "%s"
			vm_type             = "%s"
			keypair_name        = "terraform-basic"
		}

		data "outscale_tag" "web" {
			filter {
				name = "resource_id"
				values = ["${outscale_vm.basic.id}"]
			}
		}
	`, omi, vmType)
}
