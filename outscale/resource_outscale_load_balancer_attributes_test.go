package outscale

import (
	"fmt"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_LBUAttr_basic(t *testing.T) {
	var conf oscgo.AccessLog
	const (
		MIN_LB_NAME_SUFFIX int = 1
		MAX_LB_NAME_SUFFIX int = 1000
	)
	suffix := utils.RandIntRange(MIN_LB_NAME_SUFFIX, MAX_LB_NAME_SUFFIX)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer_attributes.bar2",
		Providers:     testAccProviders,
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

func testAccCheckOutscaleLBUAttrExists(n string, res *oscgo.AccessLog) resource.TestCheckFunc {
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
