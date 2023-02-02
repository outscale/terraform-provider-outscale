package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAcc_Ca_DataSource(t *testing.T) {
	dataSourceName := "data.outscale_ca.ca_data"
	dataSourcesName := "data.outscale_cas.cas_data"
	dataSourcesAllName := "data.outscale_cas.all_cas"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Ca_DataSource_Config(utils.TestCaPem),
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

func testAcc_Ca_DataSource_Config(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" { 
   ca_pem        =  %[1]q
   description   = "Ca testacc create"
}

data "outscale_ca" "ca_data" { 
  filter {
    name   = "ca_ids"
    values = ["${outscale_ca.ca_test.id}"]
  }
}

data "outscale_cas" "cas_data" { 
  filter {
    name   = "ca_ids"
    values = ["${outscale_ca.ca_test.id}"]
  }
  filter {
	name   = "descriptions"
	values = ["${outscale_ca.ca_test.description}"]
  }
  filter {
	name   = "ca_fingerprints"
	values = ["${outscale_ca.ca_test.ca_fingerprint}"]
  }
}

data "outscale_cas" "all_cas" {}
`, ca_pem)
}
