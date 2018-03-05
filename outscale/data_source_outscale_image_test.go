package outscale

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleImageDataSource_Instance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleImageDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImageDataSourceID("data.outscale_image.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.description", "Debian 9 - 4.9.51"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.block_device_mappings.#", "1"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.hypervisor", "xen"),
					resource.TestMatchResourceAttr("data.outscale_image.nat_ami", "images_set.0.image_id", regexp.MustCompile("^ami-")),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.image_type", "machine"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.is_public", "true"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.root_device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.root_device_type", "ebs"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.image_state", "available"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "images_set.0.virtualization_type", "hvm"),
				),
			},
		},
	})
}

func testAccCheckOutscaleImageDataSourceID(n string) resource.TestCheckFunc {
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

const testAccCheckOutscaleImageDataSourceConfig = `
data "outscale_image" "nat_ami" {
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
