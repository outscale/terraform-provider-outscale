package outscale

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscalePublicIPAssociation_basic(t *testing.T) {
	var a fcu.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscalePublicIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists(
						"outscale_public_ip.bar.0", &a),
					testAccCheckOutscalePublicIPAssociationExists(
						"outscale_public_ip_link.by_allocation_id", &a),
					testAccCheckOutscalePublicIPExists(
						"outscale_public_ip.bar.1", &a),
					testAccCheckOutscalePublicIPAssociationExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func TestAccOutscalePublicIPAssociation_disappears(t *testing.T) {
	var a fcu.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscalePublicIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPAssociationConfigDisappears,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists(
						"outscale_public_ip.bar", &a),
					testAccCheckOutscalePublicIPAssociationExists(
						"aws_eip_association.by_allocation_id", &a),
					testAccCheckEIPAssociationDisappears(&a),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckEIPAssociationDisappears(address *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient)
		opts := &fcu.DisassociateAddressInput{
			AssociationId: address.AssociationId,
		}
		if _, err := conn.FCU.VM.DisassociateAddress(opts); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckOutscalePublicIPAssociationExists(name string, res *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Association ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := &fcu.DescribeAddressesInput{
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("association-id"),
					Values: []*string{res.AssociationId},
				},
			},
		}
		describe, err := conn.FCU.VM.DescribeAddressesRequest(request)
		if err != nil {
			return err
		}

		if len(describe.Addresses) != 1 ||
			*describe.Addresses[0].AssociationId != *res.AssociationId {
			return fmt.Errorf("Public IP Association not found")
		}

		return nil
	}
}

func testAccCheckOutscalePublicIPAssociationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip_link" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Association ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := &fcu.DescribeAddressesInput{
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("association-id"),
					Values: []*string{aws.String(rs.Primary.ID)},
				},
			},
		}
		describe, err := conn.FCU.VM.DescribeAddressesRequest(request)
		if err != nil {
			return err
		}

		if len(describe.Addresses) > 0 {
			return fmt.Errorf("Public IP Association still exists")
		}
	}
	return nil
}

const testAccOutscalePublicIPAssociationConfig = `
resource "outscale_vm" "basic" {
	count = 2
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	subnet_id = "subnet-861fbecc"
}

resource "outscale_public_ip" "bar" {
	count = 2
}

resource "outscale_public_ip_link" "by_allocation_id" {
	allocation_id = "${outscale_public_ip.bar.0.id}"
	instance_id = "${outscale_vm.basic.0.id}"
	depends_on = ["outscale_vm.basic"]
}

resource "outscale_public_ip_link" "by_public_ip" {
	public_ip = "${outscale_public_ip.bar.1.public_ip}"
	instance_id = "${outscale_vm.basic.1.id}"
  depends_on = ["outscale_vm.basic"]
}`

const testAccOutscalePublicIPAssociationConfigDisappears = `
resource "outscale_vm" "foo" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	subnet_id = "subnet-861fbecc"
}
resource "outscale_public_ip" "bar" {
}
resource "outscale_public_ip_link" "by_allocation_id" {
	allocation_id = "${outscale_public_ip.bar.id}"
	instance_id = "${outscale_vm.foo.id}"
}`
