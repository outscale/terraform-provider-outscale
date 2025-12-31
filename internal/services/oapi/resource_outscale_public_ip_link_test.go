package oapi_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_WithPublicIPLink_basic(t *testing.T) {
	var a oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscalePublicIPLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPLinkConfig(omi, testAccVmType, utils.GetRegion(), keypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPLExists(
						"outscale_public_ip.ip_link", &a),
					testAccCheckOutscalePublicIPLinkExists(
						"outscale_public_ip_link.by_public_ip", &a),
				),
			},
		},
	})
}

func testAccCheckOutscalePublicIPLinkExists(name string, res *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		client := testacc.ConfiguredClient.OSCAPI

		request := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				LinkPublicIpIds: &[]string{res.GetLinkPublicIpId()},
			},
		}
		var response oscgo.ReadPublicIpsResponse
		err := retry.Retry(60*time.Second, func() *retry.RetryError {
			rp, httpResp, err := client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			response = rp
			return nil
		})
		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscalePublicIPLinkExists (%s)", err)
			return err
		}

		// Missing on Swagger Spec
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

func testAccCheckOutscalePublicIPLinkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip_link" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Public IP Link ID is set")
		}

		id := rs.Primary.Attributes["link_public_ip_id"]

		client := testacc.ConfiguredClient.OSCAPI

		request := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				LinkPublicIpIds: &[]string{id},
			},
		}
		var response oscgo.ReadPublicIpsResponse
		err := retry.Retry(60*time.Second, func() *retry.RetryError {
			rp, httpResp, err := client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			response = rp
			return nil
		})
		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscalePublicIPLinkDestroy (%s)", err)
			return err
		}

		if len(response.GetPublicIps()) > 0 {
			return fmt.Errorf("Public IP Link still exists")
		}
	}
	return nil
}

func testAccCheckOutscalePublicIPLExists(n string, res *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		client := testacc.ConfiguredClient.OSCAPI

		// Missing on Swagger Spec
		if strings.Contains(rs.Primary.ID, "reservation") {
			req := oscgo.ReadPublicIpsRequest{
				Filters: &oscgo.FiltersPublicIp{
					LinkPublicIpIds: &[]string{rs.Primary.ID},
				},
			}
			var resp oscgo.ReadPublicIpsResponse
			err := retry.Retry(60*time.Second, func() *retry.RetryError {
				rp, httpResp, err := client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				resp = rp
				return nil
			})
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
					PublicIpIds: &[]string{rs.Primary.ID},
				},
			}

			var response oscgo.ReadPublicIpsResponse
			var statusCode int
			err := retry.Retry(120*time.Second, func() *retry.RetryError {
				var err error
				rp, httpResp, err := client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
				if err != nil {
					if httpResp.StatusCode == http.StatusNotFound {
						return retry.RetryableError(err)
					}
					return utils.CheckThrottling(httpResp, err)
				}
				response = rp
				statusCode = httpResp.StatusCode
				return nil
			})
			if err != nil {
				if statusCode == http.StatusNotFound {
					return nil
				}
				return err
			}

			if len(response.GetPublicIps()) != 1 ||
				response.GetPublicIps()[0].GetPublicIpId() != rs.Primary.ID {
				return fmt.Errorf("PublicIP not found")
			}
			*res = response.GetPublicIps()[0]
		}

		return nil
	}
}

func testAccOutscalePublicIPLinkConfig(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_link" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "vm_link" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = [outscale_security_group.sg_link.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_public_ip" "ip_link" {}

		resource "outscale_public_ip_link" "by_public_ip" {
			public_ip = outscale_public_ip.ip_link.public_ip
			vm_id     = outscale_vm.vm_link.id
		}
	`, omi, vmType, region, keypair, sgName)
}
