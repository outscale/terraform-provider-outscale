package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVolumeAttachment_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var i fcu.Instance
	var v fcu.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volumes_link.ebs_att", "device", "/dev/sdh"),
					testAccCheckInstanceExists(
						"outscale_vm.web", &i),
					testAccCheckVolumeExists(
						"outscale_volume.example", &v),
					testAccCheckVolumeAttachmentExists(
						"outscale_volumes_link.ebs_att", &i, &v),
				),
			},
		},
	})
}

func testAccCheckVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		log.Printf("\n\n----- This is never called")
		if rs.Type != "outscale_volumes_link" {
			continue
		}
	}
	return nil
}

func testAccCheckVolumeAttachmentExists(n string, i *fcu.Instance, v *fcu.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for _, b := range i.BlockDeviceMappings {
			if rs.Primary.Attributes["device"] == *b.DeviceName {
				if b.Ebs.VolumeId != nil && rs.Primary.Attributes["volume_id"] == *b.Ebs.VolumeId {
					// pass
					return nil
				}
			}
		}

		return fmt.Errorf("Error finding instance/volume")
	}
}

const testAccVolumeAttachmentConfig = `
resource "outscale_vm" "web" {
	image_id      = "ami-880caa66"
	instance_type = "t1.micro"
	tag {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
  availability_zone = "eu-west-2a"
	size = 1
}
resource "outscale_volumes_link" "ebs_att" {
   device = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	instance_id = "${outscale_vm.web.id}"
}
`
