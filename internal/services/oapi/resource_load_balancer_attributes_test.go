package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_LBUAttributes_Basic(t *testing.T) {
	resourceName := "outscale_load_balancer_attributes.lbuattr"
	name := acctest.RandomWithPrefix("test-lbuattr")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUAttributesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health_check.#", "1"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnoresWith("access_log", "health_check")...),
		},
	})
}

func TestAccOthers_LBUAttributes_(t *testing.T) {
	resourceName := "outscale_load_balancer_attributes.lbuattr"
	name := acctest.RandomWithPrefix("test-lbuattr")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccLBUAttributesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health_check.#", "1"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnoresWith("access_log", "health_check")...),
		},
	})
}

func testAccLBUAttributesConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
  subregion_names = [var.subregion]
  load_balancer_name       = "%s"
  listeners {
    backend_port           = 8000
    backend_protocol       = "HTTP"
    load_balancer_port     = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_attributes" "lbuattr" {
	access_log {
		is_enabled = false
		osu_bucket_prefix = "donustestbucket"
	}

	health_check  {
        healthy_threshold   = 10
        check_interval      = 30
        path                = "/index.html"
        port                = 8080
        protocol            = "HTTPS"
        timeout             = 5
        unhealthy_threshold = 5
    }

	load_balancer_name = outscale_load_balancer.lbu.id
}
`, name)
}

func testAccLBUAttributesAccessLogConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
  subregion_names = [var.subregion]
  load_balancer_name       = "%s"
  listeners {
    backend_port           = 8000
    backend_protocol       = "HTTP"
    load_balancer_port     = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_attributes" "lbuattr" {
	access_log {
		publication_interval = 5
		is_enabled = false
		osu_bucket_prefix = "donustestbucket"
	}

	load_balancer_name = outscale_load_balancer.lbu.id
}
`, name)
}

func testAccLBUAttributesHealthCheckConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu" {
  subregion_names = [var.subregion]
  load_balancer_name       = "%s"
  listeners {
    backend_port           = 8000
    backend_protocol       = "HTTP"
    load_balancer_port     = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_attributes" "lbuattr" {
	health_check  {
        healthy_threshold   = 10
        check_interval      = 30
        path                = "/index.html"
        port                = 8080
        protocol            = "HTTPS"
        timeout             = 5
        unhealthy_threshold = 5
    }

	load_balancer_name = outscale_load_balancer.lbu.id
}
`, name)
}
