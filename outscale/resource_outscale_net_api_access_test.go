package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleLinOAPIAccess_basic(t *testing.T) {
	t.Skip()
	var conf fcu.VpcEndpoint

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinOAPIAccessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleLinOAPIAccessConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinOAPIAccessExists("outscale_net_api_access.link", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleLinOAPIAccessExists(n string, res *fcu.VpcEndpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No api_access id is set")
		}

		return nil
	}
}

func testAccCheckOutscaleLinOAPIAccessDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_net_api_access" {
			continue
		}

		id := rs.Primary.Attributes["net_id"]

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.FCU.VM.DescribeVpcEndpoints(&fcu.DescribeVpcEndpointsInput{
				VpcEndpointIds: []*string{aws.String(id)},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if strings.Contains(fmt.Sprint(err), "InvalidVpcEndpointId.NotFound") {
			return nil
		}

		if err != nil {
			return err
		}

	}
	return nil
}

const testAccOutscaleLinOAPIAccessConfig = `
resource "outscale_net" "foo" {
	ip_ranges = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"
}

resource "outscale_net_api_access" "link" {
	net_id = "${outscale_net.foo.id}"
	route_table_id = [
		"${outscale_route_table.foo.id}"
	]
	prefix_list_name = "com.outscale.eu-west-2.osu"
}
`
