package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscalePublicIP_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var conf fcu.Address

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},
		},
	})
}

func TestAccOutscalePublicIP_instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var conf fcu.Address

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},

			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},
		},
	})
}

// // This test is an expansion of TestAccOutscalePublicIP_instance, by testing the
// // associated Private PublicIPs of two instances
func TestAccOutscalePublicIP_associated_user_private_ip(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var one fcu.Address
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig_associated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},

			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig_associated_switch(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},
		},
	})
}

func testAccCheckOutscalePublicIPDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip" {
			continue
		}

		if strings.Contains(rs.Primary.ID, "allocation") {
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

func testAccCheckOutscalePublicIPAttributes(conf *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.PublicIp == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckOutscalePublicIPExists(n string, res *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		if strings.Contains(rs.Primary.ID, "allocation") {
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

const testAccOutscalePublicIPConfig = `
resource "outscale_public_ip" "bar" {}
`

func testAccOutscalePublicIPInstanceConfig(r int) string {
	return fmt.Sprintf(`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}
resource "outscale_public_ip" "bar" {}
`, r)
}

func testAccOutscalePublicIPInstanceConfig2(r int) string {
	return fmt.Sprintf(`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}
resource "outscale_public_ip" "bar" {}
`, r)
}

func testAccOutscalePublicIPInstanceConfig_associated(r int) string {
	return fmt.Sprintf(`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	cidr_block = "10.0.0.0/16"
	vpc_id = "${outscale_lin.vpc.id}"
}

resource "outscale_vm" "foo" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
  private_ip_address = "10.0.0.12"
  subnet_id  = "${outscale_subnet.subnet.id}"
}
resource "outscale_vm" "bar" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
  private_ip_address = "10.0.0.19"
  subnet_id  = "${outscale_subnet.subnet.id}"
}
resource "outscale_public_ip" "bar" {}
`, r)
}

func testAccOutscalePublicIPInstanceConfig_associated_switch(r int) string {
	return fmt.Sprintf(`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	cidr_block = "10.0.0.0/16"
	vpc_id = "${outscale_lin.vpc.id}"
}

resource "outscale_vm" "foo" {
 image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
  private_ip_address = "10.0.0.12"
  subnet_id  = "${outscale_subnet.subnet.id}"
}
resource "outscale_vm" "bar" {
  image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
  private_ip_address = "10.0.0.19"
  subnet_id  = "${outscale_subnet.subnet.id}"
}
resource "outscale_public_ip" "bar" {}
`, r)
}
