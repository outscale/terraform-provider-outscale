package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleLinAPIAccess_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var conf fcu.VpcEndpoint

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinAPIAccessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleLinAPIAccessConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinAPIAccessExists("outscale_lin_api_access.link", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleLinAPIAccessExists(n string, res *fcu.VpcEndpoint) resource.TestCheckFunc {
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

func testAccCheckOutscaleLinAPIAccessDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_lin_api_access" {
			continue
		}

		id := rs.Primary.Attributes["vpc_id"]

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

const testAccOutscaleLinAPIAccessConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_lin_api_access" "link" {
	vpc_id = "${outscale_lin.foo.id}"
	route_table_id = [
		"${outscale_route_table.foo.id}"
	]
	service_name = "com.outscale.eu-west-2.osu"
}
`
