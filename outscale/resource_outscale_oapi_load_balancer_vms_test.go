package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleOAPILBUAttachment_basic(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	testCheckInstanceAttached := func(count int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(conf.Instances) != count {
				return fmt.Errorf("backend_vm_id count does not match")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILBUAttachmentConfig1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(1),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPILBUAttachmentConfig2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(2),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPILBUAttachmentConfig3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(2),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPILBUAttachmentConfig4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(0),
				),
			},
		},
	})
}

// add one attachment
const testAccOutscaleOAPILBUAttachmentConfig1 = `
resource "outscale_load_balancer" "bar" {
	load_balancer_name = "load-test"
	
	availability_zones = ["eu-west-2a"]
    listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}

resource "outscale_vm" "foo1" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = ["${outscale_vm.foo1.id}"]
}
`

// add a second attachment
const testAccOutscaleOAPILBUAttachmentConfig2 = `
resource "outscale_load_balancer" "bar" {
	load_balancer_name = "load-test"
  availability_zones = ["eu-west-2a"]

    listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}

resource "outscale_vm" "foo1" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_vm" "foo2" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = ["${outscale_vm.foo1.id}"]
}

resource "outscale_load_balancer_vms" "foo2" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = ["${outscale_vm.foo2.id}"]
}
`

// swap attachments between resources
const testAccOutscaleOAPILBUAttachmentConfig3 = `
resource "outscale_load_balancer" "bar" {
	load_balancer_name = "load-test"
  availability_zones = ["eu-west-2a"]

    listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}

resource "outscale_vm" "foo1" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_vm" "foo2" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = ["${outscale_vm.foo2.id}"]
}

resource "outscale_load_balancer_vms" "foo2" {
  load_balancer_name      = "${outscale_load_balancer.bar.id}"
  backend_vm_id = ["${outscale_vm.foo1.id}"]
}
`

// destroy attachments
const testAccOutscaleOAPILBUAttachmentConfig4 = `
resource "outscale_load_balancer" "bar" {
	load_balancer_name = "load-test"
  availability_zones = ["eu-west-2a"]

    listeners {
    instance_port = 8000
    instance_protocol = "HTTP"
    load_balancer_port = 80
    protocol = "HTTP"
  }
}
`
