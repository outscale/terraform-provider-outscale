package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_WithVolumeAttachment_Basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, utils.TestAccVmType, utils.GetRegion(), keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volume_link.ebs_att", "device_name", "/dev/sdh"),
				),
			},
		},
	})
}

func TestAccVM_WithVolumeAttachment_Basic_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.3", testAccOAPIVolumeAttachmentConfig(omi, utils.TestAccVmType, utils.GetRegion(), keypair)),
	})
}

func TestAccVM_ImportVolumeAttachment_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"

	resourceName := "outscale_volume_link.ebs_att"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(), CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, utils.TestAccVmType, utils.GetRegion(), keypair),
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
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

func testAccOAPIVolumeAttachmentConfig(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vol_link" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sg_volumes_link"
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
	`, omi, vmType, region, keypair)
}
