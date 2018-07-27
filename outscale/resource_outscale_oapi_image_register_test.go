package outscale

import (
	"fmt"
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

func TestAccOutscaleOAPIImageRegister_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIImageRegisterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIImageRegisterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImageRegisterExists("outscale_image_register.outscale_image_register"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIImageRegisterDestroy(s *terraform.State) error {

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

func testAccCheckOutscaleOAPIImageRegisterExists(n string) resource.TestCheckFunc {
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

func testAccOutscaleOAPIImageRegisterConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_vm" {
    count = 1
    image_id                    = "ami-880caa66"
    type               = "c4.large"
}
resource "outscale_image_register" "outscale_image_register" {
    name        = "image_%d"
    vm_id = "${outscale_vm.outscale_vm.id}"
}`, r)
}
