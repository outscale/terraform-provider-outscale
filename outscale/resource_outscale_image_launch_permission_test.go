package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_WithImageLaunchPermission_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := "terraform-basic"

	imageID := ""

	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []r.TestStep{
			// Scaffold everything
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, utils.TestAccVmType, region, accountID, keypair, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Drop just launch permission to test destruction
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, utils.TestAccVmType, region, accountID, keypair, false, rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleImageLaunchPermissionDestroyed(accountID, &imageID),
				),
			},
			// Re-add everything so we can test when AMI disappears
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, utils.TestAccVmType, region, accountID, keypair, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Here we delete the AMI to verify the follow-on refresh after this step
			// should not error.
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, utils.TestAccVmType, region, accountID, keypair, true, rInt),
				Check: r.ComposeTestCheckFunc(
					testAccOutscaleImageDisappears(&imageID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVM_ImageLaunchPermissionDestruction_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := "terraform-basic"

	var imageID string
	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []r.TestStep{
			// Scaffold everything
			{
				Config: testAccOutscaleImageLaunchPermissionCreateConfig(omi, utils.TestAccVmType, region, keypair, rInt, true, false),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
					testAccOutscaleImageLaunchPermissionExists(accountID, &imageID),
				),
			},
			{
				Config: testAccOutscaleImageLaunchPermissionCreateConfig(omi, utils.TestAccVmType, region, keypair, rInt, true, true),
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

func testAccOutscaleImageLaunchPermissionExists(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccConfiguredClient.OSCAPI
		if has, err := hasOAPILaunchPermission(client, *imageID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("launch permission does not exist for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testAccOutscaleImageLaunchPermissionDestroyed(accountID string, imageID *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccConfiguredClient.OSCAPI
		if has, err := hasOAPILaunchPermission(client, *imageID); err != nil {
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
		client := testAccConfiguredClient.OSCAPI
		req := oscgo.DeleteImageRequest{
			ImageId: aws.StringValue(imageID),
		}

		err := r.Retry(5*time.Minute, func() *r.RetryError {
			var err error
			_, httpResp, err := client.ImageApi.DeleteImage(context.Background()).DeleteImageRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return ResourceOutscaleImageWaitForDestroy(*imageID, client, 5*utils.DeleteDefaultTimeout)
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

func testAccOutscaleImageLaunchPermissionConfig(omi, vmType, region, accountID, keypair string, includeLaunchPermission bool, r int) string {
	base := fmt.Sprintf(`
		resource "outscale_security_group" "sg_perm" {
			security_group_name = "sgLPerm"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "outscale_instance" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[5]s"
			security_group_ids = [outscale_security_group.sg_perm.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_image" "outscale_image" {
			image_name        = "terraform test-123-%[4]d"
			vm_id = outscale_vm.outscale_instance.vm_id
			no_reboot   = "true"
		}
	`, omi, vmType, region, r, keypair)

	if !includeLaunchPermission {
		return base
	}

	return base + fmt.Sprintf(`
		resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
			image_id    = outscale_image.outscale_image.image_id
			permission_additions {
				account_ids = ["%s"]
			}
		}
	`, accountID)
}

func testAccOutscaleImageLaunchPermissionCreateConfig(omi, vmType, region, keypair string, r int, includeAddtion, includeRemoval bool) string {
	base := fmt.Sprintf(`
		resource "outscale_security_group" "sg_perm" {
			security_group_name = "sgLPerm"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}
		resource "outscale_vm" "outscale_instance" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[5]s"
			security_group_ids = [outscale_security_group.sg_perm.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_image" "outscale_image" {
			image_name = "terraform test-123-%[4]d"
			vm_id      = outscale_vm.outscale_instance.vm_id
			no_reboot  = "true"
		}
	`, omi, vmType, region, r, keypair)

	if includeAddtion {
		return base + `
				resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
					image_id = outscale_image.outscale_image.image_id

					permission_additions {
						account_ids = ["520679080430"]
					}
				}
			`
	}

	if includeRemoval {
		return base + `
					resource "outscale_image_launch_permission" "outscale_image_launch_permission_two" {
						image_id = outscale_image_launch_permission.outscale_image_launch_permission.image_id

						permission_removals {
							account_ids = ["520679080430"]
						}
					}
				`
	}
	return base
}
