package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleOAPIDSLBUAttr_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 10)

	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSOutscaleOAPILBUAttrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					resource.TestCheckResourceAttr("data.outscale_load_balancer_attributes.test", "load_balancer_attributes.0.access_log.enabled", "false"),
				)},
		},
	})
}

func testAccDSOutscaleOAPILBUAttrConfig(r int) string {
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
	enabled = "false"
			s3_bucket_name = "donustestbucket"
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}

data "outscale_load_balancer_attributes" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`, r)
}
