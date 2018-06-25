package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSLoadBalancerTags_basic(t *testing.T) {
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
				Config: testAccDSODSutscaleOAPILBUDSTagsConfig,
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

const testAccDSODSutscaleOAPILBUDSTagsConfig = `
resource "outscale_load_balancer" "bar" {
  sub_region = ["eu-west-2a"]
	load_balancer_name = "foobar-terraform-elb-aaaaaa"
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }

	tag {
		bar = "baz"
	}
}

resource "outscale_load_balancer_tags" "tags" {
	load_balancer_name = ["${outscale_load_balancer.bar.id}"]
	tag = [{
		key = "bar2" 
		value = "baz2"
	}]
}

data "outscale_load_balancer_tags" "testds" {
	load_balancer_name = ["${outscale_load_balancer.bar.id}"]
}
`
