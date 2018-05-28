package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleLoadBalancerAccessLogs_basic(t *testing.T) {
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
				Config: testAccDSOutscaleLBUDSAccessLogsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUDSAccessLogsExists("data.outscale_load_balancer_access_logs.test"),
					resource.TestCheckResourceAttr(
						"data.outscale_load_balancer_access_logs.test", "enabled", "false"),
				)},
		},
	})
}

func testAccCheckOutscaleLBUDSAccessLogsExists(n string) resource.TestCheckFunc {
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

const testAccDSOutscaleLBUDSAccessLogsConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones = ["eu-west-2a"]
	load_balancer_name               = "foobar-terraform-elb-ds"
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

data "outscale_load_balancer_access_logs" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
}
`
