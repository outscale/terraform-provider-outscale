package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleImageTask_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleImageTaskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleImageTaskExists("outscale_image_tasks.outscale_image_tasks"),
				),
			},
		},
	})
}

func testAccCheckOutscaleImageTaskExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image task id is set")
		}

		return nil
	}
}

var testAccOutscaleImageTaskConfig = `
resource "outscale_vm" "outscale_vm" {
    count = 1

    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"

}

resource "outscale_image" "outscale_image" {
    name            = "image_${outscale_vm.outscale_vm.id}"
    instance_id     = "${outscale_vm.outscale_vm.id}"
}

resource "outscale_image_tasks" "outscale_image_tasks" {
    count = 1

    image_id = "${outscale_image.outscale_image.image_id}"
}
`
