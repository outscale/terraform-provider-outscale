package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIPublicIP_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	var conf oapi.PublicIps

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIPublicIP_instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	var conf fcu.Address
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},

			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},
		},
	})
}

// // This test is an expansion of TestAccOutscalePublicIP_instance, by testing the
// // associated Private PublicIPs of two instances
func TestAccOutscaleOAPIPublicIP_associated_user_private_ip(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
	var one oapi.PublicIps

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfigAssociated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscaleOAPIPublicIPAttributes(&one),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfigAssociatedSwitch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscaleOAPIPublicIPAttributes(&one),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip" {
			continue
		}

		if strings.Contains(rs.Primary.ID, "reservation") {
			req := &fcu.DescribeAddressesInput{
				AllocationIds: []*string{aws.String(rs.Primary.ID)},
			}

			var describe *fcu.DescribeAddressesOutput
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				var err error
				describe, err = conn.FCU.VM.DescribeAddressesRequest(req)

				return resource.RetryableError(err)
			})

			if err != nil {
				// Verify the error is what we want
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return nil
				}

				return err
			}

			if len(describe.Addresses) > 0 {
				return fmt.Errorf("still exists")
			}
		} else {
			req := &fcu.DescribeAddressesInput{
				PublicIps: []*string{aws.String(rs.Primary.ID)},
			}

			var describe *fcu.DescribeAddressesOutput
			err := resource.Retry(60*time.Second, func() *resource.RetryError {
				var err error
				describe, err = conn.FCU.VM.DescribeAddressesRequest(req)

				return resource.RetryableError(err)
			})

			if err != nil {
				// Verify the error is what we want
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return nil
				}

				return err
			}

			if len(describe.Addresses) > 0 {
				return fmt.Errorf("still exists")
			}
		}
	}

	return nil
}

func testAccCheckOutscaleOAPIPublicIPAttributes(conf *oapi.PublicIps) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.PublicIp == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIPublicIPExists(n string, res *oapi.PublicIps) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		if strings.Contains(rs.Primary.ID, "link") {
			req := &fcu.DescribeAddressesInput{
				AllocationIds: []*string{aws.String(rs.Primary.ID)},
			}
			describe, err := conn.FCU.VM.DescribeAddressesRequest(req)

			if err != nil {
				return err
			}

			if len(describe.Addresses) != 1 ||
				*describe.Addresses[0].AllocationId != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = *describe.Addresses[0]

		} else {
			req := &fcu.DescribeAddressesInput{
				PublicIps: []*string{aws.String(rs.Primary.ID)},
			}

			var describe *fcu.DescribeAddressesOutput
			err := resource.Retry(120*time.Second, func() *resource.RetryError {
				var err error
				describe, err = conn.FCU.VM.DescribeAddressesRequest(req)

				if err != nil {
					if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}

				return nil
			})

			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return nil
				}

				return err
			}

			if err != nil {

				// Verify the error is what we want
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
					return nil
				}

				return err
			}

			if len(describe.Addresses) != 1 ||
				*describe.Addresses[0].PublicIp != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = *describe.Addresses[0]
		}

		return nil
	}
}

const testAccOutscaleOAPIPublicIPConfig = `
resource "outscale_public_ip" "bar" {}
`

const testAccOutscaleOAPIPublicIPInstanceConfig = `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
}
resource "outscale_public_ip" "bar" {}
`

const testAccOutscaleOAPIPublicIPInstanceConfig2 = `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
}
resource "outscale_public_ip" "bar" {}
`

const testAccOutscaleOAPIPublicIPInstanceConfigAssociated = `
resource "outscale_vm" "foo" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
  private_ip_address = "10.0.0.12"
  subnet_id  = "subnet-861fbecc"
}
resource "outscale_vm" "bar" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
  private_ip_address = "10.0.0.19"
  subnet_id  = "subnet-861fbecc"
}
resource "outscale_public_ip" "bar" {}
`

const testAccOutscaleOAPIPublicIPInstanceConfigAssociatedSwitch = `
resource "outscale_vm" "foo" {
 image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
  private_ip_address = "10.0.0.12"
  subnet_id  = "subnet-861fbecc"
}
resource "outscale_vm" "bar" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
  private_ip_address = "10.0.0.19"
  subnet_id  = "subnet-861fbecc"
}
resource "outscale_public_ip" "bar" {}
`
