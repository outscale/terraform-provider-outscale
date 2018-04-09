package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVMSDataSource_basic(t *testing.T) {
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
				Config: testAccVMSDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "instance_id.#", "2"),
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "reservation_set.#", "2"),
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "reservation_set.0.instances_set.0.group_set.#", "1"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccVMSDataSourceConfig = `
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_vm" "basic2" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

data "outscale_vms" "basic_web" {
	instance_id = ["${outscale_vm.basic.id}", "${outscale_vm.basic2.id}"]
}`
