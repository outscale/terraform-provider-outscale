package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIVMSDataSource_basic(t *testing.T) {
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
				Config: testAccOAPIVMSDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"data.outscale_vm.basic_web", "type", "t2.micro"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccOAPIVMSDataSourceConfig = `
resource "outscale_vm" "basic" {
	image_id               = "ami-5c450b62"
	vm_type                = "c4.large"
	keypair_name           = "testkp"
	security_group_ids     = ["sg-9752b7a6"]
}

data "outscale_vms" "basic_web" {
	filter {
    name = "vm_ids"
    values = ["${outscale_vm.basic.id}"]
  }
}`
