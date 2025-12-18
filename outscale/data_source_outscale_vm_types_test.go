package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVM_TypesDataSource_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVMTypesConfig(omi, "tinav5.c2r2p2", sgName),
			},
		},
	})
}

func testAccDataSourceOutscaleVMTypesConfig(omi, vmType, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vms_types" {
			security_group_name = "%[3]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic_types" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
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
	`, omi, vmType, sgName)
}
