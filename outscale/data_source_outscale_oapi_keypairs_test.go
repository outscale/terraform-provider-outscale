package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIKeypairsDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIKeypairsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeypairsDataSourceID("data.outscale_keypairs.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypairs.nat_ami", "keypairs.0.keypair_name", "TestKey"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIKeypairsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Key Pair data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPIKeypairsDataSourceConfig = `
resource "outscale_keypair" "a_key_pair" {
	keypair_name   = "TestKey"
}

data "outscale_keypairs" "nat_ami" {
	keypair_names = ["${outscale_keypair.a_key_pair.id}"]
}
`
