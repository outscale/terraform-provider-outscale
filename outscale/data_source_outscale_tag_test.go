package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPITagDataSource(t *testing.T) {
	//t.Skip()
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPITagDataSourceConfig(omi, "t2.micro"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "key", "Name"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "value", "test-vm"),
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

			tags {
				key = "Name"
				value = "test-vm"
			}
		}

		data "outscale_tag" "web" {
			filter {
				name = "resource_ids"
				values = ["${outscale_vm.basic.id}"]
			}
		}
	`, omi, vmType)
}
