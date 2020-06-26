package outscale

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIPublicIPLink_basic(t *testing.T) {
	var a oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIPublicIPLinkConfig(omi, "tinav4.c2r2p2", region, keypair, sgId),
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

func testAccCheckOutscaleOAPIPublicIPLinkExists(name string, res *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		request := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				LinkPublicIpIds: &[]string{res.GetLinkPublicIpId()},
			},
		}
		response, _, err := conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(request)})

		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkExists (%s)", err)
			return err
		}

		//Missing on Swagger Spec
		if len(response.GetPublicIps()) != 1 ||
			response.GetPublicIps()[0].GetLinkPublicIpId() != res.GetLinkPublicIpId() {
			return fmt.Errorf("Public IP Link not found")
		}

		if len(response.GetPublicIps()) != 1 {
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

		request := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				LinkPublicIpIds: &[]string{id},
			},
		}
		response, _, err := conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(request)})

		log.Printf("[DEBUG] ERROR testAccCheckOutscaleOAPIPublicIPLinkDestroy (%s)", err)

		if err != nil {
			return err
		}

		if len(response.GetPublicIps()) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIPublicIPLExists(n string, res *oscgo.PublicIp) resource.TestCheckFunc {
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
			req := oscgo.ReadPublicIpsRequest{
				Filters: &oscgo.FiltersPublicIp{
					LinkPublicIpIds: &[]string{rs.Primary.ID},
				},
			}
			resp, _, err := conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(req)})

			if err != nil {
				return err
			}

			if len(resp.GetPublicIps()) != 1 ||
				resp.GetPublicIps()[0].GetLinkPublicIpId() != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = resp.GetPublicIps()[0]

		} else {
			req := oscgo.ReadPublicIpsRequest{
				Filters: &oscgo.FiltersPublicIp{
					PublicIps: &[]string{rs.Primary.ID},
				},
			}

			var response oscgo.ReadPublicIpsResponse
			err := resource.Retry(120*time.Second, func() *resource.RetryError {
				var err error
				response, _, err = conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(req)})

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

			if len(response.GetPublicIps()) != 1 ||
				response.GetPublicIps()[0].GetPublicIp() != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = response.GetPublicIps()[0]
		}

		return nil
	}
}

func testAccOutscaleOAPIPublicIPLinkConfig(omi, vmType, region, keypair, sgId string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = "${outscale_net.net.id}"
		}

		resource "outscale_vm" "vm" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = ["%[5]s"]
			placement_subregion_name = "%[3]sa"
		}
		
		resource "outscale_public_ip" "ip" {}
		
		resource "outscale_public_ip_link" "by_public_ip" {
			public_ip = "${outscale_public_ip.ip.public_ip}"
			vm_id     = "${outscale_vm.vm.id}"
		}
	`, omi, vmType, region, keypair, sgId)
}
