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

func TestAccOutscaleOAPIPublicIPLink_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	var a fcu.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIPublicIPLinkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPLExists(
						"outscale_public_ip.bar", &a),
					testAccCheckOutscaleOAPIPublicIPLinkExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPLinkExists(name string, res *fcu.Address) resource.TestCheckFunc {
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

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkExists (%s)", err)

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

func testAccCheckOutscaleOAPIPublicIPLinkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip_link" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		fmt.Printf("%#v", rs.Primary.Attributes)

		id := rs.Primary.Attributes["link_id"]

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

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkDestroy (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.Addresses) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIPublicIPLExists(n string, res *fcu.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		if strings.Contains(rs.Primary.ID, "reservation") {
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

const testAccOutscaleOAPIPublicIPLinkConfig = `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	subnet_id = "subnet-861fbecc"
}

resource "outscale_public_ip" "bar" {}

resource "outscale_public_ip_link" "by_public_ip" {
	public_ip = "${outscale_public_ip.bar.public_ip}"
	vm_id = "${outscale_vm.basic.id}"
  depends_on = ["outscale_vm.basic", "outscale_public_ip.bar"]
}`
