package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccVM_WithImagesDataSource_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	imageName := fmt.Sprintf("image-test-%d", acctest.RandInt())
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleImagesDataSourceConfig(omi, utils.TestAccVmType, utils.GetRegion(), imageName, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImagesDataSourceID("data.outscale_images.nat_ami"),
				),
			},
		},
	})
}

func testAccCheckOutscaleImagesDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find AMI data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AMI data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleImagesDataSourceConfig(omi, vmType, region, imageName, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_imgs_data" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic_one" {
			image_id			           = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	           = "terraform-basic"
			placement_subregion_name = "%[3]sa"
			security_group_ids = [outscale_security_group.sg_imgs_data.security_group_id]
		}

		resource "outscale_vm" "basic_two" {
			image_id			           = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	           = "terraform-basic"
			placement_subregion_name = "%[3]sa"
			security_group_ids = [outscale_security_group.sg_imgs_data.security_group_id]
		}

		resource "outscale_image" "image_one" {
			image_name = "%[4]s-one"
			vm_id = outscale_vm.basic_one.id
		}

		resource "outscale_image" "image_two" {
			image_name = "%[4]s-two"
			vm_id = outscale_vm.basic_two.id
		}

		data "outscale_images" "nat_ami" {
			filter {
				name = "image_ids"
				values = [outscale_image.image_one.id, outscale_image.image_two.id]
			}
		}
	`, omi, vmType, region, imageName, sgName)
}
