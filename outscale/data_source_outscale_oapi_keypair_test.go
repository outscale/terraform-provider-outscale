package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIKeypairDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeypairDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeypairDataSourceID("data.outscale_keypair.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypair.nat_ami", "keypair_name", "TestKey"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIKeypairDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find key pair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Key Pair data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPIKeypairDataSourceConfig = `
data "outscale_keypair" "nat_ami" {
	keypair_name = "TestKey"
}
`
