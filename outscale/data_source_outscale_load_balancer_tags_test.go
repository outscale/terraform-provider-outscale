package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIDSLoadBalancerTags_basic(t *testing.T) {
	r := acctest.RandString(4)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getTestAccDSODSutscaleOAPILBUDSTagsConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckODSutscaleOAPILBUDSTagsExists("data.outscale_load_balancer_tags.testds"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer_tags.testds", "tag.#", "1"),
				)},
		},
	})
}

func testAccCheckODSutscaleOAPILBUDSTagsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU Tag DS ID is set")
		}

		return nil
	}
}

func getTestAccDSODSutscaleOAPILBUDSTagsConfig(r string) string {
	return fmt.Sprintf(`
		resource "outscale_load_balancer" "bar" {
			subregion_names    = ["eu-west-2a"]
			load_balancer_name = "foobar-terraform-elb-%s"
		
			listeners {
				backend_port           = 8000
				backend_protocol       = "HTTP"
				load_balancer_port     = 80
				load_balancer_protocol = "HTTP"
			}
		
			tag {
				bar = "baz"
			}
		}
		
		resource "outscale_load_balancer_tags" "tags" {
			load_balancer_name = ["${outscale_load_balancer.bar.id}"]
		
			tag = [{
				key   = "bar2"
				value = "baz2"
			}]
		}
		
		data "outscale_load_balancer_tags" "testds" {
			load_balancer_name = ["${outscale_load_balancer.bar.id}"]
		}
	`, r)
}
