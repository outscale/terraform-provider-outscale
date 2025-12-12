package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_WithTagDataSource(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPITagDataSourceConfig(omi, utils.TestAccVmType),
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
func testAccOAPITagDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_tag" {
			security_group_name = "sg_tag"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basicTag" {
			image_id            = "%s"
			vm_type             = "%s"
			keypair_name        = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vm_tag.security_group_id]
			tags {
				key = "Name"
				value = "test-vm"
			}
		}

		data "outscale_tag" "web" {
			filter {
				name = "resource_ids"
				values = [outscale_vm.basicTag.id]
			}
		}
	`, omi, vmType)
}
