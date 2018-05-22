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

func TestAccOutscaleOAPIENI_basic(t *testing.T) {
	var conf fcu.NetworkInterface

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIENIExists(n string, res *fcu.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		dnir := &fcu.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{aws.String(rs.Primary.ID)},
		}

		var describeResp *fcu.DescribeNetworkInterfacesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.VM.DescribeNetworkInterfaces(dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(describeResp.NetworkInterfaces) != 1 ||
			*describeResp.NetworkInterfaces[0].NetworkInterfaceId != rs.Primary.ID {
			return fmt.Errorf("ENI not found")
		}

		*res = *describeResp.NetworkInterfaces[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIENIAttributes(conf *fcu.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if conf.Attachment != nil {
			return fmt.Errorf("expected attachment to be nil")
		}

		if *conf.AvailabilityZone != "eu-west-2a" {
			return fmt.Errorf("expected availability_zone to be eu-west-2a, but was %s", *conf.AvailabilityZone)
		}

		return nil
	}
}

const testAccOutscaleOAPIENIConfig = `
resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    lin_id              = "${outscale_lin.outscale_lin.lin_id}"
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

`
