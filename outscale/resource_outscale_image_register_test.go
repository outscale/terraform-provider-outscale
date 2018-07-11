package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleImageRegister_basic(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleImageRegisterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleImageRegisterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImageRegisterExists("outscale_image_register.outscale_image_register"),
				),
			},
		},
	})
}

func testAccCheckOutscaleImageRegisterDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_image_register" {
			continue
		}
		amiID := rs.Primary.ID
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		diReq := &fcu.DescribeImagesInput{
			ImageIds: []*string{aws.String(amiID)},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.VM.DescribeImages(diReq)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidAMIID.NotFound") {
				return nil
			}
			return fmt.Errorf("[DEBUG TES] Error register image %s", err)
		}

	}

	return nil
}

func testAccCheckOutscaleImageRegisterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role name is set")
		}

		return nil
	}
}

func testAccOutscaleImageRegisterConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_volume" "outscale_volume" {
  availability_zone = "eu-west-2a"
  size              = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
  volume_id = "${outscale_volume.outscale_volume.volume_id}"
}

resource "outscale_image_register" "outscale_image_register" {
  name = "registeredImageFromSnapshot-%d"

  root_device_name = "/dev/sda1"

  block_device_mapping {
    ebs {
	  snapshot_id = "${outscale_snapshot.outscale_snapshot.snapshot_id}"
	}
	device_name = "/dev/sda1"
  }
}`, r)
}
