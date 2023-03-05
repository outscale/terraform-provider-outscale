package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccVM_WithImageDataSource_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	imageName := fmt.Sprintf("image-test-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIImageDataSourceBasicConfig(omi, "tinav4.c2r2p2", utils.GetRegion(), imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImageDataSourceID("data.outscale_image.omi"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIImageDataSourceID(n string) resource.TestCheckFunc {
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

func testAccCheckOutscaleOAPIImageDataSourceBasicConfig(omi, vmType, region, imageName string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basicIm" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%[3]sa"
		}
		
		resource "outscale_image" "image" {
			image_name = "%[4]s"
			vm_id      = outscale_vm.basicIm.id
		}
		
		data "outscale_image" "omi" {
			filter {
				name   = "image_ids"
				values = [outscale_image.image.id]
			}
		}
	`, omi, vmType, region, imageName)
}
