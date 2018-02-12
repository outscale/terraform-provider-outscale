package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVMDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "instance_type", "t2.micro"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccVMDataSourceConfig = `
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
  }
}`
