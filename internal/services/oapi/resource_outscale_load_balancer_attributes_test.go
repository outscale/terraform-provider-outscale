package oapi_test

import (
	"fmt"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_LBUAttr_basic(t *testing.T) {
	var conf osc.AccessLog
	const (
		MIN_LB_NAME_SUFFIX int = 1
		MAX_LB_NAME_SUFFIX int = 1000
	)
	suffix := utils.RandIntRange(MIN_LB_NAME_SUFFIX, MAX_LB_NAME_SUFFIX)

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName: "outscale_load_balancer_attributes.bar2",
		Providers:     testacc.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUAttrConfig(utils.GetRegion(), suffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUAttrExists("outscale_load_balancer_attributes.bar2", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleLBUAttrExists(n string, res *osc.AccessLog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no lbu attr id is set")
		}

		return nil
	}
}

func testAccOutscaleLBUAttrConfig(region string, suffix int) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "barAtt" {
  subregion_names = ["%sa"]
  load_balancer_name       = "foobar-terraform-elb-%d"
  listeners {
    backend_port           = 8000
    backend_protocol       = "HTTP"
    load_balancer_port     = 80
    load_balancer_protocol = "HTTP"
  }
  tags {
       key = "test_baz"
       value = "baz"
  }
}

resource "outscale_load_balancer_attributes" "bar2" {
	access_log {
		is_enabled = "false"
		osu_bucket_prefix = "donustestbucket"
	}
	load_balancer_name = outscale_load_balancer.barAtt.id
}
`, region, suffix)
}
