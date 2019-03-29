package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscalePublicIPLink_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var a fcu.Address
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscalePublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPLinkConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists(
						"outscale_public_ip.bar", &a),
					testAccCheckOutscalePublicIPLinkExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func testAccCheckEIPLinkDisappears(address *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient)
		opts := &fcu.DisassociateAddressInput{
			AssociationId: address.AssociationId,
		}
		if _, err := conn.FCU.VM.DisassociateAddress(opts); err != nil {
			fmt.Printf("\n [DEBUG] ERROR testAccCheckEIPLinkDisappears (%s)", err)

			return err
		}
		return nil
	}
}

func testAccCheckOutscalePublicIPLinkExists(name string, res *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf("%#v", s.RootModule().Resources)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
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

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscalePublicIPLinkExists (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.Addresses) != 1 ||
			*describe.Addresses[0].AssociationId != *res.AssociationId {
			return fmt.Errorf("Public IP Link not found")
		}

		return nil
	}
}

func testAccCheckOutscalePublicIPLinkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip_link" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		fmt.Printf("%#v", rs.Primary.Attributes)

		id := rs.Primary.Attributes["association_id"]

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := &fcu.DescribeAddressesInput{
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("association-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}
		describe, err := conn.FCU.VM.DescribeAddressesRequest(request)

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscalePublicIPLinkDestroy (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.Addresses) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccOutscalePublicIPLinkConfig(r int) string {
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

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "${outscale_keypair.a_key_pair.key_name}"
	subnet_id = "${outscale_subnet.subnet.id}"
}

resource "outscale_public_ip" "bar" {
}

resource "outscale_public_ip_link" "by_public_ip" {
	public_ip = "${outscale_public_ip.bar.public_ip}"
	instance_id = "${outscale_vm.basic.id}"
  depends_on = ["outscale_vm.basic", "outscale_public_ip.bar"]
}`, r)
}
