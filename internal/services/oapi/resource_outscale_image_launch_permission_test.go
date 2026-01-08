package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccVM_WithImageLaunchPermission_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	imageID := ""

	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Scaffold everything
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, testAccVmType, region, accountID, keypair, true, rInt, sgName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
			// Drop just launch permission to test destruction
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, testAccVmType, region, accountID, keypair, false, rInt, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleImageLaunchPermissionDestroyed(accountID, &imageID),
				),
			},
			// Re-add everything so we can test when AMI disappears
			{
				Config: testAccOutscaleImageLaunchPermissionConfig(omi, testAccVmType, region, accountID, keypair, true, rInt, sgName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
		},
	})
}

func TestAccVM_ImageLaunchPermissionDestruction_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	var imageID string
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Scaffold everything
			{
				Config: testAccOutscaleImageLaunchPermissionCreateConfig(omi, testAccVmType, region, keypair, rInt, true, false, sgName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
					testAccOutscaleImageLaunchPermissionExists(accountID, &imageID),
				),
			},
			{
				Config: testAccOutscaleImageLaunchPermissionCreateConfig(omi, testAccVmType, region, keypair, rInt, true, true, sgName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceOAPILPIGetAttr("outscale_image.outscale_image", "id", &imageID),
				),
			},
		},
	})
}

func testCheckResourceOAPILPIGetAttr(name, key string, value *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary == nil {
			return fmt.Errorf("no primary instance: %s", name)
		}

		*value = rs.Primary.Attributes[key]

		return nil
	}
}

func testAccOutscaleImageLaunchPermissionExists(accountID string, imageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testacc.ConfiguredClient.OSCAPI
		if has, err := oapihelpers.ImageHasLaunchPermission(client, *imageID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("launch permission does not exist for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testAccOutscaleImageLaunchPermissionDestroyed(accountID string, imageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testacc.ConfiguredClient.OSCAPI
		if has, err := oapihelpers.ImageHasLaunchPermission(client, *imageID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("launch permission still exists for '%s' on '%s'", accountID, *imageID)
		}
		return nil
	}
}

func testCheckResourceGetAttr(name, key string, value *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("no primary instance: %s", name)
		}

		*value = is.Attributes[key]
		return nil
	}
}

func testAccOutscaleImageLaunchPermissionConfig(omi, vmType, region, accountID, keypair string, includeLaunchPermission bool, r int, sgName string) string {
	base := fmt.Sprintf(`
		resource "outscale_security_group" "sg_perm" {
			security_group_name = "%[6]s"
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
	`, omi, vmType, region, r, keypair, sgName)

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

func testAccOutscaleImageLaunchPermissionCreateConfig(omi, vmType, region, keypair string, r int, includeAddtion, includeRemoval bool, sgName string) string {
	base := fmt.Sprintf(`
		resource "outscale_security_group" "sg_perm" {
			security_group_name = "%[6]s"
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
	`, omi, vmType, region, r, keypair, sgName)

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
