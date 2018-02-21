package outscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscalePublicIP_basic(t *testing.T) {
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
	var conf fcu.Address

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},

			resource.TestStep{
				Config: testAccOutscalePublicIPInstanceConfig2,
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
// func TestAccOutscalePublicIP_associated_user_private_ip(t *testing.T) {
// 	var one fcu.Address

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:      func() { testAccPreCheck(t) },
// 		IDRefreshName: "outscale_public_ip.bar",
// 		Providers:     testAccProviders,
// 		CheckDestroy:  testAccCheckOutscalePublicIPDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: testAccOutscalePublicIPInstanceConfig_associated,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
// 					testAccCheckOutscalePublicIPAttributes(&one),
// 					testAccCheckOutscalePublicIPAssociated(&one),
// 				),
// 			},

// 			resource.TestStep{
// 				Config: testAccOutscalePublicIPInstanceConfig_associated_switch,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
// 					testAccCheckOutscalePublicIPAttributes(&one),
// 					testAccCheckOutscalePublicIPAssociated(&one),
// 				),
// 			},
// 		},
// 	})
// }

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
			describe, err := conn.FCU.VM.DescribeAddressesRequest(req)
			if err != nil {
				// Verify the error is what we want
				if ae, ok := err.(awserr.Error); ok && ae.Code() == "InvalidAllocationID.NotFound" || ae.Code() == "InvalidAddress.NotFound" {
					continue
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
			describe, err := conn.FCU.VM.DescribeAddressesRequest(req)
			if err != nil {
				// Verify the error is what we want
				if ae, ok := err.(awserr.Error); ok && ae.Code() == "InvalidAllocationID.NotFound" || ae.Code() == "InvalidAddress.NotFound" {
					continue
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

func testAccCheckOutscalePublicIPAssociated(conf *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conf.AssociationId == nil || *conf.AssociationId == "" {
			return fmt.Errorf("empty association_id")
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
			describe, err := conn.FCU.VM.DescribeAddressesRequest(req)
			if err != nil {
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
resource "outscale_public_ip" "bar" {
}
`

const testAccOutscalePublicIPInstanceConfig = `
resource "outscale_vm" "foo" {
	# us-west-2
	ami = "ami-4fccb37f"
	instance_type = "m1.small"
}
resource "outscale_public_ip" "bar" {
	instance = "${outscale_vm.foo.id}"
}
`

const testAccOutscalePublicIPInstanceConfig2 = `
resource "outscale_vm" "bar" {
	# us-west-2
	ami = "ami-4fccb37f"
	instance_type = "m1.small"
}
resource "outscale_public_ip" "bar" {
	instance = "${outscale_vm.bar.id}"
}
`

// const testAccOutscalePublicIPInstanceConfig_associated = `
// resource "aws_vpc" "default" {
//   cidr_block           = "10.0.0.0/16"
//   enable_dns_hostnames = true
//   tags {
//     Name = "default"
//   }
// }
// resource "aws_internet_gateway" "gw" {
//   vpc_id = "${aws_vpc.default.id}"
//   tags {
//     Name = "main"
//   }
// }
// resource "aws_subnet" "tf_test_subnet" {
//   vpc_id                  = "${aws_vpc.default.id}"
//   cidr_block              = "10.0.0.0/24"
//   map_public_ip_on_launch = true
//   depends_on = ["aws_internet_gateway.gw"]
//   tags {
//     Name = "tf_test_subnet"
//   }
// }
// resource "outscale_vm" "foo" {
//   # us-west-2
//   ami           = "ami-5189a661"
//   instance_type = "t2.micro"
//   private_ip = "10.0.0.12"
//   subnet_id  = "${aws_subnet.tf_test_subnet.id}"
//   tags {
//     Name = "foo instance"
//   }
// }
// resource "outscale_vm" "bar" {
//   # us-west-2
//   ami = "ami-5189a661"
//   instance_type = "t2.micro"
//   private_ip = "10.0.0.19"
//   subnet_id  = "${aws_subnet.tf_test_subnet.id}"
//   tags {
//     Name = "bar instance"
//   }
// }
// resource "outscale_public_ip" "bar" {
//   vpc = true
//   instance                  = "${outscale_vm.bar.id}"
//   associate_with_private_ip = "10.0.0.19"
// }
// `
// const testAccOutscalePublicIPInstanceConfig_associated_switch = `
// resource "aws_vpc" "default" {
//   cidr_block           = "10.0.0.0/16"
//   enable_dns_hostnames = true
//   tags {
//     Name = "default"
//   }
// }
// resource "aws_internet_gateway" "gw" {
//   vpc_id = "${aws_vpc.default.id}"
//   tags {
//     Name = "main"
//   }
// }
// resource "aws_subnet" "tf_test_subnet" {
//   vpc_id                  = "${aws_vpc.default.id}"
//   cidr_block              = "10.0.0.0/24"
//   map_public_ip_on_launch = true
//   depends_on = ["aws_internet_gateway.gw"]
//   tags {
//     Name = "tf_test_subnet"
//   }
// }
// resource "outscale_vm" "foo" {
//   # us-west-2
//   ami           = "ami-5189a661"
//   instance_type = "t2.micro"
//   private_ip = "10.0.0.12"
//   subnet_id  = "${aws_subnet.tf_test_subnet.id}"
//   tags {
//     Name = "foo instance"
//   }
// }
// resource "outscale_vm" "bar" {
//   # us-west-2
//   ami = "ami-5189a661"
//   instance_type = "t2.micro"
//   private_ip = "10.0.0.19"
//   subnet_id  = "${aws_subnet.tf_test_subnet.id}"
//   tags {
//     Name = "bar instance"
//   }
// }
// resource "outscale_public_ip" "bar" {
//   vpc = true
//   instance                  = "${outscale_vm.foo.id}"
//   associate_with_private_ip = "10.0.0.12"
// }
// `

// const testAccOutscalePublicIPInstanceConfig_associated_update = `
// resource "outscale_vm" "bar" {
// 	# us-west-2
// 	ami = "ami-4fccb37f"
// 	instance_type = "m1.small"
// }
// resource "outscale_public_ip" "bar" {
// 	instance = "${outscale_vm.bar.id}"
// }
// `
