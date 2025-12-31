package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccVM_DataSource_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	datasourcceName := "data.outscale_vm.basic_web"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMDataSourceConfig(omi, oapi.TestAccVmType, utils.GetRegion(), sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourcceName, "image_id", omi),
					resource.TestCheckResourceAttr(datasourcceName, "vm_type", oapi.TestAccVmType),
					resource.TestCheckResourceAttr(datasourcceName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccOAPIVMDataSourceConfig(omi, vmType, region, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_data" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "outscale_vm_data" {
			image_id     = "%[1]s"
			vm_type      = "%[2]s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_vm_data.security_group_id]
			tags {
				key   = "name"
				value = "Terraform-VM"
			}

		}

		data "outscale_vm" "basic_web" {
			filter {
				name   = "vm_ids"
				values = [outscale_vm.outscale_vm_data.vm_id]
			}
		}
	`, omi, vmType, region, sgName)
}
