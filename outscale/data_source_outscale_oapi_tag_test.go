package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPITagDataSource(t *testing.T) {
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
				Config: testAccOAPITagDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "key", "foo"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "value", "bar"),
					resource.TestCheckResourceAttr(
						"data.outscale_tag.web", "resource_type", "instance"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccOAPITagDataSourceConfig = `
resource "outscale_vm" "basic" {
  image_id = "ami-8a6a0120"
	instance_type = "m1.small"
	tag = {
		foo = "bar"
	}
}

data "outscale_tag" "web" {
	filter {
    name = "resource-id"
    values = ["${outscale_vm.basic.id}"]
	}
}`
