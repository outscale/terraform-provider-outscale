package oapi_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_WithPublicIPLink_basic(t *testing.T) {
	var a osc.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
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

func testAccCheckOutscalePublicIPLinkExists(name string, res *osc.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no public ip link id is set")
		}

		client := testacc.ConfiguredClient.OSC

		request := osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				LinkPublicIpIds: &[]string{ptr.From(res.LinkPublicIpId)},
			},
		}
		response, err := client.ReadPublicIps(context.Background(), request, options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscalePublicIPLinkExists (%s)", err)
			return err
		}

		// Missing on Swagger Spec
		if response.PublicIps == nil || len(*response.PublicIps) != 1 ||
			ptr.From((*response.PublicIps)[0].LinkPublicIpId) != ptr.From(res.LinkPublicIpId) {
			return fmt.Errorf("public ip link not found")
		}

		if len(*response.PublicIps) != 1 {
			return fmt.Errorf("public ip link not found")
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
			return fmt.Errorf("no public ip link id is set")
		}

		id := rs.Primary.Attributes["link_public_ip_id"]

		client := testacc.ConfiguredClient.OSC

		request := osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				LinkPublicIpIds: &[]string{id},
			},
		}
		response, err := client.ReadPublicIps(context.Background(), request, options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			log.Printf("[DEBUG] ERROR testAccCheckOutscalePublicIPLinkDestroy (%s)", err)
			return err
		}

		if len(ptr.From(response.PublicIps)) > 0 {
			return fmt.Errorf("public ip link still exists")
		}
	}
	return nil
}

func testAccCheckOutscalePublicIPLExists(n string, res *osc.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no publicip id is set")
		}

		client := testacc.ConfiguredClient.OSC

		// Missing on Swagger Spec
		if strings.Contains(rs.Primary.ID, "reservation") {
			req := osc.ReadPublicIpsRequest{
				Filters: &osc.FiltersPublicIp{
					LinkPublicIpIds: &[]string{rs.Primary.ID},
				},
			}
			resp, err := client.ReadPublicIps(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))
			if err != nil {
				return err
			}

			if resp.PublicIps == nil || len(*resp.PublicIps) != 1 ||
				ptr.From((*resp.PublicIps)[0].LinkPublicIpId) != rs.Primary.ID {
				return fmt.Errorf("publicip not found")
			}
			*res = (*resp.PublicIps)[0]

		} else {
			req := osc.ReadPublicIpsRequest{
				Filters: &osc.FiltersPublicIp{
					PublicIpIds: &[]string{rs.Primary.ID},
				},
			}

			response, err := client.ReadPublicIps(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))
			if err != nil {
				return err
			}

			ips := ptr.From(response.PublicIps)
			if len(ips) != 1 || ips[0].PublicIpId != rs.Primary.ID {
				return fmt.Errorf("publicip not found")
			}
			*res = ips[0]
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
