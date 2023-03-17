package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Ca_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_ca.ca_data"
	dataSourcesName := "data.outscale_cas.cas_data"
	dataSourcesAllName := "data.outscale_cas.all_cas"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:         testAccProviders,
		ExternalProviders: providerScottwinklerShell(),
		Steps: []resource.TestStep{
			{
				Config: testAcc_Ca_DataSource_Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "ca_fingerprint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ca_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),

					resource.TestCheckResourceAttrSet(dataSourcesName, "cas.#"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "filter.#"),

					resource.TestCheckResourceAttrSet(dataSourcesAllName, "cas.#"),
				),
			},
		},
	})
}

func testAcc_Ca_DataSource_Config() string {
	return fmt.Sprintf(`

resource "shell_script" "ca_gen" {
	lifecycle_commands {
		create = <<-EOF
			openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout datasource_ca.key -days 2 -out datasource_ca.pem -subj '/CN=domain.com'
		EOF
		read   = <<-EOF
			echo "{\"filename\":  \"datasource_ca.pem\"}"
		EOF
		delete = "rm -f datasource_ca.pem datasource_ca.key"
	}
	working_directory = "${path.module}/."
}

resource "outscale_ca" "ca_test" { 
	ca_pem      = file(shell_script.ca_gen.output.filename)
	description   = "Ca testacc create"
}

data "outscale_ca" "ca_data" { 
	filter {
		name   = "ca_ids"
		values = [outscale_ca.ca_test.id]
	}
}

data "outscale_cas" "cas_data" { 
	filter {
		name   = "ca_ids"
		values = [outscale_ca.ca_test.id]
	}
	filter {
		name   = "descriptions"
		values = [outscale_ca.ca_test.description]
	}
	filter {
		name   = "ca_fingerprints"
		values = [outscale_ca.ca_test.ca_fingerprint]
	}
}

data "outscale_cas" "all_cas" {
	depends_on = [
		outscale_ca.ca_test
	]
}
`)
}
