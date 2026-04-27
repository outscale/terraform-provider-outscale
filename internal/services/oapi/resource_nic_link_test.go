package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_NicLink_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_nic_link.outscale_nic_link"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNicLinkConfigBasic(omi, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "vm_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_id"),
				),
			},
			testacc.ImportStepWithStateIdFunc(resourceName, nicLinkStateIDFunc(resourceName), testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_NicLink_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccOutscaleNicLinkConfigBasic(omi, sgName),
		),
	})
}

func nicLinkStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return rs.Primary.Attributes["nic_id"], nil
	}
}

func testAccOutscaleNicLinkConfigBasic(omi, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_security_group" "security_group_nic" {
			net_id              = outscale_net.net.id
			security_group_name = "%s"
			description         = "testacc-nic-link"
		}

		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = [outscale_security_group.security_group_nic.security_group_id]
			placement_subregion_name = "%[4]sa"
			subnet_id                = outscale_subnet.outscale_subnet.id

			lifecycle { ignore_changes = [state] }
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[4]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.net.id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids   = [outscale_security_group.security_group_nic.security_group_id]
		}

		resource "outscale_nic_link" "outscale_nic_link" {
			device_number = 1
			vm_id         = outscale_vm.vm.id
			nic_id        = outscale_nic.outscale_nic.id
		}
	`, sgName, omi, testAccVmType, utils.GetRegion())
}
