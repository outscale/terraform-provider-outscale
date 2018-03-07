package outscale

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVolumeAttachment_basic(t *testing.T) {
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
						"outscale_volume_link.ebs_att", "device", "/dev/sdh"),
					testAccCheckInstanceExists(
						"outscale_vm.web", &i),
					testAccCheckVolumeExists(
						"outscale_volume.example", &v),
					testAccCheckVolumeAttachmentExists(
						"outscale_volume_link.ebs_att", &i, &v),
				),
			},
		},
	})
}

func TestAccOutscaleVolumeAttachment_skipDestroy(t *testing.T) {
	var i fcu.Instance
	var v fcu.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfigSkipDestroy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_volume_link.ebs_att", "device", "/dev/sdh"),
					testAccCheckInstanceExists(
						"outscale_vm.web", &i),
					testAccCheckVolumeExists(
						"outscale_volume.example", &v),
					testAccCheckVolumeAttachmentExists(
						"outscale_volume_link.ebs_att", &i, &v),
				),
			},
		},
	})
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

func testAccCheckVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		log.Printf("\n\n----- This is never called")
		if rs.Type != "outscale_volume_link" {
			continue
		}
	}
	return nil
}

const testAccVolumeAttachmentConfig = `
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	instance_type = "t1.micro"
	tags {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
	size = 1
	availability_zone = "eu-west-2a"
}
resource "outscale_volume_link" "ebs_att" {
  device = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	instance_id = "${outscale_vm.web.id}"
}
`

const testAccVolumeAttachmentConfigSkipDestroy = `
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	instance_type = "t1.micro"
	tags {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
	size = 1
	tags {
		Name = "TestVolume"
	}
}
data "outscale_volume" "ebs_volume" {
    filter {
	name = "size"
	values = ["${outscale_volume.example.size}"]
    }
    filter {
	name = "availability-zone"
	values = ["${outscale_volume.example.availability_zone}"]
    }
    filter {
	name = "tag:Name"
	values = ["TestVolume"]
    }
}
resource "outscale_volume_link" "ebs_att" {
  	device = "/dev/sdh"
	volume_id = "${data.outscale_volume.ebs_volume.id}"
	instance_id = "${outscale_vm.web.id}"
	skip_destroy = true
}
`
