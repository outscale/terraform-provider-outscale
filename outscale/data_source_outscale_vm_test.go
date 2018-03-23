package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVMDataSource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfig(rInt),
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

func testAccVMDataSourceConfig(r int) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	security_group = ["sg-6ed31f3e"]
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}

data "outscale_vm" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vm.basic.id}"]
	}
}

output "datasource_arch" {
	value = "${data.outscale_vm.basic_web.owner_id}"
}
`, r)
}
