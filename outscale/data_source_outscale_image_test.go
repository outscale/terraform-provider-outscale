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

func TestAccOutscaleImageDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleImageDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImageDataSourceID("data.outscale_image.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "block_device_mappings.#", "1"),
					resource.TestMatchResourceAttr("data.outscale_image.nat_ami", "image_id", regexp.MustCompile("^ami-")),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "image_type", "machine"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "product_codes.#", "0"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "root_device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "root_device_type", "ebs"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "image_state", "available"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "state_reason.code", "UNSET"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "state_reason.message", "UNSET"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "tag_set.#", "0"),
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
		name = "name"
		values = ["tf-testing-3273*"]
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
