package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleLBUAttr_basic(t *testing.T) {
	var conf lbu.LoadBalancerAttributes

	r := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer_attributes.bar2",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUAttrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckState("outscale_load_balancer_attributes.bar2"),
					testAccCheckOutscaleLBUAttrExists("outscale_load_balancer_attributes.bar2", &conf),
				)},
			{
				Config: testAccOutscaleLBUAttrConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckState("outscale_load_balancer_attributes.bar2"),
					testAccCheckOutscaleLBUAttrExists("outscale_load_balancer_attributes.bar2", &conf),
				)},
		},
	})
}

func testAccCheckOutscaleLBUAttrExists(n string, res *lbu.LoadBalancerAttributes) resource.TestCheckFunc {
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

func testAccOutscaleLBUAttrConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-%d"
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

resource "outscale_load_balancer_attributes" "bar2" {
	access_log_enabled = false
	access_log_s3_bucket_name = "donustestbucket"
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`, r)
}

func testAccOutscaleLBUAttrConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-%d"
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

resource "outscale_load_balancer_attributes" "bar2" {
	access_log_enabled = true
	access_log_s3_bucket_name = "donustestbucket"
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`, r)
}
