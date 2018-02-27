package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleTag_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfig_basic(),
				Check:  resource.ComposeTestCheckFunc(
				// resource.TestCheckResourceAttr(),
				),
			},
		},
	})
}

func testAccTagConfig_basic() string {
	return `
resource "outscale_vm" "basic" {
	Tag_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_tag" "foo" {
	name = "tf-testing-%d"
	instance_id = "${outscale_vm.basic.id}"
}
	`
}
