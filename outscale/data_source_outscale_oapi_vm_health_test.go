package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVMHealthDataSource_basic(t *testing.T) {
	t.Skip()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMHealthDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_vm_health.web"),
					// resource.TestCheckResourceAttrSet("data.outscale_vm_health.web", "instance_states"),
				),
			},
		},
	})
}

func testAccCheckOAPIOutscaleHealthDataSourceID(n string) resource.TestCheckFunc {
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

const testAccOAPIVMHealthDataSourceConfig = `
resource "outscale_load_balancer" "bar" {
  sub_regions = ["eu-west-2a"]
	load_balancer_name = "foobar-terraform-elb"
  listener {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }

	tag {
		bar = "baz"
	}
}

resource "outscale_vm" "foo1" {
  image_id = "ami-880caa66"
	type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = [{
		vm_id = "${outscale_vm.foo1.id}"
	}]
}

data "outscale_vm_health" "web" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}`
