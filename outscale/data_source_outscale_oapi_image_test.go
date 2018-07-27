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

func TestAccOutscaleOAPIImageDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIImageDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImageDataSourceID("data.outscale_image.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "block_device_mapping.#", "1"),
					resource.TestMatchResourceAttr("data.outscale_image.nat_ami", "image_id", regexp.MustCompile("^ami-")),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "type", "machine"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "product_code.#", "0"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "root_device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "root_device_type", "ebs"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "state", "available"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "state_comment.code", "UNSET"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "state_comment.message", "UNSET"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "tag.#", "0"),
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

const testAccCheckOutscaleOAPIImageDataSourceConfig = `
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_1"
  group_description = "Used in the terraform acceptance tests"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	type = "t2.micro"
	security_group = ["${outscale_firewall_rules_set.web.id}"]
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}

resource "outscale_image" "foo" {
	name = "tf-testing-%d"
	vm_id = "${outscale_vm.basic.id}"
}

data "outscale_image" "nat_ami" {
	image_id = "${outscale_image.foo.id}"
}
`
