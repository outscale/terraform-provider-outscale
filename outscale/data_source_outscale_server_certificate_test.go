package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ServerCertificate_Datasource(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("acc-test")
	dataSourceName := "data.outscale_server_certificate.test"
	dataSourcesName := "data.outscale_server_certificates.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:         testAccProviders,
		ExternalProviders: providerScottwinklerShell(),
		Steps: []resource.TestStep{
			{
				Config: testAcc_ServerCertificate_Datasource_Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "server_certificates.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "path"),
				),
			},
		},
	})
}

func testAcc_ServerCertificate_Datasource_Config(name string) string {
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
   path        = "/datasource/"
}

data "outscale_server_certificate" "test" {
	filter {
		name = "paths"
		values = [outscale_server_certificate.test.path]
	}
}

data "outscale_server_certificates" "test" {
	filter {
		name = "paths"
		values = [outscale_server_certificate.test.path]
	}
}
	`, name)
}
