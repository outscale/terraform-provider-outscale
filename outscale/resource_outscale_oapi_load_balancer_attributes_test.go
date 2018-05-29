package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleOAPILBUAttr_basic(t *testing.T) {
	var conf lbu.LoadBalancerAttributes

	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer_attributes.bar2",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILBUAttrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUAttrExists("outscale_load_balancer_attributes.bar2", &conf),
				)},
		},
	})
}

func testAccCheckOutscaleOAPILBUAttrExists(n string, res *lbu.LoadBalancerAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU Attr ID is set")
		}

		return nil
	}
}

func testAccOutscaleOAPILBUAttrConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  sub_regions = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-%d"
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }

	tags {
		bar = "baz"
	}

}

resource "outscale_load_balancer_attributes" "bar2" {
	is_enabled = "false"
			osu_bucket_prefix = "donustestbucket"
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`, r)
}
