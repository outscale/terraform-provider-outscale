package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIImageLaunchPermission_Basic(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	imageID := ""
	accountID := "520679080430"

	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []r.TestStep{
			// Scaffold everything
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionConfig(omi, "c4.large", region, accountID, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Drop just launch permission to test destruction
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionConfig(omi, "c4.large", region, accountID, false, rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleOAPIImageLaunchPermissionDestroyed(accountID, &imageID),
				),
			},
			// Re-add everything so we can test when AMI disappears
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionConfig(omi, "c4.large", region, accountID, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Here we delete the AMI to verify the follow-on refresh after this step
			// should not error.
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionConfig(omi, "c4.large", region, accountID, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleOAPIImageDisappears(&imageID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccOutscaleOAPIImageLaunchPermissionDestruction_Basic(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	var imageID string
	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []r.TestStep{
			// Scaffold everything
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionCreateConfig(omi, "c4.large", region, rInt, true, false),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			r.TestStep{
				Config: testAccOutscaleOAPIImageLaunchPermissionCreateConfig(omi, "c4.large", region, rInt, true, true),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
		},
	})
}

func testCheckResourceOAPILPIGetAttr(name, key string, value *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary == nil {
			return fmt.Errorf("No primary instance: %s", name)
		}

		*value = rs.Primary.Attributes[key]

		return nil
	}
}

func testAccOutscaleOAPIImageLaunchPermissionExists(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		if has, err := hasOAPILaunchPermission(conn, *imageID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("launch permission does not exist for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testAccOutscaleOAPIImageLaunchPermissionDestroyed(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		if has, err := hasOAPILaunchPermission(conn, *imageID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("launch permission still exists for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

// testAccOutscaleOAPIImageDisappears is technically a "test check function" but really it
// exists to perform a side effect of deleting an AMI out from under a resource
// so we can test that Terraform will react properly
func testAccOutscaleOAPIImageDisappears(imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		req := oscgo.DeleteImageRequest{
			ImageId: aws.StringValue(imageID),
		}

		err := r.Retry(5*time.Minute, func() *r.RetryError {
			var err error
			_, _, err = conn.ImageApi.DeleteImage(context.Background(), &oscgo.DeleteImageOpts{DeleteImageRequest: optional.NewInterface(req)})
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

		return resourceOutscaleOAPIImageWaitForDestroy(*imageID, conn)
	}
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

func testAccOutscaleOAPIImageLaunchPermissionConfig(omi, vmType, region, accountID string, includeLaunchPermission bool, r int) string {
	base := fmt.Sprintf(`
		resource "outscale_vm" "outscale_instance" {
			image_id           = "%s"
			vm_type            = "%s"
			keypair_name       = "terraform-basic"
			security_group_ids = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}
		
		resource "outscale_image" "outscale_image" {
			image_name        = "terraform test-123-%d"
			vm_id = "${outscale_vm.outscale_instance.id}"
			no_reboot   = "true"
		}
	`, omi, vmType, region, r)

	if !includeLaunchPermission {
		return base
	}

	return base + fmt.Sprintf(`
		resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
			image_id    = "${outscale_image.outscale_image.image_id}"
			permission_additions {
				account_ids = ["%s"]
			}
		}
	`, accountID)
}

func testAccOutscaleOAPIImageLaunchPermissionCreateConfig(omi, vmType, region string, r int, includeAddtion, includeRemoval bool) string {
	base := fmt.Sprintf(`
		resource "outscale_vm" "outscale_instance" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}
		
		resource "outscale_image" "outscale_image" {
			image_name = "terraform test-123-%d"
			vm_id      = "${outscale_vm.outscale_instance.id}"
			no_reboot  = "true"
		}
	`, omi, vmType, region, r)

	if includeAddtion {
		return base + `
			resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
				image_id = "${outscale_image.outscale_image.image_id}"
			
				permission_additions {
					account_ids = ["520679080430"]
				}
			}
		`
	}

	if includeRemoval {
		return base + `
			resource "outscale_image_launch_permission" "outscale_image_launch_permission_two" {
				image_id = "${outscale_image_launch_permission.outscale_image_launch_permission.image_id}"

				permission_removals {
					account_ids = ["520679080430"]
				}
			}
		`
	}
	return base
}
