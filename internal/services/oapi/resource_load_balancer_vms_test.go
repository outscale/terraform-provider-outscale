package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccVM_LbuBackends_Basic(t *testing.T) {
	resourceName := "outscale_load_balancer_vms.backend_test"
	sgName := acctest.RandomWithPrefix("testacc-sg")
	lbName := acctest.RandomWithPrefix("testacc-lbu")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUAttachmentConfig1(lbName, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_vm_ids.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_vm_ids.#", "2"),
				),
			},
			{
				Config: testAccLBUAttachmentAddUpdate(lbName, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_vm_ids.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_vm_ids.#", "1"),
				),
			},
			{
				Config: testAccLBUAttachmentBackendIps(lbName, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_ips.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_ips.#", "2"),
				),
			},
		},
	})
}

func TestAccVM_LbuBackends_Migration(t *testing.T) {
	sgName := acctest.RandomWithPrefix("testacc-sg")
	lbName := acctest.RandomWithPrefix("testacc-lbu")

	testacc.MigrationTest(t, "1.1.3",
		testAccLBUAttachmentConfig1(lbName, sgName),
		testAccLBUAttachmentAddUpdate(lbName, sgName),
	)
}

// add one attachment
func testAccLBUAttachmentConfig1(lbName, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = [var.subregion]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = var.image_id
  vm_type = var.vm_type
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]

  lifecycle { ignore_changes = [state] }
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1[0].vm_id, outscale_vm.foo1[1].vm_id]
}
`, lbName, sgName)
}

func testAccLBUAttachmentAddUpdate(lbName, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = [var.subregion]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = var.image_id
  vm_type = var.vm_type
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]

  lifecycle { ignore_changes = [state] }
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1[0].vm_id]
}
`, lbName, sgName)
}

func testAccLBUAttachmentBackendIps(lbName, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = [var.subregion]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = var.image_id
  vm_type = var.vm_type
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]

  lifecycle { ignore_changes = [state] }
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_ips = [outscale_vm.foo1[0].public_ip, outscale_vm.foo1[1].public_ip]
}
`, lbName, sgName)
}
