package outscale

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPIPublicIPLink_basic(t *testing.T) {
	var a oapi.PublicIp
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIPublicIPLinkConfig(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPLExists(
						"outscale_public_ip.ip", &a),
					testAccCheckOutscaleOAPIPublicIPLinkExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPLinkExists(name string, res *oapi.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := oapi.ReadPublicIpsRequest{
			Filters: oapi.FiltersPublicIp{
				LinkPublicIpIds: []string{res.LinkPublicIpId},
			},
		}
		describe, err := conn.OAPI.POST_ReadPublicIps(request)

		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkExists (%s)", err)
			return err
		}

		//Missing on Swagger Spec
		if len(describe.OK.PublicIps) != 1 ||
			describe.OK.PublicIps[0].LinkPublicIpId != res.LinkPublicIpId {
			return fmt.Errorf("Public IP Link not found")
		}

		if len(describe.OK.PublicIps) != 1 {
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

		id := rs.Primary.Attributes["link_id"]

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := oapi.ReadPublicIpsRequest{
			Filters: oapi.FiltersPublicIp{
				LinkPublicIpIds: []string{id},
			},
		}
		describe, err := conn.OAPI.POST_ReadPublicIps(request)

		log.Printf("[DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkDestroy (%s)", err)

		if err != nil {
			return err
		}

		if len(describe.OK.PublicIps) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIPublicIPLExists(n string, res *oapi.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		// Missing on Swagger Spec
		if strings.Contains(rs.Primary.ID, "reservation") {
			req := oapi.ReadPublicIpsRequest{
				Filters: oapi.FiltersPublicIp{
					LinkPublicIpIds: []string{rs.Primary.ID},
				},
			}
			resp, err := conn.OAPI.POST_ReadPublicIps(req)

			if err != nil {
				return err
			}

			describe := resp.OK

			if len(describe.PublicIps) != 1 ||
				describe.PublicIps[0].LinkPublicIpId != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = describe.PublicIps[0]

		} else {
			req := oapi.ReadPublicIpsRequest{
				Filters: oapi.FiltersPublicIp{
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

func testAccOutscaleOAPIPublicIPLinkConfig(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}
		
		resource "outscale_public_ip" "ip" {}
		
		resource "outscale_public_ip_link" "by_public_ip" {
			public_ip = "${outscale_public_ip.ip.public_ip}"
			vm_id     = "${outscale_vm.vm.id}"
		}
	`, omi, vmType, region)
}
