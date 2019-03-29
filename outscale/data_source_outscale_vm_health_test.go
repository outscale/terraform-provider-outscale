package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVMHealthDataSource_basic(t *testing.T) {
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
				Config: testAccVMHealthDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_vm_health.web"),
					// resource.TestCheckResourceAttrSet("data.outscale_vm_health.web", "instance_states"),
				),
			},
		},
	})
}

func testAccCheckOutscaleHealthDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vm health data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vm health data source ID not set")
		}
		return nil
	}
}

const testAccVMHealthDataSourceConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name = "foobar-terraform-elb"
  listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }

	tag {
		bar = "baz"
	}
}

resource "outscale_vm" "foo1" {
  image_id = "ami-880caa66"
	instance_type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  instances = [{
		instance_id = "${outscale_vm.foo1.id}"
	}]
}

data "outscale_vm_health" "web" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}`
