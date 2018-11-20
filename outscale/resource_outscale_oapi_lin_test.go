package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPILin_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
	var conf oapi.Nets

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinExists("outscale_net.vpc", &conf),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPILinExists(n string, res *oapi.Nets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp *oapi.POST_ReadNetsResponses
		conn := testAccProvider.Meta().(*OutscaleClient)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.OAPI.POST_ReadNets(oapi.ReadNetsRequest{
				Filters: oapi.Filters_6{NetIds: []string{rs.Primary.ID}},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return err
		}

		if resp.OK == nil {
			return fmt.Errorf("Net not found")
		}

		if len(resp.OK.Nets) != 1 ||
			resp.OK.Nets[0].NetId != rs.Primary.ID {
			return fmt.Errorf("Net not found")
		}

		*res = resp.OK.Nets[0]

		return nil
	}
}

//Missing on Swagger Spec
// func testAccCheckOutscaleOAPILinDestroyed(s *terraform.State) error {
// 	conn := testAccProvider.Meta().(*OutscaleClient)

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "outscale_net" {
// 			continue
// 		}

// 		// Try to find an internet gateway
// 		var resp *oapi.ReadGate
// 		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			var err error
// 			resp, err = conn.FCU.VM.DescribeInternetGateways(&fcu.DescribeInternetGatewaysInput{
// 				InternetGatewayIds: []*string{aws.String(rs.Primary.ID)},
// 			})

// 			if err != nil {
// 				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
// 					return resource.RetryableError(err)
// 				}
// 				return resource.NonRetryableError(err)
// 			}

// 			return resource.RetryableError(err)
// 		})

// 		if resp == nil {
// 			return nil
// 		}

// 		if err == nil {
// 			if len(resp.InternetGateways) > 0 {
// 				return fmt.Errorf("still exist")
// 			}
// 			return nil
// 		}

// 		// Verify the error is what we want
// 		ec2err, ok := err.(awserr.Error)
// 		if !ok {
// 			return err
// 		}
// 		if ec2err.Code() != "InvalidVPC.NotFound" {
// 			return err
// 		}
// 	}

// 	return nil
// }

const testAccOutscaleOAPILinConfig = `
resource "outscale_net" "vpc" {
	ip_range = "10.0.0.0/16"

	tags {
		key = "Name" 
		value = "outscale_net"
	}	
}
`
