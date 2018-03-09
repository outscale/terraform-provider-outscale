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

func TestAccOutscaleOAPIVolumeAttachment_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var i fcu.Instance
	var v fcu.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfig,
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

func TestAccOutscaleOAPIVolumeAttachment_skipDestroy(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var i fcu.Instance
	var v fcu.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVolumeAttachmentConfigSkipDestroy,
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

func testAccCheckOAPIVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		log.Printf("\n\n----- This is never called")
		if rs.Type != "outscale_volume_link" {
			continue
		}
	}
	return nil
}

const testAccOAPIVolumeAttachmentConfig = `
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	instance_type = "t1.micro"
	tag = {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
	size = 1
	sub_region_name = "eu-west-2a"
	tag = {
		Name = "HelloWorld Volume"
	}
}
resource "outscale_volume_link" "ebs_att" {
  device_name = "/dev/sdh"
	volume_id = "${outscale_volume.example.id}"
	vm_id = "${outscale_vm.web.id}"
}
`

const testAccOAPIVolumeAttachmentConfigSkipDestroy = `
resource "outscale_vm" "web" {
	image_id = "ami-8a6a0120"
	instance_type = "t1.micro"
	tag = {
		Name = "HelloWorld"
	}
}
resource "outscale_volume" "example" {
	size = 1
	tag = {
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
	values = ["${outscale_volume.example.sub_region_name}"]
    }
    filter {
	name = "tag:Name"
	values = ["TestVolume"]
    }
}
resource "outscale_volume_link" "ebs_att" {
  	device_name = "/dev/sdh"
	volume_id = "${data.outscale_volume.ebs_volume.id}"
	vm_id = "${outscale_vm.web.id}"
	skip_destroy = true
}
`
