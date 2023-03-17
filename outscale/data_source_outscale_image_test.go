package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Image_DataSource(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	imageName := fmt.Sprintf("image-test-%d", acctest.RandInt())

	dataSourceName := "data.outscale_image.omi"
	dataSourcesName := "data.outscale_images.omis"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Image_DataSource_Config(omi, "tinav4.c2r2p2", region, imageName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet(dataSourceName, "image_id"),
					resource.TestCheckResourceAttr(dataSourcesName, "images.#", "2"),
				),
			},
		},
	})
}

func testAcc_Image_DataSource_Config(omi, vmType, region, imageName string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic_one" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%[3]sa"
		}
		
		resource "outscale_vm" "basic_two" {
			image_id			     = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	         = "terraform-basic"
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_image" "image_one" {
			image_name = "%[4]s-one"
			vm_id = outscale_vm.basic_one.id
		}

		resource "outscale_image" "image_two" {
			image_name = "%[4]s-two"
			vm_id = outscale_vm.basic_two.id
		}
		
		data "outscale_image" "omi" {
			filter {
				name   = "image_ids"
				values = [outscale_image.image_one.id]
			}
		}

		data "outscale_images" "omis" {
			filter {
				name = "image_ids"
				values = [outscale_image.image_one.id, outscale_image.image_two.id]
			}
		}
	`, omi, vmType, region, imageName)
}
