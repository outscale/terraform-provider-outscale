package oapi_test

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_LBU_Basic(t *testing.T) {
	lbResourceName := "outscale_load_balancer.lbu"
	name := acctest.RandomWithPrefix("test-lbu")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lbResourceName, "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(lbResourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttr(lbResourceName, "secured_cookies", "true"),
				),
			},
			testacc.ImportStep(lbResourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_LBU_PublicIp(t *testing.T) {
	resourceName := "outscale_load_balancer.lbu"
	name := acctest.RandomWithPrefix("test-lbu")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUPublicIpConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttr(resourceName, "secured_cookies", "false"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_LBU_UpdateListeners(t *testing.T) {
	resourceName := "outscale_load_balancer.lbu"
	name := acctest.RandomWithPrefix("test-lbu")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUConfigListeners(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "secured_cookies", "false"),
				),
			},
			{
				Config: testAccLBUConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "listeners.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "secured_cookies", "false"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func testAccLBUConfigListeners(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
  load_balancer_name = "%s"
  subregion_names    = [var.subregion]
  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
  listeners {
    backend_port           = 8080
    backend_protocol       = "HTTP"
    load_balancer_protocol = "HTTP"
    load_balancer_port     = 8080
  }
  tags {
    key   = "testacc"
    value = "v1"
  }
}
`, name)
}

func testAccLBUConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
  load_balancer_name = "%s"
  subregion_names    = [var.subregion]
  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
  listeners {
    backend_port           = 90
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 90
  }
  tags {
    key   = "testacc"
    value = "v2"
  }
}
`, name)
}

func TestAccOthers_LBU_Migration(t *testing.T) {
	name := acctest.RandomWithPrefix("test-lbu")

	testacc.MigrationTest(t, "1.6.0",
		testAccLBUConfig(name),
		testAccLBUPublicIpConfig(name),
	)
}

func testAccLBUConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
	subregion_names = [var.subregion]
	load_balancer_name = "%s"
	secured_cookies = true

	listeners {
		backend_port = 8000
		backend_protocol = "HTTP"
		load_balancer_port = 80
		load_balancer_protocol = "HTTP"
	}

	tags {
		key = "testacc"
		value = "lbu-basic"
	}
}
`, name)
}

func testAccLBUPublicIpConfig(name string) string {
	return fmt.Sprintf(`
	resource "outscale_public_ip" "my_public_ip" {
	}

	resource "outscale_load_balancer" "lbu" {
		subregion_names = [var.subregion]
		load_balancer_name = "%s"

		listeners {
		  backend_port           = 80
		  backend_protocol       = "HTTP"
		  load_balancer_protocol = "HTTP"
		  load_balancer_port     = 80
		}

		public_ip = outscale_public_ip.my_public_ip.public_ip

		tags {
		  key = "testacc"
		  value = "lbu-public-ip"
		}
	  }
`, name)
}
