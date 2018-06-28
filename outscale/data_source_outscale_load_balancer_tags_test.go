package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleDSLoadBalancerTags_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	r := acctest.RandString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getTestAccDSODSutscaleLBUDSTagsConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckODSutscaleLBUDSTagsExists("data.outscale_load_balancer_tags.testds"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer_tags.testds", "tag_descriptions.#", "1"),
				)},
		},
	})
}

func testAccCheckODSutscaleLBUDSTagsExists(n string) resource.TestCheckFunc {
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

func getTestAccDSODSutscaleLBUDSTagsConfig(r string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name = "foobar-terraform-elb-%s"
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

resource "outscale_load_balancer_tags" "tags" {
	load_balancer_names = ["${outscale_load_balancer.bar.id}"]
	tags = [{
		key = "bar2" 
		value = "baz2"
	}]
}

data "outscale_load_balancer_tags" "testds" {
	load_balancer_names = ["${outscale_load_balancer.bar.id}"]
}
`, r)
}
