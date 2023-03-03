package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIServerCertificate_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_server_certificate.test"
	rName := acctest.RandomWithPrefix("acc-test")
	rNameUpdated := acctest.RandomWithPrefix("acc-test")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:         testAccProviders,
		ExternalProviders: providerScottwinklerShell(),
		CheckDestroy:      testAccCheckOutscaleServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIServerCertificateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleServerCertificateExists(resourceName),
				),
			},
			{
				Config: testAccOutscaleOAPIServerCertificateConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleServerCertificateExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id", "body", "private_key"},
			},
		},
	})
}

func testAccCheckOutscaleServerCertificateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set")
		}
		exists := false
		var resp oscgo.ReadServerCertificatesResponse
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(oscgo.ReadServerCertificatesRequest{}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetServerCertificates()) == 0 {
			return fmt.Errorf("Server Certificate not found (%s)", rs.Primary.ID)
		}

		for _, server := range resp.GetServerCertificates() {
			if server.GetId() == rs.Primary.ID {
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("Server Certificate not found (%s)", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckOutscaleServerCertificateDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_server_certificate_link" {
			continue
		}

		exists := false

		var resp oscgo.ReadServerCertificatesResponse
		var err error
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(oscgo.ReadServerCertificatesRequest{}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return fmt.Errorf("Server Certificate reading (%s)", rs.Primary.ID)
		}

		for _, server := range resp.GetServerCertificates() {
			if server.GetId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("Server Certificate still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleOAPIServerCertificateConfig(name string) string {
	return fmt.Sprintf(`

	resource "shell_script" "ca_gen" {
		lifecycle_commands {
			create = <<-EOF
				openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout resource_ca.key -days 2 -out resource_ca.pem -subj '/CN=domain.com'
			EOF
			read   = <<-EOF
				echo "{\"certfile\":  \"resource_ca.pem\",
				       \"keyfile\":  \"resource_ca.key\"}"
			EOF
			delete = "rm -f resource_ca.pem resource_ca.key"
		}
		working_directory = "${path.module}/."
	}
	
	resource "outscale_server_certificate" "test" { 
		name        = "%s"
		body        =  file(shell_script.ca_gen.output.certfile)
		private_key =  file(shell_script.ca_gen.output.keyfile)
	}
	`, name)
}
