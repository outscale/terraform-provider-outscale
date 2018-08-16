package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPISubNet_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	var conf fcu.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleLinDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubNetExists("outscale_subnet.subnet", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISubNetExists(n string, res *fcu.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet id is set")
		}
		var resp *fcu.DescribeSubnetsOutput
		conn := testAccProvider.Meta().(*OutscaleClient)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.FCU.VM.DescribeSubNet(&fcu.DescribeSubnetsInput{
				SubnetIds: []*string{aws.String(rs.Primary.ID)},
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
		if len(resp.Subnets) != 1 ||
			*resp.Subnets[0].SubnetId != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*res = *resp.Subnets[0]

		return nil
	}
}

func testAccCheckOutscaleOAPISubNetDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_subnet" {
			continue
		}

		// Try to find an internet gateway
		var resp *fcu.DescribeSubnetsOutput
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.FCU.VM.DescribeSubNet(&fcu.DescribeSubnetsInput{
				SubnetIds: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		})

		if resp == nil {
			return nil
		}

		if err == nil {
			if len(resp.Subnets) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidSubnet.NotFound" {
			return err
		}
	}

	return nil
}

const testAccOutscaleOAPISubnetConfig = `
resource "outscale_net" "vpc" {
	ip_range = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	ip_range = "10.0.0.0/16"
	lin_id = "${outscale_lin.vpc.id}"
}

`
