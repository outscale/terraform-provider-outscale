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

func TestAccOutscaleOAPILBUTags_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILBUTagsConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUTagsExists("outscale_load_balancer_tags.tags"),
				)},
		},
	})
}

func testAccCheckOutscaleOAPILBUTagsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU ID is set")
		}

		return nil
	}
}

func testAccOutscaleOAPILBUTagsConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  sub_regions = ["eu-west-2a"]
	load_balancer_name = "foobar-terraform-elb-%d"
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
`, r)
}
