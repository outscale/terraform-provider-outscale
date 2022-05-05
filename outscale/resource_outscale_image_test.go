package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIImage_basic(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

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
						"outscale_image.foo", "block_device_mappings.0.bsu.delete_on_vm_deletion", "true"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "state_comment.state_code", ""),
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

		resp, _, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
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

		resp, _, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(filterReq).Execute()
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
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name		 = "terraform-basic"
			placement_subregion_name = "%sa"
		}

		resource "outscale_image" "foo" {
			image_name  = "tf-testing-%d"
			vm_id       = outscale_vm.basic.id
			no_reboot   = "true"
			description = "terraform testing"
		}
	`, omi, vmType, region, rInt)
}
