package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVMDataSource_basic(t *testing.T) {
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
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
	}
}

output "datasource_arch" {
	value = "${data.outscale_vm.basic_web.owner_id}"
}`
