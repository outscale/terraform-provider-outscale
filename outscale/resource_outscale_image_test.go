package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIImage_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()

	var ami oscgo.Image
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIImageConfigBasic(omi, "tinav4.c2r2p2", region, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIImageExists("outscale_image.foo", &ami),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "image_name", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.bsu.0.delete_on_vm_deletion", "true"),
				),
			},
		},
	})
}

func testAccCheckOAPIImageDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_image" {
			continue
		}

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadImagesResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetImages()) > 0 {
			return fmt.Errorf("Image still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOAPIImageExists(n string, ami *oscgo.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("OMI Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OMI ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadImagesResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetImages()) < 1 {
			return fmt.Errorf("Image not found (%s)", rs.Primary.ID)
		}

		ami = &resp.GetImages()[0]

		return nil
	}
}

func testAccOAPIImageConfigBasic(omi, vmType, region string, rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			placement_subregion_name = "%[3]sa"
		}
		resource "outscale_volume" "snap_volume" {
			subregion_name = "%[3]sa"
			size = 40
		}
		resource "outscale_snapshot" "snap_image" {
			volume_id = outscale_volume.snap_volume.volume_id
                }
		resource "outscale_image" "foo" {
			image_name  = "tf-testing-%d"
			description = "terraform testing"
                        root_device_name="/dev/sda1"
                        architecture = "x86_64"
                        block_device_mappings {
                          bsu  {
                            snapshot_id = outscale_snapshot.snap_image.snapshot_id
                            delete_on_vm_deletion = true
                          }
                         device_name = "/dev/sda1"
                      }
		}
	`, omi, vmType, region, rInt)
}
