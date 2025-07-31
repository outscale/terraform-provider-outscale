package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVM_WithVolumeAttachment_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
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

func TestAccVM_ImportVolumeAttachment_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resourceName := "outscale_volume_link.ebs_att"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(), CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, utils.TestAccVmType, utils.GetRegion(), keypair),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOAPIVolumeAttachmentImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOAPIVolumeAttachmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
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
