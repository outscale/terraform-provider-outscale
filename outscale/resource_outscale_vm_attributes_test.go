package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVMAttr_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm_attributes.outscale_vm_attributes", "ebs_optimized", "true"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVMAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return nil
	}
}

func testAccCheckOutscaleVMConfig_basic() string {
	return `
resource "outscale_vm" "outscale_vm" {
  count = 1

  image_id                = "ami-880caa66"
  instance_type           = "c4.large"
  ebs_optimized = false
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
  instance_id             = "${outscale_vm.outscale_vm.0.id}"
  attribute               = "ebsOptimized"
  ebs_optimized = true
}`
}
