package oapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_Image_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	sgName := acctest.RandomWithPrefix("testacc-sg")

	var ami oscgo.Image
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIImageConfigBasic(omi, testAccVmType, region, rInt, sgName),
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
	client := testacc.ConfiguredClient.OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_image" {
			continue
		}

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadImagesResponse
		err := retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := client.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetImages()) > 0 {
			return fmt.Errorf("image still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOAPIImageExists(n string, ami *oscgo.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("omi not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no omi id is set")
		}

		client := testacc.ConfiguredClient.OSCAPI

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{rs.Primary.ID}},
		}
		var resp oscgo.ReadImagesResponse
		err := retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := client.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetImages()) < 1 {
			return fmt.Errorf("image not found (%s)", rs.Primary.ID)
		}

		ami = &resp.GetImages()[0]

		return nil
	}
}

func testAccOAPIImageConfigBasic(omi, vmType, region string, rInt int, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_img" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"
 			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			placement_subregion_name = "%[3]sa"
			security_group_ids   = [outscale_security_group.sg_img.security_group_id]
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
	`, omi, vmType, region, rInt, sgName)
}
