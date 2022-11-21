package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/aws/aws-sdk-go/aws"
	r "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccImageLaunchPermission_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	imageID := ""

	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []r.TestStep{
			// Scaffold everything
			{
				Config: testAccImageLaunchPermissionConfig(omi, "tinav4.c2r2p2", region, accountID, true, rInt, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceLPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Drop just launch permission to test destruction
			{
				Config: testAccImageLaunchPermissionConfig(omi, "tinav4.c2r2p2", region, accountID, false, rInt, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testAccImageLaunchPermissionDestroyed(accountID, &imageID),
				),
			},
			// Re-add everything so we can test when AMI disappears
			{
				Config: testAccImageLaunchPermissionConfig(omi, "tinav4.c2r2p2", region, accountID, true, rInt, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceLPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Here we delete the AMI to verify the follow-on refresh after this step
			// should not error.
			{
				Config: testAccImageLaunchPermissionConfig(omi, "tinav4.c2r2p2", region, accountID, true, rInt, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testAccImageDisappears(&imageID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccImageLaunchPermissionDestruction_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	var imageID string
	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []r.TestStep{
			// Scaffold everything
			{
				Config: testAccImageLaunchPermissionCreateConfig(omi, "tinav4.c2r2p2", region, rInt, true, false, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceLPIGetAttr("outscale_image.outscale_image", "id", &imageID),
					testAccImageLaunchPermissionExists(accountID, &imageID),
				),
			},
			{
				Config: testAccImageLaunchPermissionCreateConfig(omi, "tinav4.c2r2p2", region, rInt, true, true, keypair, sgId),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceLPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
		},
	})
}

func testCheckResourceLPIGetAttr(name, key string, value *string) r.TestCheckFunc {
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

func testAccImageLaunchPermissionExists(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*Client).OSCAPI
		if has, err := hasLaunchPermission(conn, *imageID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("launch permission does not exist for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testAccImageLaunchPermissionDestroyed(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*Client).OSCAPI
		if has, err := hasLaunchPermission(conn, *imageID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("launch permission still exists for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

// testAccImageDisappears is technically a "test check function" but really it
// exists to perform a side effect of deleting an AMI out from under a resource
// so we can test that Terraform will react properly
func testAccImageDisappears(imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*Client).OSCAPI
		req := oscgo.DeleteImageRequest{
			ImageId: aws.StringValue(imageID),
		}

		err := r.Retry(5*time.Minute, func() *r.RetryError {
			var err error
			_, httpResp, err := conn.ImageApi.DeleteImage(context.Background()).DeleteImageRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceImageWaitForDestroy(*imageID, conn)
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

func testAccImageLaunchPermissionConfig(omi, vmType, region, accountID string, includeLaunchPermission bool, r int, keypair, sgId string) string {
	base := fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = "${outscale_net.net.id}"
		}

		resource "outscale_vm" "outscale_instance" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[5]s"
			security_group_ids = ["%[6]s"]
			placement_subregion_name = "%[3]sa"
		}
		
		resource "outscale_image" "outscale_image" {
			image_name        = "terraform test-123-%[4]d"
			vm_id = "${outscale_vm.outscale_instance.id}"
			no_reboot   = "true"
		}
	`, omi, vmType, region, r, keypair, sgId)

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

func testAccImageLaunchPermissionCreateConfig(omi, vmType, region string, r int, includeAddtion, includeRemoval bool, keypair, sgId string) string {
	base := fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = "${outscale_net.net.id}"
		}

		resource "outscale_vm" "outscale_instance" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[5]s"
			security_group_ids       = ["%[6]s"]
			placement_subregion_name = "%[3]sa"
		}
		
		resource "outscale_image" "outscale_image" {
			image_name = "terraform test-123-%[4]d"
			vm_id      = "${outscale_vm.outscale_instance.id}"
			no_reboot  = "true"
		}
	`, omi, vmType, region, r, keypair, sgId)

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
