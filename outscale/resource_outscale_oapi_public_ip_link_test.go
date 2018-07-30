package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPIPublicIPLink_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	var a oapi.PublicIps

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIPublicIPLinkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckState("outscale_public_ip_link.by_public_ip"),
					testAccCheckOutscaleOAPIPublicIPLExists(
						"outscale_public_ip.bar", &a),
					testAccCheckOutscaleOAPIPublicIPLinkExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPLinkExists(name string, res *oapi.PublicIps) resource.TestCheckFunc {
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

		request := oapi.ReadPublicIpsRequest{
			Filters: oapi.ReadPublicIpsFilters{
				LinkIds: []string{res.LinkId},
			},
		}
		describe, err := conn.OAPI.POST_ReadPublicIps(request)

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkExists (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.OK.PublicIps) != 1 ||
			describe.OK.PublicIps[0].ReservationId != res.ReservationId {
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

		request := oapi.ReadPublicIpsRequest{
			Filters: oapi.ReadPublicIpsFilters{
				LinkIds: []string{id},
			},
		}
		describe, err := conn.OAPI.POST_ReadPublicIps(request)

		fmt.Printf("\n [DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkDestroy (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.OK.PublicIps) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIPublicIPLExists(n string, res *oapi.PublicIps) resource.TestCheckFunc {
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
			req := oapi.ReadPublicIpsRequest{
				Filters: oapi.ReadPublicIpsFilters{
					ReservationIds: []string{rs.Primary.ID},
				},
			}
			resp, err := conn.OAPI.POST_ReadPublicIps(req)

			if err != nil {
				return err
			}

			describe := resp.OK

			if len(describe.PublicIps) != 1 ||
				describe.PublicIps[0].ReservationId != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = describe.PublicIps[0]

		} else {
			req := oapi.ReadPublicIpsRequest{
				Filters: oapi.ReadPublicIpsFilters{
					PublicIps: []string{rs.Primary.ID},
				},
			}

			var describe *oapi.ReadPublicIpsResponse
			err := resource.Retry(120*time.Second, func() *resource.RetryError {
				var err error
				resp, err := conn.OAPI.POST_ReadPublicIps(req)

				if err != nil {
					if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}
				describe = resp.OK
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

			if len(describe.PublicIps) != 1 ||
				describe.PublicIps[0].PublicIp != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = describe.PublicIps[0]
		}

		return nil
	}
}

const testAccOutscaleOAPIPublicIPLinkConfig = `
#resource "outscale_vm" "basic" {
#	image_id = "ami-8a6a0120"
#	instance_type = "t2.micro"
#	key_name = "terraform-basic"
#	subnet_id = "subnet-861fbecc"
#}

resource "outscale_public_ip" "bar" {}

resource "outscale_public_ip_link" "by_public_ip" {
	public_ip = "${outscale_public_ip.bar.public_ip}"
	#vm_id = "${outscale_vm.basic.id}"
	vm_id = "i-ccdf0eeb"
	#depends_on = ["outscale_vm.basic", "outscale_public_ip.bar"]
}`
