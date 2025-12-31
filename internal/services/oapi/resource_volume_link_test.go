package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_WithVolumeAttachment_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, testAccVmType, utils.GetRegion(), keypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volume_link.ebs_att", "device_name", "/dev/sdh"),
				),
			},
		},
	})
}

func TestAccVM_ImportVolumeAttachment_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"

	resourceName := "outscale_volume_link.ebs_att"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(), CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, testAccVmType, utils.GetRegion(), keypair, sgName),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccVM_WithVolumeAttachment_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps: testacc.FrameworkMigrationTestSteps("1.1.3",
			testAccOAPIVolumeAttachmentConfig(omi, testAccVmType, utils.GetRegion(), keypair, sgName),
		),
	})
}

func testAccCheckOAPIVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_volume_link" {
			continue
		}
	}
	return nil
}

func testAccOAPIVolumeAttachmentConfig(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vol_link" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "%[5]s"
		}
		resource "outscale_vm" "web" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			placement_subregion_name = "%[3]sb"
			security_group_ids = [outscale_security_group.sg_vol_link.security_group_id]
		}

		resource "outscale_volume" "volume" {
			subregion_name = "%[3]sb"
			volume_type    = "standard"
			size           = 100
		}

		resource "outscale_volume_link" "ebs_att" {
			device_name = "/dev/sdh"
			volume_id   = outscale_volume.volume.id
			vm_id       = outscale_vm.web.id
		}
	`, omi, vmType, region, keypair, sgName)
}
