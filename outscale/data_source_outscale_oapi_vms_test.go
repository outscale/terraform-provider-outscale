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

	if oapi == false {
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
						"data.outscale_vms.basic_web", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"data.outscale_vms.basic_web", "instance_type", "t2.micro"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccOAPIVMSDataSourceConfig = `
resource "outscale_vms" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

data "outscale_vms" "basic_web" {
	filter {
    name = "instance-id"
    values = ["${outscale_vms.basic.id}"]
  }
}`
