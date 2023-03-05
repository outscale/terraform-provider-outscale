package outscale

import (
	"fmt"
	"os"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIVolumeAttachment_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	var i oscgo.Vm
	var v oscgo.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, "tinav4.c2r2p2", utils.GetRegion(), keypair, sgId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volumes_link.ebs_att", "device_name", "/dev/sdh"),
					testAccCheckOutscaleOAPIVMExists("outscale_vm.web", &i),
					testAccCheckOAPIVolumeAttachmentExists(
						"outscale_volumes_link.ebs_att", &i, &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolumeAttachment_importBasic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	resourceName := "outscale_volumes_link.ebs_att"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig(omi, "tinav4.c2r2p2", utils.GetRegion(), keypair, sgId),
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

func testAccCheckOAPIVolumeAttachmentExists(n string, i *oscgo.Vm, v *oscgo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for _, b := range i.GetBlockDeviceMappings() {
			if rs.Primary.Attributes["device_name"] == b.GetDeviceName() {
				if rs.Primary.Attributes["volume_id"] == b.Bsu.GetVolumeId() {
					// pass
					return nil
				}
			}
		}

		return fmt.Errorf("Error finding instance/volume")
	}
}

func testAccOAPIVolumeAttachmentConfig(omi, vmType, region, keypair, sgId string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = "${outscale_net.net.id}"
		}

		resource "outscale_vm" "web" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = ["%[5]s"]
			placement_subregion_name = "%[3]sb"
		}

		resource "outscale_volume" "volume" {
			subregion_name = "%[3]sb"
			volume_type    = "standard"
			size           = 100
		}

		resource "outscale_volumes_link" "ebs_att" {
			device_name = "/dev/sdh"
			volume_id   = "${outscale_volume.volume.id}"
			vm_id       = "${outscale_vm.web.id}"
		}
	`, omi, vmType, region, keypair, sgId)
}
