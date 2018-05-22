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

func TestAccOutscaleOAPIImage_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var ami fcu.Image
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIImageConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIImageExists("outscale_image.foo", &ami),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "name", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.device_virtual_name", "/dev/sda1"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.ebs.delete_on_vm_termination", "true"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "state_reason.code", "UNSET"),
				),
			},
		},
	})
}

func testAccCheckOAPIImageDestroy(s *terraform.State) error {
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
			return fmt.Errorf("OMI %s still exists in the state: %s", *resp.Images[0].ImageId, *state)
		}
	}
	return nil
}

func testAccCheckOAPIImageExists(n string, ami *fcu.Image) resource.TestCheckFunc {
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

func testAccOAPIImageConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_image" "foo" {
	name = "tf-testing-%d"
	vm_id = "${outscale_vm.basic.id}"
}
	`, rInt)
}
