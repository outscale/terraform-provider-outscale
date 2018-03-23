package outscale

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleImagesDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleImagesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImagesDataSourceID("data.outscale_images.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.description", "Debian 9 - 4.9.51"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.block_device_mappings.#", "1"),
					resource.TestMatchResourceAttr("data.outscale_images.nat_ami", "image_set.0.image_id", regexp.MustCompile("^ami-")),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.image_type", "machine"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.is_public", "true"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.root_device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.root_device_type", "ebs"),
					resource.TestCheckResourceAttr("data.outscale_images.nat_ami", "image_set.0.image_state", "available"),
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

const testAccCheckOutscaleImagesDataSourceConfig = `
data "outscale_images" "nat_ami" {
	filter {
		name = "architecture"
		values = ["x86_64"]
	}
	filter {
		name = "virtualization-type"
		values = ["hvm"]
	}
	filter {
		name = "root-device-type"
		values = ["ebs"]
	}
	filter {
		name = "block-device-mapping.volume-type"
		values = ["standard"]
	}
}
`
