package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_LBUListenerRule_Basic(t *testing.T) {
	resourceName := "outscale_load_balancer_listener_rule.lburule"
	name := acctest.RandomWithPrefix("test-lburule")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUListenerRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "listener_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "listener.#", "1"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_LBUListenerRule_Migration(t *testing.T) {
	name := acctest.RandomWithPrefix("test-lburule")

	testacc.MigrationTest(t, "1.7.0", testAccLBUListenerRuleConfig(name))
}

func testAccLBUListenerRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_vm" "vm" {
  image_id     = var.image_id
  vm_type      = var.vm_type
  keypair_name = var.keypair_name
}

resource "outscale_load_balancer" "lbu" {
  load_balancer_name = "%s"
  subregion_names    = [var.subregion]

  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms01" {
  load_balancer_name = outscale_load_balancer.lbu.id
  backend_vm_ids     = [outscale_vm.vm.vm_id]
}

resource "outscale_load_balancer_listener_rule" "lburule" {
  listener {
    load_balancer_name = outscale_load_balancer.lbu.id
    load_balancer_port = 80
  }

  listener_rule {
    action             = "forward"
    listener_rule_name = "terraform-listener-rule"
    host_name_pattern  = "testhost"
    path_pattern       = "testpath"
    priority           = 10
  }

  vm_ids = [outscale_vm.vm.id]
}
`, name)
}
