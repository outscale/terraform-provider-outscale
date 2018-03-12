package outscale

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIVolumeAttachment_basic(t *testing.T) {
	var i fcu.Instance
	var v fcu.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volume_link.ebs_att", "device_name", "/dev/sdh"),
					testAccCheckInstanceExists(
						"outscale_vm.web", &i),
					testAccCheckOAPIVolumeExists(
						"outscale_volume.example", &v),
					testAccCheckOAPIVolumeAttachmentExists(
						"outscale_volume_link.ebs_att", &i, &v),
				),
			},
		},
	})
}

func testAccCheckOAPIVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		log.Printf("\n\n----- This is never called")
		if rs.Type != "outscale_volume_link" {
			continue
		}
	}
	return nil
}

func testAccCheckOAPIVolumeAttachmentExists(n string, i *fcu.Instance, v *fcu.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for _, b := range i.BlockDeviceMappings {
			if rs.Primary.Attributes["device_name"] == *b.DeviceName {
				if b.Ebs.VolumeId != nil && rs.Primary.Attributes["volume_id"] == *b.Ebs.VolumeId {
					// pass
					return nil
				}
			}
		}

		return fmt.Errorf("Error finding instance/volume")
	}
}

const testAccOAPIVolumeAttachmentConfig = `
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	type = "t1.micro"
	tag {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
  sub_region_name = "eu-west-2a"
	size = 1
}
resource "outscale_volume_link" "ebs_att" {
  device_name = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	vm_id = "${outscale_vm.web.id}"
}
`
