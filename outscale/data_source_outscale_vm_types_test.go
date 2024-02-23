package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVM_TypesDataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVMTypesConfig(omi, "tinav4.c1r1p1"),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMTypesConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vms_types" {
			security_group_name = "sg_vm_type"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic_types" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vms_types.security_group_id]
		}

		data "outscale_vm_types" "vm_types" {
			filter {
				name = "bsu_optimized"
				values = ["true"]
			}
		}

		data "outscale_vm_types" "all-types" { }
	`, omi, vmType)
}
