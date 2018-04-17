package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleImage_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var ami fcu.Image
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists("outscale_image.foo", &ami),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "name", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.ebs.delete_on_termination", "true"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "state_reason.code", "UNSET"),
				),
			},
		},
	})
}

func testAccCheckImageDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_image" {
			continue
		}

		// Try to find the OMI
		log.Printf("OMI-ID: %s", rs.Primary.ID)
		DescribeAmiOpts := &fcu.DescribeImagesInput{
			ImageIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeImagesOutput
		var err error
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			resp, err = conn.FCU.VM.DescribeImages(DescribeAmiOpts)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if resp == nil {
			return nil
		}

		if err != nil {
			return err
		}

		if len(resp.Images) > 0 {
			state := resp.Images[0].State
			return fmt.Errorf("OMI %s still exists in the state: %s.", *resp.Images[0].ImageId, *state)
		}
	}
	return nil
}

func testAccCheckImageExists(n string, ami *fcu.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("OMI Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OMI ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)
		opts := &fcu.DescribeImagesInput{
			ImageIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeImagesOutput
		var err error
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			resp, err = conn.FCU.VM.DescribeImages(opts)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}
		if len(resp.Images) == 0 {
			return fmt.Errorf("OMI not found")
		}
		*ami = *resp.Images[0]
		return nil
	}
}

func testAccCheckAmiBlockDevice(ami *fcu.Image, blockDevice *fcu.BlockDeviceMapping, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		devices := make(map[string]*fcu.BlockDeviceMapping)
		for _, device := range ami.BlockDeviceMappings {
			devices[*device.DeviceName] = device
		}

		// Check if the block device exists
		if _, ok := devices[n]; !ok {
			return fmt.Errorf("block device doesn't exist: %s", n)
		}

		*blockDevice = *devices[n]
		return nil
	}
}

func testAccCheckAmiEbsBlockDevice(bd *fcu.BlockDeviceMapping, ed *fcu.EbsBlockDevice) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Test for things that ed has, don't care about unset values
		cd := bd.Ebs
		if ed.VolumeType != nil {
			if *ed.VolumeType != *cd.VolumeType {
				return fmt.Errorf("Volume type mismatch. Expected: %s Got: %s",
					*ed.VolumeType, *cd.VolumeType)
			}
		}
		if ed.DeleteOnTermination != nil {
			if *ed.DeleteOnTermination != *cd.DeleteOnTermination {
				return fmt.Errorf("DeleteOnTermination mismatch. Expected: %t Got: %t",
					*ed.DeleteOnTermination, *cd.DeleteOnTermination)
			}
		}
		if ed.Encrypted != nil {
			if *ed.Encrypted != *cd.Encrypted {
				return fmt.Errorf("Encrypted mismatch. Expected: %t Got: %t",
					*ed.Encrypted, *cd.Encrypted)
			}
		}
		// Integer defaults need to not be `0` so we don't get a panic
		if ed.Iops != nil && *ed.Iops != 0 {
			if *ed.Iops != *cd.Iops {
				return fmt.Errorf("IOPS mismatch. Expected: %d Got: %d",
					*ed.Iops, *cd.Iops)
			}
		}
		if ed.VolumeSize != nil && *ed.VolumeSize != 0 {
			if *ed.VolumeSize != *cd.VolumeSize {
				return fmt.Errorf("Volume Size mismatch. Expected: %d Got: %d",
					*ed.VolumeSize, *cd.VolumeSize)
			}
		}

		return nil
	}
}

func testAccImageConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_1"
  group_description = "Used in the terraform acceptance tests"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	security_group = ["${outscale_firewall_rules_set.web.id}"]
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}

resource "outscale_image" "foo" {
	name = "tf-testing-%d"
	instance_id = "${outscale_vm.basic.id}"
}
	`, rInt, rInt)
}
