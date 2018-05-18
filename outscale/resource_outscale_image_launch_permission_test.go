package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleImageLaunchPermission_Basic(t *testing.T) {
	imageID := ""
	accountID := "520679080430"

	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []r.TestStep{
			// Scaffold everything
			r.TestStep{
				Config: testAccOutscaleImageLaunchPermissionConfig(rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Drop just launch permission to test destruction
			r.TestStep{
				Config: testAccOutscaleImageLaunchPermissionConfig(rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleImageLaunchPermissionDestroyed(accountID, &imageID),
				),
			},
			// Re-add everything so we can test when AMI disappears
			r.TestStep{
				Config: testAccOutscaleImageLaunchPermissionConfig(rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Here we delete the AMI to verify the follow-on refresh after this step
			// should not error.
			r.TestStep{
				Config: testAccOutscaleImageLaunchPermissionConfig(rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleImageDisappears(&imageID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckResourceGetAttr(name, key string, value *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s", name)
		}

		*value = is.Attributes[key]
		return nil
	}
}

func testAccOutscaleImageLaunchPermissionExists(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		if has, err := hasLaunchPermission(conn, *imageID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("launch permission does not exist for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testAccOutscaleImageLaunchPermissionDestroyed(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		if has, err := hasLaunchPermission(conn, *imageID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("launch permission still exists for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

// testAccOutscaleImageDisappears is technically a "test check function" but really it
// exists to perform a side effect of deleting an AMI out from under a resource
// so we can test that Terraform will react properly
func testAccOutscaleImageDisappears(imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		req := &fcu.DeregisterImageInput{
			ImageId: aws.String(*imageID),
		}

		err := r.Retry(5*time.Minute, func() *r.RetryError {
			var err error
			_, err = conn.VM.DeregisterImage(req)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return r.RetryableError(err)
				}
				return r.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceOutscaleImageWaitForDestroy(*imageID, conn)
	}
}

func testAccOutscaleImageLaunchPermissionConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_instance" {
    count = 1
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
}

resource "outscale_image" "outscale_image" {
    name        = "terraform test-123-%d"
    instance_id = "${outscale_vm.outscale_instance.id}"
		no_reboot   = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id    = "${outscale_image.outscale_image.image_id}"
    launch_permission {
        add {
            user_id = "520679080430"
				}
				remove {
            user_id = "520679080430"
        }
		}
}
`, r)
}
