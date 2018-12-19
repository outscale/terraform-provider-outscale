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
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIENIDataSource_basic(t *testing.T) {
	var conf fcu.NetworkInterface

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIDataSourceExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIDataSourceAttributes(&conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIENIDataSourceExists(n string, res *fcu.NetworkInterface) resource.TestCheckFunc {
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

func testAccCheckOutscaleOAPIENIDataSourceAttributes(conf *fcu.NetworkInterface) resource.TestCheckFunc {
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

func testAccCheckOutscaleOAPIENIDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		dnir := &fcu.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{aws.String(rs.Primary.ID)},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.VM.DescribeNetworkInterfaces(dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccCheckOutscaleOAPINICDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		dnir := &oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{NicIds: []string{rs.Primary.ID}},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(*dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || describeResp.OK == nil {
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
					return nil
				}
				errString = err.Error()
			} else if describeResp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
			} else if describeResp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
			} else if describeResp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
			}
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		if len(describeResp.OK.Nics) > 0 {
			return fmt.Errorf("Nic with id %s is not destroyed yet", rs.Primary.ID)
		}
	}

	return nil
}

const testAccOutscaleOAPIENIDataSourceConfig = `
resource "outscale_net" "outscale_net" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_net.outscale_net.vpc_id}"
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

data "outscale_nic" "nic" {
		network_interface_id = "NICID"
		subnet_id = "1"
}
`
